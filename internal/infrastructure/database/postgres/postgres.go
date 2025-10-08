package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/pdkonovalov/auditoria_bot/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(cfg *config.Config) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDatabase, cfg.PostgresSslMode,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	var pool *pgxpool.Pool
	maxAttempts := 5
	backoff := time.Second * 2

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			if err = pool.Ping(context.Background()); err == nil {
				break
			}
			pool.Close()
		}

		if attempt == maxAttempts {
			return nil, fmt.Errorf("Failed to connect to database after %d attempts: %w", maxAttempts, err)
		}

		time.Sleep(backoff)
		backoff *= 2
	}
	_, err = pool.Exec(context.Background(), fmt.Sprintf("set time zone '%s';", cfg.TelegramBotTimezone))
	if err != nil {
		return nil, fmt.Errorf("Failed set timezone: %s", err)
	}

	_, err = pool.Exec(context.Background(),
		`create table if not exists users (
		user_id bigint primary key,
		username text,
		firstname text,
		lastname text,
		contact_info text,
		state text,
		context jsonb,
		admin bool);

		create table if not exists events (
		event_id text primary key,
		title text,
		photo jsonb,
		text text,
		entities jsonb,
		time timestamptz,
		offline bool,
		offline_paid bool,
		online bool,
		online_paid bool,
		payment_details_firstname text,
		payment_details_lastname text,
		payment_details_account text,
		created_at timestamptz,
		created_by bigint references users(user_id),
		updated_at timestamptz,
		updated_by bigint references users(user_id),
		draft bool);

		create table if not exists bookings (
		event_id text references events(event_id) on delete cascade,
		user_id bigint references users(user_id) on delete cascade,
		payment bool,
		payment_photo jsonb,
		payment_document jsonb,
		text text,
		offline bool,
		online bool,
		check_in bool,
		check_in_at timestamptz,
		check_in_by bigint references users(user_id),
		created_at timestamptz,
		updated_at timestamptz,
		draft bool default true);`)
	if err != nil {
		return nil, fmt.Errorf("Failed init tables: %s", err)
	}

	return pool, nil
}
