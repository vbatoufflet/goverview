package main

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vbatoufflet/goverview/pkg/logger"
)

func (s *Server) handleStatic(c *gin.Context) {
	var ct string

	// Get static file path
	p, _ := c.Params.Get("path")
	if p == "" {
		p = "index.html"
	}

	// Try to get file data from filesystem
	localPath := path.Join(s.config.StaticDir, p)
	if _, err := os.Stat(localPath); err == nil {
		logger.Debug("server", "serving %q from filesystem", p)
		c.File(localPath)
		return
	}

	// Get file data from built-in assets
	data, err := Asset(strings.TrimPrefix(p, "/"))
	if err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}

	switch path.Ext(p) {
	case ".css":
		ct = "text/css"

	case ".js":
		ct = "text/javascript"

	default:
		ct = http.DetectContentType(data)
	}

	logger.Debug("server", "serving %q from built-in assets", p)
	c.Data(http.StatusOK, ct, data)
}
