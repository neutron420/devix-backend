package chat

import (
	"context"
	"errors"
	"time"

	"devix-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetConversation(ctx context.Context, u1, u2 uuid.UUID) (*models.Conversation, error) {
	var conv models.Conversation
	err := r.db.WithContext(ctx).Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		u1, u2, u2, u1,
	).First(&conv).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &conv, err
}

func (r *Repository) CreateConversation(ctx context.Context, u1, u2 uuid.UUID) (*models.Conversation, error) {
	conv := &models.Conversation{
		ID:        uuid.New(),
		User1ID:   u1,
		User2ID:   u2,
		LastMsgAt: time.Now(),
	}
	err := r.db.WithContext(ctx).Create(conv).Error
	return conv, err
}

func (r *Repository) CreateMessage(ctx context.Context, msg *models.Message) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(msg).Error; err != nil {
			return err
		}
		return tx.Model(&models.Conversation{}).
			Where("id = ?", msg.ConversationID).
			Update("last_msg_at", msg.CreatedAt).Error
	})
}

func (r *Repository) GetMessages(ctx context.Context, convID uuid.UUID, limit, offset int) ([]models.Message, error) {
	var msgs []models.Message
	err := r.db.WithContext(ctx).
		Where("conversation_id = ?", convID).
		Order("created_at desc").
		Limit(limit).Offset(offset).
		Find(&msgs).Error
	return msgs, err
}

func (r *Repository) GetConversations(ctx context.Context, userID uuid.UUID) ([]models.Conversation, error) {
	var convs []models.Conversation
	err := r.db.WithContext(ctx).
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("last_msg_at desc").
		Find(&convs).Error
	return convs, err
}

func (r *Repository) MarkAsRead(ctx context.Context, convID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Message{}).
		Where("conversation_id = ? AND sender_id != ? AND is_read = ?", convID, userID, false).
		Update("is_read", true).Error
}

func (r *Repository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Message{}).
		Joins("JOIN conversations ON conversations.id = messages.conversation_id").
		Where("(conversations.user1_id = ? OR conversations.user2_id = ?) AND messages.sender_id != ? AND messages.is_read = ?", userID, userID, userID, false).
		Count(&count).Error
	return count, err
}
