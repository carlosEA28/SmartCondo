package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/carlosEA28/smartcondo/internal/config"
	"github.com/carlosEA28/smartcondo/internal/database"
	"github.com/carlosEA28/smartcondo/internal/logger"
	mailtrapProvider "github.com/carlosEA28/smartcondo/internal/providers/mailtrap"
	"github.com/carlosEA28/smartcondo/internal/repositories"
	"github.com/carlosEA28/smartcondo/internal/services"
	"github.com/rs/zerolog"
)

func main() {
	log := logger.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	pagamentoRepo := repositories.NewGormPagamentoRepository(db)
	notificacaoRepo := repositories.NewGormNotificacaoRepository(db)
	mailtrap := mailtrapProvider.New(
		cfg.Mailtrap.Host,
		cfg.Mailtrap.Port,
		cfg.Mailtrap.Username,
		cfg.Mailtrap.Password,
		cfg.Mailtrap.FromEmail,
	)

	emailService := services.NewEmailService(pagamentoRepo, notificacaoRepo, mailtrap, cfg.Mailtrap.FromEmail)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	now := time.Now()
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
	if now.After(nextRun) {
		nextRun = nextRun.Add(24 * time.Hour)
	}

	log.Info().Time("next_run", nextRun).Msg("cron scheduler started, next run at 8:00 AM")

	firstTimer := time.NewTimer(time.Until(nextRun))
	defer firstTimer.Stop()

	for {
		select {
		case <-firstTimer.C:
			runJob(ctx, emailService, log)
			ticker.Reset(24 * time.Hour)
		case <-ticker.C:
			runJob(ctx, emailService, log)
		case <-sigCh:
			log.Info().Msg("shutting down cron scheduler")
			return
		}
	}
}

func runJob(ctx context.Context, emailService *services.EmailService, log zerolog.Logger) {
	log.Info().Msg("starting inadimplente email job")
	if err := emailService.SendInadimplenteEmails(ctx); err != nil {
		log.Error().Err(err).Msg("inadimplente email job failed")
		return
	}
	log.Info().Msg("inadimplente email job completed")
}
