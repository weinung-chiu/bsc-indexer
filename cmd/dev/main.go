package main

import (
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
	//repo := portto.NewInMemoryStore()

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

	go i.Run()

	fmt.Println("Press Ctrl+C to interrupt...")
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	for {
		select {
		case <-done:
			//todo: set context here

			fmt.Println("")

			//block, _ := i.GetBlock(14276602)
			//if block == nil {
			//	log.Println("block not found")
			//} else {
			//	log.Printf("got block %d\n", block.Number().Uint64())
			//	log.Printf("block hash : %s\n", block.Hash().String())
			//	log.Printf("block time : %d\n", block.Number().Uint64())
			//	log.Printf("parent hash : %s\n", block.ParentHash().String())
			//}
			//block, _ = i.GetBlock(14274329)
			//if block == nil {
			//	log.Println("block not found")
			//} else {
			//	log.Printf("got block %d\n", block.Number().Uint64())
			//	log.Printf("block hash : %s\n", block.Hash().String())
			//	log.Printf("block time : %d\n", block.Number().Uint64())
			//	log.Printf("parent hash : %s\n", block.ParentHash().String())
			//}

			fmt.Println("Bye Bye...")
			os.Exit(1)
		}
	}
}
