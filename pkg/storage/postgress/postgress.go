package postgress

import (
	models "Skillfactory/36-GoNews/pkg/storage/models"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

// Объект хранилища
type Storage struct {
	Db *pgxpool.Pool
}

// Вспомогательный метод конструктор для новой таблицы. Применяется для тестирования работы приложения с БД.
func (s *Storage) NewTable(ctx context.Context) error {
	_, err := s.Db.Exec(ctx, `
DROP TABLE IF EXISTS articles;
CREATE TABLE articles (
  id BIGSERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  published BIGINT,
  link TEXT NOT NULL UNIQUE 
);`)
	if err != nil {
		log.Fatalf("Error!Cant create new table:  %v\n", err)
		return err
	}
	return nil
}

// Конструктор для объекта хранилища.
func New() (*Storage, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		return nil, err
	}
	pwd := os.Getenv("DBPASSWORD")

	connString := "postgres://postgres:" + pwd + "@localhost:5432/gonews"

	db, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Cant create new instance of DB: %v\n", err)
		return nil, err
	}
	s := Storage{
		Db: db,
	}
	return &s, nil
}

// Конструктор для хранилища, применяемого для тестов.
func NewTest() (*Storage, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		return nil, err
	}
	pwd := os.Getenv("DBPASSWORDT")

	connString := "postgres://postgres:" + pwd + "@localhost:5432/newstest"

	db, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Cant create new instance of DB: %v\n", err)
		return nil, err
	}
	s := Storage{
		Db: db,
	}
	return &s, nil
}

// Метод возврата статей. На вход применяте количество статей, необходимых к возврату. Возвращает слайс с объектами или ошибку.
func (s *Storage) GetArticles(n int) ([]models.Article, error) {
	if n < 1 {
		err := fmt.Errorf("Error!Invalid count of article - got %v", n)
		log.Println(err)
		e := errors.New("Invalid count of articles!")
		return nil, e
	}
	if n == 0 {
		n = 10
	}
	q := strconv.Itoa(n)

	rows, err := s.Db.Query(context.Background(), `SELECT * FROM articles ORDER BY published DESC LIMIT $1`, q)
	if err != nil {
		log.Printf("Cant read data from database: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	articles := []models.Article{}

	for rows.Next() {
		article := models.Article{}
		err = rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.PubTime,
			&article.Link,
		)
		if err != nil {
			return nil, fmt.Errorf("unable scan row: %w", err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// Метод добавлений новых статей в базу данных. На входи принимает слайс с объектами статей, возвращает ошибку при необходимости.
func (s *Storage) AddArticle(articles []models.Article) error {
	for _, a := range articles {
		_, err := s.Db.Exec(context.Background(), `INSERT INTO articles 
		(title,content,published,link) VALUES ($1,$2,$3,$4);`,
			a.Title, a.Content, a.PubTime, a.Link)
		if err != nil {
			log.Fatalf("Cant add data in database! %v\n", err)
			return err
		}
	}
	return nil
}
