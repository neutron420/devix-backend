package chat

import "time"

type CreateMessageRequest struct {
	ReceiverID string `json:"receiver_id" binding:"required"`
	Content    string `json:"content" binding:"required,max=2000"`
}

type MessageResponse struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	SenderID       string    `json:"sender_id"`
	Content        string    `json:"content"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
}

type ConversationResponse struct {
	ID        string          `json:"id"`
	OtherUser PublicUserResp `json:"other_user"`
	LastMsgAt time.Time       `json:"last_msg_at"`
	Unread    int             `json:"unread_count"`
}

type PublicUserResp struct {
	ID          string  `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}
