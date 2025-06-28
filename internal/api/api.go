package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/log-ingestion-service/internal/storage"
	"github.com/yourusername/log-ingestion-service/internal/tracker"
)

// API handles HTTP requests
type API struct {
	router  *gin.Engine
	storage *storage.Storage
	tracker *tracker.Tracker
}

// New creates a new API instance
func New(storage *storage.Storage, tracker *tracker.Tracker) *API {
	router := gin.Default()
	api := &API{
		router:  router,
		storage: storage,
		tracker: tracker,
	}

	api.setupRoutes()
	return api
}

// setupRoutes configures the API routes
func (a *API) setupRoutes() {
	apiGroup := a.router.Group("/api")
	{
		apiGroup.GET("/logs", a.getLogs)
		apiGroup.GET("/logs/:id", a.getLogByID)
		apiGroup.GET("/status", a.getStatus)
	}
}

// Run starts the API server
func (a *API) Run(addr string) error {
	return a.router.Run(addr)
}

// getLogs returns all logs
func (a *API) getLogs(c *gin.Context) {
	logs, err := a.storage.GetPosts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// getLogByID returns a log by its ID
func (a *API) getLogByID(c *gin.Context) {
	id := c.Param("id")
	log, err := a.storage.GetPostByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, log)
}

// getStatus returns the latest ingestion status
func (a *API) getStatus(c *gin.Context) {
	status, err := a.tracker.GetLatestStatus(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}
