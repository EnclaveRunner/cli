package policy

import (
	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

// NewCmd returns the "policy" command group.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage RBAC policies",
	}
	cmd.AddCommand(
		newListCmd(),
		newCreateCmd(),
		newDeleteCmd(),
	)

	return cmd
}

// addPolicyFlags adds the shared --role, --resource-group, --method flags.
func addPolicyFlags(cmd *cobra.Command, methodHelp string) {
	cmd.Flags().String("role", "", "Role name (required)")
	cmd.Flags().String("resource-group", "", "Resource group name (required)")
	cmd.Flags().String("method", "", methodHelp+" (required)")
	_ = cmd.MarkFlagRequired("role")
	_ = cmd.MarkFlagRequired("resource-group")
	_ = cmd.MarkFlagRequired("method")
}

// policyFromFlags builds an enclave.Policy from the shared flags.
func policyFromFlags(cmd *cobra.Command) enclave.Policy {
	role, _ := cmd.Flags().GetString("role")
	rg, _ := cmd.Flags().GetString("resource-group")
	method, _ := cmd.Flags().GetString("method")

	return enclave.Policy{
		Role:          role,
		ResourceGroup: rg,
		Method:        enclave.PolicyMethod(method),
	}
}
