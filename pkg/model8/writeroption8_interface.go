package model8

import "github.com/projectdiscovery/nuclei/v3/pkg/output"

type Model8Writer8Interface interface {
	GetWriterOptions8() []output.WriterOptions
	SetDefaultWriterOptions8() (string, error)
}
