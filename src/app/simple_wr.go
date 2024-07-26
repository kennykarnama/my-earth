package app

import (
	"context"
	"log/slog"
	"time"
)

type SimpleWeatherRefresher struct {
	interval time.Duration
	locSvc   *LocationSvc
	ctx      context.Context
}

func NewSimpleWeatherRefresher(ctx context.Context, locSvc *LocationSvc, interval time.Duration) *SimpleWeatherRefresher {
	return &SimpleWeatherRefresher{
		interval: interval,
		locSvc:   locSvc,
		ctx:      ctx,
	}
}

func (s *SimpleWeatherRefresher) Watch() {
	t := time.NewTicker(s.interval)

	for {
		select {
		case <-t.C:
			slog.Info("refresh weather")
			_, err := s.locSvc.RefreshWeather(s.ctx)
			if err != nil {
				slog.Error("refresh weather", slog.Any("err", err))
			}
		case <-s.ctx.Done():
			slog.Info("refresher stopped")
			t.Stop()

			return
		}
	}
}
