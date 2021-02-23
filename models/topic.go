package models

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/websublime/courier/storage"
	"github.com/websublime/courier/storage/namespace"
)

type Topic struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Topic      string     `json:"name" db:"topic"`
	Audience   *Audience  `json:"audience" belongs_to:"audience"`
	AudienceID uuid.UUID  `json:"audienceId" db:"audience_id"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt  nulls.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

func (Topic) TableName() string {
	tableName := "topics"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "." + tableName
	}

	return tableName
}

func NewTopic(name string) (*Topic, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "Error generating unique id")
	}

	topic := &Topic{
		ID:    uid,
		Topic: name,
	}

	return topic, nil
}

func FindTopicsByAUdience(tx *storage.Connection, audienceID uuid.UUID) ([]*Topic, error) {
	topics := []*Topic{}
	q := tx.Q().Where("audience_id = ?", audienceID)

	err := q.All(topics)

	return topics, err
}

func FindTopicByID(tx *storage.Connection, id uuid.UUID) (*Topic, error) {
	topic := &Topic{}
	if err := tx.Q().Where("id = ?", id).First(topic); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, err
		}
		return nil, errors.Wrap(err, "error finding channel")
	}

	return topic, nil
}
