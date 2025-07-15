package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/moabdelazem/app/internal/models"
)

// Ensure MockGuestBookService implements GuestBookServiceInterface
var _ GuestBookServiceInterface = (*MockGuestBookService)(nil)

// MockGuestBookService implements a mock version of the guest book service for testing
type MockGuestBookService struct {
	messages []models.GuestBookMessage
	nextID   int
}

func NewMockGuestBookService() *MockGuestBookService {
	return &MockGuestBookService{
		messages: []models.GuestBookMessage{
			{
				ID:        1,
				Name:      "John Doe",
				Email:     "john.doe@example.com",
				Message:   "Hello, this is a test message!",
				CreatedAt: time.Now().Add(-2 * time.Hour),
				UpdatedAt: time.Now().Add(-2 * time.Hour),
			},
			{
				ID:        2,
				Name:      "Jane Smith",
				Email:     "jane.smith@example.com",
				Message:   "Another test message for the guest book.",
				CreatedAt: time.Now().Add(-1 * time.Hour),
				UpdatedAt: time.Now().Add(-1 * time.Hour),
			},
		},
		nextID: 3,
	}
}

func (m *MockGuestBookService) InitializeDatabase(ctx context.Context) error {
	return nil
}

func (m *MockGuestBookService) CreateMessage(ctx context.Context, msg *models.CreateGuestBookMessage) (*models.GuestBookMessage, error) {
	if err := m.validateCreateMessage(msg); err != nil {
		return nil, err
	}

	newMessage := models.GuestBookMessage{
		ID:        m.nextID,
		Name:      msg.Name,
		Email:     msg.Email,
		Message:   msg.Message,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.messages = append(m.messages, newMessage)
	m.nextID++

	return &newMessage, nil
}

func (m *MockGuestBookService) GetMessages(ctx context.Context, page, pageSize int) ([]models.GuestBookMessage, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	total := len(m.messages)
	offset := (page - 1) * pageSize

	if offset >= total {
		return []models.GuestBookMessage{}, total, nil
	}

	end := offset + pageSize
	if end > total {
		end = total
	}

	// Return messages in reverse order (newest first)
	result := make([]models.GuestBookMessage, 0, end-offset)
	for i := total - 1; i >= 0; i-- {
		if len(result) >= pageSize {
			break
		}
		if i < total-offset {
			result = append(result, m.messages[i])
		}
	}

	return result, total, nil
}

func (m *MockGuestBookService) GetMessageByID(ctx context.Context, idStr string) (*models.GuestBookMessage, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, fmt.Errorf("invalid message ID")
	}

	for _, msg := range m.messages {
		if msg.ID == id {
			return &msg, nil
		}
	}

	return nil, fmt.Errorf("guest book message not found")
}

func (m *MockGuestBookService) validateCreateMessage(msg *models.CreateGuestBookMessage) error {
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
