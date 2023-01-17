package maleodiscord

import (
	"context"

	"github.com/tigorlazuardi/maleo/bucket"
)

type Hook interface {
	PreMessageHook(ctx context.Context, web *WebhookContext) context.Context
	PostMessageHook(ctx context.Context, web *WebhookContext, err error)
	PreBucketUploadHook(ctx context.Context, web *WebhookContext) context.Context
	PostBucketUploadHook(ctx context.Context, web *WebhookContext, results []bucket.UploadResult)
}

var _ Hook = (*NoopHook)(nil)

type NoopHook struct{}

func (n NoopHook) PreMessageHook(ctx context.Context, _ *WebhookContext) context.Context {
	return ctx
}
func (n NoopHook) PostMessageHook(context.Context, *WebhookContext, error) {}
func (n NoopHook) PreBucketUploadHook(ctx context.Context, _ *WebhookContext) context.Context {
	return ctx
}
func (n NoopHook) PostBucketUploadHook(context.Context, *WebhookContext, []bucket.UploadResult) {}
