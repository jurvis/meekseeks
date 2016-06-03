package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/jurvis/meeseeks"
	"github.com/kardianos/osext"
	"github.com/tucnak/telebot"
)

func main() {
	pwd, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatalf("error getting executable folder: %s", err)
	}

	configJSON, err := ioutil.ReadFile(path.Join(pwd, "config.json"))
	if err != nil {
		log.Fatalf("error reading config file! Boo: %s", err)
	}

	logger := log.New(os.Stdout, "[meeseeks]", 0)
	meeseek := meeseeks.InitMeeseek(configJSON, logger)

	messages := make(chan telebot.Message)
	meeseek.Listen(messages, 1*time.Second)

	for message := range messages {
		meeseek.Router(message)
	}
}
