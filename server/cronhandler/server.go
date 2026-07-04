package cronhandler

import (
	"context"

	"github.com/robfig/cron/v3"
)

type loanHandlerProvider interface {
	// RunDelinquentCheck runs the delinquent loan check process.
	RunDelinquentCheck(ctx context.Context) error
}

type CronHandler struct {
	loanHandler loanHandlerProvider
	cron        *cron.Cron
}

func NewCronHandler(cron *cron.Cron, loanHandler loanHandlerProvider) *CronHandler {
	return &CronHandler{
		cron:        cron,
		loanHandler: loanHandler,
	}
}

func (h *CronHandler) ImplementTasks() error {
	_, err := h.cron.AddFunc("0 0 0 * * *", func() { // run every day at midnight
		h.loanHandler.RunDelinquentCheck(context.Background())
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *CronHandler) Start() {
	h.cron.Start()
}

func (h *CronHandler) Stop() {
	h.cron.Stop()
}
