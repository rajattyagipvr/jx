package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/jenkins-x/jx/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var (
	// DeployConfigPath is the name of the deployment configuration file.
	// Note that we don't use a file ending with .yml or .yaml to avoid being validated by
	// tools like Config Sync as being a valid CRD in the kubernetes cluster
	DeployConfigPath = filepath.Join(".jenkins-x", "Deployfile")
)

// +exported

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeployConfig defines how applications should be deployed into a remote environment if the remote environment git
// repository does not use a `jx-apps.yml` top level folder
// +k8s:openapi-gen=true
type DeployConfig struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Spec holds the desired state of the Action from the client
	// +optional
	Spec DeployConfigSec `json:"spec"`
}

// DeployConfigSec
type DeployConfigSec struct {
	// KptPath if using kpt to deploy applications into a GitOps repository specify the folder to deploy into.
	// For example if the root directory contains a Config Sync git layout we may want applications to be deployed into the
	// `namespaces/myapps` folder. If the `myconfig` folder is used as the root of the Config Sync configuration you may want
	// to configure something like `myconfig/namespaces/mysystem` or whatever.
	KptPath string `json:"kptPath,omitempty"`

	// Namespace specifies the namespace to deploy applications if using kpt. If specified this value will be used instead
	// of the Environment.Spec.Namespace in the Environment CRD
	Namespace string `json:"namespace,omitempty"`
}

// LoadDeployConfig loads the deploy configuration if present or returns nil if its not
func LoadDeployConfig(projectDir string) (*DeployConfig, string, error) {
	fileName := DeployConfigPath
	if projectDir != "" {
		fileName = filepath.Join(projectDir, fileName)
	}
	return LoadDeployConfigFile(fileName)
}

// LoadDeployConfigFile loads a specific deploy YAML configuration file or returns nil if it does not exist
func LoadDeployConfigFile(fileName string) (*DeployConfig, string, error) {
	exists, err := util.FileExists(fileName)
	if err != nil || !exists {
		return nil, "", err
	}
	config := DeployConfig{}
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return &config, "", fmt.Errorf("Failed to load file %s due to %s", fileName, err)
	}
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return &config, "", fmt.Errorf("Failed to unmarshal YAML file %s due to %s", fileName, err)
	}
	return &config, fileName, nil
}
