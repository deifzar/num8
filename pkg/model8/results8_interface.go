package model8

import "github.com/projectdiscovery/nuclei/v3/pkg/output"

type Model8Results8Interface interface {
	AddResultEvent(output.ResultEvent) []output.ResultEvent
	SetResultEventFromOutputfilename() error
	SetOutputfilename(string)
	GetOutputfilename() string
	GetResultEvent() []output.ResultEvent
}
