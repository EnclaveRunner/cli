package cmd

import (
	"cli/config"
	"fmt"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type ConfigModel struct {
	configTable table.Model
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display current configuration",
	Long:  `Display the current configuration loaded from files, environment variables, and flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(newModel())
		if err := p.Start(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start config command")
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func newModel() *ConfigModel {
	const columnParamID = "parameter"
	const columnValueID = "value"

	values := iterateStruct(config.Cfg)

	largestParam := 0
	largestValue := 0
	for _, pair := range values {
		if len(pair[0]) > largestParam {
			largestParam = len(pair[0])
		}
		if len(pair[1]) > largestValue {
			largestValue = len(pair[1])
		}
	}

	columns := []table.Column{
		table.NewFlexColumn(columnParamID, "Parameter", 1),
		table.NewColumn(columnValueID, "Value", largestValue),
	}

	rows := make([]table.Row, len(values))
	for i, pair := range values {
		rows[i] = table.NewRow(table.RowData{
			columnParamID: pair[0],
			columnValueID: pair[1],
		})
	}

	tableModel := table.New(columns).WithRows(rows)

	return &ConfigModel{
		configTable: tableModel,
	}
}

func iterateStruct(v any) [][]string {
	val := reflect.ValueOf(v)

	return parseValues(val, "")
}

func parseValues(val reflect.Value, prefix string) [][]string {
	// Handle pointers
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		values := [][]string{}
		typ := val.Type()
		for i := range val.NumField() {
			field := val.Field(i)
			name := typ.Field(i).Name
			values = append(values, parseValues(field, prefix+name+".")...)
		}

		return values
	default:
		return [][]string{
			{prefix[:len(prefix)-1], fmt.Sprintf("%v", val.Interface())},
		}
	}
}

func (m *ConfigModel) Init() tea.Cmd {
	return nil
}

func (m *ConfigModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case tea.WindowSizeMsg:
		m.configTable.WithTargetWidth(msg.Width - 2)
	}

	m.configTable, cmd = m.configTable.Update(msg)

	return m, cmd
}

func (m *ConfigModel) View() string {
	body := strings.Builder{}

	body.WriteString("Current configuration of the enclave CLI\n")

	pad := lipgloss.NewStyle().Padding(1)

	configTable := pad.Render(m.configTable.View())

	body.WriteString(configTable)

	return body.String()
}
