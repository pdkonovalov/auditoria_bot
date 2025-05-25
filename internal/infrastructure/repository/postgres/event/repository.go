package event

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/config"
	"github.com/pdkonovalov/auditoria_bot/internal/domain/entity"
	domain "github.com/pdkonovalov/auditoria_bot/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	pool              *pgxpool.Pool
	eventIDLen        int
	eventIDCharset    string
	eventIDTypeNumber bool
}

func New(cfg *config.Config, pool *pgxpool.Pool) (domain.EventRepository, error) {
	r := repository{
		pool:       pool,
		eventIDLen: cfg.EventIDLen,
	}
	switch cfg.EventIDCharset {
	case config.EventIDCharsetLetters:
		r.eventIDCharset = "abcdefghijklmnopqrstuvwxyz"
	case config.EventIDCharsetNumbers:
		r.eventIDTypeNumber = true
	default:
		r.eventIDCharset = cfg.EventIDCharset
	}
	return &r, nil
}

func (r *repository) Create(e *entity.Event) (string, error) {
	var query string
	var exists bool
	for {
		eventID, err := r.generateEventID()
		if err != nil {
			return "", fmt.Errorf("Failed generate event id: %s", err)
		}
		query =
			`SELECT EXISTS (
        	SELECT 1 FROM events 
        	WHERE event_id = $1)`
		err = r.pool.QueryRow(context.Background(), query, eventID).Scan(&exists)
		if !exists {
			e.EventID = eventID
			break
		}
	}
	query =
		`INSERT INTO events
		(event_id, 
		title, 
		photo, 
		text, 
		entities, 
		time, 
		offline, 
		offline_paid, 
		online, 
		online_paid, 
		payment_details_firstname,
		payment_details_lastname,
		payment_details_account,
		created_at, 
		created_by,
		updated_at, 
		updated_by,
		draft)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`
	_, err := r.pool.Exec(context.Background(), query,
		e.EventID,
		e.Title,
		e.Photo,
		e.Text,
		e.Entities,
		e.Time,
		e.Offline,
		e.OfflinePaid,
		e.Online,
		e.OnlinePaid,
		e.PaymentDetailsFirstName,
		e.PaymentDetailsLastName,
		e.PaymentDetailsAccount,
		e.CreatedAt,
		e.CreatedBy,
		e.UpdatedAt,
		e.UpdatedBy,
		e.Draft)
	if err != nil {
		return "", err
	}
	return e.EventID, nil
}

func (r *repository) Delete(eventID string) (bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM events 
        WHERE event_id = $1)`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, eventID).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	query =
		`DELETE FROM events
		WHERE event_id = $1`
	_, err = r.pool.Exec(context.Background(), query, eventID)
	return true, err
}

func (r *repository) Get(eventID string) (*entity.Event, bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM events 
        WHERE event_id = $1)`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, eventID).Scan(&exists)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	query =
		`SELECT 
		event_id, 
		title, 
		photo, 
		text, 
		entities, 
		time, 
		offline, 
		offline_paid, 
		online, 
		online_paid,
		payment_details_firstname,
		payment_details_lastname,
		payment_details_account,
		created_at, 
		created_by,
		updated_at, 
		updated_by,
		draft
    	FROM events
    	WHERE event_id = $1`
	e := entity.Event{}
	err = r.pool.QueryRow(context.Background(), query, eventID).Scan(
		&e.EventID,
		&e.Title,
		&e.Photo,
		&e.Text,
		&e.Entities,
		&e.Time,
		&e.Offline,
		&e.OfflinePaid,
		&e.Online,
		&e.OnlinePaid,
		&e.PaymentDetailsFirstName,
		&e.PaymentDetailsLastName,
		&e.PaymentDetailsAccount,
		&e.CreatedAt,
		&e.CreatedBy,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.Draft,
	)
	if err != nil {
		return nil, false, err
	}
	return &e, true, nil
}

