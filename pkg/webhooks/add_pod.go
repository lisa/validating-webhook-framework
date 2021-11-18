package webhooks

import (
	"github.com/lisa/validating-webhook-framework/pkg/webhooks/pod"
)

func init() {
	Register(pod.WebhookName, func() Webhook { return pod.NewWebhook() })
}
