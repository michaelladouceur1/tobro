package store

import "tobro/db"

func (s *Store) InitSession() (*db.SessionDataDBModel, error) {
	session, _ := s.GetSession()
	if session != nil {
		return session, nil
	}
	return s.client.SessionDataDB.CreateOne(db.SessionDataDB.PortName.Set(""), db.SessionDataDB.PortID.Set("")).Exec(s.ctx)
}

func (s *Store) GetSession() (*db.SessionDataDBModel, error) {
	return s.client.SessionDataDB.FindFirst().Exec(s.ctx)
}

func (s *Store) UpdateSession(portName string, portID string) (*db.SessionDataDBModel, error) {
	return s.client.SessionDataDB.FindUnique(db.SessionDataDB.ID.Equals(0)).Update(
		db.SessionDataDB.PortName.Set(portName),
		db.SessionDataDB.PortID.Set(portID),
	).Exec(s.ctx)
}
