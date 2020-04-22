// +build integration

package gitresolver_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx/pkg/gits"
	"github.com/jenkins-x/jx/pkg/jenkinsfile/gitresolver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testGitURL = "https://github.com/jenkins-x/jxr-packs-kubernetes.git"
	testTag    = "v0.0.5"
)

func TestInitBuildPackIntegration(t *testing.T) {
	type TestCase struct {
		Name       string
		Callback   func(*testing.T, gits.Gitter, string) error
		ShouldFail bool
	}

	testCases := []TestCase{
		{
			Name:     "masterTwice",
			Callback: testMasterBranchTwice,
		},
		{
			Name:     "tagTwice",
			Callback: testTagTwice,
		},
		{
			Name:     "masterThenTagTwice",
			Callback: testMasterThenTagThenMasterThenTag,
		},
	}

	gitter := gits.NewGitCLI()

	tmpDir, err := ioutil.TempDir("", "test-init-buildpack-")
	require.NoError(t, err, "could not make temp dir")
	t.Logf("runnings tests in dir %s", tmpDir)

	for _, tc := range testCases {
		homeDir := filepath.Join(tmpDir, tc.Name)

		oldEnvHome := os.Getenv("JX_HOME")
		err = os.Setenv("JX_HOME", homeDir)
		assert.NoError(t, err)

		t.Logf("running test %s", tc.Name)
		err = tc.Callback(t, gitter, tc.Name)
		if tc.ShouldFail {
			assert.Error(t, err, "expected error for test %s", tc.Name)
		} else {
			assert.NoError(t, err, "should not have failed for test %s", tc.Name)
		}

		defer func() {
			os.Setenv("JX_HOME", oldEnvHome)
			if os.Getenv("JX_REMOVE_TMP") != "false" {
				err := os.RemoveAll(homeDir)
				if err != nil {
					t.Logf("Error cleaning up tmpdir because %v", err)
				}
			}
		}()
	}
}

func testMasterBranchTwice(t *testing.T, gitter gits.Gitter, name string) error {
	dir, _ := assertInitBuildPack(t, gitter, "master", name)
	dir2, _ := assertInitBuildPack(t, gitter, "master", name)
	assert.Equal(t, dir, dir2, "should have returned the same dir on both runs for %s", name)
	return nil
}

func testTagTwice(t *testing.T, gitter gits.Gitter, name string) error {
	dir, _ := assertInitBuildPack(t, gitter, testTag, name)
	dir2, _ := assertInitBuildPack(t, gitter, testTag, name)
	assert.Equal(t, dir, dir2, "should have returned the same dir on both runs for %s", name)
	return nil
}

func testMasterThenTagThenMasterThenTag(t *testing.T, gitter gits.Gitter, name string) error {
	dir, _ := assertInitBuildPack(t, gitter, "master", name)
	dir2, _ := assertInitBuildPack(t, gitter, testTag, name)
	assert.Equal(t, dir, dir2, "should have returned the same dir on both runs for %s", name)
	dir2, _ = assertInitBuildPack(t, gitter, "master", name)
	assert.Equal(t, dir, dir2, "should have returned the same dir on both runs for %s", name)
	dir2, _ = assertInitBuildPack(t, gitter, testTag, name)
	assert.Equal(t, dir, dir2, "should have returned the same dir on both runs for %s", name)
	return nil
}

func assertInitBuildPack(t *testing.T, gitter gits.Gitter, ref string, name string) (string, error) {
	t.Logf("trying InitBuildPack with URL %s and ref %s for %s", testGitURL, ref, name)
	dir, err := gitresolver.InitBuildPack(gitter, testGitURL, ref)
	require.NoError(t, err)
	require.NotEmpty(t, dir, "should have found dir for run on empty packs dir for %s", name)

	// assert the file exists
	fileName := filepath.Join(dir, "go", "pipeline.yaml")
	assert.FileExists(t, fileName, "with git ref %s for test %s", ref, name)
	t.Logf("file exists for ref %s name %s for test %s", ref, fileName, name)
	return dir, err
}
