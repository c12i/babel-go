package web

import (
	"log"
	"net/http"

	"github.com/c12i/babel-go/internal/library"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	lib    *library.Library
	logger *log.Logger
}

func NewHandler(lib *library.Library, logger *log.Logger) *Handler {
	return &Handler{
		lib:    lib,
		logger: logger,
	}
}

func (h *Handler) SearchForm(c *gin.Context) {
	c.HTML(http.StatusOK, "search.html", gin.H{
		"title": "Search",
	})
}

func (h *Handler) Home(c *gin.Context) {
	h.logger.Println("Serving home page")

	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Library of Babel",
	})
}

func (h *Handler) SearchPost(c *gin.Context) {
	text := c.PostForm("text")

	if text == "" {
		h.logger.Println("Empty search query")
		c.HTML(http.StatusBadRequest, "search.html", gin.H{
			"title": "Search",
			"error": "Please enter text to search",
		})
		return
	}

	h.logger.Printf("Searching for: %q", text)

	locations, err := h.lib.SearchPaginated(text, 0, 10)
	if err != nil {
		h.logger.Printf("Search failed: %v", err)
		c.HTML(http.StatusInternalServerError, "search.html", gin.H{
			"title": "Search",
			"error": "Search failed",
		})
		return
	}

	totalCount := h.lib.GetOccurrenceCount(text)

	c.HTML(http.StatusOK, "search.html", gin.H{
		"title":     "Search Results",
		"query":     text,
		"locations": locations,
		"total":     totalCount,
	})
}

func (h *Handler) Browse(c *gin.Context) {
	locationStr := c.PostForm("location")

	if locationStr == "" {
		h.logger.Println("No location provided")
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"title": "Browse",
			"error": "No location specified",
		})
		return
	}

	location, err := library.LocationFromString(locationStr)
	if err != nil {
		h.logger.Printf("Invalid location: %s - %v", locationStr, err)
		c.HTML(http.StatusBadRequest, "browse.html", gin.H{
			"title": "Browse",
			"error": "Invalid location format",
		})
		return
	}

	h.logger.Printf("Browsing: %s", location.String())

	content, err := h.lib.Browse(location)
	if err != nil {
		h.logger.Printf("Browse failed: %v", err)
		c.HTML(http.StatusInternalServerError, "browse.html", gin.H{
			"title": "Browse",
			"error": "Failed to load page",
		})
		return
	}

	c.HTML(http.StatusOK, "browse.html", gin.H{
		"title":    "Page Content",
		"location": location,
		"content":  content,
	})
}
