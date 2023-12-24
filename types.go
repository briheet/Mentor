package main

import (
	"github.com/google/uuid"
	"time"
)

type CreateMemberRequest struct {
	FirstName   string      `json:"firstname"`
	LastName    string		`json:"lastname"`
	Tech        string		`json:"tech"`
	About       string		`json:"about"`
	Discord     string		`json:"discord"`
	Linkedin    string      `json:"linkedin"`
}

type Member struct {
	ID          uuid.UUID   `json:"id"`
	FirstName   string      `json:"firstname"`
	LastName    string		`json:"lastname"`
	Tech        string		`json:"tech"`
	About       string		`json:"about"`
	Discord     string		`json:"discord"`
	Linkedin    string      `json:"linkedin"`
	CreatedAt   time.Time   `json:"createdAt"`
}

func NewMember (firstname, lastname, tech, about, discord, linkedin string) *Member {
	return &Member{
		ID:          uuid.New(),
		FirstName:   firstname,
		LastName:    lastname,
		Tech:        tech,
		About:       about,
		Discord:     discord,
		Linkedin:    linkedin,
		CreatedAt:   time.Now().UTC(),
	}
} 