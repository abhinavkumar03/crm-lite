package processor

import (
	"context"

	rdto "github.com/abhinavkumar03/crm-lite/backend/internal/record/dto"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/engine"
	"github.com/abhinavkumar03/crm-lite/backend/internal/workflow/entity"
)

// guardedRecords wraps record mutations so workflow-caused CUD carries
// reentrancy metadata (source=workflow, depth, exclude same workflow).
type guardedRecords struct {
	inner      engine.RecordMutator
	workflowID string
	depth      int
}

func (g *guardedRecords) Get(ctx context.Context, orgID, moduleID, id, userID string, expand bool) (*rdto.RecordResponse, error) {
	return g.inner.Get(ctx, orgID, moduleID, id, userID, expand)
}

func (g *guardedRecords) Create(ctx context.Context, orgID, moduleID, userID string, req rdto.CreateRecordRequest) (*rdto.RecordResponse, error) {
	ctx = engine.WithMutationMeta(ctx, engine.MutationMeta{
		Source: entity.SourceWorkflow, Depth: g.depth + 1, ExcludeWorkflowID: g.workflowID,
	})
	return g.inner.Create(ctx, orgID, moduleID, userID, req)
}

func (g *guardedRecords) Update(ctx context.Context, orgID, moduleID, id, userID string, req rdto.UpdateRecordRequest) (*rdto.RecordResponse, error) {
	ctx = engine.WithMutationMeta(ctx, engine.MutationMeta{
		Source: entity.SourceWorkflow, Depth: g.depth + 1, ExcludeWorkflowID: g.workflowID,
	})
	return g.inner.Update(ctx, orgID, moduleID, id, userID, req)
}

func withGuardedRecords(deps engine.ActionDeps, workflowID string, depth int) engine.ActionDeps {
	if deps.Records == nil {
		return deps
	}
	out := deps
	out.Records = &guardedRecords{inner: deps.Records, workflowID: workflowID, depth: depth}
	return out
}
