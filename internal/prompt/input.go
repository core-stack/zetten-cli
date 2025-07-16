package prompt

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

type PromptInputOpts struct {
	DefaultValue string
}

type CreatePromptInputOpts func(*PromptInputOpts)

func WithDefaultValue(defaultValue string) CreatePromptInputOpts {
	return func(pio *PromptInputOpts) {
		pio.DefaultValue = defaultValue
	}
}

func PromptInput(label string, opts ...CreatePromptInputOpts) (string, error) {
	options := &PromptInputOpts{}
	for _, opt := range opts {
		opt(options)
	}

	prompt := promptui.Prompt{
		Label:   label,
		Default: options.DefaultValue,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("❌ Error reading %s: %v\n", label, err)
	}
	return strings.TrimSpace(result), nil
}

func PromptConfirm(label string, defaultYes bool) (bool, error) {
	defaultStr := "n"
	if defaultYes {
		defaultStr = "y"
	}

	prompt := promptui.Prompt{
		Label:     fmt.Sprintf("%s (y/n)", label),
		Default:   defaultStr,
		AllowEdit: true,
		Validate: func(input string) error {
			switch strings.ToLower(input) {
			case "y", "n":
				return nil
			default:
				return fmt.Errorf("type 'y' or 'n'")
			}
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return defaultYes, nil
	}

	return strings.ToLower(result) == "y", nil
}
