package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/handler"
	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/model"
	"github.com/matveevsa/ai-for-developers-project-386/backend/internal/store"
)

func main() {
	s := store.New()

	demoTypes := []model.EventTypeCreate{
		{Name: "15-минутная встреча", Duration: 15, Description: "Короткая встреча для быстрых вопросов"},
		{Name: "30-минутная встреча", Duration: 30, Description: "Стандартная встреча для обсуждения деталей"},
		{Name: "Часовая консультация", Duration: 60, Description: "Подробное обсуждение проекта"},
	}
	createdTypes := make([]model.EventType, 0, len(demoTypes))
	for _, d := range demoTypes {
		et := model.EventType{
			ID:          s.NextEventID(),
			OwnerID:     "owner-1",
			Name:        d.Name,
			Description: d.Description,
			Duration:    d.Duration,
		}
		s.CreateEventType(et)
		createdTypes = append(createdTypes, et)
	}

	now := time.Now().UTC()
	today := now.Truncate(24 * time.Hour)

	pastBookings := []struct {
		eventType  model.EventType
		guestName  string
		guestEmail string
		daysAgo    int
		hour       int
		minute     int
	}{
		{createdTypes[0], "Алексей Смирнов", "alex@example.com", 2, 10, 0},
		{createdTypes[0], "Мария Петрова", "maria@example.com", 2, 14, 30},
		{createdTypes[1], "Дмитрий Козлов", "dmitry@example.com", 1, 11, 0},
		{createdTypes[2], "Елена Соколова", "elena@example.com", 3, 9, 0},
	}
	for _, pb := range pastBookings {
		slotStart := today.AddDate(0, 0, -pb.daysAgo).Add(time.Duration(pb.hour)*time.Hour + time.Duration(pb.minute)*time.Minute)
		slotID := fmt.Sprintf("%s-%d", pb.eventType.ID, slotStart.Unix())
		s.CreateBooking(model.Booking{
			ID:          s.NextBookingID(),
			SlotID:      slotID,
			EventTypeID: pb.eventType.ID,
			GuestName:   pb.guestName,
			GuestEmail:  pb.guestEmail,
		})
	}

	upcomingBookings := []struct {
		eventType  model.EventType
		guestName  string
		guestEmail string
		daysAhead  int
		hour       int
		minute     int
	}{
		{createdTypes[0], "Ольга Новикова", "olga@example.com", 2, 9, 0},
		{createdTypes[0], "Сергей Морозов", "sergey@example.com", 2, 14, 0},
		{createdTypes[1], "Анна Белова", "anna@example.com", 2, 10, 0},
		{createdTypes[1], "Павел Романов", "pavel@example.com", 3, 15, 30},
		{createdTypes[2], "Татьяна Орлова", "tatiana@example.com", 3, 9, 0},
		{createdTypes[2], "Николай Григорьев", "nikolay@example.com", 4, 11, 0},
	}
	for _, ub := range upcomingBookings {
		slotStart := today.AddDate(0, 0, ub.daysAhead).Add(time.Duration(ub.hour)*time.Hour + time.Duration(ub.minute)*time.Minute)
		slotID := fmt.Sprintf("%s-%d", ub.eventType.ID, slotStart.Unix())
		s.CreateBooking(model.Booking{
			ID:          s.NextBookingID(),
			SlotID:      slotID,
			EventTypeID: ub.eventType.ID,
			GuestName:   ub.guestName,
			GuestEmail:  ub.guestEmail,
		})
	}

	mux := http.NewServeMux()
	handler.RegisterAll(mux, s)

	if _, err := os.Stat("/frontend"); err == nil {
		mux.Handle("GET /", http.FileServer(http.Dir("/frontend")))
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Backend starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
