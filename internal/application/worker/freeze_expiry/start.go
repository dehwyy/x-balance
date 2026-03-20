package freezeexpiry

import (
	"context"

	tlog "github.com/dehwyy/tracerfx/pkg/tracer/log"
)

func (w *Worker) Start(ctx context.Context) {
	pubsub := w.client.Subscribe(ctx, "__keyevent@0__:expired")
	defer func() {
		if err := pubsub.Close(); err != nil {
			tlog.FromContext(ctx).Error("failed to close pubsub", "err", err)
		}
	}()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			w.handleExpiry(ctx, msg.Payload)
		}
	}
}
