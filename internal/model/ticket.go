package model

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	EventID     uuid.UUID  `json:"event_id" db:"event_id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	Email       string     `json:"email" db:"email"`
	Name        string     `json:"name" db:"name"`
	TicketCode  string     `json:"ticket_code" db:"ticket_code"`
	Status      string     `json:"status" db:"status"`
	Price       float64    `json:"price" db:"price"`
	ReservedAt  time.Time  `json:"reserved_at" db:"reserved_at"`
	CheckedInAt *time.Time `json:"checked_in_at,omitempty" db:"checked_in_at"`
	CheckedInBy *uuid.UUID `json:"checked_in_by,omitempty" db:"checked_in_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateTicketRequest struct {
	EventID uuid.UUID `json:"event_id" binding:"required"`
	UserID  uuid.UUID `json:"user_id" binding:"required"`
	Email   string    `json:"email" binding:"required,email"`
	Name    string    `json:"name" binding:"required"`
}

const (
	TicketStatusReserved  = "reserved"
	TicketStatusConfirmed = "confirmed"
	TicketStatusCancelled = "cancelled"
	TicketStatusUsed      = "used"
)
