package watchers

import (
	"github.com/dkotik/oakacs/v1"
	"go.uber.org/zap"
)

// ZapLogger converts ACS events into log records.
func ZapLogger(logger *zap.Logger) Watcher {
	return Watcher(func() (chan (oakacs.Event), error) {
		if logger == nil {
			var err error
			logger, err = zap.NewProduction()
			if err != nil {
				return nil, err
			}
		}
		events := make(chan (oakacs.Event), 64) // a small buffer
		// go func() {
		// 	defer logger.Sync()
		// 	for event := range events {
		// 		switch event.Type {
		// 		default:
		// 			logger.Info(event.Type.String(),
		// 				zap.Time("ts", event.When),
		// 				zap.String("session", event.Session.String()),
		// 			)
		// 		}
		// 	}
		// }()
		return events, nil
	})
}

// ZapSugarLogger converts ACS events into log records with nice terminal colors.
func ZapSugarLogger() Watcher {
	return Watcher(func() (chan (oakacs.Event), error) {
		// var err error
		// dev, err := zap.NewDevelopment()
		// if err != nil {
		// 	return nil, err
		// }
		// logger := dev.Sugar()

		events := make(chan (oakacs.Event), 64) // a small buffer
		// go func() {
		// 	defer logger.Sync()
		// 	for event := range events {
		// 		switch event.Type {
		// 		default:
		// 			logger.Info(event.Type.String(),
		// 				zap.Time("ts", event.When),
		// 				zap.String("session", event.Session.String()),
		// 			)
		// 		}
		// 	}
		// }()
		return events, nil
	})
}
