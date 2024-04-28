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
	GetAllTicket() ([]model.Ticket, error)
	GetHistoryTicketByTicketId(ticketId string) ([]model.HistoryTicket, error)
	GetTicketById(ticketId string) (model.Ticket, error)
	UpdateUserTicket(userId, ticketId string, historyTicket model.HistoryTicket) error
	UpdateEditTicket(editTicket model.EditTicketRequest, ticketId string, historyTicket model.HistoryTicket) error
	UpdateStatusTicket(status, ticketId string, historyTicket model.HistoryTicket) error
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

func (r *ticketRepository) GetAllTicket() ([]model.Ticket, error) {
	query := `SELECT * FROM tickets`
	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		ticket := model.Ticket{}
		err = rows.Scan(&ticket.Id, &ticket.UserId, &ticket.Title, &ticket.Description, &ticket.Status, &ticket.Point, &ticket.CreatedAt, &ticket.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (r *ticketRepository) GetHistoryTicketByTicketId(ticketId string) ([]model.HistoryTicket, error) {
	query := `SELECT * FROM history_ticket WHERE ticket_id = $1`
	rows, err := r.db.Query(context.Background(), query, ticketId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var historyTickets []model.HistoryTicket
	for rows.Next() {
		historyTicket := model.HistoryTicket{}
		err = rows.Scan(&historyTicket.Id, &historyTicket.TicketId, &historyTicket.Date, &historyTicket.Title,
			&historyTicket.User, &historyTicket.CreatedAt, &historyTicket.UpdatedAt)
		if err != nil {
			return nil, err
		}
		historyTickets = append(historyTickets, historyTicket)
	}

	return historyTickets, nil
}

func (r *ticketRepository) GetTicketById(ticketId string) (model.Ticket, error) {
	query := `SELECT * FROM tickets WHERE id = $1`
	row := r.db.QueryRow(context.Background(), query, ticketId)

	ticket := model.Ticket{}
	err := row.Scan(&ticket.Id, &ticket.UserId, &ticket.Title, &ticket.Description, &ticket.Status, &ticket.Point, &ticket.CreatedAt, &ticket.UpdatedAt)
	if err != nil {
		return model.Ticket{}, err
	}

	return ticket, nil
}

func (r *ticketRepository) UpdateUserTicket(userId, ticketId string, historyTicket model.HistoryTicket) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `UPDATE tickets SET user_id = $1, updated_at = $2 WHERE id = $3`
	_, err = tx.Exec(context.Background(), query, userId, historyTicket.UpdatedAt, ticketId)
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

func (r *ticketRepository) UpdateEditTicket(editTicket model.EditTicketRequest, ticketId string, historyTicket model.HistoryTicket) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `UPDATE tickets SET title = $1, description = $2, point = $3, updated_at = $4 WHERE id = $5`
	_, err = tx.Exec(context.Background(), query, editTicket.Title, editTicket.Description, editTicket.Point, historyTicket.UpdatedAt, ticketId)
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

func (r *ticketRepository) UpdateStatusTicket(status, ticketId string, historyTicket model.HistoryTicket) error {
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `UPDATE tickets SET status = $1, updated_at = $2 WHERE id = $3`
	_, err = tx.Exec(context.Background(), query, status, historyTicket.UpdatedAt, ticketId)
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
