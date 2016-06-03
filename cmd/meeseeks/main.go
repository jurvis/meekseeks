package main

import (
	"log"
	"os"
	"time"

	"github.com/jurvis/meeseeks"
	"github.com/tucnak/telebot"
)

func main() {
	logger := log.New(os.Stdout, "[meeseeks]", 0)
	meeseek := meeseeks.InitMeeseek(logger)

	messages := make(chan telebot.Message)
	meeseek.Listen(messages, 1*time.Second)

	for message := range messages {
		meeseek.Router(message)
	}
}
