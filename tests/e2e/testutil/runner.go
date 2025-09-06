package testutil

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Static errors for better error handling.
var (
	ErrInvalidBinaryPath   = errors.New("invalid binary path")
	ErrBinaryNotFound      = errors.New("binary not found")
	ErrBinaryIsDirectory   = errors.New("binary path is a directory")
	ErrBinaryNotExecutable = errors.New("binary is not executable")
)

const (
	defaultRunnerTimeout = 30 * time.Second
	coverageDirMode      = 0o750
)

// FlowRunner handles subprocess execution of mono-expenses binary.
type FlowRunner struct {
	t           *testing.T
	flowFile    string
	configDir   string
	workDir     string
	timeout     time.Duration
	binaryPath  string
	coverageDir string
	stdout      bytes.Buffer
	stderr      bytes.Buffer
	envVars     map[string]string
	command     string
	args        []string
}

// findProjectRoot searches for the project root by looking for go.mod.
func findProjectRoot(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	projectRoot := wd
	for {
		_, err := os.Stat(filepath.Join(projectRoot, "go.mod"))
		if err == nil {
			return projectRoot
		}

		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			t.Fatalf("Reached filesystem root without finding go.mod")
		}

		projectRoot = parent
	}
}

// NewFlowRunner creates a new flow runner.
func NewFlowRunner(t *testing.T) *FlowRunner {
	t.Helper()

	projectRoot := findProjectRoot(t)
	binaryPath := filepath.Join(projectRoot, "bin", "mono-expenses-e2e")

	return &FlowRunner{
		t:           t,
		flowFile:    "",
		configDir:   "",
		workDir:     "",
		timeout:     defaultRunnerTimeout,
		binaryPath:  binaryPath, // Use absolute path to coverage-instrumented binary
		coverageDir: "",
		stdout:      bytes.Buffer{},
		stderr:      bytes.Buffer{},
		envVars:     make(map[string]string),
		command:     "",
		args:        []string{},
	}
}

// SetFlowFile sets the flow file to execute.
func (r *FlowRunner) SetFlowFile(flowFile string) {
	r.flowFile = flowFile
}

// SetConfigDir sets the config directory.
func (r *FlowRunner) SetConfigDir(configDir string) {
	r.configDir = configDir
}

// SetWorkDir sets the working directory.
func (r *FlowRunner) SetWorkDir(workDir string) {
	r.workDir = workDir
}

// SetTimeout sets the execution timeout.
func (r *FlowRunner) SetTimeout(timeout time.Duration) {
	r.timeout = timeout
}

// SetBinaryPath sets the path to the binary to execute.
func (r *FlowRunner) SetBinaryPath(binaryPath string) {
	r.binaryPath = binaryPath
}

// SetEnvVars sets environment variables for the test execution.
func (r *FlowRunner) SetEnvVars(envVars map[string]string) {
	r.envVars = envVars
}

// SetCommand sets the command to execute.
func (r *FlowRunner) SetCommand(command string) {
	r.command = command
}

// SetArgs sets the arguments for the command.
func (r *FlowRunner) SetArgs(args []string) {
	r.args = args
}

// Execute runs the flow and returns the result.
func (r *FlowRunner) Execute() *FlowTestResult {
	start := time.Now()

	// Reset buffers for this execution
	r.stdout.Reset()
	r.stderr.Reset()

	// Setup coverage collection
	r.setupCoverage()

	// Create context with timeout for command execution
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	// Create and configure command with all setup in one place
	cmd, err := r.createAndConfigureCommand(ctx)
	if err != nil {
		return &FlowTestResult{
			ExitCode: 1,
			Stdout:   "",
			Stderr:   fmt.Sprintf("Command setup failed: %v", err),
			Error:    err,
			Duration: time.Since(start),
		}
	}

	return r.executeCommand(cmd, start)
}

// CleanupCoverage removes coverage files for a test (optional cleanup).
func (r *FlowRunner) CleanupCoverage() {
	if r.coverageDir != "" {
		err := os.RemoveAll(r.coverageDir)
		if err != nil {
			r.t.Logf("Warning: Failed to cleanup coverage directory %s: %v", r.coverageDir, err)
		}
	}
}

// determineExitCode extracts exit code from command error.
func (r *FlowRunner) determineExitCode(err error) int {
	exitCode := 0

	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			exitCode = exitError.ExitCode()
		} else {
			// Other error (e.g., binary not found, timeout)
			exitCode = -1
		}
	}

	return exitCode
}

// setupCoverage creates a unique coverage directory for this test.
func (r *FlowRunner) setupCoverage() {
	if r.t == nil {
		return
	}

	projectRoot := findProjectRoot(r.t)

	// Create unique coverage directory for this test using absolute path
	testName := r.t.Name()
	coverageBase := filepath.Join(projectRoot, "coverage", "e2e")
	r.coverageDir = filepath.Join(coverageBase, testName)

	// Convert to absolute path to ensure it's unambiguous
	absCoverageDir, err := filepath.Abs(r.coverageDir)
	if err != nil {
		r.t.Fatalf("Failed to get absolute path for coverage directory %s: %v", r.coverageDir, err)
	}

	r.coverageDir = absCoverageDir

	// Create coverage directory with parent directories
	err = os.MkdirAll(r.coverageDir, coverageDirMode)
	if err != nil {
		r.t.Fatalf("Failed to create coverage directory %s: %v", r.coverageDir, err)
	}

	r.t.Logf("Created coverage directory: %s", r.coverageDir)
}

