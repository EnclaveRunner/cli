package task

import (
	"cli/internal/client"
	"cli/internal/output"
	"fmt"
	"os"
	"time"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs <id>",
		Short: "Get logs for a task",
		Args:  cobra.ExactArgs(1),
		RunE:  runLogs,
	}
	cmd.Flags().
		String("level", "", "Filter by log level (trace, debug, info, warn, error)")
	cmd.Flags().String("issuer", "", "Filter by issuer")
	cmd.Flags().String("since", "", "Include logs after this time (RFC3339)")
	cmd.Flags().String("until", "", "Include logs before this time (RFC3339)")

	return cmd
}

func runLogs(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(
		output.ParseFormat(cfg.Output),
		output.TaskLogColumns,
		os.Stdout,
	)

	var opts []enclave.TaskLogOption

	if v, _ := cmd.Flags().GetString("level"); v != "" {
		opts = append(opts, enclave.FilterLogByLevel(v))
	}
	if v, _ := cmd.Flags().GetString("issuer"); v != "" {
		opts = append(opts, enclave.FilterLogByIssuer(v))
	}

	since, _ := cmd.Flags().GetString("since")
	until, _ := cmd.Flags().GetString("until")
	if since != "" || until != "" {
		var from, to time.Time
		var err error
		if since != "" {
			from, err = time.Parse(time.RFC3339, since)
			if err != nil {
				return fmt.Errorf("invalid --since: %w", err)
			}
		}
		if until != "" {
			to, err = time.Parse(time.RFC3339, until)
			if err != nil {
				return fmt.Errorf("invalid --until: %w", err)
			}
		}
		if from.IsZero() {
			from = time.Unix(0, 0)
		}
		if to.IsZero() {
			to = time.Now()
		}
		opts = append(opts, enclave.FilterLogByTimeRange(from, to))
	}

	logs, err := c.GetTaskLogs(cmd.Context(), args[0], opts...)
	if err != nil {
		return fmt.Errorf("get task logs: %w", err)
	}

	return printer.Print(logs)
}
