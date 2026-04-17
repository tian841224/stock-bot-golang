package tgbot

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tian841224/stock-bot/internal/infrastructure/config"
	logger "github.com/tian841224/stock-bot/internal/infrastructure/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type stubTelegramProcessor struct {
	release chan struct{}
	result  chan error
}

func (s *stubTelegramProcessor) ProcessUpdate(ctx context.Context, _ *tgbotapi.Update) error {
	<-s.release
	s.result <- ctx.Err()
	return nil
}

func TestWebhook_DetachesRequestContextForBackgroundProcessing(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	log, err := logger.NewLogger()
	if err != nil {
		t.Fatalf("NewLogger() error = %v", err)
	}

	processor := &stubTelegramProcessor{
		release: make(chan struct{}),
		result:  make(chan error, 1),
	}

	handler := NewTgHandler(&config.Config{}, processor, log)

	body := bytes.NewBufferString(`{"message":{"chat":{"id":123},"text":"/k 2330"}}`)
	reqCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	req := httptest.NewRequest(http.MethodPost, "/telegram", body).WithContext(reqCtx)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(recorder)
	c.Request = req

	handler.Webhook(c)

	if recorder.Code != http.StatusOK {
		t.Fatalf("Webhook() status = %d, want %d", recorder.Code, http.StatusOK)
	}

	cancel()
	close(processor.release)

	select {
	case err := <-processor.result:
		if err != nil {
			t.Fatalf("ProcessUpdate received canceled context: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for background processor")
	}
}
