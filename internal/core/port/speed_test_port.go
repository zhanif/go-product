package port

import "time"

type SpeedTestRepository interface {
	Save(method, path string, duration time.Duration) error
}

type SpeedTestService interface {
	WriteLog(method, path string, duration time.Duration) error
}
