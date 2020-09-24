package repository

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"

	"github.com/opentracing/opentracing-go"
)

// healthcheck operation constant
const (
	// Operation name
	HealthCheckSlaveDBOperation  = "Repository.HealthCheck.SlaveDB"
	HealthCheckMasterDBOperation = "Repository.HealthCheck.MasterDB"
)

type IHealthCheckRepository interface {
	HealthCheckMasterDB(ctx context.Context) (isOk bool, err error)
	HealthCheckSlaveDB(ctx context.Context) (isOk bool, err error)
}

// HealthCheckRepository type
type HealthCheckRepository struct {
	MasterDB *sqlx.DB
	SlaveDB  *sqlx.DB
}

// NewHealthCheckRepository creates new HealthCheckRepository.
func NewHealthCheckRepository(masterDB *sqlx.DB, slaveDB *sqlx.DB) (*HealthCheckRepository, error) {
	if masterDB == nil {
		return nil, errors.New("the master DB connection is nil")
	}

	if slaveDB == nil {
		return nil, errors.New("the slave DB connection is nil")
	}

	r := &HealthCheckRepository{}
	r.MasterDB = masterDB
	r.SlaveDB = slaveDB
	return r, nil
}

// HealthCheckMasterDB used for MasterDB healthcheck
func (r *HealthCheckRepository) HealthCheckMasterDB(ctx context.Context) (isOk bool, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, HealthCheckMasterDBOperation)
	defer span.Finish()

	err = r.MasterDB.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}

// HealthCheckSlaveDB used for SlaveDB healthcheck
func (r *HealthCheckRepository) HealthCheckSlaveDB(ctx context.Context) (isOk bool, err error) {
	span, _ := opentracing.StartSpanFromContext(ctx, HealthCheckSlaveDBOperation)
	defer span.Finish()

	err = r.SlaveDB.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}
