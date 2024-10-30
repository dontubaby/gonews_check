package rss

import (
	models "Skillfactory/36-GoNews/pkg/storage/models"
	"log"
	"reflect"
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestParse(t *testing.T) {
	validSource := "https://habr.com/ru/rss/hub/go/all/?fl=ru"
	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(validSource)

	if err != nil {
		t.Fatalf("Error parsing URL to feed - %v", err)
	}
	var article models.Article
	var articles []models.Article
	for _, item := range feed.Items {
		article, err = FeedItemToArticle(item)
		if err != nil {
			log.Println(err)
		}
		articles = append(articles, article)
	}
	if len(articles) < 1 {
		t.Fatalf("Error parsing feed to items - %v", err)
	}

}

func TestFeedItemToArticle(t *testing.T) {
	type args struct {
		item *gofeed.Item
	}
	tests := []struct {
		name string
		args args
		want models.Article
	}{
		{
			name: "Valid data",

			args: args{

				item: &gofeed.Item{
					Title:       "Test Title 1",
					Description: "Test Description 1",
					Published:   "Wed, 30 Oct 2024 11:18:07 GMT",
					Link:        "https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L96",
				},
			},

			want: models.Article{
				Title:     "Test Title 1",
				Content:   "Test Description 1",
				PubTime: 1730287087,
				Link:      "https://github.com/mmcdole/gofeed/blob/v1.3.0/parser.go#L96",
			},
		},

		{
			name: "Empty data",

			args: args{

				item: &gofeed.Item{
					Title:     "",
					Content:   "",
					Published: "",
					Link:      "",
				},
			},

			want: models.Article{
				Title:     "",
				Content:   "",
				PubTime: -62135596800,
				Link:      "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := FeedItemToArticle(tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FeedItemToArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}
