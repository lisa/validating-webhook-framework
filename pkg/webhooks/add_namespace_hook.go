package webhooks

import (
	"github.com/openshift/validating-webhook-framework/pkg/webhooks/namespace"
)

func init() {
	Register(namespace.WebhookName, func() Webhook { return namespace.NewWebhook() })
}
