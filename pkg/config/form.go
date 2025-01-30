package config

import (
	"context"
	"fmt"
	"sort"

	"github.com/charmbracelet/huh"
	ollama "github.com/jmorganca/ollama/api"
)

type ConfigForm struct {
	form *huh.Form

	model   string
	shell   string
	explain string
	suggest string
}

func (c *ConfigForm) Run(api *ollama.Client) error {

	// Get available models from Ollama
	models, err := api.List(context.Background())
	if err != nil {
		fmt.Printf("Error fetching models: %v\n", err)
		return err
	}

	// Create model options from available Ollama models
	modelOptions := make([]huh.Option[string], 0, len(models.Models))
	sort.Slice(models.Models, func(i, j int) bool {
		return models.Models[i].Name < models.Models[j].Name
	})
	for _, model := range models.Models {
		modelOptions = append(modelOptions, huh.NewOption(
			fmt.Sprintf("%s (%.2f GB)", model.Name, float64(model.Size)/(1024*1024*1024)),
			model.Name,
		))
	}

	c.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("LLM").
				Description("Set a default Large Language Model.").
				Options(
					modelOptions...,
				).
				Value(&c.model),

			huh.NewSelect[string]().
				Title("Shell").
				Description("Overrides platform's shell for suggestions").
				Options(
					huh.NewOption("Automatic", "auto"),
					huh.NewOption("Powershell (Windows)", "powershell"),
					huh.NewOption("Bash (Linux)", "bash"),
					huh.NewOption("Zsh (macOS)", "zsh"),
				).
				Value(&c.shell),

			huh.NewSelect[string]().
				Title("Suggestion Preference").
				Description("Sets preference for command suggestions").
				Options(
					huh.NewOption("Precise", "stable"),
					huh.NewOption("Balanced", "balanced"),
					huh.NewOption("Creative", "creative"),
				).
				Value(&c.suggest),

			huh.NewSelect[string]().
				Title("Explain Preference").
				Description("Sets preference for command explanations").
				Options(
					huh.NewOption("Precise", "stable"),
					huh.NewOption("Balanced", "balanced"),
					huh.NewOption("Creative", "creative"),
				).
				Value(&c.explain),
		),
	)

	return c.form.WithTheme(huh.ThemeBase16()).Run()
}
