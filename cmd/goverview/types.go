package main

import (
	"time"

	"github.com/facette/natsort"
)

type hostResponse struct {
	Name         string      `json:"name"`
	State        int64       `json:"state"`
	StateChanged time.Time   `json:"state_changed"`
	Comments     commentList `json:"comments"`
	InDowntime   bool        `json:"in_downtime"`
	Acknowledged bool        `json:"acknowledged"`
	Links        [][2]string `json:"links"`
	Node         string      `json:"node"`
}

type hostResponseList []hostResponse

func (l hostResponseList) Len() int {
	return len(l)
}

func (l hostResponseList) Less(i, j int) bool {
	return natsort.Compare(l[i].Name, l[j].Name)
}

func (l hostResponseList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type serviceResponse struct {
	Name         string      `json:"name"`
	State        int64       `json:"state"`
	StateChanged time.Time   `json:"state_changed"`
	Comments     commentList `json:"comments"`
	InDowntime   bool        `json:"in_downtime"`
	Acknowledged bool        `json:"acknowledged"`
	Links        [][2]string `json:"links"`
	Output       string      `json:"output"`
}

type serviceResponseList []serviceResponse

func (l serviceResponseList) Len() int {
	return len(l)
}

func (l serviceResponseList) Less(i, j int) bool {
	return natsort.Compare(l[i].Name, l[j].Name)
}

func (l serviceResponseList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type searchResponse struct {
	hostResponse
	Services serviceResponseList `json:"services"`
}

type searchResponseList []searchResponse

func (l searchResponseList) Len() int {
	return len(l)
}

func (l searchResponseList) Less(i, j int) bool {
	return natsort.Compare(l[i].hostResponse.Name, l[j].hostResponse.Name)
}

func (l searchResponseList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type nodeResponse struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type nodeResponseList []nodeResponse

func (l nodeResponseList) Len() int {
	return len(l)
}

func (l nodeResponseList) Less(i, j int) bool {
	return natsort.Compare(l[i].Name, l[j].Name)
}

func (l nodeResponseList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type groupResponse struct {
	Hosts    []string `json:"hosts"`
	Services []string `json:"services"`
}

type commentEntry struct {
	Author  string    `json:"author"`
	Content string    `json:"content"`
	Type    int64     `json:"type"`
	Date    time.Time `json:"date"`
}

type commentList []commentEntry

func (l commentList) Len() int {
	return len(l)
}

func (l commentList) Less(i, j int) bool {
	return l[i].Date.Before(l[j].Date)
}

func (l commentList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
