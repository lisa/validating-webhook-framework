package webhooks

import (
	"github.com/lisa/validating-webhook-framework/pkg/webhooks/clusterlogging"
)

func init() {
	Register(clusterlogging.WebhookName, func() Webhook { return clusterlogging.NewWebhook() })
}
