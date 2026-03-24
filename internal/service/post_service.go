package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"body-smith-be/internal/model"
	"body-smith-be/internal/repository"
)

var (
	ErrInvalidPostInput = errors.New("title and content are required")
	ErrPostNotFound     = errors.New("post not found")
)

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

type PostService interface {
	Create(ctx context.Context, req model.CreatePostRequest) (*model.Post, error)
	ListAdmin(ctx context.Context, page, perPage int) (*model.PostListResponse, error)
	ListPublic(ctx context.Context, page, perPage int) (*model.PostListResponse, error)
	GetBySlug(ctx context.Context, slug string) (*model.Post, error)
	Update(ctx context.Context, id int64, req model.UpdatePostRequest) (*model.Post, error)
	Delete(ctx context.Context, id int64) error
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) Create(ctx context.Context, req model.CreatePostRequest) (*model.Post, error) {
	if err := validatePostInput(req.Title, req.Content); err != nil {
		return nil, err
	}

	slug, err := s.uniqueSlug(ctx, req.Title, nil)
	if err != nil {
		return nil, err
	}

	post := &model.Post{
		Title:           strings.TrimSpace(req.Title),
		Slug:            slug,
		Content:         strings.TrimSpace(req.Content),
		Thumbnail:       cleanOptionalString(req.Thumbnail),
		MetaTitle:       strings.TrimSpace(req.MetaTitle),
		MetaDescription: strings.TrimSpace(req.MetaDescription),
	}

	if post.MetaTitle == "" {
		post.MetaTitle = post.Title
	}

	return s.postRepo.Create(ctx, post)
}

func (s *postService) ListAdmin(ctx context.Context, page, perPage int) (*model.PostListResponse, error) {
	return s.list(ctx, page, perPage, true)
}

func (s *postService) ListPublic(ctx context.Context, page, perPage int) (*model.PostListResponse, error) {
	return s.list(ctx, page, perPage, false)
}

func (s *postService) GetBySlug(ctx context.Context, slug string) (*model.Post, error) {
	post, err := s.postRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if post == nil {
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (s *postService) Update(ctx context.Context, id int64, req model.UpdatePostRequest) (*model.Post, error) {
	if err := validatePostInput(req.Title, req.Content); err != nil {
		return nil, err
	}

	existing, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrPostNotFound
	}

	slug, err := s.uniqueSlug(ctx, req.Title, &id)
	if err != nil {
		return nil, err
	}

	existing.Title = strings.TrimSpace(req.Title)
	existing.Slug = slug
	existing.Content = strings.TrimSpace(req.Content)
	existing.Thumbnail = cleanOptionalString(req.Thumbnail)
	existing.MetaTitle = strings.TrimSpace(req.MetaTitle)
	existing.MetaDescription = strings.TrimSpace(req.MetaDescription)
	if existing.MetaTitle == "" {
		existing.MetaTitle = existing.Title
	}

	return s.postRepo.Update(ctx, existing)
}

func (s *postService) Delete(ctx context.Context, id int64) error {
	existing, err := s.postRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrPostNotFound
	}

	return s.postRepo.Delete(ctx, id)
}

func (s *postService) list(ctx context.Context, page, perPage int, isAdmin bool) (*model.PostListResponse, error) {
	pagination := model.NewPagination(page, perPage)

	var (
		items []model.PostListItem
		total int64
		err   error
	)

	if isAdmin {
		items, total, err = s.postRepo.List(ctx, pagination)
	} else {
		items, total, err = s.postRepo.ListPublic(ctx, pagination)
	}
	if err != nil {
		return nil, err
	}

	return &model.PostListResponse{
		Items:      items,
		Page:       pagination.Page,
		PerPage:    pagination.PerPage,
		Total:      total,
		TotalPages: model.TotalPages(total, pagination.PerPage),
	}, nil
}

func (s *postService) uniqueSlug(ctx context.Context, title string, excludeID *int64) (string, error) {
	base := slugify(title)
	if base == "" {
		base = "post"
	}

	candidate := base
	for i := 1; ; i++ {
		exists, err := s.postRepo.SlugExists(ctx, candidate, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

func validatePostInput(title, content string) error {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
		return ErrInvalidPostInput
	}
	return nil
}

func slugify(input string) string {
	slug := strings.ToLower(strings.TrimSpace(input))
	slug = nonSlugChars.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}

func cleanOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
