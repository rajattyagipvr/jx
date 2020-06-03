// +build unit

package config_test

import (
	"path"
	"testing"

	"github.com/jenkins-x/jx/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadDeployNotExists(t *testing.T) {
	dir := path.Join("test_data", "jx-apps-phase-bad")
	require.DirExists(t, dir)

	cfg, fileName, err := config.LoadDeployConfig(dir)
	require.NoError(t, err, "for dir %s", dir)
	require.Nil(t, cfg, "should not have a DeployConfig returned for dir %s", dir)
	assert.Empty(t, fileName, "file name for dir %s", dir)
}

func TestLoadDeployExists(t *testing.T) {
	dir := path.Join("test_data", "deploy-config-remote-env")
	require.DirExists(t, dir)

	cfg, fileName, err := config.LoadDeployConfig(dir)
	require.NoError(t, err, "for dir %s", dir)
	require.NotNil(t, cfg, "no DeployConfig returned for dir %s", dir)
	assert.NotEmpty(t, fileName, "file name for dir %s", dir)
	assert.Equal(t, "namespaces/apps", cfg.Spec.KptPath, "spec.kptPath for dir %s", dir)
	assert.Equal(t, "myapps", cfg.Spec.Namespace, "spec.namespace for dir %s", dir)

	t.Logf("loaded file %s with deploy config %#v", fileName, cfg)
}
