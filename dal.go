package main

import (
	"tobro/db"
)

type DAL struct {
	dbClient *db.PrismaClient
}

type Circuit struct {
	ID    int
	Name  string
	Board Board
}

func NewDAL() *DAL {
	return &DAL{
		dbClient: db.NewClient(),
	}
}

func (d *DAL) Connect() error {
	return d.dbClient.Connect()
}

func (d *DAL) Disconnect() {
	d.dbClient.Disconnect()
}

func (d *DAL) CreateCircuit(name string, boardType SupportedBoards) (*Circuit, error) {
	createdBoard, err := d.dbClient.Circuit.CreateOne(db.Circuit.Name.Equals(name), db.Circuit.Board.Set(string(boardType))).Exec(dbCtx)
	if err != nil {
		return nil, err
	}

	return &Circuit{
		ID:   createdBoard.ID,
		Name: createdBoard.Name,
		Board: Board{
			Pins: []Pin{},
		},
	}, nil
}
