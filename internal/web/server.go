package web

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"

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
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"formatNumber": func(n int) string {
			// convert number to string
			str := strconv.Itoa(n)
			// add commas every 3 digits from right to left
			if len(str) <= 3 {
				return str
			}
			var result []byte
			for i, digit := range str {
				if i > 0 && (len(str)-i)%3 == 0 {
					result = append(result, ',')
				}
				result = append(result, byte(digit))
			}
			return string(result)
		},
		"toPowerOf2": func(n int) string {
			if n <= 0 {
				return "0"
			}
			// calculate log2 and round to nearest integer
			exponent := int(math.Round(math.Log2(float64(n))))
			return fmt.Sprintf("2^%d", exponent)
		},
	}
	router.SetFuncMap(funcMap)

	// load templates from main directory and partials subdirectory
	mainTemplates, err := filepath.Glob("web/templates/*.tmpl")
	if err != nil {
		log.Panicf("failed to read main templates: %v", err)
	}
	partialTemplates, err := filepath.Glob("web/templates/partials/*.tmpl")
	if err != nil {
		log.Panicf("failed to read partial templates: %v", err)
	}
	allTemplates := append(mainTemplates, partialTemplates...)
	router.LoadHTMLFiles(allTemplates...)

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
	router.GET("/random", handler.RandomPage)

	return &Server{
		router: router,
		logger: logger,
	}
}

func (s *Server) Start() error {
	s.logger.Println("Web server starting on :8080")
	return s.router.Run(":8080")
}
