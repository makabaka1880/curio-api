// Created by Sean L. on Apr. 8.
// Last Updated by Sean L. on Apr. 8.
//
// curio-api
// persistence/schema.go
//
// Makabaka1880, 2026. All rights reserved.

package persistence

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"type:text;not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`

	SubmissionStarts time.Time `gorm:"type:date;not null" json:"submission_starts"`
	EventDate        time.Time `gorm:"type:date;not null" json:"event_date"`
}

type Submission struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	EventID        uuid.UUID `gorm:"type:uuid;not null" json:"event_id"`
	Name           string    `gorm:"type:text;not null" json:"name"`
	Description    string    `gorm:"type:text;not null" json:"description"`
	SubmissionDate time.Time `gorm:"type:date;not null" json:"submission_date"`
	Token          string    `gorm:"type:varchar(6);not null" json:"token"`
}

type Vote struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IpAddress    string    `gorm:"type:text;not null"`
	SubmissionID uuid.UUID `gorm:"type:uuid;index"`
	Value        int       `gorm:"type:smallint"`
}
