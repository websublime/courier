package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/websublime/courier/storage/namespace"
)

type Message struct {
	ID        uuid.UUID       `json:"id" db:"id"`
	Message   json.RawMessage `json:"message" db:"message"`
	Topic     *Topic          `json:"topic,omitempty" belongs_to:"topics"`
	TopicID   uuid.UUID       `json:"topicId" db:"topic_id"`
	CreatedAt time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time       `json:"updatedAt" db:"updated_at"`
	DeletedAt nulls.Time      `json:"deletedAt,omitempty" db:"deleted_at"`
}

func (Message) TableName() string {
	tableName := "messages"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "." + tableName
	}

	return tableName
}

func NewMessage(msg json.RawMessage, topicID uuid.UUID) (*Message, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "Error generating unique id")
	}

	message := &Message{
		ID:      uid,
		Message: msg,
		TopicID: topicID,
	}

	return message, nil
}
