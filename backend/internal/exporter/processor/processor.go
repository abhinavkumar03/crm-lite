// Package processor runs persisted export jobs in the worker. The generation
// logic lives in the export service (shared with synchronous downloads), so this
// is a thin adapter that satisfies jobs.ExportProcessor.
package processor

import (
	"context"

	"go.uber.org/zap"

	"github.com/abhinavkumar03/crm-lite/backend/internal/exporter/service"
)

type Processor struct {
	svc    *service.Service
	logger *zap.Logger
}

func New(svc *service.Service, logger *zap.Logger) *Processor {
	return &Processor{svc: svc, logger: logger}
}

func (p *Processor) Process(ctx context.Context, orgID, id string) error {
	if err := p.svc.RunJob(ctx, orgID, id); err != nil {
		p.logger.Error("export: run job", zap.Error(err), zap.String("id", id))
		return err
	}
	p.logger.Info("export: job processed", zap.String("id", id))
	return nil
}
