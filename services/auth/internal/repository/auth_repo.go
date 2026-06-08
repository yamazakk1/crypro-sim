package repository

import (
	"context"
	"log"

	"crypto-simulator/services/auth/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	log.Println("auth-repo: initialized")
	return &AuthRepo{pool: pool}
}

func (r *AuthRepo) Create(ctx context.Context, id, username, email string, passwordHash []byte) error {
	log.Printf("auth-repo: Create called: username=%s, email=%s", username, email)

	query := `INSERT INTO users (id, username, email, password_hash, role, balance_usdt)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  ON CONFLICT (username) DO NOTHING
			  RETURNING 1`

	var res int
	err := r.pool.QueryRow(ctx, query, id, username, email, string(passwordHash), models.UserRole, 10000).Scan(&res)
	if err == pgx.ErrNoRows {
		log.Printf("auth-repo: Create: username already exists: %s", username)
		return ErrUsernameAlreadyExists
	}
	if err != nil {
		log.Printf("auth-repo: Create: db error: %v", err)
		return err
	}

	log.Printf("auth-repo: Create: success, username=%s", username)
	return nil
}

func (r *AuthRepo) EmailExists(ctx context.Context, email string) (bool, error) {
	log.Printf("auth-repo: EmailExists called: email=%s", email)

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, email).Scan(&exists)
	if err != nil {
		log.Printf("auth-repo: EmailExists: db error: %v", err)
		return false, err
	}

	log.Printf("auth-repo: EmailExists: result=%v", exists)
	return exists, nil
}

func (r *AuthRepo) UsernameExists(ctx context.Context, username string) (bool, error) {
	log.Printf("auth-repo: UsernameExists called: username=%s", username)

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, username).Scan(&exists)
	if err != nil {
		log.Printf("auth-repo: UsernameExists: db error: %v", err)
		return false, err
	}

	log.Printf("auth-repo: UsernameExists: result=%v", exists)
	return exists, nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	log.Printf("auth-repo: GetUserByEmail called: email=%s", email)

	query := `SELECT id, username, password_hash, role, balance_usdt
	FROM users WHERE email = $1`

	var user models.User
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.UserID, &user.Username, &user.Passwordhash, &user.Role, &user.Balance_usdt)
	if err != nil {
		log.Printf("auth-repo: GetUserByEmail: db error: %v", err)
		return nil, err
	}

	user.Email = email
	log.Printf("auth-repo: GetUserByEmail: success, username=%s", user.Username)
	return &user, nil
}

func (r *AuthRepo) GetMeById(ctx context.Context, userId string) (*models.GetMeUser, error) {
	log.Printf("auth-repo: GetMeById called: user_id=%s", userId)

	query := `SELECT username, email, balance_usdt, role 
	FROM users WHERE id = $1`

	var user models.GetMeUser
	err := r.pool.QueryRow(ctx, query, userId).Scan(&user.Username, &user.Email, &user.Balance_usdt, &user.Role)
	if err != nil {
		log.Printf("auth-repo: GetMeById: db error: %v", err)
		return nil, err
	}

	log.Printf("auth-repo: GetMeById: success, username=%s", user.Username)
	return &user, nil
}