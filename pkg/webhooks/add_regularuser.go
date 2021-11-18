package webhooks

import (
	"github.com/openshift/validating-webhook-framework/pkg/webhooks/regularuser"
)

func init() {
	Register(regularuser.WebhookName, func() Webhook { return regularuser.NewWebhook() })
}
