package DB

import (
	"Skillfactory/36-GoNews/pkg/storage/models"
	"log"
)

type DbInterface interface {
	GetArticles(int) ([]models.Article, error)
	AddArticle([]models.Article) error
}

func GetAll(n int, db DbInterface) ([]models.Article, error) {
	result, err := db.GetArticles(n)
	if err != nil {
		log.Fatalf("Error when GET articles from server: %v\n", err)
		return nil, err
	}

	return result, nil
}

func Add(db DbInterface, articles []models.Article) error {
	err := db.AddArticle(articles)
	if err != nil {
		log.Fatalf("Error when ADD article to database: %v\n", err)
		return err
	}
	return nil
}

