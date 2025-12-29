package web

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

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
	c.HTML(http.StatusOK, "search.tmpl", gin.H{
		"title": "Search",
	})
}

func (h *Handler) Home(c *gin.Context) {
	h.logger.Println("serving home page")

	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"title": "Library of Babel",
	})
}

func (h *Handler) SearchPost(c *gin.Context) {
	text := c.PostForm("text")
	pageStr := c.DefaultPostForm("page", "1")

	if text == "" {
		h.logger.Println("empty search query")
		c.HTML(http.StatusBadRequest, "search.tmpl", gin.H{
			"title": "Search",
			"error": "Please enter text to search",
		})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	h.logger.Printf("searching for: %q (page %d)", text, page)

	const resultsPerPage = 20
	offset := (page - 1) * resultsPerPage

	locations, err := h.lib.SearchPaginated(text, offset, resultsPerPage)
	if err != nil {
		h.logger.Printf("search failed: %v", err)
		c.HTML(http.StatusInternalServerError, "search.tmpl", gin.H{
			"title": "Search",
			"error": "Search failed",
		})
		return
	}

	totalCount := h.lib.GetOccurrenceCount(text)
	totalPages := (totalCount + resultsPerPage - 1) / resultsPerPage

	c.HTML(http.StatusOK, "search.tmpl", gin.H{
		"title":       "Search Results",
		"query":       text,
		"locations":   locations,
		"total":       totalCount,
		"currentPage": page,
		"totalPages":  totalPages,
		"hasNext":     page < totalPages,
		"hasPrev":     page > 1,
	})
}

func (h *Handler) BrowseForm(c *gin.Context) {
	c.HTML(http.StatusOK, "browse.tmpl", gin.H{
		"title": "Browse",
	})
}

func (h *Handler) Browse(c *gin.Context) {
	locationStr := c.PostForm("location")
	query := c.PostForm("query")

	if locationStr == "" {
		h.logger.Println("no location provided")
		c.HTML(http.StatusBadRequest, "browse.tmpl", gin.H{
			"title": "Browse",
			"error": "No location specified",
		})
		return
	}

	location, err := library.LocationFromString(locationStr)
	if err != nil {
		h.logger.Printf("invalid location: %s - %v", locationStr, err)
		c.HTML(http.StatusBadRequest, "browse.tmpl", gin.H{
			"title": "Browse",
			"error": "Invalid location format",
		})
		return
	}

	h.logger.Printf("browsing: %s", location.String())

	content, err := h.lib.Browse(location)
	if err != nil {
		h.logger.Printf("browse failed: %v", err)
		c.HTML(http.StatusInternalServerError, "browse.tmpl", gin.H{
			"title": "Browse",
			"error": "Failed to load page",
		})
		return
	}

	formattedContent := formatPageContent(content)

	var displayContent template.HTML
	if query != "" {
		displayContent = template.HTML(highlightText(formattedContent, query)) //nolint:gosec
	} else {
		displayContent = template.HTML(html.EscapeString(formattedContent)) //nolint:gosec
	}

	c.HTML(http.StatusOK, "browse.tmpl", gin.H{
		"title":          "Page Content",
		"location":       location,
		"displayContent": displayContent,
		"hasQuery":       query != "",
		"nextLocation":   location.Next(),
		"prevLocation":   location.Previous(),
	})
}

// formatPageContent formats content to have 80 characters per line (40 lines total)
func formatPageContent(content string) string {
	const charsPerLine = 80
	var formatted strings.Builder

	runes := []rune(content)
	for i := 0; i < len(runes); i += charsPerLine {
		end := min(i+charsPerLine, len(runes))
		formatted.WriteString(string(runes[i:end]))
		if end < len(runes) {
			formatted.WriteString("\n")
		}
	}

	return formatted.String()
}

// Wraps the query text in the content with HTML mark tags for highlighting
// it handles multi-line matches by treating any whitespace in the query as matching any whitespace in the content
func highlightText(content, query string) string {
	escapedContent := html.EscapeString(content)
	escapedQuery := html.EscapeString(query)

	// create a regex pattern that treats any whitespace in the query as matching any whitespace (\s+) in the content
	// This allows matching across line breaks
	pattern := regexp.QuoteMeta(escapedQuery)
	pattern = regexp.MustCompile(`\s+`).ReplaceAllString(pattern, `\s+`)

	// use case-insensitive matching
	re := regexp.MustCompile(`(?i)` + pattern)

	// replace all matches with highlighted version
	highlighted := re.ReplaceAllStringFunc(escapedContent, func(match string) string {
		return "<mark class=\"bg-amber-500/30 text-amber-200 font-bold\">" + match + "</mark>"
	})

	return highlighted
}

func (h *Handler) RandomPage(c *gin.Context) {
	h.logger.Println("generating random page")
	location := library.RandomLocation()

	h.logger.Printf("random location: %s", location.String())

	// get the content at this location
	content, err := h.lib.Browse(location)
	if err != nil {
		h.logger.Printf("browse failed for random page: %v", err)
		c.HTML(http.StatusInternalServerError, "browse.tmpl", gin.H{
			"title": "Browse",
			"error": "Failed to load random page",
		})
		return
	}

	formattedContent := formatPageContent(content)
	displayContent := template.HTML(html.EscapeString(formattedContent)) //nolint:gosec

	c.HTML(http.StatusOK, "browse.tmpl", gin.H{
		"title":          "Random Page",
		"location":       location,
		"displayContent": displayContent,
		"nextLocation":   location.Next(),
		"prevLocation":   location.Previous(),
	})
}
