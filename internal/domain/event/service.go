package event

import (
	"errors"
	"gotickets/internal/domain/event/dto"
)

var ErrInvalidTicketCount = errors.New("total tickets cannot be lower than sold tickets")

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) CreateEvent(req dto.CreateEventRequest) (*dto.EventResponse, error) {
	event := Event{
		Title:            req.Title,
		Description:      req.Description,
		Location:         req.Location,
		StartsAt:         req.StartsAt,
		TotalTickets:     req.TotalTickets,
		AvailableTickets: req.TotalTickets,
		Price:            req.Price,
	}

	if err := s.repo.Create(&event); err != nil {
		return nil, err
	}

	return event.ToResponse(), nil

}

func (s *service) GetEvents() ([]dto.EventResponse, error) {
	events, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	// responses := make([]dto.EventResponse, len(events))

	var responses []dto.EventResponse

	for _, event := range events {
		responses = append(responses, *event.ToResponse())
	}

	return responses, nil
}

func (s *service) GetEventByID(eventId uint) (*dto.EventResponse, error) {
	event, err := s.repo.GetByID(eventId)

	if err != nil {
		return nil, err
	}

	return event.ToResponse(), nil
}

func (s *service) UpdateEvent(eventId uint, req dto.UpdateEventRequest) (*dto.EventResponse, error) {
	event, err := s.repo.GetByID(eventId) // getting existing event by the id first
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		event.Title = req.Title
	}

	if req.Description != "" {
		event.Description = req.Description
	}

	if req.Location != "" {
		event.Location = req.Location
	}

	if !req.StartsAt.IsZero() {
		event.StartsAt = req.StartsAt
	}

	if req.Price != 0 {
		event.Price = req.Price
	}

	if req.TotalTickets != 0 {
		soldTickets := event.TotalTickets - event.AvailableTickets
		if req.TotalTickets < soldTickets {
			return nil, ErrInvalidTicketCount
		}
		event.TotalTickets = req.TotalTickets
		event.AvailableTickets = req.TotalTickets - soldTickets
	}

	if err := s.repo.Update(event); err != nil {
		return nil, err
	}

	return event.ToResponse(), nil
}

func (s *service) DeleteEvent(eventId uint) error {
	return s.repo.Delete(eventId)
}
