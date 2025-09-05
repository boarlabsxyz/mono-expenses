package testutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	// File permissions for directories and files.
	dirPermissions  = 0o750
	filePermissions = 0o600
)

// createStandardFlows returns common flow definitions used across tests.
func createStandardFlows() map[string]string {
	return map[string]string{
		"single-step.json": `{
			"id": "single-step",
			"name": "Single Step Flow",
			"initialStep": "step1",
			"steps": {
				"step1": {
					"type": "prompt",
					"prompt": "Hello World"
				}
			}
		}`,
		"multi-step.json": `{
			"id": "multi-step",
			"name": "Multi Step Flow",
			"initialStep": "step1",
			"steps": {
				"step1": {
					"type": "prompt",
					"prompt": "Step 1",
					"nextStep": "step2"
				},
				"step2": {
					"type": "prompt",
					"prompt": "Step 2"
				}
			}
		}`,
		"with-conditions.json": `{
			"id": "with-conditions",
			"name": "Conditional Flow",
			"initialStep": "condition1",
			"steps": {
				"condition1": {
					"type": "condition",
					"condition": "true",
					"yes": "step1",
					"no": "step2"
				},
				"step1": {
					"type": "prompt",
					"prompt": "True branch"
				},
				"step2": {
					"type": "prompt",
					"prompt": "False branch"
				}
			}
		}`,
	}
}

// GetStandardFlows returns a copy of the standard flows map.
func GetStandardFlows() map[string]string {
	return createStandardFlows()
}

// CreateFlowsDirectory creates a standard .flows/flows directory structure
// and returns the path to the flows directory.
func CreateFlowsDirectory(t *testing.T, tempDir string) string {
	t.Helper()

	flowsDir := filepath.Join(tempDir, ".flows", "flows")
	require.NoError(t, os.MkdirAll(flowsDir, dirPermissions), "Should create flows directory")

	return flowsDir
}

// WriteFlowFiles writes a map of flow files to the specified directory.
func WriteFlowFiles(t *testing.T, flowsDir string, flows map[string]string) {
	t.Helper()

	for filename, content := range flows {
		filePath := filepath.Join(flowsDir, filename)
		require.NoError(t, os.WriteFile(filePath, []byte(content), filePermissions), "Should write flow file: %s", filename)
	}
}

// WriteStandardFlows writes all standard flows to the specified directory.
func WriteStandardFlows(t *testing.T, flowsDir string) {
	t.Helper()
	WriteFlowFiles(t, flowsDir, createStandardFlows())
}

// SetupTestWithFlows creates a temporary directory with standard flows
// and returns the temp directory and flows directory paths.
func SetupTestWithFlows(t *testing.T) (string, string) {
	t.Helper()

	tempDir := t.TempDir()
	flowsDir := CreateFlowsDirectory(t, tempDir)
	WriteStandardFlows(t, flowsDir)

	return tempDir, flowsDir
}

// SetupTestWithCustomFlows creates a temporary directory with custom flows
// and returns the temp directory and flows directory paths.
func SetupTestWithCustomFlows(t *testing.T, flows map[string]string) (string, string) {
	t.Helper()

	tempDir := t.TempDir()
	flowsDir := CreateFlowsDirectory(t, tempDir)
	WriteFlowFiles(t, flowsDir, flows)

	return tempDir, flowsDir
}

// SetupEmptyTest creates a temporary directory with just the .flows structure
// but no flow files, and returns the temp directory and flows directory paths.
func SetupEmptyTest(t *testing.T) (string, string) {
	t.Helper()

	tempDir := t.TempDir()
	flowsDir := CreateFlowsDirectory(t, tempDir)

	return tempDir, flowsDir
}

// CreateSingleFlow creates a simple flow with the given ID and name.
func CreateSingleFlow(id, name, prompt string) string {
	return `{
		"id": "` + id + `",
		"name": "` + name + `",
		"initialStep": "step1",
		"steps": {
			"step1": {
				"type": "prompt",
				"prompt": "` + prompt + `"
			}
		}
	}`
}
