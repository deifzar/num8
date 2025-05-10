package model8

import (
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
)

type Model8Scan8Interface interface {
	AddTarget(string) []string
	AddResultEvent(output.ResultEvent) []output.ResultEvent
	SetResultEventFromOutputfilename() error
	ValidRemoteDirectory(string) bool
	SetTemplateSources8(bool)
	SetTemplateFilters8ForNewTemplates() (bool, string)
	SetOutputfilename(string)
}
