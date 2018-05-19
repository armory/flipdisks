package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
	"strings"
)

func main() {
	// api := slack.New("xoxb-16939015395-366505293908-nSI2CGuND5Uehi5ZvygNEUfO")
	api := slack.New("xoxb-367584659222-366507826004-5dx77oPn5VasWfesv86rKW6h")
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			fmt.Printf("Message: %+v\n\n", ev)
			fmt.Printf("GOT MESSAGE: %+v\n", ev.Msg.Text)
			flipbotUserId := "UASEXQA04"
			if strings.Contains(ev.Msg.Text, flipbotUserId) {
				flipTableWordList := []string{
					"flip",
					"table",
				}

				curseWords := []string{
					"fuck",
					"god",
					"damn",
					"ass",
				}


				for _, word := range curseWords {
					if strings.Contains(ev.Msg.Text, word) {
						rtm.SendMessage(rtm.NewOutgoingMessage("Yo, watch your language, you dick head...", "CATMSGF7H"))
						break
					}
				}

				for _, word := range flipTableWordList {
					if strings.Contains(ev.Msg.Text, word) {
						rtm.SendMessage(rtm.NewOutgoingMessage("Flipping the table", "CATMSGF7H"))
						break
					}
				}
			}

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:
		}
	}
}
