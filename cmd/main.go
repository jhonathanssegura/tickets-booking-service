package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonathanssegura/ticket-reservation/internal/handler"
)

func main() {
	r := gin.Default()
	r.POST("/tickets", handler.ReserveTicket)
	r.Run(":8080")
}
