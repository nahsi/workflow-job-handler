package main

import (
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/webhooks/v6/github"
)

type Config struct {
	RedisDSN      string `env:"REDIS_DSN"`
	WebhookPath   string `env:"WEBHOOK_PATH"`
	WebhookSecret string `env:"WEBHOOK_SECRET"`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	hook, _ := github.New(github.Options.Secret(cfg.WebhookSecret))

	http.HandleFunc(cfg.WebhookPath, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.WorkflowJobEvent)
		if err != nil {
			fmt.Printf("%+v\n", err)
		}
		p := payload.(github.WorkflowJobPayload)
		redis := InitStorage(cfg)
		redis.Put(intoWorkFlowJob(&p))
	})

	http.ListenAndServe(":8123", nil)
}
