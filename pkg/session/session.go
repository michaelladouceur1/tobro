package session

import (
	"encoding/json"
	"os"
)

type Session struct {
	path string
	Port string `json:"port"`
}

func NewSession(path string) (*Session, error) {
	return initSession(&Session{
		path: path,
		Port: "",
	})
}

func (s *Session) UpdatePort(port string) error {
	s.Port = port
	return s.save()
}

func initSession(s *Session) (*Session, error) {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		file, err := os.Create(s.path)
		if err != nil {
			return nil, err
		}

		if err := json.NewEncoder(file).Encode(s); err != nil {
			return nil, err
		}
	}

	file, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(file).Decode(s); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Session) save() error {
	file, err := os.OpenFile(s.path, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(file).Encode(s); err != nil {
		return err
	}

	return nil
}
