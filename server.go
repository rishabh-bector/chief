package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"bector.dev/chief/config"
	"bector.dev/chief/pipeline"
	"github.com/gin-gonic/gin"
)

const STATUS_STOPPED = "STOPPED"
const STATUS_RUNNING = "RUNNING"

// ChiefServer is the main chief server
type ChiefServer struct {
	pipelines []pipeline.Pipeline
	config    *config.ChiefConfig
	server    *http.Server
}

func NewChiefServer(config *config.ChiefConfig) ChiefServer {
	return ChiefServer{
		config: config,
	}
}

func (s *ChiefServer) Start() error {
	// Router setup
	router := gin.Default()
	router.GET("/new", s.NewPipeline)
	router.GET("/kill", s.Kill)
	router.GET("/status", s.Status)

	s.server = &http.Server{
		Addr:    "localhost:2222",
		Handler: router,
	}

	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Info("Server closed under request")
		} else {
			log.Fatal("Server closed unexpectedly")
		}
	}

	return nil
}

func (s *ChiefServer) NewPipeline(c *gin.Context) {

}

func (s *ChiefServer) Kill(c *gin.Context) {
	s.server.Close()
}

type ServerStatus struct {
	Status string `json:"status"`
}

func (s *ChiefServer) Status(c *gin.Context) {
	c.JSON(200, ServerStatus{
		Status: STATUS_RUNNING,
	})
}
