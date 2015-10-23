package poller

import (
	"fmt"
	"strings"
	"sync"
	"time"

	livestatus "github.com/vbatoufflet/go-livestatus"
	"github.com/vbatoufflet/goverview/pkg/logger"
)

const (
	_ uint = iota
	workerCmdRefresh
	workerCmdShutdown
)

// Worker represents a worker instance in a poller
type Worker struct {
	Catalog    *Catalog
	Config     NodeConfig
	cs         chan uint
	wg         *sync.WaitGroup
	refreshing bool
	stopping   bool
}

// Start starts the worker processing
func (w *Worker) Start() {
	t := time.NewTicker(time.Duration(w.Config.PollInterval) * time.Second)
	tc := t.C

	for {
		select {
		case _ = <-tc:
			w.poll()

		case cmd := <-w.cs:
			switch cmd {
			case workerCmdRefresh:
				w.poll()

			case workerCmdShutdown:
				t.Stop()
				goto stop
			}
		}
	}

stop:
	w.wg.Done()
}

// Shutdown stops the worker
func (w *Worker) Shutdown() {
	if w.stopping {
		return
	}

	logger.Debug("poller", "shutting down %q polling worker", w.Config.Name)

	// Send shutdown command
	w.cs <- workerCmdShutdown
}

// Refresh emits a refresh command on the worker command channel
func (w *Worker) Refresh() {
	w.cs <- workerCmdRefresh
}

func (w *Worker) poll() {
	var (
		l        *livestatus.Livestatus
		q        *livestatus.Query
		resp     *livestatus.Response
		services map[string][]livestatus.Record
		err      error
	)

	// Skip refresh if already in refreshing state
	if w.refreshing {
		logger.Warning("poller", "worker is already refreshing data from %q node", w.Config.Name)
		return
	}

	w.refreshing = true

	logger.Debug("poller", "refreshing data from %q node", w.Config.Name)

	c := NewCatalog(w)

	// Create new LiveStatus instance
	args := strings.SplitN(w.Config.Address, ":", 2)
	if len(args) != 2 {
		logger.Error("poller", "invalid node address %q", w.Config.Address)
		goto end
	}

	l = livestatus.NewLivestatus(args[0], args[1])

	// Query services data
	q = l.Query("services")
	q.Columns(
		"host_name",
		"description",
		"state",
		"last_state_change",
		"comments_with_extra_info",
		"plugin_output",
		"groups",
	)
	q.Filter("state_type > 0")

	resp, err = q.Exec()
	if err != nil {
		logger.Error("poller", "unable to retrieve services from %q node: %s", w.Config.Name, err)
		goto end
	}

	services = make(map[string][]livestatus.Record)
	for _, r := range resp.Records {
		name, err := r.GetString("host_name")
		if err != nil {
			logger.Error("poller", "unable to retrieve \"host_name\" from %q node: %s",
				w.Config.Name, err)
			goto end
		}

		if _, ok := services[name]; !ok {
			services[name] = make([]livestatus.Record, 0)
		}

		services[name] = append(services[name], r)
	}

	// Query hosts data
	q = l.Query("hosts")
	q.Columns(
		"host_name",
		"state",
		"last_state_change",
		"comments_with_extra_info",
		"groups",
	)

	resp, err = q.Exec()
	if err != nil {
		logger.Error("poller", "unable to retrieve hosts from %q node: %s", w.Config.Name, err)
		goto end
	}

	// Fill catalog with retrieved data
	for _, r := range resp.Records {
		h := &Host{
			Catalog: c,
		}

		h.Name, _ = r.GetString("host_name")
		h.State, _ = r.GetInt("state")
		h.StateChanged, _ = r.GetTime("last_state_change")
		h.Comments, _ = getComments(r)
		h.Links, _ = getLinks(h.Name, "", w.Config.Links)
		h.Groups, _ = getStringSlice(r, "groups")
		h.Services = make([]*Service, 0)

		if _, ok := services[h.Name]; ok {
			for _, sr := range services[h.Name] {
				s := &Service{Host: h}

				s.Name, _ = sr.GetString("description")
				s.State, _ = sr.GetInt("state")
				s.StateChanged, _ = sr.GetTime("last_state_change")
				s.Comments, _ = getComments(sr)
				s.Links, _ = getLinks(h.Name, s.Name, w.Config.Links)
				s.Output, _ = sr.GetString("plugin_output")
				s.Groups, _ = getStringSlice(sr, "groups")

				h.Services = append(h.Services, s)
			}
		}

		c.Hosts = append(c.Hosts, h)
	}

	// Replace catalog
	w.Catalog = c

end:
	// Reset refreshing state
	w.refreshing = false
}

// NewWorker returns a new worker instance
func NewWorker(c NodeConfig, wg *sync.WaitGroup) *Worker {
	return &Worker{
		Config: c,
		cs:     make(chan uint),
		wg:     wg,
	}
}

func getComments(r livestatus.Record) ([]Comment, error) {
	data, err := r.GetSlice("comments_with_extra_info")
	if err != nil {
		return []Comment{}, err
	}

	result := []Comment{}
	for _, entry := range data {
		slice, ok := entry.([]interface{})
		if !ok || len(slice) != 5 {
			logger.Error("poller", "skipping comment entry %v", entry)
			continue
		}

		ct, ok := slice[3].(float64)
		if !ok {
			logger.Error("poller", "unable to parse comment type: skipping entry %v", entry)
			continue
		}

		cd, ok := slice[4].(float64)
		if !ok {
			logger.Error("poller", "unable to parse comment date: skipping entry %v", entry)
			continue
		}

		result = append(result, Comment{
			Author:  fmt.Sprintf("%v", slice[1]),
			Content: fmt.Sprintf("%v", slice[2]),
			Type:    int64(ct),
			Date:    time.Unix(int64(cd), 0),
		})
	}

	return result, nil
}

func getLinks(host, service string, config []LinkConfig) ([][2]string, error) {
	result := make([][2]string, len(config))

	for i, l := range config {
		if service == "" {
			result[i] = [2]string{l.Label, l.HostURL}
			result[i][1] = strings.Replace(result[i][1], "%h", host, -1)
		} else {
			result[i] = [2]string{l.Label, l.ServiceURL}
			result[i][1] = strings.Replace(result[i][1], "%h", host, -1)
			result[i][1] = strings.Replace(result[i][1], "%s", service, -1)
		}
	}

	return result, nil
}

func getStringSlice(r livestatus.Record, key string) ([]string, error) {
	data, err := r.GetSlice(key)
	if err != nil {
		return []string{}, err
	}

	result := []string{}
	for _, entry := range data {
		value, ok := entry.(string)
		if ok {
			result = append(result, value)
		}
	}

	return result, nil
}
