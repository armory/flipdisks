package slackbot

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/armory/flipdisks/controller/pkg/flipboard"
	"github.com/armory/flipdisks/controller/pkg/github"
	"github.com/armory/flipdisks/controller/pkg/options"
	"github.com/nlopes/slack"
)

type Slack struct {
	token             string
	countdownDate     string
	githubEmojiLookup github.EmojiLookup
	RTM               *slack.RTM
}

func NewSlack(token string, c string, g github.EmojiLookup) *Slack {
	api := slack.New(token)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	return &Slack{
		token:             token,
		countdownDate:     c,
		githubEmojiLookup: g,
		RTM:               rtm,
	}
}

func (s *Slack) StartSlackListener(board *flipboard.Flipboard, flipboardMsgChn chan options.FlipboardMessageOptions) {
	var oldStopper chan struct{}

	for msg := range s.RTM.IncomingEvents {
		switch event := msg.Data.(type) {
		case *slack.MessageEvent:
			// close the hold handleSlackMsg
			if oldStopper != nil {
				close(oldStopper)
			}

			stopper := make(chan struct{})
			go s.handleSlackMsg(event, board, flipboardMsgChn, stopper, s.countdownDate, s.githubEmojiLookup)
			oldStopper = stopper

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

func (s *Slack) handleSlackMsg(slackEvent *slack.MessageEvent, board *flipboard.Flipboard, flipboardMsgChn chan options.FlipboardMessageOptions, stopper chan struct{}, countdownDate string, githubEmojiLookup github.EmojiLookup) {
	rawMsg := slackEvent.Msg.Text
	if slackEvent.SubMessage != nil {
		rawMsg = slackEvent.SubMessage.Text
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
		msg.Message = s.renderSlackEmojis(msg.Message, githubEmojiLookup)

		board.Enqueue(&msg)
	}
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

func (s *Slack) renderSlackEmojis(msg string, githubEmojiLookup github.EmojiLookup) string {
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
				emojiImgUrl = githubEmojiLookup[emojiName]
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

