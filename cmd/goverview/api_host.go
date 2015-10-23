package main

import (
	"net/http"
	"os"
	"sort"

	"github.com/fatih/set"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleAPIHost(c *gin.Context) {
	name, _ := c.Params.Get("name")

	// Get hosts list by their name
	hosts, err := s.poller.GetHostByName(name)
	if err == os.ErrNotExist {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	// Prepare hosts response
	result := make(hostResponseList, 0)

	for _, h := range hosts {
		hr := hostResponse{
			Name:         h.Name,
			State:        h.State,
			StateChanged: h.StateChanged,
			// TODO: fix downtimes
			// InDowntime:   len(h.Downtimes) > 0,
		}

		if h.Catalog.Worker.Config.Label != "" {
			hr.Node = h.Catalog.Worker.Config.Label
		} else {
			hr.Node = h.Catalog.Worker.Config.Name
		}

		result = append(result, hr)
	}

	sort.Sort(result)

	c.Header("Cache-Control", "private, max-age=0")
	c.JSON(http.StatusOK, result)
}

func (s *Server) handleAPIHosts(c *gin.Context) {
	hs := set.New(set.ThreadSafe)

	// Get hosts list
	hosts, _ := s.poller.GetHosts()
	for _, h := range hosts {
		hs.Add(h.Name)
	}

	// Get unique set of hosts names
	result := set.StringSlice(hs)
	sort.Strings(result)

	c.Header("Cache-Control", "private, max-age=0")
	c.JSON(http.StatusOK, result)
}
