package user

import "time"

type Role string

const (
	RoleAdmin  Role = "admin"
	RoleClient Role = "client"
)

type User struct {
	id             string
	name           string
	email          string
	passwordHash   string
	role           Role
	failedAttempts int
	lockedUntil    *time.Time // nullable — nil quando nao esta bloqueado
	createdAt      time.Time
	updatedAt      time.Time
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
func Reconstitute(
	id, name, email, passwordHash string,
	role Role,
	failedAttempts int,
	lockedUntil *time.Time,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		id:             id,
		name:           name,
		email:          email,
		passwordHash:   passwordHash,
		role:           role,
		failedAttempts: failedAttempts,
		lockedUntil:    lockedUntil,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

// Getters
func (u *User) ID() string             { return u.id }
func (u *User) Name() string           { return u.name }
func (u *User) Email() string          { return u.email }
func (u *User) PasswordHash() string   { return u.passwordHash }
func (u *User) Role() Role             { return u.role }
func (u *User) FailedAttempts() int    { return u.failedAttempts }
func (u *User) LockedUntil() *time.Time { return u.lockedUntil }
func (u *User) CreatedAt() time.Time   { return u.createdAt }
func (u *User) UpdatedAt() time.Time   { return u.updatedAt }

// Setters — apenas campos mutaveis apos criacao
func (u *User) SetID(id string)     { u.id = id }
func (u *User) SetName(name string) { u.name = name; u.touch() }

// --- Comportamentos de seguranca (account lockout) ---

// IsLocked retorna true se o usuario esta temporariamente bloqueado.
// Retorna false quando o lockedUntil ja expirou (auto-unlock).
func (u *User) IsLocked() bool {
	if u.lockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.lockedUntil)
}

// RegisterFailedLogin incrementa o contador de tentativas falhas e,
// se atingir o limite, bloqueia a conta por lockDuration.
// A decisao de "quando bloquear" fica no dominio, nao no use case.
func (u *User) RegisterFailedLogin(maxAttempts int, lockDuration time.Duration) {
	u.failedAttempts++
	if u.failedAttempts >= maxAttempts {
		until := time.Now().Add(lockDuration)
		u.lockedUntil = &until
	}
	u.touch()
}

// ResetFailedLogins limpa o contador e remove qualquer bloqueio.
// Deve ser chamado apos login bem-sucedido.
func (u *User) ResetFailedLogins() {
	u.failedAttempts = 0
	u.lockedUntil = nil
	u.touch()
}

func (u *User) touch() { u.updatedAt = time.Now() }
