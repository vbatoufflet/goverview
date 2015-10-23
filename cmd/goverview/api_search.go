package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/fatih/set"
	"github.com/gin-gonic/gin"
	"github.com/vbatoufflet/goverview/pkg/poller"
	"github.com/vbatoufflet/goverview/pkg/slice"
)

type searchQuery struct {
	Nodes        []string `json:"nodes"`
	States       []int64  `json:"states"`
	Groups       []string `json:"groups"`
	Downtimes    bool     `json:"downtimes"`
	Acknowledges bool     `json:"acknowledges"`
	Filter       string   `json:"filter"`
}

func (s *Server) handleAPISearch(c *gin.Context) {
	var hgs, sgs set.Interface

	// Parse search query
	q := searchQuery{}
	if err := c.BindJSON(&q); err != nil {
		fmt.Printf("ERROR: unable to parse JSON data: %s\n", err)
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	// Get groups set
	gs := set.New(set.ThreadSafe)
	for _, g := range q.Groups {
		gs.Add(g)
	}

	// Prepare hosts response
	result := make(searchResponseList, 0)

	hosts, _ := s.poller.GetHosts()
	for _, h := range hosts {
		// Check for node and name filters
		if len(q.Nodes) > 0 && slice.StringIndexOf(q.Nodes, h.Catalog.Worker.Config.Name) == -1 ||
			q.Filter != "" && !strings.Contains(strings.ToLower(h.Name), strings.ToLower(q.Filter)) {
			continue
		}

		// Prepare search response
		sr := searchResponse{
			hostResponse: hostResponse{
				Name:         h.Name,
				State:        h.State,
				StateChanged: h.StateChanged,
				Comments:     make(commentList, 0),
				Links:        h.Links,
			},
			Services: make(serviceResponseList, 0),
		}

		// Fill comments and sort them by date
		for _, c := range h.Comments {
			if c.Type == poller.CommentTypeDowntime {
				// Check for downtimes filter
				if !q.Downtimes {
					goto nextHost
				}

				sr.hostResponse.InDowntime = true
			} else if c.Type == poller.CommentTypeAcknowledgement {
				// Check for acknowledges filter
				if !q.Acknowledges {
					goto nextHost
				}

				sr.hostResponse.Acknowledged = true
			}

			sr.hostResponse.Comments = append(sr.Comments, commentEntry{
				Author:  c.Author,
				Content: c.Content,
				Type:    c.Type,
				Date:    c.Date,
			})
		}

		sort.Sort(sr.Comments)

		// Add extra node information
		if h.Catalog.Worker.Config.Label != "" {
			sr.Node = h.Catalog.Worker.Config.Label
		} else {
			sr.Node = h.Catalog.Worker.Config.Name
		}

		// Prepare set for groups filtering
		if !gs.IsEmpty() {
			hgs = set.New(set.ThreadSafe)
			for _, g := range h.Groups {
				hgs.Add(g)
			}
		}

		// Fill services list
		for _, svc := range h.Services {
			// Check for states filter
			if len(q.States) > 0 && slice.Int64IndexOf(q.States, svc.State) == -1 {
				continue
			}

			// Check for groups filter
			if !gs.IsEmpty() {
				if len(svc.Groups) > 0 {
					sgs = set.New(set.ThreadSafe)
					for _, g := range svc.Groups {
						sgs.Add(g)
					}
				} else {
					sgs = hgs
				}

				if set.Intersection(gs, sgs).IsEmpty() {
					continue
				}
			}

			svcr := serviceResponse{
				Name:         svc.Name,
				State:        svc.State,
				StateChanged: svc.StateChanged,
				Comments:     make(commentList, 0),
				Links:        svc.Links,
				Output:       svc.Output,
			}

			// Fill comments and sort them by date
			for _, c := range svc.Comments {
				if c.Type == poller.CommentTypeDowntime {
					// Check for downtimes filter
					if !q.Downtimes {
						goto nextService
					}

					svcr.InDowntime = true
				} else if c.Type == poller.CommentTypeAcknowledgement {
					// Check for acknowledges filter
					if !q.Acknowledges {
						goto nextService
					}

					svcr.Acknowledged = true
				}

				svcr.Comments = append(svcr.Comments, commentEntry{
					Author:  c.Author,
					Content: c.Content,
					Type:    c.Type,
					Date:    c.Date,
				})
			}

			sort.Sort(svcr.Comments)

			sr.Services = append(sr.Services, svcr)
		nextService:
		}

		if len(sr.Services) == 0 {
			continue
		}

		sort.Sort(sr.Services)

		result = append(result, sr)
	nextHost:
	}

	sort.Sort(result)

	c.Header("Cache-Control", "private, max-age=0")
	c.JSON(http.StatusOK, result)
}
