package model8

import (
	"deifzar/num8/pkg/log8"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
)

type WriteOptions8 struct {
	WriteOptions []output.WriterOptions
}

func NewModel8WriterOptions8() Model8Writer8Interface {
	return &WriteOptions8{
		WriteOptions: nil,
	}
}

func (wo8 *WriteOptions8) GetWriterOptions8() []output.WriterOptions {
	return wo8.WriteOptions
}

func (wo8 *WriteOptions8) SetDefaultWriterOptions8() (string, error) {

	currenttime := time.Now()
	suffix := fmt.Sprintf("%d-%d-%d-%d-%d-%d", currenttime.Year(), currenttime.Month(), currenttime.Day(), currenttime.Hour(), currenttime.Minute(), currenttime.Second())
	// var writer io.WriteCloser
	file1, errs := os.CreateTemp("./tmp", "result-"+suffix)
	// writer, errs := os.Create("./tmp/temp-1111.txt")
	if errs != nil {
		log8.BaseLogger.Debug().Msg(errs.Error())
		return "", errs
	}
	// file2, errs := os.Create("./tmp/errorsink-" + suffix)
	// if errs != nil {
	// 	log8.BaseLogger.Debug().Msg(errs.Error())
	// 	return "", errs
	// }
	// file3, errs := os.Create("./tmp/tracesink-" + suffix)
	// if errs != nil {
	// 	log8.BaseLogger.Debug().Msg(errs.Error())
	// 	return "", errs
	// }
	writer1 := io.WriteCloser(file1)
	// writer2 := io.WriteCloser(file2)
	// writer3 := io.WriteCloser(file3)
	// writer := io.WriteCloser(os.Stdout)

	wo8.WriteOptions = append(wo8.WriteOptions, output.WithJson(true, true))
	wo8.WriteOptions = append(wo8.WriteOptions, output.WithTimestamp(true))
	wo8.WriteOptions = append(wo8.WriteOptions, output.WithNoMetadata(false))
	wo8.WriteOptions = append(wo8.WriteOptions, output.WithMatcherStatus(true))
	wo8.WriteOptions = append(wo8.WriteOptions, output.WithStoreResponse(true, "./tmp"))
	wo8.WriteOptions = append(wo8.WriteOptions, output.WithAurora(aurora.NewAurora(true)))
	wo8.WriteOptions = append(wo8.WriteOptions, output.WithWriter(writer1))
	// wo8.WriteOptions = append(wo8.WriteOptions, output.WithErrorSink(writer2))
	// wo8.WriteOptions = append(wo8.WriteOptions, output.WithTraceSink(writer3))
	filename := file1.Name()
	return filename, nil
}
