package cloud

import (
	"sort"
	"strings"
)

const (
	AKS        = "aks"
	ALIBABA    = "alibaba"
	AWS        = "aws"
	EKS        = "eks"
	GKE        = "gke"
	ICP        = "icp"
	IKS        = "iks"
	KUBERNETES = "kubernetes"
	KIND       = "kind"
	OKE        = "oke"
	OPENSHIFT  = "openshift"
	PKS        = "pks"
)

// KubernetesProviders list of all available Kubernetes providers
var KubernetesProviders = []string{GKE, OKE, AKS, AWS, EKS, KIND, KUBERNETES, IKS, OPENSHIFT, JX_INFRA, PKS, ICP, ALIBABA}

// KubernetesProviderOptions returns all the Kubernetes providers as a string
func KubernetesProviderOptions() string {
	values := []string{}
	values = append(values, KubernetesProviders...)
	sort.Strings(values)
	return strings.Join(values, ", ")
}
