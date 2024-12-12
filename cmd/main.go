package main

import (
	"bkh-ecom/internal/app"
	"bkh-ecom/internal/config"
	"bkh-ecom/internal/logger"
	"bkh-ecom/internal/repository"
	"bkh-ecom/internal/service"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGQUIT, syscall.SIGTERM)
	defer stop()

	app.InitLogger(ctx)

	conf, err := config.NewConfig(ctx)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{Panic: "Error with reading config"})
	}

	db, err := connectToPostgres(ctx, conf)
	if err != nil {
		logger.ErrorKV(ctx, logger.Data{Panic: "Error connecting to DB", Detail: fmt.Sprintf("%+v", conf.DB)})
	}

	dao := repository.NewDAO(db)
	clickService := service.NewClickService(ctx, dao)

	httpSrv := app.ServerInit(ctx, clickService)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.WarnKV(gCtx, logger.Data{Msg: "Postgres connected!"})
		return db.Ping(ctx)
	})
	g.Go(func() error {
		fmt.Println("Server started!")
		addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
		return httpSrv.Listen(addr)
	})
	g.Go(func() error {
		<-gCtx.Done()
		fmt.Println("Server is shut down.")
		return httpSrv.ShutdownWithTimeout(time.Second * 5)
	})
	g.Go(func() error {
		<-gCtx.Done()
		db.Close()
		return nil
	})

	if err = g.Wait(); err != nil {
		logger.ErrorKV(gCtx, logger.Data{
			Error:  err,
			Msg:    "Error group wait got error",
			Detail: fmt.Sprintf("exit reason: %v", err),
		})
	}

	logger.WarnKV(gCtx, logger.Data{Msg: "Server was gracefully shut down."})
}

func connectToPostgres(ctx context.Context, config *config.Entity) (*pgxpool.Pool, error) {

	connString := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable&pool_max_conns=%v",
		config.DB.User, config.DB.Pass, config.DB.Host, config.DB.Port, config.DB.Name, config.DB.PoolSize)
	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}
	pgxConfig.MaxConnLifetime = time.Minute * 10

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