func (r *repository) GetFirstAfter(start time.Time, userID *int64) (*entity.Event, bool, error) {
	var query string
	var exists bool
	if userID == nil {
		query =
			`SELECT EXISTS (
        	SELECT 1 FROM events 
        	WHERE time > $1 AND draft IS FALSE
		)`
		err := r.pool.QueryRow(context.Background(), query, start).Scan(&exists)
		if err != nil {
			return nil, false, err
		}
		if !exists {
			return nil, false, nil
		}
		query =
			`SELECT
		event_id, 
		title, 
		photo, 
		text, 
		entities, 
		time, 
		offline, 
		offline_paid, 
		online, 
		online_paid, 
		payment_details_firstname,
		payment_details_lastname,
		payment_details_account,
		created_at, 
		created_by,
		updated_at, 
		updated_by,
		draft
    	FROM events
    	WHERE time > $1 AND draft IS FALSE
		ORDER BY time ASC, title ASC
		LIMIT 1`
		e := entity.Event{}
		err = r.pool.QueryRow(context.Background(), query, start).Scan(
			&e.EventID,
			&e.Title,
			&e.Photo,
			&e.Text,
			&e.Entities,
			&e.Time,
			&e.Offline,
			&e.OfflinePaid,
			&e.Online,
			&e.OnlinePaid,
			&e.PaymentDetailsFirstName,
			&e.PaymentDetailsLastName,
			&e.PaymentDetailsAccount,
			&e.CreatedAt,
			&e.CreatedBy,
			&e.UpdatedAt,
			&e.UpdatedBy,
			&e.Draft,
		)
		if err != nil {
			return nil, false, err
		}
		return &e, true, nil
	}
	query =
		`SELECT EXISTS (
        	SELECT 1 
			FROM events e JOIN bookings b ON e.event_id = b.event_id
        	WHERE
			e.time > $1 
			AND e.draft IS FALSE 
			AND b.draft IS FALSE 
			AND b.user_id = $2
		)`
	err := r.pool.QueryRow(context.Background(), query, start, *userID).Scan(&exists)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	query =
		`SELECT
		e.event_id, 
		e.title, 
		e.photo, 
		e.text, 
		e.entities, 
		e.time, 
		e.offline, 
		e.offline_paid, 
		e.online, 
		e.online_paid,
		e.payment_details_firstname,
		e.payment_details_lastname,
		e.payment_details_account,
		e.created_at, 
		e.created_by,
		e.updated_at, 
		e.updated_by,
		e.draft
    	FROM events e JOIN bookings b ON e.event_id = b.event_id
    	WHERE
		e.time > $1 
		AND e.draft IS FALSE
		AND b.draft IS FALSE 
		AND b.user_id = $2
		ORDER BY e.time ASC, e.title ASC
		LIMIT 1`
	e := entity.Event{}
	err = r.pool.QueryRow(context.Background(), query, start, *userID).Scan(
		&e.EventID,
		&e.Title,
		&e.Photo,
		&e.Text,
		&e.Entities,
		&e.Time,
		&e.Offline,
		&e.OfflinePaid,
		&e.Online,
		&e.OnlinePaid,
		&e.PaymentDetailsFirstName,
		&e.PaymentDetailsLastName,
		&e.PaymentDetailsAccount,
		&e.CreatedAt,
		&e.CreatedBy,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.Draft,
	)
	if err != nil {
		return nil, false, err
	}
	return &e, true, nil
}

