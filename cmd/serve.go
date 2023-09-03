/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/bamdadam/rubbit/internal/db/rdb"
	"github.com/bamdadam/rubbit/internal/http/handler"
	"github.com/bamdadam/rubbit/internal/rabbit"
	redisHandler "github.com/bamdadam/rubbit/internal/store/rdb"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "serve",
	Short: "start a rabbit server",
	Long:  `start a rabbit server which can send normal or delayed messages`,
	Run: func(cmd *cobra.Command, args []string) {
		rdb, err := rdb.New(context.Background())
		if err != nil {
			log.Fatal("can't make redis connection: ", err)
		}
		rh := redisHandler.New(rdb.Client)
		rb, err := rabbit.InitRabbitHandler()
		if err != nil {
			log.Fatal("error while running rabbit server: ", err)
		}
		app := fiber.New(
			fiber.Config{
				AppName: "Rubbit Server",
			},
		)
		g := app.Group("/")
		hl := handler.Handler{
			RH:  rb,
			RDB: rh,
		}
		hl.RegisterHandler(g)
		err = app.Listen("0.0.0.0:8080")
		if err != nil {
			log.Fatal("error while running fiber server: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
