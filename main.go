package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/juankhusuma/scele-bot/handlers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Token := os.Getenv("DC_TOKEN")
	dc, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal(err)
	}

	dc.AddHandler(handlers.MessageCreate)
	dc.AddHandler(handlers.InteractionCreate)

	err = dc.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Bot alive!")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dc.Close()
}
