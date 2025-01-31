package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"DeliveryTimePrediction/internal/config"

	"github.com/jackc/pgx/v4"
)

type Storage struct {
	conn *pgx.Conn
}

func New(ctx context.Context, cfg config.PostgresConfig) (*Storage, error) {
	dbPath := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	deadline := time.After(cfg.Timeout)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			conn, err := pgx.Connect(ctx, dbPath)
			if err != nil {
				log.Printf("error connecting to database: %v", err)
				continue
			}
			if err = conn.Ping(ctx); err != nil {
				log.Printf("error pinging database: %v", err)
				continue
			}
			log.Println("Successful database connection")
			return &Storage{conn: conn}, nil

		case <-deadline:
			return nil, fmt.Errorf("unable to connect to database")
		}
	}
}

func (s *Storage) GetResult(ctx context.Context, id string) (float64, error) {
	request := `SELECT delivery_time FROM results WHERE id = $1`
	var deliveryTime float64

	err := s.conn.QueryRow(ctx, request, id).Scan(&deliveryTime)
	return deliveryTime, err
}

func (s *Storage) InsertResult(ctx context.Context, id string, deliveryTime float64) error {
	request := `INSERT INTO results (id, delivery_time) VALUES ($1, $2)`
	_, err := s.conn.Exec(ctx, request, id, deliveryTime)
	return err
}
