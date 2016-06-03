package meeseeks

import "strings"

// SayHello says hi to you
func (m *Meeseeks) SayHello(msg *message) {
	m.SendMessage(msg.Chat, "Hello there, "+msg.Sender.FirstName+"!", nil)
}

// Clear returns a message that clears out the folder
func (m *Meeseeks) Clear(msg *message) {
	m.SendMessage(msg.Chat, "Lol, sure."+strings.Repeat("\n", 41)+"Cleared.", nil)
}

// Start returns some help text.
func (m *Meeseeks) Start(msg *message) {
	m.SendMessage(msg.Chat, `Hi there! I can help you with the following things:

    /Hello - say Hello to the bot
    /Clear - clears out your NSFW crap
    /urbandict - does an Urban Dictionary search

    Give these commands a try!`, nil)
}

// Help returns some help text
func (m *Meeseeks) Help(msg *message) {
	m.SendMessage(msg.Chat, `Some commands:

    /Hello - say Hello to the bot
    /Clear - clears out your NSFW crap
    /urbandict - does an Urban Dictionary search
  `, nil)
}
