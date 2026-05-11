package chat

import (
	"context"
	"time"

	apperrors "devix-backend/internal/errors"
	"devix-backend/internal/models"
	"devix-backend/internal/modules/user"
	"devix-backend/internal/modules/websocket"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Service struct {
	repo        *Repository
	userService *user.Service
	wsService   *websocket.Service
	log         zerolog.Logger
}

func NewService(repo *Repository, userService *user.Service, wsService *websocket.Service, log zerolog.Logger) *Service {
	return &Service{
		repo:        repo,
		userService: userService,
		wsService:   wsService,
		log:         log.With().Str("module", "chat").Logger(),
	}
}

func (s *Service) SendMessage(ctx context.Context, senderID uuid.UUID, req *CreateMessageRequest) (*MessageResponse, error) {
	receiverID, err := uuid.Parse(req.ReceiverID)
	if err != nil {
		return nil, apperrors.BadRequest("Invalid receiver ID")
	}

	if senderID == receiverID {
		return nil, apperrors.BadRequest("You cannot message yourself")
	}

	// 1. Get or Create Conversation
	conv, err := s.repo.GetConversation(ctx, senderID, receiverID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}
	if conv == nil {
		conv, err = s.repo.CreateConversation(ctx, senderID, receiverID)
		if err != nil {
			return nil, apperrors.Internal(err)
		}
	}

	// 2. Create Message
	msg := &models.Message{
		ID:             uuid.New(),
		ConversationID: conv.ID,
		SenderID:       senderID,
		Content:        req.Content,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.CreateMessage(ctx, msg); err != nil {
		return nil, apperrors.Internal(err)
	}

	res := &MessageResponse{
		ID:             msg.ID.String(),
		ConversationID: msg.ConversationID.String(),
		SenderID:       msg.SenderID.String(),
		Content:        msg.Content,
		IsRead:         msg.IsRead,
		CreatedAt:      msg.CreatedAt,
	}

	// 3. Real-time Notification
	s.wsService.NotifyUser(ctx, receiverID, "new_message", res)

	return res, nil
}

func (s *Service) GetConversations(ctx context.Context, userID uuid.UUID) ([]ConversationResponse, error) {
	convs, err := s.repo.GetConversations(ctx, userID)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	res := make([]ConversationResponse, 0, len(convs))
	for _, c := range convs {
		otherUserID := c.User1ID
		if otherUserID == userID {
			otherUserID = c.User2ID
		}

		otherUser, _ := s.userService.GetPublicProfileByID(ctx, otherUserID)
		
		cr := ConversationResponse{
			ID:        c.ID.String(),
			LastMsgAt: c.LastMsgAt,
		}
		if otherUser != nil {
			cr.OtherUser = PublicUserResp{
				ID:          otherUser.ID,
				Username:    otherUser.Username,
				DisplayName: otherUser.DisplayName,
				AvatarURL:   otherUser.AvatarURL,
			}
		}
		res = append(res, cr)
	}

	return res, nil
}

func (s *Service) GetMessages(ctx context.Context, userID, convID uuid.UUID, limit, offset int) ([]MessageResponse, error) {
	// Mark as read first
	_ = s.repo.MarkAsRead(ctx, convID, userID)

	msgs, err := s.repo.GetMessages(ctx, convID, limit, offset)
	if err != nil {
		return nil, apperrors.Internal(err)
	}

	res := make([]MessageResponse, 0, len(msgs))
	for _, m := range msgs {
		res = append(res, MessageResponse{
			ID:             m.ID.String(),
			ConversationID: m.ConversationID.String(),
			SenderID:       m.SenderID.String(),
			Content:        m.Content,
			IsRead:         m.IsRead,
			CreatedAt:      m.CreatedAt,
		})
	}
	return res, nil
}

func (s *Service) NotifyTyping(userID, otherUserID uuid.UUID, isTyping bool) {
	s.wsService.BroadcastTyping(userID, otherUserID, isTyping)
}

func (s *Service) MarkAsRead(ctx context.Context, userID, convID uuid.UUID) error {
	if err := s.repo.MarkAsRead(ctx, convID, userID); err != nil {
		return err
	}

	// Notify the sender that their messages were read
	s.wsService.NotifyUser(ctx, userID, "chat:read", map[string]string{
		"conversation_id": convID.String(),
	})

	return nil
}
