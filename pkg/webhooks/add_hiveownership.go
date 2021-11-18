package webhooks

import (
	"github.com/lisa/validating-webhook-framework/pkg/webhooks/hiveownership"
)

func init() {
	Register(hiveownership.WebhookName, func() Webhook { return hiveownership.NewWebhook() })
}
