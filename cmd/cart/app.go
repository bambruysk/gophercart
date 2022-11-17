package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

type Application struct {
	signalChannel chan os.Signal
	exitSignal    chan struct{}

	logger *logrus.Logger

	onStart    func()
	onShutdown func() error
}

func NewApplication(logger *logrus.Logger, onStart func(), onShutdown func() error) *Application {

	app := &Application{logger: logger, onStart: onStart, onShutdown: onShutdown}

	app.signalChannel = make(chan os.Signal, 1)
	app.exitSignal = make(chan struct{})

	return app
}

func (a *Application) initSignals() {
	signal.Notify(a.signalChannel, syscall.SIGTERM)
	signal.Notify(a.signalChannel, syscall.SIGINT)
	signal.Notify(a.signalChannel, syscall.SIGKILL)
	for s := range a.signalChannel {
		switch s {
		case syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL:
			close(a.signalChannel)
			a.logger.Warnf("We got %s, shutdown application...", s)
			close(a.exitSignal)
			return
		}
	}
}

func (a *Application) Run() {
	a.logger.Info("start application")
	a.initSignals()
	<-a.exitSignal
}
