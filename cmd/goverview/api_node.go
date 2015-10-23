package main

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func (s *Server) handleAPINodes(c *gin.Context) {
	result := make(nodeResponseList, 0)

	for _, n := range s.config.Poller.Nodes {
		result = append(result, nodeResponse{
			Name:  n.Name,
			Label: n.Label,
		})
	}

	sort.Sort(result)

	c.Header("Cache-Control", "private, max-age=0")
	c.JSON(http.StatusOK, result)
}
