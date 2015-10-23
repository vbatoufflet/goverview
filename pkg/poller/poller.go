package poller

import (
	"os"
	"sync"

	"github.com/fatih/set"
	"github.com/vbatoufflet/goverview/pkg/logger"
)

// Poller represents a poller instance
type Poller struct {
	config   Config
	Workers  []*Worker
	wg       *sync.WaitGroup
	stopping bool
}

// Shutdown stops the poller
func (p *Poller) Shutdown() {
	if p.stopping {
		return
	}

	logger.Info("poller", "shutting down poller workers")

	for _, w := range p.Workers {
		w.Shutdown()
	}

	p.wg.Wait()
}

// Refresh refreshes the poller catalog content
func (p *Poller) Refresh() {
	logger.Info("poller", "refreshing poller catalog")

	for _, w := range p.Workers {
		w.Refresh()
	}
}

// GetHosts returns the list of hosts present in the poller
func (p *Poller) GetHosts() ([]*Host, error) {
	result := []*Host{}

	for _, w := range p.Workers {
		if w.Catalog == nil {
			continue
		}

		result = append(result, w.Catalog.Hosts...)
	}

	return result, nil
}

// GetHostByName returns an host based on its name, or nil if not found
func (p *Poller) GetHostByName(name string) ([]*Host, error) {
	result := []*Host{}

	for _, w := range p.Workers {
		if w.Catalog == nil {
			continue
		}

		h, err := w.Catalog.GetHostByName(name)
		if err == os.ErrNotExist {
			continue
		}

		result = append(result, h)
	}

	if len(result) == 0 {
		return nil, os.ErrNotExist
	}

	return result, nil
}

// NewPoller returns a new poller instance
func NewPoller(c Config) *Poller {
	p := &Poller{
		config:  c,
		Workers: make([]*Worker, 0),
		wg:      &sync.WaitGroup{},
	}

	logger.Info("poller", "starting poller workers")

	names := set.New(set.ThreadSafe)

	for _, n := range c.Nodes {
		if n.Name == "" {
			logger.Error("poller", "missing mandatory `name' node parameter")
			continue
		} else if n.Address == "" {
			logger.Error("poller", "missing mandatory `remote_addr' node parameter")
			continue
		} else if names.Has(n.Name) {
			logger.Error("poller", "duplicate name `%s' found in node definition", n.Name)
			continue
		}

		w := NewWorker(n, p.wg)

		p.Workers = append(p.Workers, w)
		p.wg.Add(1)

		go w.Start()
		w.Refresh()

		names.Add(n.Name)
	}

	return p
}
