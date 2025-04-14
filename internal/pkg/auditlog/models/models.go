package models

import "time"

type Log struct {
	UserID      int       `json:"user_id"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
	Error       string    `json:"error"`
}

type Task struct {
	UserID          int       `json:"user_id" db:"user_id"`
	Action          string    `json:"action" db:"action"`
	Description     string    `json:"description" db:"description"`
	Timestamp       time.Time `json:"timestamp" db:"timestamp"`
	Error           string    `json:"error" db:"error"`
	Attempts        int       `json:"attempts" db:"attempts"`
	Status          string    `json:"status" db:"status"`
	Created_at      time.Time `json:"created_at" db:"created_at"`
	Updated_at      time.Time `json:"updated_at" db:"updated_at"`
	Complited_at    time.Time `json:"complited_at" db:"complited_at"`
	Processing_from time.Time `json:"processing_from" db:"processing_from"`
}
