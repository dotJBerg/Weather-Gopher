package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Thank you, claude for helping me style :)
var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		MarginBottom(1)
	
	subtitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#2D3748")).
		Padding(0, 1)

	infoStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#2D3748")).
		MarginTop(1)

	highlightStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000")).
		Bold(true)
)

type Model struct {
	location string
	loading  bool
	spinner  spinner.Model
	forecast []ForecastData
	err      error
	quitting bool
}

func InitialModel(location string, showForecast bool) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56D4"))

	return Model{
		location: location,
		loading:  true,
		spinner:  s,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchForecastCmd(),
	)
}

func (m Model) fetchForecastCmd() tea.Cmd {
	return func() tea.Msg {
		forecast, err := GetForecast(m.location)
		if err != nil {
			return fetchErrorMsg{err}
		}
		return forecastDataMsg{forecast}
	}
}

type forecastDataMsg struct {
	forecast []ForecastData
}

type fetchErrorMsg struct {
	err error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "r":
			if !m.loading {
				m.loading = true
				m.forecast = []ForecastData{}
				m.err = nil
				
				return m, tea.Batch(
					m.spinner.Tick,
					m.fetchForecastCmd(),
				)
			}
		}

	case forecastDataMsg:
		m.forecast = msg.forecast
		m.err = nil
		m.loading = false
		return m, nil

	case fetchErrorMsg:
		m.err = msg.err
		m.forecast = []ForecastData{}
		m.loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if m.quitting {
		return "Thanks for using Weather-Gopher!\n"
	}

	var s strings.Builder

	s.WriteString(titleStyle.Render("Weather-Gopher"))
	s.WriteString("\n")

	if m.loading {
		s.WriteString(fmt.Sprintf("\n %s Loading forecast for %s...\n",
			m.spinner.View(), m.location))
		return s.String()
	}

	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err)))
		s.WriteString("\nPress 'r' to retry or 'q' to quit\n")
		return s.String()
	}

	if len(m.forecast) > 0 {
		s.WriteString(subtitleStyle.Render(fmt.Sprintf("5-Day Forecast for %s", m.location)))
		s.WriteString("\n\n")

		for i, day := range m.forecast {
			weekday := day.Date.Format("Monday")
			date := day.Date.Format("Jan 2")

			if i > 0 {
				s.WriteString("\n")
			}

			s.WriteString(highlightStyle.Render(fmt.Sprintf("%s (%s):\n", weekday, date)))
			s.WriteString(fmt.Sprintf("  Conditions: %s\n", day.Description))
			s.WriteString(fmt.Sprintf("  Temperature: %d°F to %d°F\n", day.TempMin, day.TempMax))
			s.WriteString(fmt.Sprintf("  Humidity: %d%%\n", day.Humidity))
			s.WriteString(fmt.Sprintf("  Wind: %d mph\n", day.WindSpeed))
		}
	} else {
		s.WriteString("No forecast data available.\n")
	}

	s.WriteString("\n")
	s.WriteString(infoStyle.Render("Press 'r' to refresh, 'q' to quit"))

	return s.String()
}

func StartUI(location string, showForecast bool) error {
	p := tea.NewProgram(InitialModel(location, showForecast))
	_, err := p.Run()
	return err
}

