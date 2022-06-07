package postgres

import (
	"database/sql"
	"log"

	"github.com/Power9-Alpha/laser"
)

type Service struct {
	db     *sql.DB
	create *sql.Stmt
	getOne *sql.Stmt
	getAll *sql.Stmt
	update *sql.Stmt
	delete *sql.Stmt
}

func (s *Service) Init(db *sql.DB) error {
	s.db = db
	var err error
	if s.create, err = s.db.Prepare(SQLServiceInsert); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if s.getOne, err = s.db.Prepare(SQLServiceSelect); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if s.getAll, err = s.db.Prepare(SQLServiceSelect); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if s.update, err = s.db.Prepare(SQLServiceUpdate); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	} else if s.delete, err = s.db.Prepare(SQLServiceDelete); err != nil {
		log.Printf("failed to create prepared statement: %s", err)
		return err
	}
	return err
}

func (s *Service) Insert(service *laser.Service) error {
	_, err := s.create.Exec(service.Name, service.Technology, service.PointOfContact)
	// @note: we could capture the result if SQL RETURNING such as to capture
	//        the id or entire model
	return err
}

func (s *Service) SelectOne(id string) (*laser.Service, error) {
	service := laser.Service{}
	if err := s.getOne.QueryRow(id).Scan(&service.Name, &service.Technology, &service.PointOfContact); err == sql.ErrNoRows {
		return nil, nil // @note: caller must check for nil model to address 404
	} else if err != nil {
		log.Printf("Select Failed: %s (%s)", id, err)
		return nil, err
	}
	return &service, nil
}

func (s *Service) Select() ([]laser.Service, error) {
	rows, err := s.getAll.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []laser.Service{}

	for rows.Next() {
		service := laser.Service{}
		if err := rows.Scan(&service.Name, &service.Technology, &service.PointOfContact); err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func (s *Service) Update(service *laser.Service) error {
	_, err := s.update.Exec(service.Name, service.Technology, service.PointOfContact)
	return err
}

func (s *Service) Delete(id string) error {
	_, err := s.delete.Exec(id)
	return err
}
