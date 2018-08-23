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

func (s *FlipBoardDisplayOptions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type optsDefaults FlipBoardDisplayOptions

	// todo, make this so that we don't have 2 defaults..
	raw := optsDefaults{
		DisplayTime: 5000,
		Inverted:    false,
		BWThreshold: 140, // magic
		Fill:        "",
		Align:       "center center",
	}

	if err := unmarshal(&raw); err != nil {
		return err
	}

	raw.XAlign, raw.YAlign = GetAlignOptions(raw.Align)

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
