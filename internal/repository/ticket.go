package repository

import (
	"context"
	"github.com/gemm123/vkrf-ticket/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ticketRepository struct {
	db *pgxpool.Pool
}

type TicketRepository interface {
	CreateTicket(ticket model.Ticket, historyTicket model.HistoryTicket) error
}

func NewTicketRepository(db *pgxpool.Pool) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) CreateTicket(ticket model.Ticket, historyTicket model.HistoryTicket) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `INSERT INTO tickets (id, user_id, title, description, status, point, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(context.Background(), query,
		ticket.Id, ticket.UserId, ticket.Title, ticket.Description, ticket.Status, ticket.Point,
		ticket.CreatedAt, ticket.UpdatedAt)
	if err != nil {
		return err
	}

	query2 := `INSERT INTO history_ticket (id, ticket_id, date, title, "user", created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(context.Background(), query2, historyTicket.Id, historyTicket.TicketId,
		historyTicket.Date, historyTicket.Title, historyTicket.User, historyTicket.CreatedAt, historyTicket.UpdatedAt)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
