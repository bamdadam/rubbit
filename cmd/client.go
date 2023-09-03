/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/bamdadam/rubbit/internal/db/rdb"
	"github.com/bamdadam/rubbit/internal/rabbit"
	redisHandler "github.com/bamdadam/rubbit/internal/store/rdb"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "start a rabbit client",
	Long:  `stary a rabbit client which reads messages based on its topic and saves them into redis`,
	Run: func(cmd *cobra.Command, args []string) {
		rdb, err := rdb.New(context.Background())
		if err != nil {
			log.Fatal("can't make redis connection: ", err)
		}
		rh := redisHandler.New(rdb.Client)
		err = rabbit.InitRabbitClient(rh)
		if err != nil {
			log.Fatal("error while running rabbit client: ", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
