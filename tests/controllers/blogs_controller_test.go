package controllers_test

import (
	"encoding/json"
	"go_http_testing/controllers"
	"go_http_testing/models"
	models_test "go_http_testing/tests/models"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := models_test.SetupTestDB()
	if err != nil {
		log.Fatalf("Error setting up test database: %v", err)
	}

	code := m.Run()

	// Cleanup
	models_test.TeardownTestDB()

	os.Exit(code)
}

func TestBlogsIndex(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/blogs", nil)
	w := httptest.NewRecorder()

	// Seed the database with some test data
	blog1 := models.Blog{Title: "First Blog", Content: "Content of the first blog"}
	blog2 := models.Blog{Title: "Second Blog", Content: "Content of the second blog"}
	models.DB.Create(&blog1)
	models.DB.Create(&blog2)

	controllers.BlogsIndex(w, req)
	res := w.Result()

	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, res.StatusCode, http.StatusOK, "API should return 200 status code")

	// Read data from the body and parse the JSON
	var blogs []models.Blog
	err := json.NewDecoder(res.Body).Decode(&blogs)
	assert.NoError(t, err)

	// Check the length of the blogs array
	assert.Len(t, blogs, 2)

	// Check that the first blog matches the expected data
	assert.Equal(t, blogs[0].Title, blog2.Title)
	assert.Equal(t, blogs[0].Content, blog2.Content)

	// Check that the second blog matches the expected data
	assert.Equal(t, blogs[1].Title, blog1.Title)
	assert.Equal(t, blogs[1].Content, blog1.Content)
}

func TestBlogsIndexEmptyTable(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/blogs", nil)
	w := httptest.NewRecorder()

	// Empty the table
	models.DB.Exec("DELETE FROM blogs")

	controllers.BlogsIndex(w, req)
	res := w.Result()

	defer res.Body.Close()

	// Check the status code
	assert.Equal(t, res.StatusCode, http.StatusOK, "API should return 200 status code")

	// Read data from the body and parse the JSON
	var blogs []models.Blog
	err := json.NewDecoder(res.Body).Decode(&blogs)
	assert.NoError(t, err)

	// Check the length of the blogs array
	assert.Len(t, blogs, 0)
}

func TestBlogShow(t *testing.T) {
	// Seed the table
	blog := models.Blog{Title: "Test Blog", Content: "This is a test blog content"}
	models.DB.Create(&blog)

	t.Run("Valid Blog ID", func(t *testing.T) {
		blogID := strconv.FormatUint(uint64(blog.ID), 10)
		req := httptest.NewRequest(http.MethodGet, "/blogs/"+blogID, nil)
		w := httptest.NewRecorder()

		controllers.BlogShow(w, req)
		res := w.Result()
		defer res.Body.Close()

		// Check the status code
		assert.Equal(t, http.StatusOK, res.StatusCode)

		// Read data from the body and parse the JSON
		var returnedBlog models.Blog
		err := json.NewDecoder(res.Body).Decode(&returnedBlog)
		assert.NoError(t, err)

		// Check the returned blog matches the seeded data
		assert.Equal(t, blog.Title, returnedBlog.Title)
		assert.Equal(t, blog.Content, returnedBlog.Content)
	})

	t.Run("Invalid Blog ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/blogs/invalid", nil)
		w := httptest.NewRecorder()

		controllers.BlogShow(w, req)
		res := w.Result()
		defer res.Body.Close()

		// Check the status code
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("Nonexistent Blog ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/blogs/99999", nil)
		w := httptest.NewRecorder()

		controllers.BlogShow(w, req)
		res := w.Result()
		defer res.Body.Close()

		// Check the status code
		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
