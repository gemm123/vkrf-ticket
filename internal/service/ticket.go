package service

import (
	"context"
	grpcserver "github.com/gemm123/vkrf-ticket/internal/grpc"
	"github.com/gemm123/vkrf-ticket/internal/model"
	"github.com/gemm123/vkrf-ticket/internal/repository"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"time"
)

type ticketService struct {
	ticketRepository repository.TicketRepository
	conn             *grpc.ClientConn
}

type TicketService interface {
	CreateTicket(ticket model.TicketRequest, email string) error
}

func NewTicketService(ticketRepository repository.TicketRepository, conn *grpc.ClientConn) TicketService {
	return &ticketService{
		ticketRepository: ticketRepository,
		conn:             conn,
	}
}

func (s *ticketService) CreateTicket(ticket model.TicketRequest, email string) error {
	c := grpcserver.NewUserServiceClient(s.conn)
	userRequest := grpcserver.GetUserByEmailRequest{
		Email: email,
	}
	resp, err := c.GetUserByEmail(context.Background(), &userRequest)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}

	userId, _ := uuid.Parse(resp.User.Id)
	t := model.Ticket{
		Id:          uuid.New(),
		UserId:      userId,
		Title:       ticket.Title,
		Description: ticket.Description,
		Status:      ticket.Status,
		Point:       ticket.Point,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ht := model.HistoryTicket{
		Id:        uuid.New(),
		TicketId:  t.Id,
		Date:      t.CreatedAt.Format("02 Jan 2006"),
		Title:     "Ticket Created",
		User:      resp.User.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.ticketRepository.CreateTicket(t, ht); err != nil {
		return err
	}

	return nil
}
