package domain

import "time"

type User struct {
	Name        string // Format: users/{user_id}
	DisplayName string
	Email       string
	CreateTime  time.Time
	UpdateTime  time.Time
}
