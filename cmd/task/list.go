package task

import (
	"fmt"
	"os"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List tasks",
		RunE:  runList,
	}
	cmd.Flags().String("state", "", "Filter by state (e.g. running, failed, completed)")
	return cmd
}

func runList(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.TaskColumns, os.Stdout)

	var opts []enclave.ListTasksOption
	if v, _ := cmd.Flags().GetString("state"); v != "" {
		opts = append(opts, enclave.FilterByState(v))
	}

	tasks, err := enclave.Collect(c.ListTasks(cmd.Context(), opts...))
	if err != nil {
		return fmt.Errorf("list tasks: %w", err)
	}
	return printer.Print(tasks)
}
