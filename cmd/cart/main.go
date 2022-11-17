package main

import (
	"context"

	"github.com/sirupsen/logrus"

	"gophercart/internal/options"
	"gophercart/internal/repo"
	"gophercart/internal/usecase"
)

var (
	appCtx       context.Context
	appCtxCancel context.CancelFunc

	repository repo.Repo

	log *logrus.Logger
)

func main() {
	log = logrus.StandardLogger()
	app := NewApplication(log, OnStart, OnShutdown)
	app.Run()
}

const config = "conf/conf.toml"

func OnStart() {
	appCtx, appCtxCancel = context.WithCancel(context.Background())

	opts, err := options.NewOptionsFromFile(config)
	if err != nil {
		log.Fatalln("unable load options")
	}

	repository, err = repo.NewRepository(appCtx, opts.RepoOptions, log)

	proc := usecase.New(log, repository, nil)

	_ = proc

}

func OnShutdown() error {
	appCtxCancel()

	return nil
}
