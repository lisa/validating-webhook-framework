package webhooks

import (
	"github.com/openshift/validating-webhook-framework/pkg/webhooks/mutatingnode"
)

func init() {
	Register(mutatingnode.WebhookName, func() Webhook { return mutatingnode.NewWebhook() })
}
