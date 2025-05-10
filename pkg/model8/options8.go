package model8

import (
	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
)

type Options8 struct {
	Options []nuclei.NucleiSDKOptions
}

func NewModel8Options8() Model8Options8Interface {
	return &Options8{
		Options: []nuclei.NucleiSDKOptions{},
	}
}

func (o *Options8) AddOption(option nuclei.NucleiSDKOptions) []nuclei.NucleiSDKOptions {
	o.Options = append(o.Options, option)
	return o.Options
}

func (o *Options8) GetOptions() []nuclei.NucleiSDKOptions {
	return o.Options
}
