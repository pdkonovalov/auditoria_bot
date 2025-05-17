package user

import (
	"context"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	domain "github.com/pdkonovalov/auditoria_bot/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) (domain.UserRepository, error) {
	return &repository{pool}, nil
}

func (r *repository) Create(u *entity.User) (bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM users 
        WHERE user_id = $1)`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, u.UserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	query =
		`INSERT INTO users
		(user_id, username, firstname, lastname, contact_info, state, context, admin)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = r.pool.Exec(context.Background(), query,
		u.UserID, u.Username, u.FirstName, u.LastName, u.ContactInfo, u.State, u.Context, u.Admin)
	return false, err
}

func (r *repository) Get(userID int64) (*entity.User, bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM users 
        WHERE user_id = $1 
    )`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, userID).Scan(&exists)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	query =
		`SELECT user_id, username, firstname, lastname, contact_info, state, context, admin
    	FROM users
    	WHERE user_id = $1`
	u := entity.User{}
	err = r.pool.QueryRow(context.Background(), query, userID).Scan(
		&u.UserID,
		&u.Username,
		&u.FirstName,
		&u.LastName,
		&u.ContactInfo,
		&u.State,
		&u.Context,
		&u.Admin,
	)
	if err != nil {
		return nil, false, err
	}
	return &u, true, nil
}

func (r *repository) Update(u *entity.User) (bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM users 
        WHERE user_id = $1 
    )`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, u.UserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	query = `
    	UPDATE users 
    	SET 
			username = $2,
			firstname = $3,
			lastname = $4,
			contact_info = $5,
			state = $6,
			context = $7,
			admin = $8 
    	WHERE user_id = $1`
	_, err = r.pool.Exec(context.Background(), query,
		u.UserID,
		u.Username,
		u.FirstName,
		u.LastName,
		u.ContactInfo,
		u.State,
		u.Context,
		u.Admin,
	)
	return true, err
}
