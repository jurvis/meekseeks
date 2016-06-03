package meeseeks

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/tucnak/telebot"
)

// Meeseeks is the main struct. All response funcs bind to this
type Meeseeks struct {
	Name string
	bot  *telebot.Bot
	log  *log.Logger
	fmap FuncMap
	keys config
}

// Configuration struct for setting up Meeseeks
type config struct {
	Name           string `json:"name"`
	TelegramAPIKey string `json:"telegram_api_key"`
}

type message struct {
	Cmd  string
	Args []string
	*telebot.Message
}

// GetArgs prints out arguments for the message in a string
func (m message) GetArgString() string {
	argString := ""
	for _, s := range m.Args {
		argString += s + " "
	}
	return strings.TrimSpace(argString)
}

// A FuncMap is a map of command strings to response functions
// It is use for routing commands to responses.
type FuncMap map[string]ResponseFunc

// ResponseFunc is a handler for a bot command
type ResponseFunc func(m *message)

// InitMeeseek initialises a Meeseek
func InitMeeseek(lg *log.Logger) *Meeseeks {
	if lg == nil {
		lg = log.New(os.Stdout, "[meeseeks]", 0)
	}

	bot, err := telebot.NewBot(os.Getenv("TELEGRAM_API_KEY"))
	if err != nil {
		log.Fatalf("error creating a new bot :/ %s", err)
	}

	m := &Meeseeks{Name: os.Getenv("TELEGRAM_BOT_NAME"), bot: bot, log: lg}

	m.fmap = m.getDefaultFuncMap()

	return m
}

// Listen exposes telebot Listen API
func (m *Meeseeks) Listen(subscription chan telebot.Message, timeout time.Duration) {
	m.bot.Listen(subscription, timeout)
}

func (m *Meeseeks) getDefaultFuncMap() FuncMap {
	return FuncMap{
		"/start":     m.Start,
		"/help":      m.Help,
		"/hello":     m.SayHello,
		"/clear":     m.Clear,
		"/urbandict": m.UrbanDictSearch,
		"/ud":        m.UrbanDictSearch,
	}
}

// AddFunction adds a response function to the FuncMap
func (m *Meeseeks) AddFunction(command string, resp ResponseFunc) error {
	if !strings.HasPrefix(command, "/") {
		return fmt.Errorf("not a valid command string - it should be of the format /something")
	}
	m.fmap[command] = resp
	return nil
}

// Router routes received Telegram messages to the appropriate response functions.
func (m *Meeseeks) Router(msg telebot.Message) {
	jmsg := m.parseMessage(&msg)
	if jmsg.Cmd != "" {
		m.log.Printf("[%s][id: %d] command: %s, args: %s", time.Now().Format(time.RFC3339), jmsg.ID, jmsg.Cmd, jmsg.GetArgString())
	}
	execFn := m.fmap[jmsg.Cmd]

	if execFn != nil {
		m.GoSafely(func() { execFn(jmsg) })
	}
}

// GoSafely is a utility wrapper to recover and log panics in goroutines.
// If we use naked goroutines, a panic in any one of them crashes
// the whole program. Using GoSafely prevents this.
func (m *Meeseeks) GoSafely(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]

				m.log.Printf("PANIC: %s\n%s", err, stack)
			}
		}()

		fn()
	}()
}

func (m *Meeseeks) parseMessage(msg *telebot.Message) *message {
	cmd := ""
	args := []string{}

	if msg.IsReply() {
		r := regexp.MustCompile(`\/\w*`)
		res := r.FindString(msg.ReplyTo.Text)
		for k, _ := range m.fmap {
			if res == k {
				cmd = k
				args = strings.Split(msg.Text, " ")
				break
			}
		}
	} else if msg.Text != "" {
		msgTokens := strings.Fields(msg.Text)
		cmd, args = strings.ToLower(msgTokens[0]), msgTokens[1:]
		if strings.Contains(cmd, "@") {
			c := strings.Split(cmd, "@")
			cmd = c[0]
		}
	}

	return &message{Cmd: cmd, Args: args, Message: msg}
}

// SendMessage sends the message using the Meeseeks struct
func (m *Meeseeks) SendMessage(recipient telebot.Recipient, msg string, options *telebot.SendOptions) {
	m.bot.SendMessage(recipient, msg, options)
}
