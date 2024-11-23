package domain

import (
	"bytes"
	"encoding/gob"
	"time"
)

type User struct {
	// TODO: rename 'Name' to 'UserID' and make it into a struct
	Name        string // Format: users/{user_id}
	DisplayName string
	Email       string
	CreateTime  time.Time
	UpdateTime  time.Time
	DeleteTime  time.Time
}

func (u *User) Copy() (*User, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)

	err := enc.Encode(u)
	if err != nil {
		return nil, err
	}

	var userCopy User
	err = dec.Decode(&userCopy)
	if err != nil {
		return nil, err
	}
	return &userCopy, nil
}
