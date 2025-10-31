package store

import (
	"awesomeProject1/internal/repository"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	config                 *Config
	db                     *pgxpool.Pool
	subscriptionRepository *repository.SubscriptionRepository
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Start() error {
	db, err := pgxpool.New(context.Background(), s.config.DbConnString)
	if err != nil {
		return err
	}

	err = db.Ping(context.Background())

	if err != nil {
		return err
	}
	println("BD Connected")
	s.db = db
	return nil
}

func (s *Store) Stop() {
	s.db.Close()
}

func (s *Store) SubscriptionRepository() *repository.SubscriptionRepository {
	if s.subscriptionRepository == nil {
		s.subscriptionRepository = repository.NewSubscriptionRepository(s.db)
	}
	return s.subscriptionRepository
}
