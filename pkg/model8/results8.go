package model8

import (
	"bufio"
	"deifzar/num8/pkg/log8"
	"encoding/json"
	"os"

	"github.com/projectdiscovery/nuclei/v3/pkg/output"
)

type Results8 struct {
	Outputfilename string
	ResultEvent    []output.ResultEvent
}

func NewModel8Results8() Model8Results8Interface {
	return &Results8{
		Outputfilename: "",
		ResultEvent:    nil,
	}
}

func (r *Results8) AddResultEvent(re output.ResultEvent) []output.ResultEvent {
	r.ResultEvent = append(r.ResultEvent, re)
	return r.ResultEvent
}

func (r *Results8) SetResultEventFromOutputfilename() error {
	readFile, err := os.Open(r.Outputfilename)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return err
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var re output.ResultEvent
	for fileScanner.Scan() {
		err = json.Unmarshal(fileScanner.Bytes(), &re)
		if err != nil {
			log8.BaseLogger.Debug().Msg(err.Error())
			return err
		}
		r.AddResultEvent(re)
	}
	return nil
}

func (r *Results8) SetOutputfilename(f string) {
	r.Outputfilename = f
}

func (r *Results8) GetResultEvent() []output.ResultEvent {
	return r.ResultEvent
}

func (r *Results8) GetOutputfilename() string {
	return r.Outputfilename
}
