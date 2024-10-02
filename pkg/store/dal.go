package store

import (
	"context"
	"log"
	"tobro/db"
	"tobro/pkg/models/circuit"
	"tobro/pkg/models/pin"
)

type DAL struct {
	ctx    context.Context
	client *db.PrismaClient
}

func New() *DAL {
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

// Circuit

func (d *DAL) InitCircuit(circuit *circuit.Circuit) (*db.CircuitDBModel, error) {
	dbCircuit, err := d.client.CircuitDB.FindFirst().With(db.CircuitDB.Pins.Fetch()).Exec(d.ctx)
	if err != nil {
		if err.Error() == "Error: Record not found" {
			return d.CreateCircuit(circuit.Name, string(circuit.Board))
		}
		return nil, err
	}

	return dbCircuit, nil
}

func (d *DAL) GetCircuitByID(id int) (*db.CircuitDBModel, error) {
	return d.client.CircuitDB.FindUnique(
		db.CircuitDB.ID.Equals(id)).With(db.CircuitDB.Pins.Fetch()).Exec(d.ctx)
}

func (d *DAL) CreateCircuit(name string, board string) (*db.CircuitDBModel, error) {
	newCircuit, err := d.client.CircuitDB.CreateOne(
		db.CircuitDB.Name.Set(name),
		db.CircuitDB.Board.Set(board)).Exec(d.ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	pins, err := circuit.SupportedBoardPins(board)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	d.AddPins(newCircuit.ID, pins)

	return d.GetCircuitByID(newCircuit.ID)
}

func (d *DAL) SaveCircuit(circuit circuit.Circuit) (*db.CircuitDBModel, error) {
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

func (d *DAL) AddPin(circuitID int, pin pin.Pin) (*db.PinDBModel, error) {
	return d.client.PinDB.CreateOne(
		db.PinDB.PinNumber.Set(pin.PinNumber),
		db.PinDB.Circuit.Link(db.CircuitDB.ID.Equals(circuitID)),
		db.PinDB.Mode.Set(int(pin.Mode))).Exec(d.ctx)
}

func (d *DAL) AddPins(circuitID int, pins []pin.Pin) ([]*db.PinDBModel, error) {
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

// Sketch

func (d *DAL) CreateSketch(circuitID int, name string) (*db.SketchDBModel, error) {
	return d.client.SketchDB.CreateOne(
		db.SketchDB.Name.Set(name),
		db.SketchDB.Circuit.Link(db.CircuitDB.ID.Equals(circuitID)),
	).Exec(d.ctx)
}
