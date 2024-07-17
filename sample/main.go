package main

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"os"

	daprd "github.com/dapr/go-sdk/service/grpc"

	"github.com/dapr/go-sdk/service/common"
	"go.uber.org/zap"
)

var log *zap.Logger

func init() {
	log, _ = zap.NewDevelopment()
}

func main() {

	d := getSleepTime()
	log.Info("starting applicaton", zap.Duration("sleep", d))

	// **** SLEEP ****
	time.Sleep(d)

	log.Info("creating dapr service")
	s, err := daprd.NewService(":51051")
	if err != nil {
		log.Fatal("error creating dapr service", zap.Error(err))
		return
	}

	// Add binding  handlers
	for k, v := range map[string]common.BindingInvocationHandler{
		"h1-binding": h1,
		"h2-binding": h2,
	} {
		log.Info("add cron binding handler", zap.String("name", k))
		if err := s.AddBindingInvocationHandler(k, v); err != nil {
			log.Error("error adding cron binding handler", zap.Error(err))
		}
	}

	sc, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-sc.Done()
		s.GracefulStop()
		os.Exit(0)
	}()

	log.Info("starting dapr service")
	if err := s.Start(); err != nil {
		log.Fatal("error starting dapr service", zap.Error(err))
		return
	}
}

// getSleepTime returns the sleep time from the environment variable ZZZ_SLEEP. If not set it returns 0
func getSleepTime() time.Duration {
	e := os.Getenv("ZZZ_SLEEP")
	if e == "" {
		return 0
	}

	log.Info("sleep duration", zap.String("ZZZ_SLEEP", e))
	d, err := time.ParseDuration(e)
	if err != nil {
		log.Error("error parsing duration", zap.Error(err))
		return 0
	}
	return d
}

// h1 is a handler for the moris-cron-sample binding
func h1(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Info("h1 handler was invoked")
	return nil, nil
}

// h2 is a handler for the billing-observe-receiving binding
func h2(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	log.Info("h1 handler was invoked")
	return nil, nil
}
