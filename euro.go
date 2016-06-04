package meeseeks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/tucnak/telebot"
)

const euroTeamsURL = "http://api.football-data.org/v1/soccerseasons/424/teams"

// EuroFixturesSearch Queries http://api.football-data.org/v1/soccerseasons/424 to retrieve some options
func (m *Meeseeks) EuroFixturesSearch(msg *message) {
	if len(msg.Args) == 0 {
		resp, err := http.Get(euroTeamsURL)

		if err != nil {
			m.log.Printf("failure retrieving teams from API for query '%s': %s", euroTeamsURL, err)
		}

		jsonBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			m.log.Printf("failure reading json results for Euro search query '%s': %s", euroTeamsURL, err)
			return
		}

		teamRes := struct {
			List []struct {
				Name string `json:"name"`
			} `json:"teams"`
		}{}

		err = json.Unmarshal(jsonBody, &teamRes)
		if err != nil {
			m.log.Printf("failure unmarshalling json for Euro search query '%s': %s", euroTeamsURL, err)
			return
		}

		options := make([][]string, 12)
		count := 0
		for i := 0; i < len(teamRes.List)/2; i++ {
			options[i] = make([]string, 2)
			for j := 0; j < 2; j++ {
				options[i][j] = teamRes.List[count].Name
				count++
			}
		}

		so := &telebot.SendOptions{ReplyTo: *msg.Message, ReplyMarkup: telebot.ReplyMarkup{
			ForceReply:      true,
			Selective:       true,
			CustomKeyboard:  options,
			OneTimeKeyboard: true,
		}}
		m.SendMessage(msg.Chat, "/euro: Tell me which team you are curious about\nHere are some commands to try: \n* France\n* Germany\n\n\U0001F4A1", so)
		return
	}

}
