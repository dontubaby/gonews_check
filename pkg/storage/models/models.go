package models

type Article struct {
	ID        int    `db:"id"`
	Title     string `db:"title"`
	Content   string `db:"description"`
	PubTime int64  `db:"published"`
	Link      string `db:"link"`
}
