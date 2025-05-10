package db8

import "deifzar/num8/pkg/model8"

type Db8Historyissue8Interface interface {
	GetAllHistoryIssues() ([]model8.Historyissue8, error)
	GetAllHistoryIssuesByStatus(model8.Status) ([]model8.Historyissue8, error)
	InsertBatch([]model8.Historyissue8) error
}
