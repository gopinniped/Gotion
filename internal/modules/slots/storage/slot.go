package storage

import "github.com/jmoiron/sqlx"

type SlotStorage struct {
	DB *sqlx.DB
}

func NewSlotStorage(db *sqlx.DB) *SlotStorage {
	return &SlotStorage{DB: db}
}

func Create()
