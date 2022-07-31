package app

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/analytics/internal/domain/events"
	"gitlab.com/g6834/team41/analytics/internal/domain/statistics"
	grpc "gitlab.com/g6834/team41/analytics/internal/grpc"
	mq "gitlab.com/g6834/team41/analytics/internal/kafka"
	"gitlab.com/g6834/team41/analytics/internal/pg"
	"gitlab.com/g6834/team41/analytics/internal/ports"

	"gitlab.com/g6834/team41/analytics/internal/env"

	"gitlab.com/g6834/team41/analytics/internal/http/handlers"
	"gitlab.com/g6834/team41/analytics/internal/http/middlewares"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "gitlab.com/g6834/team41/analytics/docs"
)

type App struct {
	m          *chi.Mux
	Auth       ports.AuthService
	Statistics ports.Statistics
	Events     ports.Events
}

func NewApp() *App {
	// Create postgres connection
	logrus.Debug("Connecting to postgres...")
	db, err := pg.NewPG(env.E.C.DB.Host, env.E.C.DB.User, env.E.C.DB.Password, env.E.C.DB.Name, env.E.C.DB.SSL, env.E.C.DB.Port)
	if err != nil {
		logrus.Panic(fmt.Errorf("failed to connect to postgres: %w", err))
	}
	repo := pg.NewAnalytics(db)

	return &App{
		Auth:       grpc.NewClient(env.E.C.AuthAddress),
		Statistics: statistics.New(repo),
		Events:     events.New(repo),
		m:          chi.NewRouter(),
	}
}

func (a *App) Run() error {
	a.bindHandlers()
	a.bindSwagger()
	a.bindProfiler()

	//logrus.Info("Starting gRPC server...")
	//go grpc.StartServer(env.E.C.GrpcAddress, a.Events)

	logrus.Info("Starting MQ Consumer...")
	go mq.StartConsumer(
		context.TODO(),
		env.E.C.Kafka.Brokers,
		env.E.C.Kafka.Topic,
		env.E.C.Kafka.GroupId,
		a.Events,
	)

	logrus.Info("Starting server...")
	return http.ListenAndServe(env.E.C.HostAddress, a.m)
}

const (
	CountAccepted = "/count-accepted"
	CountDeclined = "/count-declined"
	SumReaction   = "/sum-reaction"
)

func (a *App) bindHandlers() {

	a.m.Route("/analytics", func(r chi.Router) {
		r.Use(middlewares.GetCheckAuthFunc(a.Auth))

		r.Handle(CountAccepted, handlers.CountAcceptedTask{Statistics: a.Statistics})
		r.Handle(CountDeclined, handlers.CountDeclinedTask{Statistics: a.Statistics})
		r.Handle(SumReaction, handlers.SumReaction{Statistics: a.Statistics})
	})
}

func (a *App) bindProfiler() {
	a.m.Route("/debug/pprof", func(r chi.Router) {
		r.HandleFunc("/", pprof.Index)
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
	})
}

func (a *App) bindSwagger() {
	a.m.Route("/swagger", func(r chi.Router) {
		r.HandleFunc("/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost"+env.E.C.HostAddress+"/swagger/doc.json"),
		))
	})
}
