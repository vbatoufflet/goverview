package main

import (
	"net/http"

	"github.com/facette/natsort"
	"github.com/fatih/set"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleAPIGroups(c *gin.Context) {
	gs := set.New(set.ThreadSafe)

	// Get unique sets of groups names
	hosts, _ := s.poller.GetHosts()
	for _, h := range hosts {
		for _, g := range h.Groups {
			gs.Add(g)
		}

		for _, s := range h.Services {
			for _, g := range s.Groups {
				gs.Add(g)
			}
		}
	}

	result := set.StringSlice(gs)
	natsort.Sort(result)

	c.Header("Cache-Control", "private, max-age=0")
	c.JSON(http.StatusOK, result)
}
