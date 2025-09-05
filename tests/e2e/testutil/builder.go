package testutil

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const (
	defaultTestTimeout = 30 * time.Second
	e2eGroup           = "e2e"
	envVarParts        = 2
)

// FlowTestBuilder provides a fluent API for building flow execution tests.
type FlowTestBuilder struct {
	t             *testing.T
	flowFile      string
	configDir     string
	timeout       time.Duration
	workDir       string
	expectExit    *int
	expectError   string
	expectOut     string
	envVars       map[string]string
	command       string
	args          []string
	requiredPaths []string
}

// FlowTestResult contains the results of a flow test execution.
type FlowTestResult struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Error    error
	Duration time.Duration
}

// NewFlowTest creates a new flow test builder.
func NewFlowTest(t *testing.T) *FlowTestBuilder {
	t.Helper()

	return &FlowTestBuilder{
		t:             t,
		flowFile:      "",
		configDir:     "",
		workDir:       "",
		timeout:       defaultTestTimeout, // Default timeout
		expectExit:    nil,
		expectError:   "",
		expectOut:     "",
		envVars:       make(map[string]string),
		command:       "",
		args:          []string{},
		requiredPaths: []string{},
	}
}

// WithFlow sets the flow file to execute.
func (b *FlowTestBuilder) WithFlow(flowFile string) *FlowTestBuilder {
	b.flowFile = flowFile

	return b
}

// WithCommand sets the command and arguments to execute.
func (b *FlowTestBuilder) WithCommand(command string, args ...string) *FlowTestBuilder {
	b.command = command
	b.args = args

	// Record required file path for execute command so we can validate later when workDir is known
	if strings.ToLower(command) == "execute" && len(args) > 0 {
		b.requiredPaths = append(b.requiredPaths, args[0])
	}

	return b
}

// WithConfig sets the config directory for the test.
func (b *FlowTestBuilder) WithConfig(configDir string) *FlowTestBuilder {
	b.configDir = configDir

	return b
}

// WithTimeout sets the execution timeout.
func (b *FlowTestBuilder) WithTimeout(timeout time.Duration) *FlowTestBuilder {
	b.timeout = timeout

	return b
}

// WithWorkDir sets the working directory for the test.
func (b *FlowTestBuilder) WithWorkDir(workDir string) *FlowTestBuilder {
	b.workDir = workDir

	return b
}

// ExpectExitCode sets the expected exit code.
func (b *FlowTestBuilder) ExpectExitCode(code int) *FlowTestBuilder {
	b.expectExit = &code

	return b
}

// ExpectSuccess is a convenience method for expecting exit code 0.
func (b *FlowTestBuilder) ExpectSuccess() *FlowTestBuilder {
	return b.ExpectExitCode(0)
}

// ExpectFailure is a convenience method for expecting non-zero exit code.
func (b *FlowTestBuilder) ExpectFailure() *FlowTestBuilder {
	return b.ExpectExitCode(1)
}

// ExpectOutput sets the expected output substring.
func (b *FlowTestBuilder) ExpectOutput(output string) *FlowTestBuilder {
	b.expectOut = output

	return b
}

// ExpectError sets the expected error substring.
func (b *FlowTestBuilder) ExpectError(errorMsg string) *FlowTestBuilder {
	b.expectError = errorMsg

	return b
}

// WithEnv sets an environment variable for the test execution.
func (b *FlowTestBuilder) WithEnv(key, value string) *FlowTestBuilder {
	b.envVars[key] = value

	return b
}

// WithEnvFile loads environment variables from a file.
func (b *FlowTestBuilder) WithEnvFile(path string) *FlowTestBuilder {
	file, err := os.Open(path) // #nosec G304 - path is controlled in test environment
	if err != nil {
		b.t.Fatalf("Required env file %s is missing: %v", path, err)
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			b.t.Fatalf("Failed to close env file %s: %v", path, closeErr)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", envVarParts)
		if len(parts) == envVarParts {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			b.envVars[key] = value
		}
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		b.t.Fatalf("Failed to read env file %s: %v", path, scanErr)
	}

	return b
}

