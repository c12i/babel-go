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

	// set custom template functions
	router.SetFuncMap(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	})

	router.LoadHTMLGlob("web/templates/*.html")

	// serve static files
	router.Static("/static", "./web/static")

	// healthcheck
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// library routes
	router.GET("/", handler.Home)
	router.GET("/search", handler.SearchForm)
	router.POST("/search", handler.SearchPost)
	router.GET("/browse", handler.BrowseForm)
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