// EnsureBinaryExists checks if the test binary exists and builds it if needed.
func EnsureBinaryExists(t *testing.T) {
	t.Helper()

	projectRoot := findProjectRoot(t)
	binaryPath := filepath.Join(projectRoot, "bin", "mono-expenses-e2e")

	// Check if binary exists
	_, err := os.Stat(binaryPath)
	if os.IsNotExist(err) {
		t.Logf("Binary %s not found, building...", binaryPath)

		// Build the binary with coverage - use project root we already found
		cmd := exec.CommandContext(context.Background(), "make", "build-e2e-coverage")
		cmd.Dir = projectRoot

		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Failed to build binary: %v\nOutput: %s", err, output)
		}

		t.Logf("Binary built successfully")
	}
}

// validateBinaryPath ensures the binary path is safe to execute.
func validateBinaryPath(path string) error {
	// Convert to absolute path to prevent path traversal
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidBinaryPath, err)
	}

	// Check if file exists and is executable
	info, err := os.Stat(absPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrBinaryNotFound, err)
	}

	if info.IsDir() {
		return fmt.Errorf("%w: %s", ErrBinaryIsDirectory, absPath)
	}

	// Check file permissions (executable)
	if info.Mode()&0o111 == 0 {
		return fmt.Errorf("%w: %s", ErrBinaryNotExecutable, absPath)
	}

	return nil
}

// sanitizeArgs removes potentially dangerous arguments.
func sanitizeArgs(args []string) []string {
	var sanitized []string

	for _, arg := range args {
		// Remove null bytes and other dangerous characters
		cleaned := strings.ReplaceAll(arg, "\x00", "")
		cleaned = strings.TrimSpace(cleaned)

		// Skip empty arguments
		if cleaned != "" {
			sanitized = append(sanitized, cleaned)
		}
	}

	return sanitized
}

// createAndConfigureCommand creates a fully configured command ready for execution.
func (r *FlowRunner) createAndConfigureCommand(ctx context.Context) (*exec.Cmd, error) {
	// Validate binary path before use
	err := validateBinaryPath(r.binaryPath)
	if err != nil {
		return nil, fmt.Errorf("invalid binary path: %w", err)
	}

	// Build arguments
	args := r.buildCommandArgs()

	// Sanitize arguments to prevent command injection
	sanitizedArgs := sanitizeArgs(args)

	// Create command with context
	cmd := exec.CommandContext(ctx, r.binaryPath, sanitizedArgs...)

	// Configure command properties
	r.configureCommand(cmd)

	return cmd, nil
}

// buildCommandArgs builds the command arguments array.
func (r *FlowRunner) buildCommandArgs() []string {
	args := []string{}

	if r.configDir != "" {
		args = append(args, "--config", r.configDir)
	}

	// Use explicit command if provided
	if r.command != "" {
		args = append(args, r.command)
		args = append(args, r.args...)
	} else {
		// Fallback to legacy behavior for backward compatibility
		if r.flowFile != "" {
			// Use execute command when flow file is provided
			args = append(args, "execute", r.flowFile)
		} else {
			// Use list command when no flow file is provided
			args = append(args, "list")
		}
	}

	return args
}

// configureCommand sets up all command properties.
func (r *FlowRunner) configureCommand(cmd *exec.Cmd) {
	// Configure all command properties in one place
	cmd.Stdout = &r.stdout
	cmd.Stderr = &r.stderr

	// Setup working directory if specified
	if r.workDir != "" {
		cmd.Dir = r.workDir
	}

	// Setup environment
	r.setupEnvironment(cmd)
}

// setupEnvironment configures environment variables for the command.
func (r *FlowRunner) setupEnvironment(cmd *exec.Cmd) {
	// Setup environment variables
	if cmd.Env == nil {
		cmd.Env = os.Environ()
	}

	// Add custom environment variables from the builder
	for key, value := range r.envVars {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	// Setup coverage environment
	if r.coverageDir != "" {
		cmd.Env = append(cmd.Env, "GOCOVERDIR="+r.coverageDir)
		r.t.Logf("Set GOCOVERDIR=%s", r.coverageDir)
	}
}

// executeCommand runs the command and handles the result.
func (r *FlowRunner) executeCommand(cmd *exec.Cmd, start time.Time) *FlowTestResult {
	// Execute command
	err := cmd.Run()

	// Calculate duration
	duration := time.Since(start)

	// Determine exit code
	exitCode := r.determineExitCode(err)

	return &FlowTestResult{
		ExitCode: exitCode,
		Stdout:   r.stdout.String(),
		Stderr:   r.stderr.String(),
		Error:    err,
		Duration: duration,
	}
}
