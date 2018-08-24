package options

import (
	"regexp"
)

type FlipBoardDisplayOptions struct {
	Message     string `yaml:"message"`
	DisplayTime int    `yaml:"displayTime"`
	Append      bool   `yaml:"append"`
	Align       string `yaml:"align"`
	XAlign      string
	YAlign      string
	FontSize    int    `yaml:"font-size"`
	Kerning     int    `yaml:"kerning"`
	Inverted    bool   `yaml:"inverted"`
	BWThreshold int    `yaml:"bwThreshold"`
	Fill        string `yml:"fill"`
}

func GetDefaultOptions() FlipBoardDisplayOptions {
	return FlipBoardDisplayOptions{
		DisplayTime: 10000,
		Inverted:    false,
		BWThreshold: 140, // magic
		Fill:        "",
		Align:       "center center",
	}
}

func (s *FlipBoardDisplayOptions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// make a copy of the type
	type nestedOptsDefaults FlipBoardDisplayOptions

	// cast the defaults to the copy of the type, otherwise it'll be circular
	raw := nestedOptsDefaults(GetDefaultOptions())

	if err := unmarshal(&raw); err != nil {
		return err
	}

	// uncast the nested type so it's the correct type
	*s = FlipBoardDisplayOptions(raw)
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
