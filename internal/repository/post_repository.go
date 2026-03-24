package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"body-smith-be/internal/model"
)

type PostRepository interface {
	Create(ctx context.Context, post *model.Post) (*model.Post, error)
	List(ctx context.Context, pagination model.Pagination) ([]model.PostListItem, int64, error)
	ListPublic(ctx context.Context, pagination model.Pagination) ([]model.PostListItem, int64, error)
	GetByID(ctx context.Context, id int64) (*model.Post, error)
	GetBySlug(ctx context.Context, slug string) (*model.Post, error)
	SlugExists(ctx context.Context, slug string, excludeID *int64) (bool, error)
	Update(ctx context.Context, post *model.Post) (*model.Post, error)
	Delete(ctx context.Context, id int64) error
}

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *model.Post) (*model.Post, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO posts (title, slug, content, thumbnail, meta_title, meta_description)
		VALUES (?, ?, ?, ?, ?, ?)
	`, post.Title, post.Slug, post.Content, post.Thumbnail, post.MetaTitle, post.MetaDescription)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}

func (r *postRepository) List(ctx context.Context, pagination model.Pagination) ([]model.PostListItem, int64, error) {
	var posts []model.PostListItem
	if err := r.db.SelectContext(ctx, &posts, `
		SELECT id, title, slug, content, thumbnail, meta_title, meta_description, created_at, updated_at
		FROM posts
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, pagination.PerPage, pagination.Offset); err != nil {
		return nil, 0, err
	}

	total, err := r.countPosts(ctx)
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *postRepository) ListPublic(ctx context.Context, pagination model.Pagination) ([]model.PostListItem, int64, error) {
	return r.List(ctx, pagination)
}

func (r *postRepository) GetByID(ctx context.Context, id int64) (*model.Post, error) {
	var post model.Post
	err := r.db.GetContext(ctx, &post, `
		SELECT id, title, slug, content, thumbnail, meta_title, meta_description, created_at, updated_at
		FROM posts
		WHERE id = ?
		LIMIT 1
	`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) GetBySlug(ctx context.Context, slug string) (*model.Post, error) {
	var post model.Post
	err := r.db.GetContext(ctx, &post, `
		SELECT id, title, slug, content, thumbnail, meta_title, meta_description, created_at, updated_at
		FROM posts
		WHERE slug = ?
		LIMIT 1
	`, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) SlugExists(ctx context.Context, slug string, excludeID *int64) (bool, error) {
	var count int
	query := `SELECT COUNT(1) FROM posts WHERE slug = ?`
	args := []any{slug}
	if excludeID != nil {
		query += ` AND id != ?`
		args = append(args, *excludeID)
	}

	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *postRepository) Update(ctx context.Context, post *model.Post) (*model.Post, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE posts
		SET title = ?, slug = ?, content = ?, thumbnail = ?, meta_title = ?, meta_description = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, post.Title, post.Slug, post.Content, post.Thumbnail, post.MetaTitle, post.MetaDescription, post.ID)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, post.ID)
}

func (r *postRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id = ?`, id)
	return err
}

func (r *postRepository) countPosts(ctx context.Context) (int64, error) {
	var total int64
	err := r.db.GetContext(ctx, &total, `SELECT COUNT(1) FROM posts`)
	return total, err
}
