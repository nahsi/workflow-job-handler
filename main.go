package main

import (
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-playground/webhooks/v6/github"
)

type Config struct {
	RedisDSN string `env:"REDIS_DSN"`
	Path     string `env:"WEBHOOK_PATH"`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	hook, _ := github.New(github.Options.Secret("haph8Cioteing3ne"))

	http.HandleFunc(cfg.Path, func(w http.ResponseWriter, r *http.Request) {
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
