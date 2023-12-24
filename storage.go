package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateMember(*Member) error
	DeleteMember(uuid.UUID) error
	UpdateMember(*Member) error
	GetMemberByID(uuid.UUID) (*Member, error)
	GetMembers() ([]*Member, error)
	GetMembersByTech(tech string) ([]*Member, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	//connStr := "user=postgres dbname=postgres password=briheet sslmode=disable"
	connStr := "postgres://mentor_user:9aLIh9IoMHznNjcqouMx93NlUxhbrUgA@dpg-clvakemd3nmc738av6j0-a.singapore-postgres.render.com/mentor"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
			return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	},nil
}

func (s *PostgresStore) Init() error {
	return s.createMemberTable()
}

func (s *PostgresStore) createMemberTable() error {
	resp, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS member (
		  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			tech VARCHAR(50),
			about VARCHAR(50),
			discord VARCHAR(50),
			linkedin VARCHAR(50),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)

	fmt.Printf("%+v\n", resp)

	if err != nil {
			return err
	}

	return nil
}



func (s *PostgresStore) CreateMember(mem *Member) error {
	query := `
			INSERT INTO member
			(first_name, last_name, tech, about, discord, linkedin, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.Exec(query,
		mem.FirstName, 
		mem.LastName, 
		mem.Tech, 
		mem.About, 
		mem.Discord,
		mem.Linkedin, 
		mem.CreatedAt)

	if err != nil {
			return err
	}

	return nil
}



func (s *PostgresStore) DeleteMember (id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM member WHERE id = $1", id)
	return err
}




func (s *PostgresStore) UpdateMember (mem *Member) error {
	_, err := s.db.Exec(`
		UPDATE member
		SET first_name = $2, last_name = $3, tech = $4, about = $5, discord = $6, linkedin = $7, created_at = $8
		WHERE id = $1`,
		mem.ID,
		mem.FirstName,
		mem.LastName,
		mem.Tech,
		mem.About,
		mem.Discord,
		mem.Linkedin,
		mem.CreatedAt)
	return err
}




func (s *PostgresStore) GetMemberByID (id uuid.UUID) (*Member, error) {
	member := &Member{}
	err := s.db.QueryRow("SELECT * FROM member WHERE id = $1", id).
		Scan(&member.ID, 
			&member.FirstName, 
			&member.LastName, 
			&member.Tech, 
			&member.About, 
			&member.Discord, 
			&member.Linkedin, 
			&member.CreatedAt)

	if err != nil {
		return nil, err
	}
	return member, nil
}



func (s *PostgresStore) GetMembersByTech(tech string) ([]*Member, error) {
	//query := "SELECT * FROM member WHERE tech = $1"
	query := "SELECT * FROM member WHERE tech ILIKE '%' || $1 || '%'"
  rows, err := s.db.Query(query, tech)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, err
	}
	defer rows.Close()

	members := []*Member{}
	for rows.Next() {
		member := &Member{}
		if err := rows.Scan(
			&member.ID,
			&member.FirstName,
			&member.LastName,
			&member.Tech,
			&member.About,
			&member.Discord,
			&member.Linkedin,
			&member.CreatedAt,
		); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error reading rows:", err)
		return nil, err
	}

	return members, nil
}





func (s *PostgresStore) GetMembers() ([]*Member, error) {
	rows, err := s.db.Query(`SELECT * FROM member`)
	if err != nil {
			return nil, err
	}

	defer rows.Close()

	members := []*Member{}
	for rows.Next() {
			member := new(Member)
			err := rows.Scan(
					&member.ID,
					&member.FirstName,
					&member.LastName,
					&member.Tech,
					&member.About,
					&member.Discord,
					&member.Linkedin,
					&member.CreatedAt, 
			)
			if err != nil {
					return nil, err
			}
			members = append(members, member)
	}

	//fmt.Println(members)

	return members, nil
}


