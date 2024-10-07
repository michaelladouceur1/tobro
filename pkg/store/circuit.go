package store

import (
	"log"
	"tobro/db"
	"tobro/internal/models/circuit"
	"tobro/internal/models/pin"
)

func (s *Store) InitCircuit(circuit *circuit.Circuit) (*db.CircuitDBModel, error) {
	dbCircuit, err := s.client.CircuitDB.FindFirst().With(db.CircuitDB.Pins.Fetch()).Exec(s.ctx)
	if err != nil {
		if err.Error() == "Error: Record not found" {
			return s.CreateCircuit(circuit.Name, string(circuit.Board))
		}
		return nil, err
	}

	return dbCircuit, nil
}

func (s *Store) GetCircuitByID(id int) (*db.CircuitDBModel, error) {
	return s.client.CircuitDB.FindUnique(
		db.CircuitDB.ID.Equals(id)).With(db.CircuitDB.Pins.Fetch()).Exec(s.ctx)
}

func (s *Store) CreateCircuit(name string, board string) (*db.CircuitDBModel, error) {
	newCircuit, err := s.client.CircuitDB.CreateOne(
		db.CircuitDB.Name.Set(name),
		db.CircuitDB.Board.Set(board)).Exec(s.ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	pins, err := circuit.SupportedBoardPins(board)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	s.AddPins(newCircuit.ID, pins)

	return s.GetCircuitByID(newCircuit.ID)
}

func (s *Store) SaveCircuit(circuit circuit.Circuit) (*db.CircuitDBModel, error) {
	_, err := s.client.CircuitDB.FindUnique(db.CircuitDB.ID.Equals(circuit.ID)).Update(
		db.CircuitDB.Name.Set(circuit.Name),
		db.CircuitDB.Board.Set(string(circuit.Board)),
	).Exec(s.ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	for _, pin := range circuit.Pins {
		_, err := s.client.PinDB.FindMany(db.PinDB.And(db.PinDB.PinNumber.Equals(pin.PinNumber), db.PinDB.CircuitID.Equals(circuit.ID))).Update(
			db.PinDB.Mode.Set(int(pin.Mode)),
		).Exec(s.ctx)
		if err != nil {
			log.Print(err)
			return nil, err
		}
	}

	return s.GetCircuitByID(circuit.ID)
}

func (s *Store) AddPin(circuitID int, pin pin.Pin) (*db.PinDBModel, error) {
	return s.client.PinDB.CreateOne(
		db.PinDB.PinNumber.Set(pin.PinNumber),
		db.PinDB.Circuit.Link(db.CircuitDB.ID.Equals(circuitID)),
		db.PinDB.Mode.Set(int(pin.Mode))).Exec(s.ctx)
}

func (s *Store) AddPins(circuitID int, pins []pin.Pin) ([]*db.PinDBModel, error) {
	var pinModels []*db.PinDBModel
	for _, pin := range pins {
		pinModel, err := s.AddPin(circuitID, pin)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		pinModels = append(pinModels, pinModel)
	}
	return pinModels, nil
}
