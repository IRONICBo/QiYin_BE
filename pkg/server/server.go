package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server server interface.
type Server interface {
	ListenAndServe() error
}

// InitServer init server.
func InitServer(address string, r *gin.Engine) Server {
	server := &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return server
}
