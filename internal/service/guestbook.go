package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/moabdelazem/app/internal/models"
	"github.com/moabdelazem/app/internal/repository"
)

type GuestBookService struct {
	repo *repository.GuestBookRepository
}

func NewGuestBookService(repo *repository.GuestBookRepository) *GuestBookService {
	return &GuestBookService{repo: repo}
}

func (s *GuestBookService) InitializeDatabase(ctx context.Context) error {
	return s.repo.CreateTable(ctx)
}

func (s *GuestBookService) CreateMessage(ctx context.Context, msg *models.CreateGuestBookMessage) (*models.GuestBookMessage, error) {
	if err := s.validateCreateMessage(msg); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, msg)
}

func (s *GuestBookService) GetMessages(ctx context.Context, page, pageSize int) ([]models.GuestBookMessage, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	messages, err := s.repo.GetAll(ctx, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

func (s *GuestBookService) GetMessageByID(ctx context.Context, idStr string) (*models.GuestBookMessage, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid message ID")
	}

	return s.repo.GetByID(ctx, id)
}

func (s *GuestBookService) validateCreateMessage(msg *models.CreateGuestBookMessage) error {
	if len(msg.Name) < 2 || len(msg.Name) > 100 {
		return fmt.Errorf("name must be between 2 and 100 characters")
	}

	if len(msg.Email) == 0 || len(msg.Email) > 255 {
		return fmt.Errorf("email must be between 1 and 255 characters")
	}

	if len(msg.Message) < 10 || len(msg.Message) > 1000 {
		return fmt.Errorf("message must be between 10 and 1000 characters")
	}

	return nil
}
