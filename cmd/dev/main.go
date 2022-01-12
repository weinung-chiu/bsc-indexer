package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"portto"
)

func main() {
	rpcEndpoint := "https://bsc-dataseed1.ninicoin.io/"
	//repo := portto.NewInMemoryStore()

	dsn := "root:mypasswd@tcp(127.0.0.1:3306)/bsc?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Silent),
	//})
	if err != nil {
		log.Fatal("failed to connect to db, ", err)
	}
	repo := portto.NewSQLStore(db)

	i, err := portto.NewIndexer(rpcEndpoint, repo)
	if err != nil {
		log.Fatal("failed to make new indexer, ", err)
	}

	go i.Run()

	block, _ := i.GetBlock(13000000)
	if block == nil {
		log.Println("block not found")
	} else {
		log.Printf("got block %d\n", block.Number)
		log.Printf("block hash : %s\n", block.Hash)
		log.Printf("block time : %d\n", block.Time)
		log.Printf("parent hash : %s\n", block.ParentHash)
	}

	block, _ = i.GetBlock(14274329)
	if block == nil {
		log.Println("block not found")
	} else {
		log.Printf("got block %d\n", block.Number)
		log.Printf("block hash : %s\n", block.Hash)
		log.Printf("block time : %d\n", block.Time)
		log.Printf("parent hash : %s\n", block.ParentHash)
	}
	log.Println("waiting...")
	log.Println("waiting...")
	log.Println("waiting...")
	time.Sleep(5 * time.Second)

	blocks, err := i.GetNewBlocks(3)
	if err != nil {
		log.Fatal(err)
	}
	for _, block := range blocks {
		log.Println(block.Number)
		log.Println(block.Hash)
	}
	log.Println("waiting...")
	log.Println("waiting...")
	log.Println("waiting...")
	time.Sleep(12 * time.Second)
	blocks, err = i.GetNewBlocks(3)
	if err != nil {
		log.Fatal(err)
	}
	for _, block := range blocks {
		log.Println(block.Number)
		log.Println(block.Hash)
	}

	fmt.Println("Press Ctrl+C to interrupt...")
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	for {
		select {
		case <-done:
			//todo: set context here

			fmt.Println("")

			fmt.Println("Bye Bye...")
			os.Exit(1)
		}
	}
}
