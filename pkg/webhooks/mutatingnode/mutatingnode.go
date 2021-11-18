package mutatingnode

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/lisa/validating-webhook-framework/pkg/webhooks/utils"
	admissionv1 "k8s.io/api/admission/v1"
	admissionregv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	admissionctl "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	WebhookName string = "node-mutating"
	docString   string = `Managed OpenShift Customers may use tolerations on Pods that could cause those Pods to be scheduled on infra or master nodes.`

	masterLabel string = "node-role.kubernetes.io/master"
	infraLabel  string = "node-role.kubernetes.io"
	workerLabel string = "node-role.kubernetes.io/worker"
)

var (
	adminGroups = []string{"dedicated-admins"}
	log         = logf.Log.WithName(WebhookName)

	scope = admissionregv1.AllScopes
	rules = []admissionregv1.RuleWithOperations{
		{
			Operations: []admissionregv1.OperationType{admissionregv1.OperationAll},
			Rule: admissionregv1.Rule{
				APIGroups:   []string{""},
				APIVersions: []string{"*"},
				Resources:   []string{"nodes"},
				Scope:       &scope,
			},
		},
	}
)

// ObjectSelector implements Webhook interface
func (s *MutatingNodeWebhook) ObjectSelector() *metav1.LabelSelector { return nil }

func (s *MutatingNodeWebhook) Doc() string {
	return fmt.Sprintf(docString)
}

// TimeoutSeconds implements Webhook interface
func (s *MutatingNodeWebhook) TimeoutSeconds() int32 { return 1 }

// MatchPolicy implements Webhook interface
func (s *MutatingNodeWebhook) MatchPolicy() admissionregv1.MatchPolicyType {
	return admissionregv1.Equivalent
}

// Name implements Webhook interface
func (s *MutatingNodeWebhook) Name() string { return WebhookName }

// FailurePolicy implements Webhook interface and defines how unrecognized errors and timeout errors from the admission webhook are handled. Allowed values are Ignore or Fail.
// Ignore means that an error calling the webhook is ignored and the API request is allowed to continue.
// It's important to leave the FailurePolicy set to Ignore because otherwise the pod will fail to be created as the API request will be rejected.
func (s *MutatingNodeWebhook) FailurePolicy() admissionregv1.FailurePolicyType {
	return admissionregv1.Ignore
}

// Rules implements Webhook interface
func (s *MutatingNodeWebhook) Rules() []admissionregv1.RuleWithOperations { return rules }

// GetURI implements Webhook interface
func (s *MutatingNodeWebhook) GetURI() string { return "/" + WebhookName }

// SideEffects implements Webhook interface
func (s *MutatingNodeWebhook) SideEffects() admissionregv1.SideEffectClass {
	return admissionregv1.SideEffectClassNone
}

// SyncSetLabelSelector returns the label selector to use in the SyncSet.
func (s *MutatingNodeWebhook) SyncSetLabelSelector() metav1.LabelSelector {
	return utils.DefaultLabelSelector()
}

// Validate implements Webhook interface
func (s *MutatingNodeWebhook) Validate(req admissionctl.Request) bool {
	// Check if incoming request is a node request
	// Retrieve old and new node objects
	node := &corev1.Node{}
	oldNode := &corev1.Node{}

	err := json.Unmarshal(req.Object.Raw, node)
	if err != nil {
		errMsg := "Failed to Unmarshal node object"
		log.Error(err, errMsg)
		return false
	}
	err = json.Unmarshal(req.OldObject.Raw, oldNode)
	if err != nil {
		errMsg := "Failed to Unmarshal old node object"
		log.Error(err, errMsg)
		return false
	}
	return true

}

func (s *MutatingNodeWebhook) Authorized(request admissionctl.Request) admissionctl.Response {
	return s.authorized(request)
}

func (s *MutatingNodeWebhook) authorized(request admissionctl.Request) admissionctl.Response {
	var ret admissionctl.Response

	if request.AdmissionRequest.UserInfo.Username == "system:unauthenticated" {
		// This could highlight a significant problem with RBAC since an
		// unauthenticated user should have no permissions.
		log.Info("system:unauthenticated made a webhook request. Check RBAC rules", "request", request.AdmissionRequest)
		ret = admissionctl.Denied("Unauthenticated")
		ret.UID = request.AdmissionRequest.UID
		return ret
	}

	// Check that the current user is a dedicated admin
	for _, userGroup := range request.UserInfo.Groups {

		// Retrieve old and new node objects
		node := &corev1.Node{}
		oldNode := &corev1.Node{}

		err := json.Unmarshal(request.Object.Raw, node)
		if err != nil {
			errMsg := "Failed to Unmarshal node object"
			log.Error(err, errMsg)
			ret = admissionctl.Denied("UnauthorizedAction")
			return ret
		}
		err = json.Unmarshal(request.OldObject.Raw, oldNode)
		if err != nil {
			errMsg := "Failed to Unmarshal old node object"
			log.Error(err, errMsg)
			ret = admissionctl.Denied("UnauthorizedAction")
			return ret
		}

		log.Info("Request log", "oldNode.Labels", oldNode.Labels, "node.Labels", node.Labels, "username", request.UserInfo.Username)

		if utils.SliceContains(userGroup, adminGroups) {

			// Fail on none worker nodes
			if _, ok := oldNode.Labels[workerLabel]; ok {

				// Fail on infra,worker nodes
				if val, ok := oldNode.Labels[infraLabel]; ok && val == "infra" {
					log.Info("cannot edit non-worker node")
					ret.UID = request.AdmissionRequest.UID
					ret = admissionctl.Denied("UnauthorizedAction")
					return ret
				}

				// Do not allow worker node type to change to master
				if _, ok := node.Labels[masterLabel]; ok {
					log.Info("cannot change worker node to master")
					ret.UID = request.AdmissionRequest.UID
					ret = admissionctl.Denied("UnauthorizedAction")
					return ret
				}

				// Do not allow worker node type to change to infra
				if _, ok := node.Labels[infraLabel]; ok {
					log.Info("cannot change worker node to infra")
					ret.UID = request.AdmissionRequest.UID
					ret = admissionctl.Denied("UnauthorizedAction")
					return ret
				}

				// Fail on removed worker node label
				if _, ok := node.Labels[workerLabel]; !ok {
					log.Info("cannot remove worker node label from worker node")
					ret.UID = request.AdmissionRequest.UID
					ret = admissionctl.Denied("UnauthorizedAction")
					return ret
				}
			} else {
				log.Info("cannot edit non-worker nodes")
				ret.UID = request.AdmissionRequest.UID
				ret = admissionctl.Denied("UnauthorizedAction")
				return ret
			}
		}
	}
	// Allow Operation
	msg := "New label does not infringe on node properties"
	log.Info(msg)
	ret = admissionctl.Patched(msg)
	ret.UID = request.AdmissionRequest.UID
	return ret

}

type MutatingNodeWebhook struct {
	mu sync.Mutex
	s  runtime.Scheme
}

func NewWebhook() *MutatingNodeWebhook {
	scheme := runtime.NewScheme()
	admissionv1.AddToScheme(scheme)
	corev1.AddToScheme(scheme)
	return &MutatingNodeWebhook{
		s: *scheme,
	}
}
