package config

var (
	// ProtectedNamespaces comes from
	// raw.githubusercontent.com/openshift/managed-cluster-config/master/deploy/osd-managed-resources/addons-namespaces.ConfigMap.yaml
	ProtectedNamespaces = []string{
		"acm",
		"addon-dba-operator",
		"codeready-workspaces-operator",
		"codeready-workspaces-operator-qe",
		"openshift-logging",
		"openshift-storage",
		"prow",
		"redhat-addon-operator",
		"redhat-gpu-operator",
		"redhat-kas-fleetshard-operator",
		"redhat-kas-fleetshard-operator-qe",
		"redhat-managed-kafka-operator",
		"redhat-managed-kafka-operator-qe",
		"redhat-ocm-addon-test-operator",
		"redhat-ods-operator",
		"redhat-reference-addon",
		"redhat-rhmi-operator",
		"redhat-rhoam-operator",
		"redhat-rhoami-operator",
	}
)
