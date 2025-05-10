package db8

import (
	"database/sql"
	"deifzar/num8/pkg/log8"
	"deifzar/num8/pkg/model8"

	"github.com/gofrs/uuid/v5"
	_ "github.com/lib/pq"
)

type Db8Endpoint8 struct {
	Db *sql.DB
}

func NewDb8Endpoint8(db *sql.DB) Db8Endpoint8Interface {
	return &Db8Endpoint8{Db: db}
}

func (m *Db8Endpoint8) GetAllEndpoints() ([]model8.Endpoint8, error) {
	query, err := m.Db.Query("SELECT id, endpoint, live, hostnameid FROM ONLY cptm8endpoint WHERE hostnameid IN (SELECT id FROM cptm8hostname WHERE enabled = true and domainid IN (SELECT id FROM cptm8domain WHERE enabled = true))")
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return []model8.Endpoint8{}, err
	}
	var endpoints []model8.Endpoint8
	if query != nil {
		for query.Next() {
			var (
				id         uuid.UUID
				endpoint   string
				live       bool
				hostnameid uuid.UUID
			)
			err := query.Scan(&id, &endpoint, &live, &hostnameid)
			if err != nil {
				log8.BaseLogger.Debug().Msg(err.Error())
				return nil, err
			}
			e := model8.Endpoint8{Id: id, Endpoint: endpoint, Live: live, Hostnameid: hostnameid}
			endpoints = append(endpoints, e)
		}
	}
	return endpoints, nil
}

func (m *Db8Endpoint8) GetAllHTTPEndpoints() ([]model8.Endpoint8, error) {
	query, err := m.Db.Query("SELECT id, endpoint, live, hostnameid FROM ONLY cptm8endpoint WHERE endpoint LIKE 'http%' AND hostnameid IN (SELECT id FROM cptm8hostname WHERE enabled = true and domainid IN (SELECT id FROM cptm8domain WHERE enabled = true))")
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return []model8.Endpoint8{}, err
	}
	var endpoints []model8.Endpoint8
	if query != nil {
		for query.Next() {
			var (
				id         uuid.UUID
				endpoint   string
				live       bool
				hostnameid uuid.UUID
			)
			err := query.Scan(&id, &endpoint, &live, &hostnameid)
			if err != nil {
				log8.BaseLogger.Debug().Msg(err.Error())
				return nil, err
			}
			e := model8.Endpoint8{Id: id, Endpoint: endpoint, Live: live, Hostnameid: hostnameid}
			endpoints = append(endpoints, e)
		}
	}
	return endpoints, nil
}

func (m *Db8Endpoint8) GetAllByDomainID(domainID uuid.UUID) ([]model8.Endpoint8, error) {
	query, err := m.Db.Query("SELECT id, endpoint, live, hostnameid FROM ONLY cptm8endpoint WHERE hostnameid IN (SELECT id FROM cptm8hostname WHERE domainid = $1 AND domainid IN (SELECT id FROM cptm8domain enabled = true))", domainID)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return []model8.Endpoint8{}, err
	}
	var endpoints []model8.Endpoint8
	if query != nil {
		for query.Next() {
			var (
				id         uuid.UUID
				endpoint   string
				live       bool
				hostnameid uuid.UUID
			)
			err := query.Scan(&id, &endpoint, &live, &hostnameid)
			if err != nil {
				log8.BaseLogger.Debug().Msg(err.Error())
				return nil, err
			}
			e := model8.Endpoint8{Id: id, Endpoint: endpoint, Live: live, Hostnameid: hostnameid}
			endpoints = append(endpoints, e)
		}
	}
	return endpoints, nil
}

