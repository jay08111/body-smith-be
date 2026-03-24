package model

import "time"

type Post struct {
	ID              int64      `db:"id" json:"id"`
	Title           string     `db:"title" json:"title"`
	Slug            string     `db:"slug" json:"slug"`
	Content         string     `db:"content" json:"content"`
	Thumbnail       *string    `db:"thumbnail" json:"thumbnail,omitempty"`
	MetaTitle       string     `db:"meta_title" json:"meta_title"`
	MetaDescription string     `db:"meta_description" json:"meta_description"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"-"`
}

type CreatePostRequest struct {
	Title           string  `json:"title"`
	Content         string  `json:"content"`
	Thumbnail       *string `json:"thumbnail"`
	MetaTitle       string  `json:"meta_title"`
	MetaDescription string  `json:"meta_description"`
}

type UpdatePostRequest struct {
	Title           string  `json:"title"`
	Content         string  `json:"content"`
	Thumbnail       *string `json:"thumbnail"`
	MetaTitle       string  `json:"meta_title"`
	MetaDescription string  `json:"meta_description"`
}

type PostListItem struct {
	ID              int64      `db:"id" json:"id"`
	Title           string     `db:"title" json:"title"`
	Slug            string     `db:"slug" json:"slug"`
	Content         string     `db:"content" json:"content"`
	Thumbnail       *string    `db:"thumbnail" json:"thumbnail,omitempty"`
	MetaTitle       string     `db:"meta_title" json:"meta_title"`
	MetaDescription string     `db:"meta_description" json:"meta_description"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"-"`
}

type PostListResponse struct {
	Items      []PostListItem `json:"items"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	Total      int64          `json:"total"`
	TotalPages int            `json:"total_pages"`
}