// Run executes the flow test and returns the result.
func (b *FlowTestBuilder) Run() *FlowTestResult {
	// Ensure binary exists before execution
	EnsureBinaryExists(b.t)

	// Track execution timing
	start := time.Now()

	// Note: Flow file is not required for list command testing
	// The current implementation only supports list command anyway

	// Set default work directory if not provided
	if b.workDir == "" {
		b.workDir = b.t.TempDir()
	}

	// Create flow runner
	runner := NewFlowRunner(b.t)

	// Configure runner
	if b.flowFile != "" {
		runner.SetFlowFile(b.flowFile)
	}

	// Set command and args if specified
	if b.command != "" {
		runner.SetCommand(b.command)
		runner.SetArgs(b.args)
	}

	runner.SetTimeout(b.timeout)
	runner.SetWorkDir(b.workDir)

	if b.configDir != "" {
		runner.SetConfigDir(b.configDir)
	}

	// Pass environment variables to runner
	runner.SetEnvVars(b.envVars)

	// Validate required paths (e.g., execute target flow) relative to workDir.
	// Enforce only for success-path tests; for failure-path tests, rely on CLI validation.
	enforce := b.expectExit != nil && *b.expectExit == 0
	for _, p := range b.requiredPaths {
		checkPath := p
		if !filepath.IsAbs(checkPath) {
			checkPath = filepath.Join(b.workDir, checkPath)
		}

		_, err := os.Stat(checkPath)
		if err != nil {
			if enforce {
				b.t.Fatalf("Required file for command not found: %s (%v)", checkPath, err)
			} else {
				b.t.Logf("Required file for command not found (expected in failure test): %s (%v)", checkPath, err)
			}
		}
	}

	// Execute the flow
	result := runner.Execute()

	// Calculate duration
	duration := time.Since(start)
	result.Duration = duration

	// Record test execution for coverage tracking
	RecordTestExecution(b.t, e2eGroup, "", duration, result.ExitCode == 0)

	// Validate expectations
	b.validateResult(result)

	return result
}

// validateResult checks the test result against expectations.
func (b *FlowTestBuilder) validateResult(result *FlowTestResult) {
	// Check exit code expectation
	if b.expectExit != nil {
		if result.ExitCode != *b.expectExit {
			b.t.Fatalf("Expected exit code %d, got %d", *b.expectExit, result.ExitCode)
		}
	}

	// Check output expectation
	if b.expectOut != "" {
		if !contains(result.Stdout, b.expectOut) {
			b.t.Fatalf("Expected output to contain %q, got: %s", b.expectOut, result.Stdout)
		}
	}

	// Check error expectation
	if b.expectError != "" {
		if !contains(result.Stderr, b.expectError) {
			b.t.Fatalf("Expected error to contain %q, got: %s", b.expectError, result.Stderr)
		}
	}
}

// contains checks if a string contains a substring (case-insensitive basic check).
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// TestDataPath returns the absolute path to a test data file.
func TestDataPath(relativePath string) string {
	root := detectProjectRoot()

	return filepath.Join(root, "tests", "e2e", "testdata", relativePath)
}

func detectProjectRoot() string {
	workingDir, err := os.Getwd()
	if err != nil {
		return "."
	}

	cur := workingDir
	for {
		_, statErr := os.Stat(filepath.Join(cur, "go.mod"))
		if statErr == nil {
			return cur
		}

		parent := filepath.Dir(cur)
		if parent == cur {
			return workingDir
		}

		cur = parent
	}
}

// FlowPath returns the absolute path to a flow file.
func FlowPath(flowName string) string {
	rel := TestDataPath(filepath.Join("flows", flowName))

	abs, err := filepath.Abs(rel)
	if err != nil {
		return rel
	}

	return abs
}
