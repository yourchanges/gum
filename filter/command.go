package filter

import (
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/gum/internal/exit"
	"github.com/charmbracelet/gum/internal/files"
	"github.com/charmbracelet/gum/internal/stdin"
	"github.com/charmbracelet/gum/style"
)

// Run provides a shell script interface for filtering through options, powered
// by the textinput bubble.
func (o Options) Run() error {
	i := textinput.New()
	i.Focus()

	i.Prompt = o.Prompt
	i.PromptStyle = o.PromptStyle.ToLipgloss()
	i.Placeholder = o.Placeholder
	i.Width = o.Width

	var choices []string
	if input, _ := stdin.Read(); input != "" {
		choices = strings.Split(strings.TrimSpace(input), "\n")
	} else {
		choices = files.List()
	}

	p := tea.NewProgram(model{
		choices:        choices,
		indicator:      o.Indicator,
		matches:        matchAll(choices),
		textinput:      i,
		indicatorStyle: o.IndicatorStyle.ToLipgloss(),
		matchStyle:     o.MatchStyle.ToLipgloss(),
		textStyle:      o.TextStyle.ToLipgloss(),
	}, tea.WithOutput(os.Stderr))

	tm, err := p.StartReturningModel()
	m := tm.(model)

	if m.aborted {
		return exit.ErrAborted
	}
	if len(m.matches) > m.selected && m.selected >= 0 {
		fmt.Println(m.matches[m.selected].Str)
	}

	return err
}

// BeforeReset hook. Used to unclutter style flags.
func (o Options) BeforeReset(ctx *kong.Context) error {
	return style.HideFlags(ctx)
}
