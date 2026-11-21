package server

import (
	"gotickets/internal/domain/booking"
	"gotickets/internal/domain/event"
	"gotickets/internal/domain/user"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func seedDatabase(db *gorm.DB) error {
	if err := db.Exec("UPDATE users SET google_id = NULL WHERE google_id = ''").Error; err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	users := []user.User{
		{Name: "Super Admin", Email: "admin@gotix.local", Password: string(passwordHash), Role: user.RoleSuperAdmin},
		{Name: "Organizer", Email: "organizer@gotix.local", Password: string(passwordHash), Role: user.RoleOrganizer},
		{Name: "Alex Customer", Email: "alex@gotix.local", Password: string(passwordHash), Role: user.RoleUser},
		{Name: "Sam Customer", Email: "sam@gotix.local", Password: string(passwordHash), Role: user.RoleUser},
	}
	for i := range users {
		var count int64
		if err := db.Model(&user.User{}).Where("email = ?", users[i].Email).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.Create(&users[i]).Error; err != nil {
				return err
			}
		}
	}

	now := time.Now().UTC()
	events := []event.Event{
		{Title: "Go Conference 2026", Description: "A practical backend engineering conference for Go developers.", Location: "Berlin", StartsAt: now.AddDate(0, 1, 7), TotalTickets: 300, AvailableTickets: 300, Price: 75},
		{Title: "Cloud Native Summit", Description: "Cloud platforms, Kubernetes, observability, and platform engineering.", Location: "Amsterdam", StartsAt: now.AddDate(0, 2, 1), TotalTickets: 180, AvailableTickets: 180, Price: 120},
		{Title: "Indie Music Night", Description: "An evening featuring emerging independent artists.", Location: "London", StartsAt: now.AddDate(0, 0, 21), TotalTickets: 500, AvailableTickets: 500, Price: 35},
	}
	for i := range events {
		var count int64
		if err := db.Model(&event.Event{}).Where("title = ?", events[i].Title).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			if err := db.Create(&events[i]).Error; err != nil {
				return err
			}
		}
	}

	var customer user.User
	if err := db.Where("email = ?", "alex@gotix.local").First(&customer).Error; err != nil {
		return err
	}
	var seededEvent event.Event
	if err := db.Where("title = ?", "Go Conference 2026").First(&seededEvent).Error; err != nil {
		return err
	}
	var bookingCount int64
	if err := db.Model(&booking.Booking{}).Where("booking_code LIKE ?", "GT-SEED-%").Count(&bookingCount).Error; err != nil {
		return err
	}
	if bookingCount == 0 {
		seedBooking := booking.Booking{
			UserID: customer.ID, EventID: seededEvent.ID, Quantity: 2,
			TotalPrice: 2 * seededEvent.Price, Status: booking.BookingConfirmed,
			BookingCode: "GT-SEED-001",
		}
		return db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&seedBooking).Error; err != nil {
				return err
			}
			seededEvent.AvailableTickets -= seedBooking.Quantity
			return tx.Save(&seededEvent).Error
		})
	}
	return nil
}
