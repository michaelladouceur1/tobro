package main

import (
	"context"
	"tobro/db"
)

type DAL struct {
	ctx    context.Context
	client *db.PrismaClient
}

type Circuit struct {
	ID    int
	Name  string
	Board Board
}

func NewDAL() *DAL {
	return &DAL{
		ctx:    context.Background(),
		client: db.NewClient(),
	}
}

func (d *DAL) Connect() error {
	return d.client.Connect()
}

func (d *DAL) Disconnect() {
	d.client.Disconnect()
}

func (d *DAL) CreateCircuit(name string, boardType SupportedBoards) (*db.CircuitModel, error) {
	return d.client.Circuit.CreateOne(
		db.Circuit.Name.Equals(name),
		db.Circuit.Board.Set(string(boardType))).Exec(d.ctx)
}

func (d *DAL) AddPin(circuitID int, pinID int, mode PinMode) (*db.PinModel, error) {
	return d.client.Pin.CreateOne(db.Pin.Pin.Set(pinID), db.Pin.Circuit.Link(db.Circuit.ID.Equals(circuitID)), db.Pin.Mode.Set(int(mode))).Exec(d.ctx)
}

func (d *DAL) AddPins(circuitID int, pins []Pin) ([]*db.PinModel, error) {
	var pinModels []*db.PinModel
	for _, pin := range pins {
		pinModel, err := d.AddPin(circuitID, pin.ID, pin.Mode)
		if err != nil {
			return nil, err
		}
		pinModels = append(pinModels, pinModel)
	}
	return pinModels, nil
}
