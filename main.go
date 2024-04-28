package main

import (
	"github.com/gemm123/vkrf-ticket/config"
	"github.com/gemm123/vkrf-ticket/internal/controller"
	"github.com/gemm123/vkrf-ticket/internal/repository"
	"github.com/gemm123/vkrf-ticket/internal/service"
	"github.com/gemm123/vkrf-ticket/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err)
	}

	db := config.InitConnPool()
	defer db.Close()

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	ticketRepository := repository.NewTicketRepository(db)

	ticketService := service.NewTicketService(ticketRepository, conn)

	tickerController := controller.NewTicketController(ticketService)

	app := fiber.New()

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World!")
	})

	api := app.Group("/api")
	v1 := api.Group("/v1", middleware.Middleware)
	v1.Get("/tickets", tickerController.GetAllTicket)
	v1.Post("/tickets/create", tickerController.CreateTicket)
	v1.Get("/tickets/:ticketId/", tickerController.GetDetailTicket)
	v1.Put("/tickets/:ticketId/edit/assignee", tickerController.UpdateUserTicket)
	v1.Put("/tickets/:ticketId/edit", tickerController.UpdateEditTicket)

	app.Listen(":3001")
}
