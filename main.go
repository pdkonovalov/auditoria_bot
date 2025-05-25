package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/pdkonovalov/auditoria_bot/internal/config"
	"github.com/pdkonovalov/auditoria_bot/internal/infrastructure/database/postgres"
	booking_repository "github.com/pdkonovalov/auditoria_bot/internal/infrastructure/repository/postgres/booking"
	event_repository "github.com/pdkonovalov/auditoria_bot/internal/infrastructure/repository/postgres/event"
	user_repository "github.com/pdkonovalov/auditoria_bot/internal/infrastructure/repository/postgres/user"
	"github.com/pdkonovalov/auditoria_bot/internal/infrastructure/validator"
	"github.com/pdkonovalov/auditoria_bot/internal/interfaces/telegram"
)

func main() {
	log.Print("read configuration...")

	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Error init config: %s", err)
	}

	log.Print("configuration successfull loaded:")
	cfg_masked, err := cfg.StringSecureMasked()
	if err != nil {
		log.Fatalf("Error print config: %s", err)
	}

	log.Print("\n" + cfg_masked)

	log.Print("connect to postgres...")

	pool, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("Error init postgres: %s", err)
	}

	log.Print("connect to postgres successfull")

	log.Print("init repositories...")

	userRepo, err := user_repository.New(pool)
	if err != nil {
		log.Fatalf("Error init user repository: %s", err)
	}

	eventRepo, err := event_repository.New(cfg, pool)
	if err != nil {
		log.Fatalf("Error init event repository: %s", err)
	}

	bookingRepo, err := booking_repository.New(pool)
	if err != nil {
		log.Fatalf("Error init booking repository: %s", err)
	}

	log.Print("init repositories successfull")

	log.Print("configure bot...")

	validator, err := validator.New()
	if err != nil {
		log.Fatalf("Error init validator: %s", err)
	}

	bot, err := telegram.New(cfg, userRepo, eventRepo, bookingRepo, validator)
	if err != nil {
		log.Fatalf("Error init bot: %s", err)
	}

	log.Print("configure bot successfull")
	log.Print("starting bot...")

	go bot.Start()

	log.Print("bot successfull started")
	log.Print("press ctrl c to shutdown")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Print("shutdown bot...")
		pool.Close()
		bot.Stop()
	}()
	wg.Wait()
	log.Print("bot stopped")
	log.Print("exit")
}
