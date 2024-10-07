package store

import (
	"context"
	"tobro/db"
)

type Store struct {
	ctx    context.Context
	client *db.PrismaClient
}

func New(client *db.PrismaClient) *Store {
	return &Store{
		ctx:    context.Background(),
		client: client,
	}
}

func (s *Store) Connect() error {
	return s.client.Connect()
}

func (s *Store) Disconnect() {
	s.client.Disconnect()
}
