package controller

import (
	"github.com/gemm123/vkrf-ticket/internal/model"
	"github.com/gemm123/vkrf-ticket/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ticketController struct {
	ticketService service.TicketService
	validate      *validator.Validate
}

type TicketController interface {
	CreateTicket(ctx *fiber.Ctx) error
	GetAllTicket(ctx *fiber.Ctx) error
	GetDetailTicket(ctx *fiber.Ctx) error
	UpdateUserTicket(ctx *fiber.Ctx) error
	UpdateEditTicket(ctx *fiber.Ctx) error
	UpdateStatusTicket(ctx *fiber.Ctx) error
	Summary(ctx *fiber.Ctx) error
	Performance(ctx *fiber.Ctx) error
}

func NewTicketController(ticketService service.TicketService, validate *validator.Validate) TicketController {
	return &ticketController{ticketService: ticketService, validate: validate}
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
	err := c.validate.Struct(ticket)
	if err != nil {
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

func (c *ticketController) GetDetailTicket(ctx *fiber.Ctx) error {
	ticketId := ctx.Params("ticketId")
	detailTicket, err := c.ticketService.GetDetailTicket(ticketId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get detail ticket",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"status":  fiber.StatusOK,
		"data":    detailTicket,
	})
}

func (c *ticketController) UpdateUserTicket(ctx *fiber.Ctx) error {
	ticketId := ctx.Params("ticketId")
	email := ctx.Locals("email").(string)
	var jsonData map[string]interface{}
	if err := ctx.BodyParser(&jsonData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	err := c.validate.Var(jsonData["email"], "required,email")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	if err := c.ticketService.UpdateUserTicket(jsonData["email"].(string), ticketId, email); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update ticket",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Ticket updated",
		"status":  fiber.StatusOK,
	})
}

func (c *ticketController) UpdateEditTicket(ctx *fiber.Ctx) error {
	ticketId := ctx.Params("ticketId")
	email := ctx.Locals("email").(string)
	editTicket := model.EditTicketRequest{}
	if err := ctx.BodyParser(&editTicket); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	if err := c.validate.Struct(editTicket); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	if err := c.ticketService.UpdateEditTicket(ticketId, email, editTicket); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update ticket",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Ticket updated",
		"status":  fiber.StatusOK,
	})
}

func (c *ticketController) UpdateStatusTicket(ctx *fiber.Ctx) error {
	ticketId := ctx.Params("ticketId")
	email := ctx.Locals("email").(string)
	var jsonData map[string]interface{}

	if err := ctx.BodyParser(&jsonData); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	if err := c.validate.Var(jsonData["status"], "required"); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
			"status":  fiber.StatusBadRequest,
			"error":   err.Error(),
		})
	}

	if err := c.ticketService.UpdateStatusTicket(ticketId, email, jsonData["status"].(string)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update ticket",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Ticket updated",
		"status":  fiber.StatusOK,
	})
}

func (c *ticketController) Summary(ctx *fiber.Ctx) error {
	email := ctx.Locals("email").(string)
	summary, err := c.ticketService.Summary(email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get summary",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"status":  fiber.StatusOK,
		"data":    summary,
	})
}

func (c *ticketController) Performance(ctx *fiber.Ctx) error {
	email := ctx.Locals("email").(string)
	performance, err := c.ticketService.Performance(email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get performance",
			"status":  fiber.StatusInternalServerError,
			"error":   err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"status":  fiber.StatusOK,
		"data":    performance,
	})
}
