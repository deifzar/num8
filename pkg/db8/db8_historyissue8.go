package db8

import (
	"database/sql"
	"deifzar/num8/pkg/log8"
	"deifzar/num8/pkg/model8"
	"time"

	_ "github.com/lib/pq"

	"github.com/gofrs/uuid/v5"
)

type Db8Historyissue8 struct {
	Db *sql.DB
}

func NewDb8Historyissue8(db *sql.DB) Db8Historyissue8Interface {
	return &Db8Historyissue8{Db: db}
}

func (m *Db8Historyissue8) GetAllHistoryIssuesByStatus(status model8.Status) ([]model8.Historyissue8, error) {
	query, err := m.Db.Query("SELECT id, endpointid, issue, url, signature, status, firsttimefound FROM cptm8historyissue WHERE status = $1", status)
	if err != nil {
		log8.BaseLogger.Debug().Stack().Err(err).Msg("")
		return []model8.Historyissue8{}, err
	}
	var historyissues []model8.Historyissue8
	if query != nil {
		for query.Next() {
			var (
				id             uuid.UUID
				endpointid     uuid.UUID
				issue          model8.Issue8
				url            string
				signature      string
				status         model8.Status
				foundFirsttime time.Time
			)
			err := query.Scan(&id, &endpointid, &issue, &url, &signature, &status, &foundFirsttime)
			if err != nil {
				log8.BaseLogger.Debug().Msg(err.Error())
				return nil, err
			}
			hi := model8.Historyissue8{Id: id, Endpointid: endpointid, Issue: issue, Url: url, Signature: signature, Status: status, FoundFirsttime: foundFirsttime}
			historyissues = append(historyissues, hi)
		}
	}
	return historyissues, nil
}

func (m *Db8Historyissue8) GetAllHistoryIssues() ([]model8.Historyissue8, error) {
	query, err := m.Db.Query("SELECT issue, url, signature, id, endpointid, status, firsttimefound FROM cptm8historyissue")
	if err != nil {
		log8.BaseLogger.Debug().Stack().Err(err).Msg("")
		return []model8.Historyissue8{}, err
	}
	var historyissues []model8.Historyissue8
	if query != nil {
		for query.Next() {
			var (
				id             uuid.UUID
				endpointid     uuid.UUID
				issue          model8.Issue8
				url            string
				signature      string
				status         model8.Status
				foundFirsttime time.Time
			)
			err := query.Scan(&id, &endpointid, &issue, &url, &signature, &status, &foundFirsttime)
			if err != nil {
				log8.BaseLogger.Debug().Msg(err.Error())
				return nil, err
			}
			hi := model8.Historyissue8{Id: id, Endpointid: endpointid, Issue: issue, Url: url, Signature: signature, Status: status, FoundFirsttime: foundFirsttime}
			historyissues = append(historyissues, hi)
		}
	}
	return historyissues, nil
}

func (m *Db8Historyissue8) InsertBatch(historyIssues8 []model8.Historyissue8) error {
	tx, err := m.Db.Begin()
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO cptm8historyissue (endpointid, url, signature, issue, status, firsttimefound) VALUES ($1, $2, $3, $4, $5, NOW()) ON CONFLICT (signature) DO UPDATE SET status = EXCLUDED.status`)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return err
	}
	defer stmt.Close()
	var err2 error
	for _, h := range historyIssues8 {
		_, err2 = stmt.Exec(h.Endpointid, h.Url, h.Signature, h.Issue, h.Status)
		if err2 != nil {
			_ = tx.Rollback()
			log8.BaseLogger.Debug().Msg(err2.Error())
			return err2
		}
	}
	err2 = tx.Commit()
	if err2 != nil {
		_ = tx.Rollback()
		log8.BaseLogger.Debug().Msg(err2.Error())
		return err
	}
	return nil
}
