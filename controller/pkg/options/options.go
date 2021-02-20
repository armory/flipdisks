package options

import (
	"fmt"
	"regexp"
	"time"

	"github.com/armory/flipdisks/pkg/virtualboard"
	"gopkg.in/yaml.v2"
)

type FlipboardMessageOptions struct {
	// Display text or image/gif url by only setting Message
	// Otherwise, you can create your own virtual board and have only that rendered.
	Message      string `yaml:"message"`
	VirtualBoard *virtualboard.VirtualBoard

	DisplayTime      int    `yaml:"displayTime"` // in ms
	Append           bool   `yaml:"append"`
	Align            string `yaml:"align"`
	XAlign           string
	YAlign           string
	FontSize         int    `yaml:"font-size"`
	Kerning          int    `yaml:"kerning"`
	Inverted         bool   `yaml:"inverted"`
	BWThreshold      int    `yaml:"bwThreshold"`
	Fill             string `yaml:"fill"`
	SendPanelByPanel bool   `yaml:"sendPanelByPanel"`
}

func GetDefaultOptions() FlipboardMessageOptions {
	return FlipboardMessageOptions{
		DisplayTime:      int(5 * (time.Second / time.Millisecond)),	// stored in ms
		Inverted:         false,
		BWThreshold:      140, // magic
		Fill:             "",
		Align:            "center center",
		SendPanelByPanel: true,
	}
}

// SetDisplayTime is a helper given a time.Duration converting to the correct DisplayTime
func (s *FlipboardMessageOptions) SetDisplayTime(duration time.Duration) {
	s.DisplayTime = int(duration.Nanoseconds() / 1e6) // because ms is "Intentionally unimplemented" https://github.com/golang/go/issues/5491
}

func (s *FlipboardMessageOptions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// make a copy of the type
	type nestedOptsDefaults FlipboardMessageOptions

	// cast the defaults to the copy of the type, otherwise it'll be circular
	raw := nestedOptsDefaults(GetDefaultOptions())

	if err := unmarshal(&raw); err != nil {
		return err
	}

	// uncast the nested type so it's the correct type
	*s = FlipboardMessageOptions(raw)
	return nil
}

func GetAlignOptions(align string) (string, string) {
	alignmentOptions := regexp.MustCompile("( |,)+").Split(align, -1)
	var XAlign, YAlign string
	XAlign = alignmentOptions[0]
	if len(alignmentOptions) > 1 {
		YAlign = alignmentOptions[1]
	}

	return XAlign, YAlign
}
func SplitMessageAndOptions(rawMsg string) []FlipboardMessageOptions {
	var messages []FlipboardMessageOptions

	playlistRegex := regexp.MustCompile(`^\n*---\n`)
	isPlaylist := playlistRegex.Match([]byte(rawMsg))
	if isPlaylist == false {
		msgAndOptions := regexp.MustCompile(`\s*--(-+)\s*`).Split(rawMsg, -1)
		var m FlipboardMessageOptions

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

func unmarshleOptions(rawOptions string) FlipboardMessageOptions {
	// reset the options for each Message
	opts := GetDefaultOptions()

	yaml.Unmarshal([]byte(rawOptions), &opts)

	return opts
}
