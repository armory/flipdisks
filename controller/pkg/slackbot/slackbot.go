package slackbot

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/armory/flipdisks/controller/pkg/github"
	"github.com/armory/flipdisks/controller/pkg/options"
	"github.com/nlopes/slack"
	"gopkg.in/yaml.v2"
)

type Slack struct {
	countdownDate     string
	githubEmojiLookup github.EmojiLookup
}

func NewSlack(c string, g github.EmojiLookup) *Slack {
	return &Slack{
		countdownDate:     c,
		githubEmojiLookup: g,
	}
}

func (s *Slack) StartSlackListener(slackToken string, flipboardMsgChn chan options.FlipBoardDisplayOptions) {
	api := slack.New(slackToken)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)
	api.SetDebug(false)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	var oldStopper chan struct{}

	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			// close the hold handleSlackMsg
			if oldStopper != nil {
				close(oldStopper)
			}

			stopper := make(chan struct{})
			go handleSlackMsg(ev, rtm, flipboardMsgChn, stopper, s.countdownDate, s.githubEmojiLookup)
			oldStopper = stopper

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		case *slack.ConnectionErrorEvent:
			fmt.Println("Connection Error!")
			fmt.Printf("%#v\n", ev.ErrorObj.Error())

		default:
			fmt.Println("Event Received: ")
			fmt.Printf("%#v\n", msg)
			fmt.Printf("%#v\n", msg.Data)
		}
	}
}

func handleSlackMsg(slackEvent *slack.MessageEvent, rtm *slack.RTM, flipboardMsgChn chan options.FlipBoardDisplayOptions, stopper chan struct{}, countdownDate string, githubEmojiLookup github.EmojiLookup) {
	rawMsg := slackEvent.Msg.Text
	if slackEvent.SubMessage != nil {
		rawMsg = slackEvent.SubMessage.Text
	}

	fmt.Printf("Raw Slack Message: %+v\n", rawMsg)

	messages := splitMessageAndOptions(rawMsg)

	fmt.Printf("%#v \n", messages)

	// do some message cleanup because of slack formatting
	for _, msg := range messages {
		if strings.ToLower(msg.Message) == "help" {
			respondWithHelpMsg(rtm, slackEvent.Msg.Channel)
			return
		}

		msg.Message = renderSlackUsernames(msg.Message, rtm)
		msg.Message = cleanupSlackEncodedCharacters(msg.Message)
		msg.Message = renderSlackEmojis(msg.Message, rtm, githubEmojiLookup)
	}

	for {
		select {
		case <-stopper:
			return // we've received a new message, let's stop looping
		default:
			for _, msg := range messages {
				fmt.Printf("Rendering Message: %+v\n", msg.Message)

				flipboardMsgChn <- msg
				fmt.Println("sleeping", msg.DisplayTime)
				time.Sleep(time.Millisecond * time.Duration(msg.DisplayTime))
				fmt.Println("end sleeping")
			}
			if countdownDate != "" {
				messages = []options.FlipBoardDisplayOptions{countdown(countdownDate)}
			}
		}
	}
}

func cleanupSlackEncodedCharacters(msg string) string {
	// replace slack tokens that are rendered to characters
	msg = strings.Replace(msg, "&lt;", "<", -1)
	msg = strings.Replace(msg, "&gt;", ">", -1)
	return msg
}

func renderSlackUsernames(msg string, rtm *slack.RTM) string {
	userIds := regexp.MustCompile("<@\\w+>").FindAllString(msg, -1)
	for _, slackFmtMsgUserId := range userIds {
		// in the Message we'll receive something like "<@U123123>", the id is actually "U123123"
		userId := strings.Replace(strings.Replace(slackFmtMsgUserId, "<@", "", 1), ">", "", 1)

		user, err := rtm.GetUserInfo(userId)
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

func renderSlackEmojis(msg string, rtm *slack.RTM, githubEmojiLookup github.EmojiLookup) string {
	var err error

	if slackEmojiLookup == nil {
		slackEmojiLookup, err = rtm.GetEmoji()
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

func respondWithHelpMsg(rtm *slack.RTM, channelId string) {
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
	rtm.SendMessage(rtm.NewOutgoingMessage(msg, channelId))
}

func splitMessageAndOptions(rawMsg string) []options.FlipBoardDisplayOptions {
	var messages []options.FlipBoardDisplayOptions

	playlistRegex := regexp.MustCompile(`^---\n`)
	isPlaylist := playlistRegex.Match([]byte(rawMsg))
	if isPlaylist == false {
		msgAndOptions := regexp.MustCompile(`\s*--(-+)\s*`).Split(rawMsg, -1)
		var m options.FlipBoardDisplayOptions

		rawOptions := ""
		if len(msgAndOptions) > 1 { // is there options?
			rawOptions = msgAndOptions[1]
		}

		m = unmarshleOptions(rawOptions)
		m.Message = msgAndOptions[0]

		messages = append(messages, m)
	} else {
		rawPlaylist := playlistRegex.Split(rawMsg, -1)[1]

		err := yaml.Unmarshal([]byte(rawPlaylist), &messages)

		if err != nil {
			fmt.Println("Could not unmarshal the yaml")
			fmt.Println(err)
		}
	}

	return messages
}

func countdown(countdownDate string) options.FlipBoardDisplayOptions {
	horizonEventTime, err := time.Parse("2006-01-02", countdownDate)
	if err != nil {
		fmt.Println(err)
	}
	t := time.Now()
	elapsed := horizonEventTime.Sub(t)
	days := int(elapsed.Hours() / 24)
	hours := int(elapsed.Hours()) % 24
	mins := int(elapsed.Minutes()) % 60
	secs := int(elapsed.Seconds()) % 60
	msg := options.FlipBoardDisplayOptions{
		Message:     fmt.Sprintf("HORIZON EVENT\n%d:%02d:%02d:%02d", days, hours, mins, secs),
		DisplayTime: 1000,
		Align:       "center center",
	}
	return msg
}

func unmarshleOptions(rawOptions string) options.FlipBoardDisplayOptions {
	// reset the options for each Message
	opts := options.GetDefaultOptions()

	yaml.Unmarshal([]byte(rawOptions), &opts)

	return opts
}
