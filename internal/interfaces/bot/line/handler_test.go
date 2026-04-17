package linebot

import (
	"context"
	"errors"
	"testing"
	"time"

	sdklinebot "github.com/line/line-bot-sdk-go/v8/linebot"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
)

type stubLineProcessor struct {
	release chan struct{}
	result  chan error
}

func (s *stubLineProcessor) ProcessTextMessage(ctx context.Context, _ *sdklinebot.Event, _ *sdklinebot.TextMessage) error {
	<-s.release
	s.result <- ctx.Err()
	return nil
}

func TestLineWebhookBackgroundContextCanOutliveRequest(t *testing.T) {
	t.Parallel()

	log, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	processor := &stubLineProcessor{
		release: make(chan struct{}),
		result:  make(chan error, 1),
	}

	handler := &LineBotHandler{
		messageProcessor: processor,
		logger:           log,
	}

	reqCtx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		defer close(done)

		reqCtx = logger.WithLogger(reqCtx, log)
		processCtx, processCancel := context.WithTimeout(logger.DetachContext(reqCtx, log), asyncProcessingTimeout)
		defer processCancel()

		reqLogger := logger.FromContext(processCtx, log).With(
			logger.String("reply_token", "reply-token"),
		)
		eventCtx := logger.WithLogger(processCtx, reqLogger)

		if err := handler.messageProcessor.ProcessTextMessage(
			eventCtx,
			&sdklinebot.Event{ReplyToken: "reply-token"},
			&sdklinebot.TextMessage{Text: "/d 2330"},
		); err != nil && !errors.Is(err, context.Canceled) {
			t.Errorf("ProcessTextMessage() unexpected error = %v", err)
		}
	}()

	cancel()
	close(processor.release)

	select {
	case err := <-processor.result:
		if err != nil {
			t.Fatalf("ProcessTextMessage received canceled context: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for background processor")
	}

	<-done
}
