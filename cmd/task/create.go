package task

import (
	"fmt"
	"os"
	"strings"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <source>",
		Short: "Create a new task",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreate,
	}
	cmd.Flags().StringSlice("args", nil, "Arguments to pass to the task")
	cmd.Flags().StringArray("env", nil, "Environment variables in KEY=VALUE format")
	cmd.Flags().String("callback", "", "Callback URL to invoke on completion")
	cmd.Flags().Int("retries", 0, "Maximum number of retries")
	cmd.Flags().String("retention", "", "Retention duration (e.g. 24h)")
	return cmd
}

func runCreate(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.TaskColumns, os.Stdout)

	var opts []enclave.CreateTaskOption

	if taskArgs, _ := cmd.Flags().GetStringSlice("args"); len(taskArgs) > 0 {
		opts = append(opts, enclave.WithArgs(taskArgs...))
	}
	if envVars, _ := cmd.Flags().GetStringArray("env"); len(envVars) > 0 {
		var envs []enclave.EnvironmentVariable
		for _, kv := range envVars {
			parts := strings.SplitN(kv, "=", 2)
			if len(parts) == 2 {
				envs = append(envs, enclave.EnvironmentVariable{Key: parts[0], Value: parts[1]})
			}
		}
		if len(envs) > 0 {
			opts = append(opts, enclave.WithEnv(envs...))
		}
	}
	if v, _ := cmd.Flags().GetString("callback"); v != "" {
		opts = append(opts, enclave.WithCallback(v))
	}
	if cmd.Flags().Changed("retries") {
		n, _ := cmd.Flags().GetInt("retries")
		opts = append(opts, enclave.WithRetries(n))
	}
	if v, _ := cmd.Flags().GetString("retention"); v != "" {
		opts = append(opts, enclave.WithRetention(v))
	}

	t, err := c.CreateTask(cmd.Context(), args[0], opts...)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return printer.Print([]any{t})
}
