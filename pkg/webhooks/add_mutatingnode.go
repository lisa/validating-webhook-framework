package webhooks

import (
	"github.com/lisa/validating-webhook-framework/pkg/webhooks/mutatingnode"
)

func init() {
	Register(mutatingnode.WebhookName, func() Webhook { return mutatingnode.NewWebhook() })
}
