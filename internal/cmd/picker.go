package cmd

import (
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/jamesonstone/kp/internal/picker"
	"github.com/jamesonstone/kp/internal/prompt"
	"github.com/spf13/cobra"
)

func (a *app) runPicker(cmd *cobra.Command, prompts []prompt.Prompt) error {
	var (
		name string
		err  error
	)
	if a.noFzf {
		name, err = a.pickNumbered(prompts)
	} else {
		items := make([]picker.Item, len(prompts))
		for i, p := range prompts {
			items[i] = promptItem(cmd.Root().CommandPath(), p)
		}
		name, err = a.pick(cmd, items)
	}
	if err != nil {
		return err
	}
	return a.runPrompt(name)
}

// pick runs the interactive picker, or the injected test runner, and returns
// the selected item ID.
func (a *app) pick(cmd *cobra.Command, items []picker.Item) (string, error) {
	run := a.pickerRunner
	if run == nil {
		run = func(items []picker.Item) (string, error) {
			return picker.Run(cmd.Context(), cmd.CommandPath(), items)
		}
	}

	id, err := run(items)
	if err != nil {
		return "", mapPickerError(err)
	}
	for _, item := range items {
		if item.ID == id {
			return id, nil
		}
	}
	return "", NewExitError(ExitUser, fmt.Errorf("invalid picker selection %q", id))
}

func (a *app) pickNumbered(prompts []prompt.Prompt) (string, error) {
	for i, p := range prompts {
		fmt.Fprintf(a.stderr, "%d\t%s\t%s\n", i+1, p.Name, p.Label)
	}

	line, err := a.inputReader.ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		if errors.Is(err, io.EOF) {
			return "", NewExitError(ExitCancel, errPickerCanceled)
		}
		return "", NewExitError(ExitUser, err)
	}

	choiceText := strings.TrimSpace(line)
	if choiceText == "" {
		return "", NewExitError(ExitCancel, errPickerCanceled)
	}

	choice, err := strconv.Atoi(choiceText)
	if err != nil {
		return "", NewExitError(ExitUser, fmt.Errorf("invalid selection %q", choiceText))
	}
	if choice < 1 || choice > len(prompts) {
		return "", NewExitError(ExitUser, fmt.Errorf("selection %d out of range", choice))
	}

	return prompts[choice-1].Name, nil
}

var errPickerCanceled = picker.ErrCanceled

var pickerFarewells = []string{
	"👋 Tiny wave goodbye—your prompts will be right here.",
	"✨ Poof! The picker took the scenic exit.",
	"🦆 The picker waddled off. Quack soon!",
	"🌙 Menu tucked in for a tiny nap.",
	"🛸 Picker beamed out. Safe travels, captain!",
	"🧠 Tiny fun fact: choosing nothing is still a choice.",
	"🌈 Nothing selected; everything is still sparkly.",
	"🪁 The picker caught a breeze. See you next time!",
}

var (
	pickerFarewellStart  = uint64(rand.IntN(len(pickerFarewells)))
	pickerFarewellCursor atomic.Uint64
)

// ExitMessage returns terminal text without changing an error's exit classification.
func ExitMessage(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, errPickerCanceled) {
		return nextPickerFarewell()
	}
	return err.Error()
}

func nextPickerFarewell() string {
	offset := pickerFarewellCursor.Add(1) - 1
	index := (pickerFarewellStart + offset) % uint64(len(pickerFarewells))
	return pickerFarewells[index]
}

func mapPickerError(err error) error {
	switch {
	case errors.Is(err, errPickerCanceled):
		return NewExitError(ExitCancel, err)
	case errors.Is(err, picker.ErrNoTerminal):
		return NewExitError(ExitConfig, fmt.Errorf("%w; run 'kp --help' or 'kp list --no-fzf'", err))
	default:
		return NewExitError(ExitUser, err)
	}
}
