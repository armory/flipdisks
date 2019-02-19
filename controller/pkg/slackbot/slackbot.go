package slackbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"text/template"

	"github.com/armory/flipdisks/controller/pkg/flipboard"
	"github.com/armory/flipdisks/controller/pkg/github"
	"github.com/armory/flipdisks/controller/pkg/options"
	"github.com/nlopes/slack"
)

type Slack struct {
	token             string
	githubEmojiLookup github.EmojiLookup
	RTM               *slack.RTM
}

func NewSlack(token string, g github.EmojiLookup) *Slack {
	api := slack.New(token)
	rtm := api.NewRTM()

	go rtm.ManageConnection()

	return &Slack{
		token:             token,
		githubEmojiLookup: g,
		RTM:               rtm,
	}
}

func (s *Slack) StartSlackListener(board *flipboard.Flipboard) {
	for msg := range s.RTM.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			go s.handleSlackMsg(event, board)

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		case *slack.ConnectionErrorEvent:
			fmt.Println("Connection Error!")
			fmt.Printf("%#v\n", event.ErrorObj.Error())

		default:
			fmt.Println("Event Received: ")
			fmt.Printf("%#v\n", msg)
			fmt.Printf("%#v\n", msg.Data)
		}
	}
}

func (s *Slack) handleSlackMsg(slackEvent *slack.MessageEvent, board *flipboard.Flipboard) {
	rawMsg := slackEvent.Msg.Text
	if slackEvent.SubMessage != nil {
		rawMsg = slackEvent.SubMessage.Text
	}

	if strings.HasPrefix(rawMsg, s.getMyUserIdFormatted()) {
		if strings.Contains(strings.ToLower(rawMsg), "ssh") {
			s.respondWithSSHConnectionString(slackEvent.Msg.Channel)
			return
		}

		if rawMsg == "" {
			return // let's just ignore it if we don't have anything to display on the board
		}

		rawMsg = s.editSettings(rawMsg, board, slackEvent)
	}

	fmt.Printf("Raw Slack Message: %+v\n", rawMsg)

	messages := options.SplitMessageAndOptions(rawMsg)

	fmt.Printf("%#v \n", messages)

	// do some message cleanup because of slack formatting
	for _, msg := range messages {
		if strings.ToLower(msg.Message) == "help" {
			s.respondWithHelpMsg(slackEvent.Msg.Channel)
			return
		}

		msg.Message = s.renderSlackUsernames(msg.Message)
		msg.Message = cleanupSlackEncodedCharacters(msg.Message)
		msg.Message = s.renderSlackEmojis(msg.Message)

		board.Enqueue(&msg)
	}
}

func (s *Slack) editSettings(rawMsg string, board *flipboard.Flipboard, event *slack.MessageEvent) string {
	// remove my userId from the message
	cleanMsg := strings.Replace(rawMsg, s.getMyUserIdFormatted(), "", -1)
	settings := strings.Split(strings.TrimSpace(cleanMsg), " ")

	settingName := settings[0]

	switch settingName {
	case "countdown":
		fallthrough
	case "countdownClock":
		val := settings[1]

		if val == "enable" || val == "enabled" {
			flipboard.EnableCountdownClock(board)
			return "enabled countdown"
		} else if val == "disable" || val == "disabled" {
			flipboard.DisableCountdownClock(board)
			return "disabled countdown"
		} else {
			if err := flipboard.SetCountdownClock(board, val); err != nil {
				return "date not set, " + err.Error()
			}
			return "setting countdown to " + val
		}
	case "help":
		s.respondWithSettingsHelpMessage(event.Msg.Channel)
		return ""
	}

	s.respondWithSettingsHelpMessage(event.Msg.Channel)
	return "received an unknown setting"
}

func (s *Slack) getMyUserIdFormatted() string {
	return fmt.Sprintf("<@%s>", s.RTM.GetInfo().User.ID)
}

func cleanupSlackEncodedCharacters(msg string) string {
	// replace slack tokens that are rendered to characters
	msg = strings.Replace(msg, "&lt;", "<", -1)
	msg = strings.Replace(msg, "&gt;", ">", -1)
	return msg
}

func (s *Slack) renderSlackUsernames(msg string) string {
	userIds := regexp.MustCompile("<@\\w+>").FindAllString(msg, -1)
	for _, slackFmtMsgUserId := range userIds {
		// in the Message we'll receive something like "<@U123123>", the id is actually "U123123"
		userId := strings.Replace(strings.Replace(slackFmtMsgUserId, "<@", "", 1), ">", "", 1)

		user, err := s.RTM.GetUserInfo(userId)
		if err != nil {
			name := user.Name
			if user.Profile.FirstName != "" {
				name = user.Profile.FirstName
			}
			msg = strings.Replace(msg, "<@"+user.ID+">", name, -1)
		}
	}
	return msg
}

