package main

import (
	"context"
	"tobro/db"
)

type DAL struct {
	ctx    context.Context
	client *db.PrismaClient
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

func (d *DAL) CreateCircuit(name string, boardType SupportedBoards) (*db.CircuitDBModel, error) {
	return d.client.CircuitDB.CreateOne(
		db.CircuitDB.Name.Equals(name),
		db.CircuitDB.Board.Set(string(boardType))).Exec(d.ctx)
}

func (d *DAL) AddPin(circuitID int, pinID int, mode PinMode) (*db.PinDBModel, error) {
	return d.client.PinDB.CreateOne(
		db.PinDB.Pin.Set(pinID),
		db.PinDB.Circuit.Link(db.CircuitDB.ID.Equals(circuitID)),
		db.PinDB.Mode.Set(int(mode))).Exec(d.ctx)
}

func (d *DAL) AddPins(circuitID int, pins []Pin) ([]*db.PinDBModel, error) {
	var pinModels []*db.PinDBModel
	for _, pin := range pins {
		pinModel, err := d.AddPin(circuitID, pin.ID, pin.Mode)
		if err != nil {
			return nil, err
		}
		pinModels = append(pinModels, pinModel)
	}
	return pinModels, nil
}
