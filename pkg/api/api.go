package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"Skillfactory/36-GoNews/pkg/storage/postgress"

	"github.com/gorilla/mux"
)

// Объект API
type Api struct {
	db *postgress.Storage
	r  *mux.Router
}

// Конуструктор объекта API
func New(db *postgress.Storage) *Api {
	api := Api{db: db, r: mux.NewRouter()}
	api.endpoints()
	return &api
}

// Init-метод для API-роутера
func (api *Api) Router() *mux.Router {
	return api.r
}

// Метод регистратор endpoint-ов и настраивающий саброутинг для файлового сервера (веб-приложения).
func (api *Api) endpoints() {
	// получить n последних новостей
	api.r.HandleFunc("/news/{n}", api.GetArticlesHandler).Methods(http.MethodGet, http.MethodOptions)
	// веб-приложение
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))

}

// Хэндлер обрабатывающий запрос по маршруту /news/{n} для возврата списка статей в JSON формате
func (api *Api) GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}

	s := mux.Vars(r)["n"]
	n, _ := strconv.Atoi(s)
	
	news, err := api.db.GetArticles(n)
	if err != nil {
		log.Printf("API get articles error - %v", err)
	}
	json.NewEncoder(w).Encode(news)
}
