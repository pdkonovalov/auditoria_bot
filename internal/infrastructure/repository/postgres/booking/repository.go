package booking

import (
	"context"

	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	domain "github.com/pdkonovalov/auditoria_bot/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) (domain.BookingRepository, error) {
	return &repository{pool}, nil
}

func (r *repository) Create(b *entity.Booking) (bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM bookings 
        WHERE event_id = $1 AND user_id = $2)`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, b.EventID, b.UserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	query =
		`INSERT INTO bookings
		(event_id, user_id, payment, payment_photo, payment_document, text, offline, online, created_at, updated_at, draft)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err = r.pool.Exec(context.Background(), query,
		b.EventID, b.UserID, b.Payment, b.PaymentPhoto, b.PaymentDocument, b.Text, b.Offline, b.Online, b.CreatedAt, b.UpdatedAt, b.Draft)
	return false, err
}

func (r *repository) Get(userID int64, eventID string) (*entity.Booking, bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM bookings 
        WHERE user_id = $1 AND event_id = $2)`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, userID, eventID).Scan(&exists)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	query =
		`SELECT event_id, user_id, payment, payment_photo, payment_document, text, offline, online, created_at, updated_at, draft
    	FROM bookings
    	WHERE user_id = $1 AND event_id = $2`
	b := entity.Booking{}
	err = r.pool.QueryRow(context.Background(), query, userID, eventID).Scan(
		&b.EventID,
		&b.UserID,
		&b.Payment,
		&b.PaymentPhoto,
		&b.PaymentDocument,
		&b.Text,
		&b.Offline,
		&b.Online,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.Draft,
	)
	if err != nil {
		return nil, false, err
	}
	return &b, true, nil
}

func (r *repository) GetByEventID(eventID string, offline, online bool) ([]*entity.Booking, error) {
	var query = `SELECT event_id, user_id, payment, payment_photo, payment_document, text, offline, online, created_at, updated_at, draft
	FROM bookings
	WHERE event_id = $1 AND draft IS FALSE AND offline = $2 AND online = $3
	ORDER BY COALESCE(updated_at, created_at) ASC`
	rows, err := r.pool.Query(context.Background(), query, eventID, offline, online)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*entity.Booking
	for rows.Next() {
		booking := &entity.Booking{}
		err := rows.Scan(&booking.EventID, &booking.UserID, &booking.Payment, &booking.PaymentPhoto, &booking.PaymentDocument, &booking.Text, &booking.Offline, &booking.Online, &booking.CreatedAt, &booking.UpdatedAt, &booking.Draft)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *repository) GetByUserID(userID int64) ([]*entity.Booking, error) {
	var query = `SELECT event_id, user_id, payment, payment_photo, payment_document, text, offline, online, created_at, updated_at, draft
		FROM bookings
		WHERE user_id = $1 AND draft IS FALSE
		ORDER BY COALESCE(updated_at, created_at) ASC`
	rows, err := r.pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*entity.Booking
	for rows.Next() {
		booking := &entity.Booking{}
		err := rows.Scan(&booking.EventID, &booking.UserID, &booking.Payment, &booking.PaymentPhoto, &booking.PaymentDocument, &booking.Text, &booking.Offline, &booking.Online, &booking.CreatedAt, &booking.UpdatedAt, &booking.Draft)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *repository) Update(b *entity.Booking) (bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM bookings 
        WHERE event_id = $1 AND user_id = $2)`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, b.EventID, b.UserID).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	query = `
    	UPDATE bookings
    	SET 
			payment = $3,
			payment_photo = $4,
			payment_document = $5,
			text =  $6,
			offline = $7,
			online = $8,
			created_at = $9,
			updated_at = $10,
			draft = $11
    	WHERE event_id = $1 AND user_id = $2`
	_, err = r.pool.Exec(context.Background(), query,
		b.EventID,
		b.UserID,
		b.Payment,
		b.PaymentPhoto,
		b.PaymentDocument,
		b.Text,
		b.Offline,
		b.Online,
		b.CreatedAt,
		b.UpdatedAt,
		b.Draft,
	)
	return true, err
}

func (r *repository) Delete(userID int64, eventID string) (bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM bookings 
        WHERE user_id = $1 AND event_id = $2)`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, userID, eventID).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	query =
		`DELETE FROM bookings
		WHERE user_id = $1 AND event_id = $2`
	_, err = r.pool.Exec(context.Background(), query, userID, eventID)
	return true, err
}
