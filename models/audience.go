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

type Topics []Topic

type Audience struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Topics    Topics     `json:"topics,omitempty" has_many:"topics"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt nulls.Time `json:"deletedAt,omitempty" db:"deleted_at"`
}

func (Audience) TableName() string {
	tableName := "audiences"

	if namespace.GetNamespace() != "" {
		return namespace.GetNamespace() + "." + tableName
	}

	return tableName
}

func NewAudience(name string) (*Audience, error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, errors.Wrap(err, "Error generating unique id")
	}

	aud := &Audience{
		ID:   uid,
		Name: name,
	}

	return aud, nil
}

func FindAudience(tx *storage.Connection, name string) (*Audience, error) {
	aud := &Audience{}
	if err := tx.Where("name = ?", name).First(aud); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, errors.Wrap(err, "Audience not found")
		}

		return nil, errors.Wrap(err, err.Error())
	}

	return aud, nil
}
