package mocks

import (
	"context"
)

type HealthCheckRepositoryHealthyMock struct {
}

func (s *HealthCheckRepositoryHealthyMock) HealthCheckMasterDB(ctx context.Context) (isOk bool, err error) {
	return true, nil
}
func (s *HealthCheckRepositoryHealthyMock) HealthCheckSlaveDB(ctx context.Context) (isOk bool, err error) {
	return true, nil
}

type HealthCheckRepositoryUnhealthyMock struct {
}

func (s *HealthCheckRepositoryUnhealthyMock) HealthCheckMasterDB(ctx context.Context) (isOk bool, err error) {
	return false, nil
}
func (s *HealthCheckRepositoryUnhealthyMock) HealthCheckSlaveDB(ctx context.Context) (isOk bool, err error) {
	return false, nil
}

type HealthCheckRepositoryPanicMock struct {
}

func (s *HealthCheckRepositoryPanicMock) HealthCheckMasterDB(ctx context.Context) (isOk bool, err error) {
	panic(true)
}
func (s *HealthCheckRepositoryPanicMock) HealthCheckSlaveDB(ctx context.Context) (isOk bool, err error) {
	panic(true)
}
