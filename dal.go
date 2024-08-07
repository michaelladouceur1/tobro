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

// func NewPinFromDBModel(model *db.PinDBModel) *Pin {

// }

func NewCircuitFromDBModel(model *db.CircuitDBModel) *Circuit {
	newCircuit := NewCircuit(model.ID, model.Name, SupportedBoards(model.Board), portServer)

	pins := model.Pins()
	for _, pin := range pins {
		cPin, err := newCircuit.GetPin(pin.PinNumber)
		if err != nil {
			continue
		}

		cPin.UpdateFromDBModel(&pin)
	}

	return newCircuit
}

func (d *DAL) Connect() error {
	return d.client.Connect()
}

func (d *DAL) Disconnect() {
	d.client.Disconnect()
}

func (d *DAL) GetCircuitByID(id int) (*Circuit, error) {
	circuit, err := d.client.CircuitDB.FindUnique(
		db.CircuitDB.ID.Equals(id)).Exec(d.ctx)
	if err != nil {
		return nil, err
	}

	return NewCircuitFromDBModel(circuit), nil
}

func (d *DAL) CreateCircuit(circuit Circuit) (*Circuit, error) {
	newCircuit, err := d.client.CircuitDB.CreateOne(
		db.CircuitDB.Name.Equals(circuit.Name),
		db.CircuitDB.Board.Set(string(circuit.Board))).Exec(d.ctx)
	if err != nil {
		return nil, err
	}

	d.AddPins(newCircuit.ID, circuit.Pins)

	return d.GetCircuitByID(newCircuit.ID)
}

func (d *DAL) AddPin(circuitID int, pin Pin) (*db.PinDBModel, error) {
	return d.client.PinDB.CreateOne(
		db.PinDB.PinNumber.Set(pin.PinNumber),
		db.PinDB.Circuit.Link(db.CircuitDB.ID.Equals(circuitID)),
		db.PinDB.Mode.Set(int(pin.Mode))).Exec(d.ctx)
}

func (d *DAL) AddPins(circuitID int, pins []Pin) ([]*db.PinDBModel, error) {
	var pinModels []*db.PinDBModel
	for _, pin := range pins {
		pinModel, err := d.AddPin(circuitID, pin)
		if err != nil {
			return nil, err
		}
		pinModels = append(pinModels, pinModel)
	}
	return pinModels, nil
}