var slackEmojiLookup map[string]string

func (s *Slack) renderSlackEmojis(msg string) string {
	var err error

	if slackEmojiLookup == nil {
		slackEmojiLookup, err = s.RTM.GetEmoji()
		if err != nil {
			log.Panicln("Could not get emojis from Slack", err)
			return msg
		}
	}

	emojis := regexp.MustCompile(":\\w+:").FindAllString(msg, -1)
	for _, slackFmtMsgEmoji := range emojis {
		// in the Message we'll receive something like ":smile:", this will actually return 😊
		emojiName := strings.Replace(strings.Replace(slackFmtMsgEmoji, ":", "", 1), ":", "", 1)

		if emojiName != "" {
			emojiImgUrl := slackEmojiLookup[emojiName]

			// follow the aliases for emojis
			for strings.Contains(emojiImgUrl, "alias:") {
				nextEmojiName := strings.Replace(emojiImgUrl, "alias:", "", -1)
				emojiImgUrl = slackEmojiLookup[nextEmojiName]
			}

			if emojiImgUrl == "" {
				emojiImgUrl = s.githubEmojiLookup[emojiName]
			}

			if emojiImgUrl == "" {
				continue
			}

			msg = strings.Replace(msg, ":"+emojiName+":", emojiImgUrl, -1)
		}
	}

	return msg
}

func (s *Slack) respondWithHelpMsg(channelId string) {
	msg := `DM me something and I'll try to display that on the board.

You can also supply options for the board by doing:
`

	msg += "```"
	msg += `
Your Message, 🚀, or img_url goes here.
---
align:        # 10 5           // set position of media; horizontally or vertically
align:        # center center  // (left,center,right)  (top,center,bottom)
inverted:     # (true/false) invert the text or image
bwThreshold:  # (0-256) set the threshold value for either "on" or "off"
fill:         # ("", true/false) leave blank for autofill, or select your own fill
`

	// we would like to add support for this in the future
	//kerning: 0	         // spacing between letters
	//font-size: 1           // ??

	msg += "```\n\n"
	s.RTM.SendMessage(s.RTM.NewOutgoingMessage(msg, channelId))
}

func (s *Slack) respondWithSettingsHelpMessage(channelId string) {
	t, _ := template.New("help").Parse("```" + `You can change settings of the bot by saying:
	@{{.Username}} <setting_name> <setting_val>

Available Settings:
---
help                   # show this help message
countdown enable       # enable the countdown clock
countdown disable      # disable the countdown clock
countdown YYYY-MM-DD   # set a new countdown date and enable it
` + "```")

	var buff bytes.Buffer
	_ = t.Execute(&buff, struct {
		Username string
	}{
		Username: s.RTM.GetInfo().User.Name,
	})

	s.RTM.SendMessage(s.RTM.NewOutgoingMessage(buff.String(), channelId))
}

func (s *Slack) respondWithSSHConnectionString(channelId string) {
	slackMessage := ""

	ngrokResp, err := http.Get("http://localhost:4040/api/tunnels")
	if err != nil || ngrokResp == nil {
		slackMessage = "local ngrok could not be reached"
		s.RTM.SendMessage(s.RTM.NewOutgoingMessage(slackMessage, channelId))
		return
	}

	defer ngrokResp.Body.Close()

	ngrokBody, err := ioutil.ReadAll(ngrokResp.Body)
	if err != nil {
		slackMessage = "could not read ngrok response"
		s.RTM.SendMessage(s.RTM.NewOutgoingMessage(slackMessage, channelId))
		return
	}

	ngrok := struct {
		Tunnels []struct {
			Name      string
			PublicUrl string `json:"public_url"`
		}
	}{}

	err = json.Unmarshal(ngrokBody, &ngrok)
	if err != nil {
		slackMessage = "could not parse ngrok response"
		s.RTM.SendMessage(s.RTM.NewOutgoingMessage(slackMessage, channelId))
		return
	}

	// yay! we can actually send the ssh connection string
	ngrokUrl, _ := url.Parse(ngrok.Tunnels[0].PublicUrl)
	if err != nil {
		slackMessage = "could not parse ngrok url"
		s.RTM.SendMessage(s.RTM.NewOutgoingMessage(slackMessage, channelId))
		return
	}

	slackMessage = fmt.Sprintf("`ssh -p %s pi@%s`\n", ngrokUrl.Port(), ngrokUrl.Hostname())
	slackMessage += fmt.Sprintf("`host: %s:%s`", ngrokUrl.Hostname(), ngrokUrl.Port())
	s.RTM.SendMessage(s.RTM.NewOutgoingMessage(slackMessage, channelId))
}
