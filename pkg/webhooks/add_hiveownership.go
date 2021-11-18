package webhooks

import (
	"github.com/openshift/validating-webhook-framework/pkg/webhooks/hiveownership"
)

func init() {
	Register(hiveownership.WebhookName, func() Webhook { return hiveownership.NewWebhook() })
}
