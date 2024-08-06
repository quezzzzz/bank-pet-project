package main

import (
	"bank"
	"bank/pkg/handler"
	"bank/pkg/repository"
	"bank/pkg/service"
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error while initialization config : %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error while loading env param : %s", err.Error())

	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("error while db init : %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(bank.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error while running server")
		}
	}()
	logrus.Println("Bank started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Println("Bank shutting down")
	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Error("error while shutting down server", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Error("error while closing connection with db", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
