package rss

import (
	"log"
	"strings"
	"time"

	models "Skillfactory/36-GoNews/pkg/storage/models"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
)

// Функция парсинга источника RSS. На вход получает строку с адресом источника, возвращает слайс объектов статей или ошибку.
func Parse(source string) ([]models.Article, error) {
	parser := gofeed.NewParser()
	var articles []models.Article
	var article models.Article
	feed, err := parser.ParseURL(source)
	if err != nil {
		log.Printf("Parsing error - %v", err)
		return nil, err
	}
	for _, item := range feed.Items {
		article, err = FeedItemToArticle(item)
		if err != nil {
			log.Println(err)
		}
		articles = append(articles, article)
	}
	return articles, nil
}

// Функция конвертер элементов ленты в объекты models.Article. На вход получает указатель объект ленты  gofeed.Item. Возвращает объект статьи или ошибку.
func FeedItemToArticle(item *gofeed.Item) (a models.Article, err error) {

	published := strings.ReplaceAll(item.Published, ",", "")
	t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", published)
	if err != nil {
		t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", published)
	} else if err != nil {
		log.Println(err)
		return models.Article{}, err
	}
	a.Content = item.Description

	a.Content = strip.StripTags(a.Content)
	unixtime := t.Unix()

	a = models.Article{
		Title:   item.Title,
		Content: a.Content,
		PubTime: unixtime,
		Link:    item.Link,
	}
	return a, nil
}
