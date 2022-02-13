package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/rs/zerolog"
)

type Config struct {
	RedisDSN      string `env:"REDIS_DSN"`
	WebhookPath   string `env:"WEBHOOK_PATH"`
	WebhookSecret string `env:"WEBHOOK_SECRET,unset"`
}

func main() {
	// Set up logging
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerologr.NameFieldName = "logger"
	zerologr.NameSeparator = "/"

	zl := zerolog.New(os.Stderr).With().Caller().Timestamp().Logger()
	var log logr.Logger = zerologr.New(&zl)

	// Parse settings from enviromental variables
	cfg := Config{}
	opts := &env.Options{RequiredIfNoDef: true}
	if err := env.Parse(&cfg, *opts); err != nil {
		log.Error(err, "error during startup")
		os.Exit(1)
	}

	hook, _ := github.New(github.Options.Secret(cfg.WebhookSecret))

	http.HandleFunc(cfg.WebhookPath, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.WorkflowJobEvent)
		if err != nil {
			log.Error(err, "failed to parse event")
		}

		redis, err := InitStorage(cfg)
		if err != nil {
			log.Error(err, "error during startup")
			os.Exit(1)
		}

		switch payload.(type) {
		case github.WorkflowJobPayload:
			p := intoWorkFlowJob(payload.(github.WorkflowJobPayload))
			if err := redis.Put(p); err != nil {
				log.Error(err, "write to Redis failed")
			}

		default:
			log.Info(
				"accepting only workflow_job events",
				"event", fmt.Sprintf("%v", payload.(github.Event)),
			)
		}
	})

	if err := http.ListenAndServe(":8123", nil); err != nil {
		log.Error(err, "error during startup")
		os.Exit(1)
	}
}
