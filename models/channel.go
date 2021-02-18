package models

import (
	"time"

	"github.com/gobuffalo/pop/nulls"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
	"github.com/websublime/courier/storage"
	"github.com/websublime/courier/storage/namespace"
)

type Channel struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	Name       string     `json:"name" db:"name"`
	Audience   *Audience  `belongs_to:"audiences"`
	AudienceID uuid.UUID  `json:"audienceId"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt  nulls.Time `json:"deleteadAt,omitempty" db:"deleted_at"`
}

func (Channel) TableName() string {
	tableName := "channels"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "." + tableName
	}

	return tableName
}

func NewChannel(name string) (*Channel, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "Error generating unique id")
	}

	channels := &Channel{
		ID:   uid,
		Name: name,
	}

	return channels, nil
}

func FindChannels(tx *storage.Connection, aud uuid.UUID) {
	tx.Q().Where("audience_id = ?", aud)
}
