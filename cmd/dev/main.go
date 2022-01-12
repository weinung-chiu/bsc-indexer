package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/gin-gonic/gin"
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

	apiService := NewAPIService(i)
	go apiService.serveAndListen(":80")

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

func NewAPIService(i *portto.Indexer) *APIService {
	return &APIService{indexer: i}
}

type APIService struct {
	indexer *portto.Indexer
}

func (s APIService) serveAndListen(addr string) {
	r := gin.Default()
	r.GET("/blocks", s.blocksHandler)
	err := r.Run(addr)
	if err != nil {
		log.Fatal("failed to run http server, ", err)
	}
}

func (s APIService) blocksHandler(c *gin.Context) {
	limitRaw := c.DefaultQuery("limit", "1")
	limit, err := strconv.Atoi(limitRaw)
	if err != nil || limit < 1 || limit > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "limit must be a number and between 1 and 10",
		})
		return
	}

	blocks, err := s.indexer.GetNewBlocks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"blocks": blocks,
	})
}
