package postgress

import (
	models "Skillfactory/36-GoNews/pkg/storage/models"
	"context"
	"testing"
	//"context"
	//"fmt"
	//"log"
	//"strconv"
	//"github.com/jackc/pgx/v4/pgxpool"
)

func TestNew(t *testing.T) {
	_, err := New()
	if err != nil {
		t.Fatalf("Error of creating DB instance - %v", err)
	}
}

func TestAddArticle(t *testing.T) {
	db, err := New()
	if err != nil {
		t.Fatalf("Error create DB instance - %v", err)
	}

	articles := []models.Article{
		{
			ID:      1,
			Title:   "Test Title",
			Content: "Some test content here",
			PubTime: 1729584999,
			Link:    "https://go.dev/play/#",
		},
		{
			ID:      2,
			Title:   "Test Title2",
			Content: "Some test content here2",
			PubTime: 1729584991,
			Link:    "https://github.com/stretchr/testify",
		},
	}

	err = db.AddArticle(articles)
	if err != nil {
		t.Fatalf("Error adding article in database - %v", err)
	}
}

func TestGetArticle(t *testing.T) {
	db, err := NewTest()
	if err != nil {
		t.Fatalf("Error create DB instance - %v", err)
	}
	db.NewTable(context.Background())

	articles := []models.Article{
		{
			ID:      1,
			Title:   "Test Title",
			Content: "Some test content here",
			PubTime: 1129584991,
			Link:    "https://go.dev/play/1",
		},
		{
			ID:      2,
			Title:   "Test Title2",
			Content: "Some test content here2",
			PubTime: 1229584991,
			Link:    "https://go.dev/play/2",
		},
		{
			ID:      3,
			Title:   "Test Title3",
			Content: "Some test content here3",
			PubTime: 1329584991,
			Link:    "https://go.dev/play/3",
		},
	}
	err = db.AddArticle(articles)
	if err != nil {
		t.Fatalf("Error adding article in database - %v", err)
	}

	news, err := db.GetArticles(3)
	if err != nil {
		t.Fatalf("Error get artciles from DB - %v", err)
	}
	_, err = db.GetArticles(-1)
	if err.Error() != "Invalid count of articles!" {
		t.Fatalf("Error! Invalid input->")
	}
	if len(news) != len(articles) {
		t.Fatalf("not all news received - %v", err)
	}
	if news[0].ID != 3 || news[1].ID != 2 || news[2].ID != 1 {
		t.Fatalf("Mismatch IDs - %v", err)
	}
	if news[2].Title != "Test Title" || news[1].Title != "Test Title2" || news[0].Title != "Test Title3" {
		t.Fatalf("Mismatch Title - %v", err)
	}
	if news[2].Content != "Some test content here" || news[1].Content != "Some test content here2" || news[0].Content != "Some test content here3" {
		t.Fatalf("Mismatch description - %v", err)
	}
	if news[2].PubTime != 1129584991 || news[1].PubTime != 1229584991 || news[0].PubTime != 1329584991 {
		t.Fatalf("Mismatch Published - %v", err)
	}
	if news[2].Link != "https://go.dev/play/1" || news[1].Link != "https://go.dev/play/2" || news[0].Link != "https://go.dev/play/3" {
		t.Fatalf("Mismatch Link - %v", err)
	}
	db.NewTable(context.Background())
}
