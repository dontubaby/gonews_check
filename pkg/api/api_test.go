package api

import (
	"Skillfactory/36-GoNews/pkg/storage/postgress"

	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetArticlesHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/news?n=3", nil)
	w := httptest.NewRecorder()

	pool, err := postgress.NewTest()
	if err != nil {
		fmt.Printf("Error DB connection - %v", err)
	}
	defer pool.Db.Close()

	api := New(pool)

	api.GetArticlesHandler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if w.Code != http.StatusOK {
		t.Errorf("Invalid http code - %v", err)
	}
}
