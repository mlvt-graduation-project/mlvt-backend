package main

import (
	"context"
	"fmt"
	"mlvt/internal/entity"
	"mlvt/internal/infra/db/mongodb"
	"mlvt/internal/infra/env"
	"mlvt/internal/infra/zap-logging/log"
	"mlvt/internal/initialize"
	"mlvt/internal/repo/traffic_repo"
	"mlvt/internal/service/traffic_service"
	"os"
	"time"
)

func main() {
	// Initialize Logger
	if err := initialize.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Logger initialization failed: %v\n", err)
		os.Exit(1)
	}

	mongoConn := mongodb.NewMongoDBClient(env.EnvConfig.MongoDBEndPoint)

	// Initialize AWS Clients
	s3Client, err := initialize.InitAWS()
	if err != nil {
		log.Errorf("AWS initialization failed: %v", err)
		os.Exit(1)
	}

	tRepo := traffic_repo.NewTrafficRepo(mongoConn)
	tService := traffic_service.NewTrafficService(tRepo, s3Client)

	// 3) Seed data
	// We'll create seeds for:
	//    - Current day (24 hours)
	//    - Current week (7 days, Mon-Sun)
	//    - Current year (12 months)
	// Each entry is spread across time so you can test your day/week/year grouping logic.

	var baseTime time.Time
	baseTime, _ = time.Parse("2006-01-02", "2024-01-21")

	err = seedDayTraffic(tService, baseTime)
	if err != nil {
		log.Errorf("failed to seed day traffic: %v", err)
		return
	}

	err = seedWeekTraffic(tService, baseTime)
	if err != nil {
		log.Errorf("failed to seed week traffic: %v", err)
		return
	}

	err = seedYearTraffic(tService, baseTime)
	if err != nil {
		log.Errorf("failed to seed year traffic: %v", err)
		return
	}

	fmt.Println("Traffic seeding completed successfully.")
}

// seedDayTraffic creates 24 traffic items within the current day, each in a distinct hour.
func seedDayTraffic(s traffic_service.TrafficService, baseTime time.Time) error {
	now := baseTime
	// Start = midnight
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	for hour := 0; hour < 24; hour++ {
		ts := startOfDay.Add(time.Duration(hour) * time.Hour).Unix()

		traffic := entity.Traffic{
			ActionType:  entity.ProcessSTTModelAction, // or whatever action you like
			Description: fmt.Sprintf("Seeding hour %02d:00", hour),
			UserID:      2, // or any test user
			Timestamp:   ts,
		}
		if _, err := s.CreateTraffic(context.Background(), traffic); err != nil {
			return fmt.Errorf("insert traffic for hour %d: %w", hour, err)
		}
	}

	fmt.Println("Seeded 24 day-traffic entries.")
	return nil
}

// seedWeekTraffic creates 7 traffic items in the current week, Monday-Sunday.
func seedWeekTraffic(s traffic_service.TrafficService, baseTime time.Time) error {
	now := baseTime
	wd := int(now.Weekday())
	if wd == 0 {
		wd = 7
	}
	// Monday of this week
	monday := now.AddDate(0, 0, -(wd - 1))
	startOfMonday := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())

	// For i in [0..6] => Monday + i days
	for i := 0; i < 7; i++ {
		day := startOfMonday.AddDate(0, 0, i)
		ts := day.Unix()
		traffic := entity.Traffic{
			ActionType:  entity.ProcessTTSModelAction,
			Description: fmt.Sprintf("Seeding day #%d of this week", i+1),
			UserID:      2,
			Timestamp:   ts,
		}
		if _, err := s.CreateTraffic(context.Background(), traffic); err != nil {
			return fmt.Errorf("insert traffic for day offset %d: %w", i, err)
		}
	}

	fmt.Println("Seeded 7 week-traffic entries.")
	return nil
}

// seedYearTraffic creates 12 traffic items in the current year, Jan–Dec.
func seedYearTraffic(s traffic_service.TrafficService, baseTime time.Time) error {
	now := baseTime
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())

	for monthOffset := 0; monthOffset < 12; monthOffset++ {
		thisMonth := startOfYear.AddDate(0, monthOffset, 0)
		ts := thisMonth.Unix()

		traffic := entity.Traffic{
			ActionType:  entity.ProcessFullPipelineModelAction,
			Description: fmt.Sprintf("Seeding month #%d", monthOffset+1),
			UserID:      2,
			Timestamp:   ts,
		}
		if _, err := s.CreateTraffic(context.Background(), traffic); err != nil {
			return fmt.Errorf("insert traffic for month offset %d: %w", monthOffset, err)
		}
	}

	fmt.Println("Seeded 12 year-traffic entries.")
	return nil
}
