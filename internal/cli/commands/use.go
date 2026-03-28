package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/farion1231/ccmux/internal/store"
)

type ExecResult struct {
	Binary string
	Args   []string
}

func RunUse(state *store.AppState, name string, extraArgs []string) error {
	overlay, err := state.Service.BuildOverlay(name)
	if err != nil {
		return err
	}
	tmpfile, err := os.CreateTemp("", "ccc-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	data, err := json.MarshalIndent(overlay, "", "  ")
	if err != nil {
		os.Remove(tmpfile.Name())
		return err
	}
	if _, err := tmpfile.Write(data); err != nil {
		os.Remove(tmpfile.Name())
		return err
	}
	tmpfile.Close()
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		os.Remove(tmpfile.Name())
		return fmt.Errorf("claude not found in PATH: %w", err)
	}
	args := []string{"claude", "--settings", tmpfile.Name()}
	args = append(args, extraArgs...)
	p, _ := state.Service.GetByName(name)
	providerName := name
	if p != nil {
		providerName = p.Name
	}
	fmt.Printf("▸ launching claude with profile '%s'\n", providerName)
	return syscall.Exec(claudePath, args, os.Environ())
}

// BuildExecResult builds an ExecResult for use by the TUI (post-exit syscall.Exec).
func BuildExecResult(state *store.AppState, name string, extraArgs []string) (*ExecResult, error) {
	overlay, err := state.Service.BuildOverlay(name)
	if err != nil {
		return nil, err
	}

	tmpfile, err := os.CreateTemp("", "ccc-*.json")
	if err != nil {
		return nil, err
	}

	data, err := json.MarshalIndent(overlay, "", "  ")
	if err != nil {
		os.Remove(tmpfile.Name())
		return nil, err
	}
	if _, err := tmpfile.Write(data); err != nil {
		os.Remove(tmpfile.Name())
		return nil, err
	}
	tmpfile.Close()

	claudePath, err := exec.LookPath("claude")
	if err != nil {
		os.Remove(tmpfile.Name())
		return nil, fmt.Errorf("claude not found in PATH: %w", err)
	}

	args := []string{"claude", "--settings", tmpfile.Name()}
	args = append(args, extraArgs...)

	return &ExecResult{Binary: claudePath, Args: args}, nil
}