func (r *repository) GetFirstBefore(end time.Time, userID *int64) (*entity.Event, bool, error) {
	var query string
	var exists bool
	if userID == nil {
		query =
			`SELECT EXISTS (
        SELECT 1 FROM events 
        WHERE time < $1 AND draft IS FALSE)`
		err := r.pool.QueryRow(context.Background(), query, end).Scan(&exists)
		if err != nil {
			return nil, false, err
		}
		if !exists {
			return nil, false, nil
		}
		query =
			`SELECT
		event_id, 
		title, 
		photo, 
		text, 
		entities, 
		time, 
		offline, 
		offline_paid, 
		online, 
		online_paid,
		payment_details_firstname,
		payment_details_lastname,
		payment_details_account,
		created_at, 
		created_by,
		updated_at, 
		updated_by,
		draft
    	FROM events
    	WHERE time < $1 AND draft IS FALSE
		ORDER BY time DESC, title DESC
		LIMIT 1`
		e := entity.Event{}
		err = r.pool.QueryRow(context.Background(), query, end).Scan(
			&e.EventID,
			&e.Title,
			&e.Photo,
			&e.Text,
			&e.Entities,
			&e.Time,
			&e.Offline,
			&e.OfflinePaid,
			&e.Online,
			&e.OnlinePaid,
			&e.PaymentDetailsFirstName,
			&e.PaymentDetailsLastName,
			&e.PaymentDetailsAccount,
			&e.CreatedAt,
			&e.CreatedBy,
			&e.UpdatedAt,
			&e.UpdatedBy,
			&e.Draft,
		)
		if err != nil {
			return nil, false, err
		}
		return &e, true, nil
	}
	query =
		`SELECT EXISTS (
        		SELECT 1 
				FROM events e JOIN bookings b ON e.event_id = b.event_id
        		WHERE
				e.time < $1 
				AND e.draft IS FALSE 
				AND b.draft IS FALSE 
				AND b.user_id = $2)`
	err := r.pool.QueryRow(context.Background(), query, end, *userID).Scan(&exists)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}
	query =
		`SELECT
		e.event_id, 
		e.title, 
		e.photo, 
		e.text, 
		e.entities, 
		e.time, 
		e.offline, 
		e.offline_paid, 
		e.online, 
		e.online_paid,
		e.payment_details_firstname,
		e.payment_details_lastname,
		e.payment_details_account,
		e.created_at, 
		e.created_by,
		e.updated_at, 
		e.updated_by,
		e.draft
    	FROM events e JOIN bookings b ON e.event_id = b.event_id
    	WHERE
		e.time < $1 
		AND e.draft IS FALSE
		AND b.draft IS FALSE 
		AND b.user_id = $2
		ORDER BY e.time DESC, e.title DESC
		LIMIT 1`
	e := entity.Event{}
	err = r.pool.QueryRow(context.Background(), query, end, *userID).Scan(
		&e.EventID,
		&e.Title,
		&e.Photo,
		&e.Text,
		&e.Entities,
		&e.Time,
		&e.Offline,
		&e.OfflinePaid,
		&e.Online,
		&e.OnlinePaid,
		&e.PaymentDetailsFirstName,
		&e.PaymentDetailsLastName,
		&e.PaymentDetailsAccount,
		&e.CreatedAt,
		&e.CreatedBy,
		&e.UpdatedAt,
		&e.UpdatedBy,
		&e.Draft,
	)
	if err != nil {
		return nil, false, err
	}
	return &e, true, nil
}

