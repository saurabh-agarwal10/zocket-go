package model

import "time"

type UserInfo struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Mobile    int       `json:"mobile"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
