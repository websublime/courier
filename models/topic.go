package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/websublime/courier/storage"
	"github.com/websublime/courier/storage/namespace"
)

type Topic struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Topic      string     `json:"topic" db:"topic"`
	Audience   *Audience  `json:"audience,omitempty" belongs_to:"audience"`
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

func ContainsTopic(topics []*Topic, value string) bool {
	for _, topic := range topics {
		if topic.Topic == value {
			return true
		}
	}

	return false
}

func FindTopicsByAudienceID(tx *storage.Connection, audienceID uuid.UUID) ([]*Topic, error) {
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
		return nil, errors.Wrap(err, "error finding topic")
	}

	return topic, nil
}

func FindTopicByName(tx *storage.Connection, name string) (*Topic, error) {
	topic := &Topic{}
	if err := tx.Q().Where("topic = ?", name).First(topic); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, err
		}
		return nil, errors.Wrap(err, fmt.Sprintf("error finding topic: %s", name))
	}

	return topic, nil
}

func FindTopicByNameAndAudienceID(tx *storage.Connection, name string, audienceID uuid.UUID) (*Topic, error) {
	topic := &Topic{}
	if err := tx.Q().Where("topic = ?", name).Where("audience_id = ?", audienceID).First(topic); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, err
		}
		return nil, errors.Wrap(err, fmt.Sprintf("error finding topic: %s for audience: %s", name, audienceID))
	}

	return topic, nil
}
