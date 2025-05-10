package model8

import (
	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
)

type Model8Options8Interface interface {
	AddOption(nuclei.NucleiSDKOptions) []nuclei.NucleiSDKOptions
	GetOptions() []nuclei.NucleiSDKOptions
}
