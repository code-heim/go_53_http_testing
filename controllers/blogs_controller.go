package controllers

import (
	"encoding/json"
	"go_http_testing/models"
	"net/http"
	"strconv"
	"strings"
)

func BlogsIndex(w http.ResponseWriter, r *http.Request) {
	blogs := models.BlogsAll()
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Encode the blogs data to JSON and write to the response
	if err := json.NewEncoder(w).Encode(blogs); err != nil {
		// If encoding fails, respond with a 500 Internal Server Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// BlogShow handles requests to show a single blog by ID
func BlogShow(w http.ResponseWriter, r *http.Request) {
	// Extract the blog ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/blogs/")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid blog ID", http.StatusBadRequest)
		return
	}

	// Fetch the blog by ID
	blog := models.BlogsFind(id)
	if blog.ID == 0 {
		http.Error(w, "Blog not found", http.StatusNotFound)
		return
	}

	// Set the Content-Type header and write the response as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(blog); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
