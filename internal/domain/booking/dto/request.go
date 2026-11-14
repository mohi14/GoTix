package dto

type CreateBookingRequest struct {
	EventID  uint `json:"event_id" validate:"required"`
	Quantity int  `json:"quantity" validate:"required,gt=0"`
}