func (r *repository) GetByDate(t time.Time, userID *int64) ([]*entity.Event, error) {
	if userID == nil {
		var query = `SELECT
		event_id, 
		title, 
		photo, 
		text, 
		entities, 
		time, 
		offline, 
		offline_paid, 
		online, 
		online_paid, 
		payment_details_firstname,
		payment_details_lastname,
		payment_details_account,
		created_at, 
		created_by,
		updated_at, 
		updated_by,
		draft
		FROM events
		WHERE DATE(time) = DATE($1) AND draft IS FALSE
		ORDER BY time ASC, title ASC`
		rows, err := r.pool.Query(context.Background(), query, t)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var events []*entity.Event
		for rows.Next() {
			e := &entity.Event{}
			err := rows.Scan(
				&e.EventID,
				&e.Title,
				&e.Photo,
				&e.Text,
				&e.Entities,
				&e.Time,
				&e.Offline,
				&e.OfflinePaid,
				&e.Online,
				&e.OnlinePaid,
				&e.PaymentDetailsFirstName,
				&e.PaymentDetailsLastName,
				&e.PaymentDetailsAccount,
				&e.CreatedAt,
				&e.CreatedBy,
				&e.UpdatedAt,
				&e.UpdatedBy,
				&e.Draft,
			)
			if err != nil {
				return nil, err
			}
			events = append(events, e)
		}

		return events, nil
	}
	var query = `SELECT
	e.event_id, 
	e.title, 
	e.photo, 
	e.text, 
	e.entities, 
	e.time, 
	e.offline, 
	e.offline_paid, 
	e.online, 
	e.online_paid,
	e.payment_details_firstname,
	e.payment_details_lastname,
	e.payment_details_account,
	e.created_at, 
	e.created_by,
	e.updated_at, 
	e.updated_by,
	e.draft
	FROM events e JOIN bookings b ON e.event_id = b.event_id
	WHERE
	DATE(e.time) = DATE($1) 
	AND e.draft IS FALSE
	AND b.draft IS FALSE 
	AND b.user_id = $2
	ORDER BY e.time ASC, e.title ASC`
	rows, err := r.pool.Query(context.Background(), query, t, *userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*entity.Event
	for rows.Next() {
		e := &entity.Event{}
		err := rows.Scan(
			&e.EventID,
			&e.Title,
			&e.Photo,
			&e.Text,
			&e.Entities,
			&e.Time,
			&e.Offline,
			&e.OfflinePaid,
			&e.Online,
			&e.OnlinePaid,
			&e.PaymentDetailsFirstName,
			&e.PaymentDetailsLastName,
			&e.PaymentDetailsAccount,
			&e.CreatedAt,
			&e.CreatedBy,
			&e.UpdatedAt,
			&e.UpdatedBy,
			&e.Draft,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *repository) Update(e *entity.Event) (bool, error) {
	var query string
	query =
		`SELECT EXISTS (
        SELECT 1 FROM events 
        WHERE event_id = $1)`
	var exists bool
	err := r.pool.QueryRow(context.Background(), query, e.EventID).Scan(&exists)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}
	query = `
    	UPDATE events 
    	SET
		title = $2,
		photo = $3,
		text = $4,
		entities = $5,
		time = $6,
		offline = $7,
		offline_paid = $8,
		online = $9,
		online_paid = $10,
		payment_details_firstname = $11,
		payment_details_lastname = $12,
		payment_details_account = $13,
		created_at = $14,
		created_by = $15,
		updated_at = $16,
		updated_by = $17,
		draft = $18
    	WHERE event_id = $1`
	_, err = r.pool.Exec(context.Background(), query,
		e.EventID,
		e.Title,
		e.Photo,
		e.Text,
		e.Entities,
		e.Time,
		e.Offline,
		e.OfflinePaid,
		e.Online,
		e.OnlinePaid,
		e.PaymentDetailsFirstName,
		e.PaymentDetailsLastName,
		e.PaymentDetailsAccount,
		e.CreatedAt,
		e.CreatedBy,
		e.UpdatedAt,
		e.UpdatedBy,
		e.Draft,
	)
	return true, err
}

func (r *repository) generateEventID() (string, error) {
	if r.eventIDTypeNumber {
		var min, max int64 = 1, 10
		for range r.eventIDLen - 1 {
			min *= 10
			max *= 10
		}
		rand, err := rand.Int(rand.Reader, big.NewInt(max-min))
		if err != nil {
			return "", err
		}
		eventID_int := int(min + rand.Int64())
		return strconv.Itoa(eventID_int), nil
	}
	eventID := make([]byte, r.eventIDLen, r.eventIDLen)
	for index := range eventID {
		rand_big, err := rand.Int(rand.Reader, big.NewInt(int64(len(r.eventIDCharset))))
		if err != nil {
			return "", err
		}
		rand_int := int(rand_big.Int64())
		eventID[index] = r.eventIDCharset[rand_int]
	}
	return string(eventID), nil
}
