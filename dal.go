package main

import (
	"context"
	"log"
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
		db.CircuitDB.ID.Equals(id)).With(db.CircuitDB.Pins.Fetch()).Exec(d.ctx)
	if err != nil {
		return nil, err
	}

	return NewCircuitFromDBModel(circuit), nil
}

func (d *DAL) CreateCircuit(circuit Circuit) (*Circuit, error) {
	newCircuit, err := d.client.CircuitDB.CreateOne(
		db.CircuitDB.Name.Set(circuit.Name),
		db.CircuitDB.Board.Set(string(circuit.Board))).Exec(d.ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	d.AddPins(newCircuit.ID, circuit.Pins)

	return d.GetCircuitByID(newCircuit.ID)
}

func (d *DAL) SaveCircuit(circuit Circuit) (*Circuit, error) {
	_, err := d.client.CircuitDB.FindUnique(db.CircuitDB.ID.Equals(circuit.ID)).Update(
		db.CircuitDB.Name.Set(circuit.Name),
		db.CircuitDB.Board.Set(string(circuit.Board)),
	).Exec(d.ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for _, pin := range circuit.Pins {
		_, err := d.client.PinDB.FindMany(db.PinDB.And(db.PinDB.PinNumber.Equals(pin.PinNumber), db.PinDB.CircuitID.Equals(circuit.ID))).Update(
			db.PinDB.Mode.Set(int(pin.Mode)),
		).Exec(d.ctx)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	}

	return d.GetCircuitByID(circuit.ID)
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
			log.Print(err)
			return nil, err
		}
		pinModels = append(pinModels, pinModel)
	}
	return pinModels, nil
}
