package user

import "time"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleClient Role = "client"
)

type User struct {
	id           string
	name         string
	email        string
	passwordHash string
	role         Role
	createdAt    time.Time
	updatedAt    time.Time
}

func New(name, email, passwordHash string, role Role) *User {
	now := time.Now()
	return &User{
		name:         name,
		email:        email,
		passwordHash: passwordHash,
		role:         role,
		createdAt:    now,
		updatedAt:    now,
	}
}

// Reconstitute restaura uma entidade a partir de dados persistidos (uso exclusivo de repositories).
func Reconstitute(id, name, email, passwordHash string, role Role, createdAt, updatedAt time.Time) *User {
	return &User{
		id:           id,
		name:         name,
		email:        email,
		passwordHash: passwordHash,
		role:         role,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// Getters
func (u *User) ID() string           { return u.id }
func (u *User) Name() string         { return u.name }
func (u *User) Email() string        { return u.email }
func (u *User) PasswordHash() string { return u.passwordHash }
func (u *User) Role() Role           { return u.role }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// Setters — apenas campos mutaveis apos criacao
func (u *User) SetID(id string)     { u.id = id }
func (u *User) SetName(name string) { u.name = name; u.touch() }

func (u *User) touch() { u.updatedAt = time.Now() }
