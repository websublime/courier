package models

import (
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
)

type Hook struct {
	Action  string      `json:"action"`
	Topic   string      `json:"topic"`
	Message interface{} `json:"message"`
}

func NewHook(action string, topic string, message interface{}) *Hook {
	hook := &Hook{
		Action:  action,
		Topic:   topic,
		Message: message,
	}

	return hook
}

func (hook *Hook) Validate() *validate.Errors {
	return validate.Validate(
		&validators.StringIsPresent{Field: hook.Action, Name: "Action", Message: "Action name is missing"},
		&validators.StringIsPresent{Field: hook.Topic, Name: "Topic", Message: "Topic name is missing"},
	)
}
