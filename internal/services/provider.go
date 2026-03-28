package services

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotFound     = errors.New("provider not found")
	ErrAmbiguous    = errors.New("ambiguous provider name")
	ErrNotSupported = errors.New("not supported in this mode")
)

type Provider struct {
	Name        string
	Description string
	Env         map[string]string
	ModelAlias  string
}

type SettingsOverlay struct {
	Env   map[string]string `json:"env,omitempty"`
	Model string            `json:"model,omitempty"`
}

type ProviderService interface {
	List() ([]Provider, error)
	GetCurrent() (*Provider, error)
	GetByName(name string) (*Provider, error)
	SetCurrent(name string) error
	BuildOverlay(name string) (*SettingsOverlay, error)

	Add(p Provider) error
	Edit(name string) error
	Remove(name string) error
	Import(name string) error
}

func ResolveName(input string, providers []Provider) (string, error) {
	inputLower := strings.ToLower(input)
	var exact, prefix []string

	for _, p := range providers {
		nameLower := strings.ToLower(p.Name)
		if nameLower == inputLower {
			exact = append(exact, p.Name)
		} else if strings.HasPrefix(nameLower, inputLower) {
			prefix = append(prefix, p.Name)
		}
	}

	if len(exact) == 1 {
		return exact[0], nil
	}
	if len(exact) > 1 {
		return "", fmt.Errorf("%w '%s', matches: %s", ErrAmbiguous, input, strings.Join(exact, ", "))
	}
	if len(prefix) == 1 {
		return prefix[0], nil
	}
	if len(prefix) > 1 {
		return "", fmt.Errorf("%w '%s', matches: %s", ErrAmbiguous, input, strings.Join(prefix, ", "))
	}
	return "", fmt.Errorf("%w: '%s'", ErrNotFound, input)
}
