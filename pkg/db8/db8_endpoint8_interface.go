package db8

import (
	"deifzar/num8/pkg/model8"

	"github.com/gofrs/uuid/v5"
)

type Db8Endpoint8Interface interface {
	GetAllEndpoints() ([]model8.Endpoint8, error)
	GetAllHTTPEndpoints() ([]model8.Endpoint8, error)
	GetAllByDomainID(uuid.UUID) ([]model8.Endpoint8, error)
	GetAllHTTPByDomainID(uuid.UUID) ([]model8.Endpoint8, error)
	GetAllByHostnameID(uuid.UUID) ([]model8.Endpoint8, error)
	GetAllHTTPByHostnameID(uuid.UUID) ([]model8.Endpoint8, error)
	GetOneEndpointByID(uuid.UUID) (model8.Endpoint8, error)
}
