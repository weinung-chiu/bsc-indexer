package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"portto"
)

func main() {
	rpcEndpoint := "https://bsc-dataseed1.ninicoin.io/"

	dsn := "root:mypasswd@tcp(127.0.0.1:3306)/bsc?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to db, ", err)
	}
	repo := portto.NewSQLStore(db)

	i, err := portto.NewIndexer(rpcEndpoint, repo)
	if err != nil {
		log.Fatal("failed to make new indexer, ", err)
	}

	apiService := portto.NewAPIService(i)
	go apiService.ListenAndServe(":80")

	ctx, cancel := context.WithCancel(context.Background())
	go i.Run(ctx)

	fmt.Println("Press Ctrl+C to interrupt...")
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	for {
		select {
		case <-done:
			cancel()
			i.StopWait()
			fmt.Println("")
			fmt.Println("Bye Bye")
			os.Exit(1)
		}
	}
}
