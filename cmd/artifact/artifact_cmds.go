package artifact

import (
	"fmt"
	"os"
	"strings"

	"cli/internal/client"
	"cli/internal/output"

	"github.com/EnclaveRunner/sdk-go/enclave"
	"github.com/spf13/cobra"
)

func newNamespaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "namespace",
		Short: "Manage artifact namespaces",
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all artifact namespaces",
		RunE:  runNamespaceList,
	}
	cmd.AddCommand(listCmd)
	return cmd
}

func runNamespaceList(cmd *cobra.Command, _ []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())

	// ListArtifactNamespaces returns Artifact objects; we only show namespace.
	nsCol := []output.Column{
		{Header: "NAMESPACE", Extract: func(r any) string {
			a := r.(enclave.Artifact)
			// Namespace entries have no Name; deduplicate by showing Namespace.
			return a.Namespace
		}},
	}
	printer := output.New(output.ParseFormat(cfg.Output), nsCol, os.Stdout)

	namespaces, err := enclave.Collect(c.ListArtifactNamespaces(cmd.Context()))
	if err != nil {
		return fmt.Errorf("list artifact namespaces: %w", err)
	}

	// Deduplicate namespace names.
	seen := map[string]bool{}
	unique := make([]enclave.Artifact, 0, len(namespaces))
	for _, a := range namespaces {
		key := a.Namespace
		if !seen[key] {
			seen[key] = true
			unique = append(unique, a)
		}
	}

	return printer.Print(unique)
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <namespace>",
		Short: "List artifacts in a namespace",
		Args:  cobra.ExactArgs(1),
		RunE:  runList,
	}
}

func runList(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.ArtifactColumns, os.Stdout)

	artifacts, err := enclave.Collect(c.ListArtifacts(cmd.Context(), args[0]))
	if err != nil {
		return fmt.Errorf("list artifacts: %w", err)
	}
	return printer.Print(artifacts)
}

func newVersionsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "versions <namespace> <name>",
		Short: "List all versions of an artifact",
		Args:  cobra.ExactArgs(2),
		RunE:  runVersions,
	}
}

func runVersions(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.ArtifactColumns, os.Stdout)

	versions, err := enclave.Collect(c.ListArtifactVersions(cmd.Context(), args[0], args[1]))
	if err != nil {
		return fmt.Errorf("list artifact versions: %w", err)
	}
	return printer.Print(versions)
}

func newUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload <namespace> <name> <file>",
		Short: "Upload an artifact",
		Args:  cobra.ExactArgs(3),
		RunE:  runUpload,
	}
}

func runUpload(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())

	f, err := os.Open(args[2])
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	result, err := c.UploadArtifact(cmd.Context(), args[0], args[1], f)
	if err != nil {
		return fmt.Errorf("upload artifact: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Uploaded. Version hash: %s\n", result.VersionHash)
	return nil
}

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <namespace> <name> <tag-or-hash>",
		Short: "Get artifact metadata by tag or hash",
		Args:  cobra.ExactArgs(3),
		RunE:  runGet,
	}
}

func runGet(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.ArtifactColumns, os.Stdout)

	namespace, name, ref := args[0], args[1], args[2]
	var a enclave.Artifact
	var err error
	if isHash(ref) {
		a, err = c.GetArtifactByHash(cmd.Context(), namespace, name, ref)
	} else {
		a, err = c.GetArtifactByTag(cmd.Context(), namespace, name, ref)
	}
	if err != nil {
		return fmt.Errorf("get artifact: %w", err)
	}
	return printer.Print([]any{a})
}

func newDownloadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download <namespace> <name> <tag-or-hash>",
		Short: "Download an artifact",
		Args:  cobra.ExactArgs(3),
		RunE:  runDownload,
	}
	cmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	return cmd
}

func runDownload(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())

	namespace, name, ref := args[0], args[1], args[2]
	var reader interface {
		Read(p []byte) (int, error)
		Close() error
	}
	var err error
	if isHash(ref) {
		reader, err = c.DownloadArtifactByHash(cmd.Context(), namespace, name, ref)
	} else {
		reader, err = c.DownloadArtifactByTag(cmd.Context(), namespace, name, ref)
	}
	if err != nil {
		return fmt.Errorf("download artifact: %w", err)
	}
	defer reader.Close()

	out, _ := cmd.Flags().GetString("output")
	var w *os.File
	if out == "" {
		w = os.Stdout
	} else {
		w, err = os.Create(out)
		if err != nil {
			return fmt.Errorf("create output file: %w", err)
		}
		defer w.Close()
	}

	buf := make([]byte, 32*1024)
	for {
		n, readErr := reader.Read(buf)
		if n > 0 {
			if _, writeErr := w.Write(buf[:n]); writeErr != nil {
				return fmt.Errorf("write output: %w", writeErr)
			}
		}
		if readErr != nil {
			break
		}
	}
	return nil
}

func newTagCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag <namespace> <name> <tag-or-hash>",
		Short: "Update tags on an artifact version",
		Args:  cobra.ExactArgs(3),
		RunE:  runTag,
	}
	cmd.Flags().StringSlice("tags", nil, "New tag list (replaces existing tags)")
	_ = cmd.MarkFlagRequired("tags")
	return cmd
}

func runTag(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.ArtifactColumns, os.Stdout)

	namespace, name, ref := args[0], args[1], args[2]
	tags, _ := cmd.Flags().GetStringSlice("tags")

	var a enclave.Artifact
	var err error
	if isHash(ref) {
		a, err = c.UpdateArtifactTagsByHash(cmd.Context(), namespace, name, ref, tags)
	} else {
		a, err = c.UpdateArtifactTagsByTag(cmd.Context(), namespace, name, ref, tags)
	}
	if err != nil {
		return fmt.Errorf("update artifact tags: %w", err)
	}
	return printer.Print([]any{a})
}

func newDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <namespace> <name> <tag-or-hash>",
		Short: "Delete an artifact version by tag or hash",
		Args:  cobra.ExactArgs(3),
		RunE:  runDelete,
	}
}

func runDelete(cmd *cobra.Command, args []string) error {
	c := client.FromContext(cmd.Context())
	cfg := client.ConfigFromContext(cmd.Context())
	printer := output.New(output.ParseFormat(cfg.Output), output.ArtifactColumns, os.Stdout)

	namespace, name, ref := args[0], args[1], args[2]
	var a enclave.Artifact
	var err error
	if isHash(ref) {
		a, err = c.DeleteArtifactByHash(cmd.Context(), namespace, name, ref)
	} else {
		a, err = c.DeleteArtifactByTag(cmd.Context(), namespace, name, ref)
	}
	if err != nil {
		return fmt.Errorf("delete artifact: %w", err)
	}
	return printer.Print([]any{a})
}

// isHash returns true if s looks like a SHA-256 hex digest (64 hex chars).
func isHash(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, r := range strings.ToLower(s) {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			return false
		}
	}
	return true
}