func (m *Db8Endpoint8) GetAllHTTPByDomainID(domainID uuid.UUID) ([]model8.Endpoint8, error) {
	query, err := m.Db.Query("SELECT id, endpoint, live, hostnameid FROM ONLY cptm8endpoint WHERE endpoint LIKE 'http%' AND hostnameid IN (SELECT id FROM cptm8hostname WHERE domainid = $1 AND domainid IN (SELECT id FROM cptm8domain enabled = true))", domainID)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return []model8.Endpoint8{}, err
	}
	var endpoints []model8.Endpoint8
	if query != nil {
		for query.Next() {
			var (
				id         uuid.UUID
				endpoint   string
				live       bool
				hostnameid uuid.UUID
			)
			err := query.Scan(&id, &endpoint, &live, &hostnameid)
			if err != nil {
				log8.BaseLogger.Debug().Msg(err.Error())
				return nil, err
			}
			e := model8.Endpoint8{Id: id, Endpoint: endpoint, Live: live, Hostnameid: hostnameid}
			endpoints = append(endpoints, e)
		}
	}
	return endpoints, nil
}

func (m *Db8Endpoint8) GetAllByHostnameID(hostnameID uuid.UUID) ([]model8.Endpoint8, error) {
	query, err := m.Db.Query("SELECT id, endpoint, live, hostnameid FROM ONLY cptm8endpoint WHERE hostnameid = $1", hostnameID)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return []model8.Endpoint8{}, err
	}
	var endpoints []model8.Endpoint8
	if query != nil {
		for query.Next() {
			var (
				id         uuid.UUID
				endpoint   string
				live       bool
				hostnameid uuid.UUID
			)
			err := query.Scan(&id, &endpoint, &live, &hostnameid)
			if err != nil {
				log8.BaseLogger.Debug().Msg(err.Error())
				return nil, err
			}
			e := model8.Endpoint8{Id: id, Endpoint: endpoint, Live: live, Hostnameid: hostnameid}
			endpoints = append(endpoints, e)
		}
	}
	return endpoints, nil
}

func (m *Db8Endpoint8) GetAllHTTPByHostnameID(hostnameID uuid.UUID) ([]model8.Endpoint8, error) {
	query, err := m.Db.Query("SELECT id, endpoint, live, hostnameid FROM ONLY cptm8endpoint WHERE endpoint LIKE 'http%' AND hostnameid = $1", hostnameID)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return []model8.Endpoint8{}, err
	}
	var endpoints []model8.Endpoint8
	if query != nil {
		for query.Next() {
			var (
				id         uuid.UUID
				endpoint   string
				live       bool
				hostnameid uuid.UUID
			)
			err := query.Scan(&id, &endpoint, &live, &hostnameid)
			if err != nil {
				log8.BaseLogger.Debug().Msg(err.Error())
				return nil, err
			}
			e := model8.Endpoint8{Id: id, Endpoint: endpoint, Live: live, Hostnameid: hostnameid}
			endpoints = append(endpoints, e)
		}
	}
	return endpoints, nil
}

func (m *Db8Endpoint8) GetOneEndpointByID(endpointID uuid.UUID) (model8.Endpoint8, error) {
	query, err := m.Db.Query("SELECT id, endpoint, live, hostnameid FROM ONLY cptm8endpoint WHERE id = $1", endpointID)
	if err != nil {
		log8.BaseLogger.Debug().Msg(err.Error())
		return model8.Endpoint8{}, err
	}
	var e model8.Endpoint8
	if query != nil {
		for query.Next() {
			var (
				id         uuid.UUID
				endpoint   string
				live       bool
				hostnameid uuid.UUID
			)
			err := query.Scan(&id, &endpoint, &live, &hostnameid)
			if err != nil {
				log8.BaseLogger.Debug().Msg(err.Error())
				return model8.Endpoint8{}, err
			}
			e = model8.Endpoint8{Id: id, Endpoint: endpoint, Live: live, Hostnameid: hostnameid}
		}
	}
	return e, nil
}

func (m *Db8Endpoint8) GetEndpointIDByEndpoint(endpoint string) (uint, error) {
	row := m.Db.QueryRow("SELECT id FROM ONLY cptm8endpoint WHERE endpoint like $1", endpoint)
	var id uint
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return 0, nil
	case nil:
		return id, nil
	default:
		return 0, err
	}
}
