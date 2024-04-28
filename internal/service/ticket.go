package service

import (
	"context"
	"fmt"
	"github.com/gemm123/vkrf-ticket/helper"
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
	GetAllTicket() ([]model.TicketResponse, error)
	GetDetailTicket(ticketId string) (model.DetailTicketResponse, error)
	UpdateUserTicket(emailAssignee, ticketId, email string) error
	UpdateEditTicket(ticketId, email string, editTicket model.EditTicketRequest) error
	UpdateStatusTicket(ticketId, email, status string) error
	Summary(email string) ([]model.SummaryResponse, error)
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

func (s *ticketService) GetAllTicket() ([]model.TicketResponse, error) {
	tickets, err := s.ticketRepository.GetAllTicket()
	if err != nil {
		return nil, err
	}

	ticketResponses := make([]model.TicketResponse, 0)
	for _, ticket := range tickets {
		c := grpcserver.NewUserServiceClient(s.conn)
		userRequest := grpcserver.GetUserByUserIdRequest{
			UserId: ticket.UserId.String(),
		}
		resp, err := c.GetUserByUserId(context.Background(), &userRequest)
		if err != nil {
			log.Printf("Error: %v", err)
			return nil, err
		}

		ticketResponse := model.TicketResponse{
			Id:          ticket.Id,
			Title:       ticket.Title,
			Description: ticket.Description,
			Status:      ticket.Status,
			Point:       ticket.Point,
			User:        resp.User.Name,
			ProfilePic:  resp.User.ProfilePic,
		}

		ticketResponses = append(ticketResponses, ticketResponse)
	}

	return ticketResponses, nil
}

func (s *ticketService) GetDetailTicket(ticketId string) (model.DetailTicketResponse, error) {
	ticket, err := s.ticketRepository.GetTicketById(ticketId)
	if err != nil {
		return model.DetailTicketResponse{}, err
	}

	c := grpcserver.NewUserServiceClient(s.conn)
	userRequest := grpcserver.GetUserByUserIdRequest{
		UserId: ticket.UserId.String(),
	}
	resp, err := c.GetUserByUserId(context.Background(), &userRequest)
	if err != nil {
		log.Printf("Error: %v", err)
		return model.DetailTicketResponse{}, err
	}

	historyTickets, err := s.ticketRepository.GetHistoryTicketByTicketId(ticketId)
	if err != nil {
		return model.DetailTicketResponse{}, err
	}

	historyTicketResponses := make([]model.HistoryTicketResponse, 0)
	for _, ht := range historyTickets {
		htr := model.HistoryTicketResponse{
			Date:  ht.Date,
			Title: ht.Title,
			User:  ht.User,
		}
		historyTicketResponses = append(historyTicketResponses, htr)
	}

	dtr := model.DetailTicketResponse{
		Id:                    ticket.Id.String(),
		Username:              resp.User.Name,
		ProfilePic:            resp.User.ProfilePic,
		Title:                 ticket.Title,
		Description:           ticket.Description,
		Status:                ticket.Status,
		Point:                 ticket.Point,
		HistoryTicketResponse: historyTicketResponses,
	}

	return dtr, nil
}

func (s *ticketService) UpdateUserTicket(emailAssignee, ticketId, email string) error {
	resp, err := helper.GetUserByEmailGrpc(s.conn, emailAssignee)
	if err != nil {
		return err
	}

	resp2, err := helper.GetUserByEmailGrpc(s.conn, email)
	if err != nil {
		return err
	}

	historyTicket := model.HistoryTicket{
		Id:        uuid.New(),
		TicketId:  uuid.MustParse(ticketId),
		Date:      time.Now().Format("02 Jan 2006"),
		Title:     fmt.Sprintf("Change Assignees to %s", resp.User.Name),
		User:      resp2.User.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.ticketRepository.UpdateUserTicket(resp.User.Id, ticketId, historyTicket); err != nil {
		return err
	}

	return nil
}

func (s *ticketService) UpdateEditTicket(ticketId, email string, editTicket model.EditTicketRequest) error {
	resp, err := helper.GetUserByEmailGrpc(s.conn, email)
	if err != nil {
		return err
	}

	historyTicket := model.HistoryTicket{
		Id:        uuid.New(),
		TicketId:  uuid.MustParse(ticketId),
		Date:      time.Now().Format("02 Jan 2006"),
		Title:     fmt.Sprintf("Edited by %s", resp.User.Name),
		User:      resp.User.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.ticketRepository.UpdateEditTicket(editTicket, ticketId, historyTicket); err != nil {
		return err
	}

	return nil
}

func (s *ticketService) UpdateStatusTicket(ticketId, email, status string) error {
	resp, err := helper.GetUserByEmailGrpc(s.conn, email)
	if err != nil {
		return err
	}

	historyTicket := model.HistoryTicket{
		Id:        uuid.New(),
		TicketId:  uuid.MustParse(ticketId),
		Date:      time.Now().Format("02 Jan 2006"),
		Title:     fmt.Sprintf("%s Change status to %s", resp.User.Name, status),
		User:      resp.User.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.ticketRepository.UpdateStatusTicket(status, ticketId, historyTicket); err != nil {
		return err
	}

	return nil
}

func (s *ticketService) Summary(email string) ([]model.SummaryResponse, error) {
	resp, err := helper.GetUserByEmailGrpc(s.conn, email)
	if err != nil {
		return nil, err
	}

	var summaryResponses []model.SummaryResponse
	result, err := s.ticketRepository.CountTicketGroupByStatus(resp.User.Id)
	if err != nil {
		return nil, err
	}
	for _, r := range result {
		summaryResponses = append(summaryResponses, model.SummaryResponse{
			Status:    r.Status,
			TotalTask: r.Count,
		})

	}

	result2, err := s.ticketRepository.SumTicketGroupByStatus(resp.User.Id)
	if err != nil {
		return nil, err
	}
	for _, r := range result2 {
		for i, sr := range summaryResponses {
			if sr.Status == r.Status {
				summaryResponses[i].Point = r.Point
			}
		}
	}

	return summaryResponses, nil
}
