package poller

import (
	"os"
	"time"
)

const (
	// CommentTypeUser represents an user comment type
	CommentTypeUser = 1
	// CommentTypeDowntime represents a downtime comment type
	CommentTypeDowntime = 2
	// CommentTypeFlapping represents a flap comment type
	CommentTypeFlapping = 3
	// CommentTypeAcknowledgement represents a acknowledge comment type
	CommentTypeAcknowledgement = 4
)

// Catalog represents the states catalog in a poller worker
type Catalog struct {
	Hosts  []*Host
	Worker *Worker
}

// GetHostByName returns an host based on its name, or nil if not found
func (c *Catalog) GetHostByName(name string) (*Host, error) {
	for _, h := range c.Hosts {
		if h.Name == name {
			return h, nil
		}
	}

	return nil, os.ErrNotExist
}

// NewCatalog returns a new states catalog instance
func NewCatalog(w *Worker) *Catalog {
	return &Catalog{
		Hosts:  make([]*Host, 0),
		Worker: w,
	}
}

// Host represents a host entry in a states catalog
type Host struct {
	Name         string
	State        int64
	StateChanged time.Time
	Comments     []Comment
	Links        [][2]string
	Services     []*Service
	Groups       []string
	Catalog      *Catalog
}

// Service represents a service entry in a states catalog host
type Service struct {
	Name         string
	State        int64
	StateChanged time.Time
	Comments     []Comment
	Links        [][2]string
	Output       string
	Groups       []string
	Host         *Host
}

// Comment represents a comment entry in a states catalog host or service
type Comment struct {
	Author  string
	Content string
	Type    int64
	Date    time.Time
}
