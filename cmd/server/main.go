package main

import (
	"context"
	"go_payment_microservice/cmd/server/middlewares"
	"go_payment_microservice/internal/logger"
	//"log"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"os/signal"
	//"runtime/pprof"
	"syscall"

	"go_payment_microservice/cmd/server/handlers"
	"go_payment_microservice/internal/clients"
	"go_payment_microservice/internal/config"
	"go_payment_microservice/internal/services"

	"github.com/gofiber/fiber/v2"
	_ "go_payment_microservice/cmd/server/docs"
)

// @title PaymentsAPI
// @version 1.0
// @description Payments API microservice

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		slog.ErrorContext(context.Background(), err.Error())
		panic(err)
	}

	logger.InitLogger(cfg)

	//if cfg.EnableCPUProfiler == "true" {
	//	f, err := os.Create("/tmp/profiles/cpu.pprof")
	//	if err != nil {
	//		log.Fatal("could not create CPU profile: ", err)
	//	}
	//	defer func(f *os.File) {
	//		err := f.Close()
	//		if err != nil {
	//			log.Fatal("could not close CPU profile: ", err)
	//		}
	//	}(f)
	//
	//	if err := pprof.StartCPUProfile(f); err != nil {
	//		log.Fatal("could not start CPU profile: ", err)
	//	}
	//	defer pprof.StopCPUProfile()
	//
	//	log.Println("CPU profiling started")
	//	// Profiling will run for the duration of the main logic
	//}

	//_, err = pyroscope.Start(pyroscope.Config{
	//	ApplicationName: "messages.api",
	//	ServerAddress:   cfg.PyroscopeServer,
	//	Logger:          nil,
	//	ProfileTypes: []pyroscope.ProfileType{
	//		pyroscope.ProfileCPU,
	//		pyroscope.ProfileAllocObjects,
	//		pyroscope.ProfileInuseObjects,
	//		pyroscope.ProfileAllocSpace,
	//		pyroscope.ProfileMutexCount,
	//		pyroscope.ProfileMutexDuration,
	//		pyroscope.ProfileBlockCount,
	//		pyroscope.ProfileBlockDuration,
	//	},
	//})
	//if err != nil {
	//	logger.GetLogger().Fatal(ctx, err)
	//}

	c, err := clients.NewClients(ctx, cfg)
	if err != nil {
		logger.GetLogger().Fatal(ctx, err)
	}

	s := services.NewServices(c, cfg)
	m := middlewares.NewMiddlewares(s, cfg)
	h := handlers.NewHandlers(cfg, s, m)

	app := fiber.New()

	h.RegisterRoutes(app)

	go func() {
		err := app.Listen(":8080")
		if err != nil {
			logger.GetLogger().Fatal(context.Background(), err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	err = app.Shutdown()
	if err != nil {
		logger.GetLogger().Fatal(context.Background(), err)
	}
}
