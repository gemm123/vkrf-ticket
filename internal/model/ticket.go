package model

import (
	"github.com/google/uuid"
	"time"
)

type Ticket struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Point       int       `json:"point"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Point       int    `json:"point"`
}

type TicketResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	User        string `json:"user"`
	ProfilePic  string `json:"profile_pic"`
}

type HistoryTicket struct {
	Id        uuid.UUID `json:"id"`
	TicketId  uuid.UUID `json:"ticket_id"`
	Date      string    `json:"date"`
	Title     string    `json:"title"`
	User      string    `json:"user"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
