package main

import (
	//"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"Skillfactory/36-GoNews/pkg/api"
	"Skillfactory/36-GoNews/pkg/rss"
	DB "Skillfactory/36-GoNews/pkg/storage"
	"Skillfactory/36-GoNews/pkg/storage/models"
	"Skillfactory/36-GoNews/pkg/storage/postgress"
)

// Объект с настройками приложения
type Config struct {
	RSSsources []string `json:"source"`
	Interval   int      `json:"interval"`
}

// Функция - конвертер JSON файла с настройками в объект с настройками приложения. TODO: перенести функцию в пакет utils
func ParseConfigFile(filename string) (Config, error) {
	var data []byte
	_, err := os.Open(filename)
	if err != nil {
		log.Printf("Open file error - %v", err)
		return Config{}, err
	}
	//defer configFile.Close()

	data, err = ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	jsonErr := json.Unmarshal(data, &config)
	if jsonErr != nil {
		log.Printf("Unmarhaling error - %v", err)
		return Config{}, err
	}
	return config, nil
}

// Функция асинхронной обработки RSS-лент. На вход принимает источник RSS, интерфейс БД, канал для записи статей, канал для записи ошибок парсинга RSS
func AsynParser(source string, db DB.DbInterface, articles chan<- []models.Article, errs chan<- error, interval int) {
	for {
		news, err := rss.Parse(source)
		if err != nil {
			errs <- err
			continue
		}
		articles <- news
		time.Sleep(time.Minute * time.Duration(interval))
	}
}

func main() {

	pool, err := postgress.New()
	if err != nil {
		log.Printf("Error DB connection - %v", err)
	}
	defer pool.Db.Close()

	api := api.New(pool)

	config, err := ParseConfigFile("config.json")
	if err != nil {
		log.Printf("Config decoding error - %v", err)
	}

	articleStream := make(chan []models.Article)
	errorStream := make(chan error)
	for _, source := range config.RSSsources {
		go AsynParser(source, pool, articleStream, errorStream, config.Interval)
	}

	go func() {
		log.Println("articleStream start working")
		for article := range articleStream {
			pool.AddArticle(article)
		}
		log.Println("articleStream end working")
	}()

	go func() {
		log.Println("errorStream start working")
		for err := range errorStream {
			log.Println("Error:", err)
		}
		log.Println("errorStream end working")
	}()

	log.Println("Server start working!")
	err = http.ListenAndServe(":80", api.Router())

	if err != nil {
		log.Println(err)

	}

}
