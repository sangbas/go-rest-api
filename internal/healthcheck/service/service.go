package service

import (
	"context"
	"fmt"
	"github.com/go-rest-api/internal/healthcheck/repository"
	"sync"

	"github.com/opentracing/opentracing-go"

	logger "github.com/sirupsen/logrus"
)

// Operation Constants
const (
	// Operation name
	APIHealthCheckOperation            = "service.healthcheck.api"
	InfrastructureHealthCheckOperation = "service.healthcheck.infrastructure"
)

// Dependency Constants
const (
	// Dependency types
	hardDependencyType = "hard"
	softDependencyType = "soft"
)

// Health Message Constants
const (
	// Health message
	HealthyMsg  = "It's healthy as hell."
	CoughingMsg = "It's getting cough. Please check the soft dependency."
	DyingMsg    = "It's dying. Please check the hard dependency."
)

type IHealthCheckService interface {
	API(ctx context.Context)
	Infrastructure(ctx context.Context) InfrastructureHealthCheckResponse
}

type healthItem struct {
	Name           string `json:"name"`
	IsHealthy      bool   `json:"is_healthy"`
	DependencyType string `json:"dependency_type"`
	Remarks        string `json:"remarks"`
}

// InfrastructureHealthCheckResponse type
type InfrastructureHealthCheckResponse struct {
	Items  []healthItem `json:"items"`
	Result string       `json:"result"`
	IsOk   bool         `json:"-"`
	Mutex  *sync.Mutex  `json:"-"`
}

// HealthCheckService type
type HealthCheckService struct {
	repo repository.IHealthCheckRepository
}

// NewHealthCheckService used for initiate HealthCheckService
func NewHealthCheckService(repo repository.IHealthCheckRepository) (*HealthCheckService, error) {
	return &HealthCheckService{
		repo: repo,
	}, nil
}

// API used for provide API service health check
func (s *HealthCheckService) API(ctx context.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx, APIHealthCheckOperation)
	defer span.Finish()
}

// Infrastructure used for provide Infrastructure service healthcheck
func (s *HealthCheckService) Infrastructure(ctx context.Context) InfrastructureHealthCheckResponse {
	span, ctx := opentracing.StartSpanFromContext(ctx, InfrastructureHealthCheckOperation)
	defer span.Finish()

	var (
		wg           sync.WaitGroup
		healthResult InfrastructureHealthCheckResponse
	)

	healthResult.Mutex = &sync.Mutex{}
	s.getMasterDBStatus(ctx, &healthResult, &wg)
	s.getSlaveDBStatus(ctx, &healthResult, &wg)
	wg.Wait()
	healthResult.examineHealth()
	return healthResult
}

// getMasterDBStatus get master database status
func (s *HealthCheckService) getMasterDBStatus(ctx context.Context, healthResponse *InfrastructureHealthCheckResponse, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logger.Error(fmt.Errorf("panic: %v", err))
				wg.Done()
			}
		}()

		item := healthItem{
			Name:           "Master Database SQL",
			DependencyType: hardDependencyType,
			IsHealthy:      true,
			Remarks:        "",
		}

		isOk, err := s.repo.HealthCheckMasterDB(ctx)
		if !isOk {
			item.IsHealthy = false
			item.Remarks = fmt.Sprint(err)
		}

		healthResponse.addItem(item)
		wg.Done()
	}()
}

// getSlaveDBStatus get slave database status
func (s *HealthCheckService) getSlaveDBStatus(ctx context.Context, healthResponse *InfrastructureHealthCheckResponse, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logger.Error(fmt.Errorf("panic: %v", err))
				wg.Done()
			}
		}()

		item := healthItem{
			Name:           "Slave Database SQL",
			DependencyType: hardDependencyType,
			IsHealthy:      true,
			Remarks:        "",
		}

		isOk, err := s.repo.HealthCheckSlaveDB(ctx)
		if !isOk {
			item.IsHealthy = false
			item.Remarks = fmt.Sprint(err)
		}

		healthResponse.addItem(item)
		wg.Done()
	}()
}

// AddItem add healths item
func (ahr *InfrastructureHealthCheckResponse) addItem(item healthItem) {
	ahr.Mutex.Lock()
	ahr.Items = append(ahr.Items, item)
	ahr.Mutex.Unlock()
}

// examineHealth examines the health based on item
func (ahr *InfrastructureHealthCheckResponse) examineHealth() {
	// Set default is healthy
	ahr.Result = HealthyMsg
	ahr.IsOk = true

	var unhealthyItemIsExist bool
	for _, v := range ahr.Items {
		if !v.IsHealthy && v.DependencyType == hardDependencyType {
			ahr.IsOk = false
			ahr.Result = DyingMsg
			unhealthyItemIsExist = true
		} else if !v.IsHealthy && v.DependencyType == softDependencyType {
			unhealthyItemIsExist = true
		}
	}

	if !ahr.IsOk {
		ahr.Result = DyingMsg
	} else {
		if unhealthyItemIsExist {
			ahr.Result = CoughingMsg
		}
	}
}
