package models

import "time"

type Message struct {
	CreatedAt time.Time 
	Content string `json:"content"`
}