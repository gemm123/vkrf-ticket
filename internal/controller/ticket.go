package controller

import (
	"github.com/gemm123/vkrf-ticket/internal/model"
	"github.com/gemm123/vkrf-ticket/internal/service"
	"github.com/gofiber/fiber/v2"
)

type ticketController struct {
	ticketService service.TicketService
}

type TicketController interface {
	CreateTicket(ctx *fiber.Ctx) error
	GetAllTicket(ctx *fiber.Ctx) error
}

func NewTicketController(ticketService service.TicketService) TicketController {
	return &ticketController{ticketService: ticketService}
}

func (c *ticketController) CreateTicket(ctx *fiber.Ctx) error {
	ticket := model.TicketRequest{}
	if err := ctx.BodyParser(&ticket); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	email := ctx.Locals("email").(string)

	if err := c.ticketService.CreateTicket(ticket, email); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create ticket",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Ticket created",
		"status":  fiber.StatusCreated,
	})
}

func (c *ticketController) GetAllTicket(ctx *fiber.Ctx) error {
	tickets, err := c.ticketService.GetAllTicket()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get all tickets",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"status":  fiber.StatusOK,
		"data":    tickets,
	})
}
