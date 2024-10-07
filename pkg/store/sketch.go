package store

import "tobro/db"

func (s *Store) CreateSketch(circuitID int, name string) (*db.SketchDBModel, error) {
	return s.client.SketchDB.CreateOne(
		db.SketchDB.Name.Set(name),
		db.SketchDB.Circuit.Link(db.CircuitDB.ID.Equals(circuitID)),
	).Exec(s.ctx)
}
