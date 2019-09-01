package main

import (
	"context"
	"fmt"
	"github.com/MihailShev/calendar-service/calendarpb"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer func() {
		err = cc.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	client := calendarpb.NewCalendarClient(cc)

	uuid := createEvent(client)
	event := getEvent(client, uuid)

	fmt.Printf("Created event: %+v\n", event)

	event = updateEvent(client, event)

	fmt.Printf("Updated event description: %+v\n", event)
}

func createEvent(client calendarpb.CalendarClient) uint64 {
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	start := time.Now()
	end := start.Add(5 * time.Minute)

	res, err := client.CreateEvent(ctx, &calendarpb.CreateEventReq{
		Event: &calendarpb.Event{
			Title:       "Test event",
			Start:       &timestamp.Timestamp{Seconds: start.Unix(), Nanos: 0},
			End:         &timestamp.Timestamp{Seconds: end.Unix(), Nanos: 0},
			Description: "Description",
			NoticeTime:  5,
			UserId:      1,
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	return res.GetUUID()
}

func getEvent(client calendarpb.CalendarClient, uuid uint64) *calendarpb.Event {
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	res, err := client.GetEvent(ctx, &calendarpb.GetEventReq{UUID: uuid})

	if err != nil {
		log.Fatal(err)
	}

	if res.Error != "" {
		log.Fatal(res.Error)
	}

	return res.GetEvent()
}

func updateEvent(client calendarpb.CalendarClient, event *calendarpb.Event) *calendarpb.Event {
	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	res, err := client.UpdateEvent(ctx, &calendarpb.UpdateEventReq{
		Event: &calendarpb.Event{
			UUID:        event.UUID,
			Title:       event.Title,
			Start:       event.Start,
			End:         event.End,
			Description: "New description",
			UserId:      event.UserId,
			NoticeTime:  event.NoticeTime,
		},
	})

	if err != nil {
		log.Fatal(event)
	}

	if res.Error != "" {
		log.Fatal(res.Error)
	}

	return res.GetEvent()
}