package service

import (
	"go-product/internal/core/port"
	"time"
)

type SpeedTestServiceImpl struct {
	repo port.SpeedTestRepository
}

func NewSpeedTestService(repo port.SpeedTestRepository) port.SpeedTestService {
	return &SpeedTestServiceImpl{repo: repo}
}

func (s *SpeedTestServiceImpl) WriteLog(method, path string, duration time.Duration) error {
	return s.repo.Save(method, path, duration)
}
