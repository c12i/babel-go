package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	logger *log.Logger
}

func NewServer(handler *Handler, logger *log.Logger) *Server {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("web/templates/*")
	router.SetHTMLTemplate(template.Must(template.ParseGlob("web/templates/*.html")))

	// healthcheck
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// library routes
	router.GET("/", handler.Home)
	router.GET("/search", handler.SearchForm)
	router.POST("/search", handler.SearchPost)
	router.POST("/browse", handler.Browse)

	return &Server{
		router: router,
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Println("Web server starting on :8080")
	return s.router.Run(":8080")
}
