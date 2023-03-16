package storage

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	TgID   int64
	TgName string
}

type Request struct {
	Query        string
	Request_time string
}

type StorageI interface {
	Create(u *User) error
	Get(tgID int64) (string, error)
	GetOrCreate(u *User) error
	CreateRequest(query string, tgId int64) error
	GetFirstRequest(tgId int64) (*Request, error)
	GetRequest(tgId int64) ([]*Request, error)
}

type storagePg struct {
	db *sqlx.DB
}

func NewStoragePg(db *sqlx.DB) StorageI {
	return &storagePg{
		db: db,
	}
}

func (s *storagePg) Create(user *User) error {
	query := `
		INSERT INTO users(
			tg_id,
			tg_name
		) VALUES($1,$2)`

	_, err := s.db.Exec(query, user.TgID, user.TgName)
	if err != nil {
		log.Printf("Error while create user: %v", err)
		return err
	}

	return nil
}

func (s *storagePg) Get(id int64) (string, error) {
	var result string

	query := `
		SELECT
			tg_id
		FROM users
		WHERE tg_id=$1
	`

	row := s.db.QueryRow(query, id)
	err := row.Scan(
		&result,
	)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *storagePg) GetOrCreate(req *User) error {
	_, err := s.Get(req.TgID)
	if errors.Is(err, sql.ErrNoRows) {
		err := s.Create(&User{
			TgID:   req.TgID,
			TgName: req.TgName,
		})
		if err != nil {
			log.Printf("Error while create user: %v", err)
			return err
		}
	} else if err != nil {
		log.Printf("Error while checking user: %v", err)
		return err
	}

	return nil
}

func (s *storagePg) CreateRequest(query string, tgId int64) error {
	_, err := s.db.Exec(`INSERT INTO 
	requests(tg_id,request)
	VALUES($1,$2)`, tgId, query)
	if err != nil {
		return err
	}
	return nil
}

func (s *storagePg) GetFirstRequest(tgId int64) (*Request, error) {
	response := &Request{}
	var request_time time.Time
	err := s.db.QueryRow(`SELECT request_time, request 
	FROM requests WHERE tg_id=$1 ORDER BY request_time limit 1`, tgId).Scan(
		&request_time,
		&response.Query,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error while get first request: %v", err)
		return &Request{}, err
	}
	response.Request_time = request_time.Format(time.RFC1123)
	return response, nil
}

func (s *storagePg) GetRequest(tgId int64) ([]*Request, error) {
	var request_time time.Time
	rows, err := s.db.Query(`SELECT 
	request,request_time 
	FROM requests 
	WHERE tg_id=$1`, tgId)
	if err != nil {
		log.Printf("Error while get request: %v", err)
		return []*Request{}, err
	}
	response := []*Request{}
	for rows.Next() {
		i := &Request{}
		err := rows.Scan(
			&i.Query,
			&request_time,
		)
		if err != nil {
			log.Printf("Error while get request in for: %v", err)
			return []*Request{}, err
		}
		i.Request_time = request_time.Format(time.RFC1123)
		response = append(response, i)
	}
	return response, nil
}
