package cmd

import (
	"bytes"
	"cli/client"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Upload, download and manage artifacts",
	Long:  `Commands for uploading, downloading and managing of artifacts.`,
}

var artifactListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"query"},
	Short:   "List artifacts",
	Long:    `List all artifacts in the Enclave system. Use flags to filter the results.`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		var source *string
		var author *string
		var name *string

		if s, err := cmd.Flags().GetString("source"); err == nil && s != "" {
			source = &s
		}
		if a, err := cmd.Flags().GetString("author"); err == nil && a != "" {
			author = &a
		}
		if n, err := cmd.Flags().GetString("name"); err == nil && n != "" {
			name = &n
		}

		c := getClient()
		l, err := c.GetArtifactListWithResponse(
			cmd.Context(),
			&client.GetArtifactListParams{
				Source: source,
				Author: author,
				Name:   name,
			},
		)

		handleResponse(l, err, "")

		tableData := [][]string{}
		for _, artifact := range *l.JSON200 {
			fqn := fmt.Sprintf(
				"%s/%s/%s",
				artifact.Fqn.Source,
				artifact.Fqn.Author,
				artifact.Fqn.Name,
			)
			tags := ""
			if len(artifact.Tags) > 0 {
				tags = strings.Join(artifact.Tags, "\n")
			}
			tableData = append(tableData, []string{
				fqn,
				artifact.VersionHash,
				tags,
				artifact.CreatedAt.Format("2006-01-02 15:04:05"),
				strconv.Itoa(artifact.Pulls),
			})
		}

		printTable(tableData, []string{
			"FQN",
			"HASH",
			"TAGS",
			"CREATED",
			"PULLS",
		})
	},
}

var artifactUploadCmd = &cobra.Command{
	Use:     "upload <fqn> <wasm-file>",
	Aliases: []string{"create", "push"},
	Short:   "Upload new artifact",
	Long:    "Upload a new artifact and assign tags",
	Args: func(cmd *cobra.Command, args []string) error {
		//nolint:mnd // 2 because two args are required
		if len(args) != 2 {
			return errors.New(
				"provide a FQN (Fully Qualified Name) and the path to a wasm file to upload",
			)
		}

		if _, err := parseFQN(args[0]); err != nil {
			return err
		}

		if !strings.HasSuffix(args[1], ".wasm") {
			return errors.New("provided file needs to be compiled WASM (.wasm)")
		}

		if _, err := os.Stat(args[1]); err != nil {
			return fmt.Errorf("provided wasm file not valid: %w", err)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fqn, _ := parseFQN(args[0])

		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		//nolint:errcheck // Close only throws when called twice (not possible here)
		defer w.Close()

		err := w.WriteField("source", fqn.Source)
		if err != nil {
			log.Fatal().Str("field", "source").Msg("Failed to write field")
		}

		err = w.WriteField("author", fqn.Author)
		if err != nil {
			log.Fatal().Str("field", "author").Msg("Failed to write field")
		}

		err = w.WriteField("name", fqn.Name)
		if err != nil {
			log.Fatal().Err(err).Str("field", "name").Msg("Failed to write field")
		}

		tagsString, err := cmd.Flags().GetString("tags")
		if err == nil {
			for i, tag := range strings.Split(tagsString, " ") {
				tag = strings.TrimSpace(tag)
				if tag == "" {
					continue
				}

				err = w.WriteField("tag", tag)
				if err != nil {
					log.Fatal().
						Err(err).
						Int("index", i).
						Str("field", "tag").
						Msg("Failed to write field")
				}
			}
		}

		fileWriter, err := w.CreateFormFile("file", "plugin.wasm")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create form file")
		}

		fileReader, err := os.Open(args[1])
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to open wasm file")
		}

		written, err := io.Copy(fileWriter, fileReader)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read wasm file")
		}

		log.Info().Int64("size", written).Msg("Read wasm file")

		err = w.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close multipart writer")
		}

		c := getClient()
		r, err := c.PostArtifactUploadWithBodyWithResponse(
			cmd.Context(),
			"multipart/form-data; boundary="+w.Boundary(),
			&b,
		)

		handleResponse(r, err, "Uploaded artifact successfully!")

		printArtifact(r.JSON201)
	},
}

var artifactDownloadCmd = &cobra.Command{
	Use:     "download <fqn> <output-file>",
	Aliases: []string{"pull"},
	Short:   "Download an artifact",
	Long:    "Download an artifact by its FQN (Fully Qualified Name)",
	Args: func(cmd *cobra.Command, args []string) error {
		//nolint:mnd // Two arguments are expected
		if len(args) != 2 {
			return errors.New(
				"provide a FQN (Fully Qualified Name) and the output file path",
			)
		}

		if _, _, err := parseFQNWithIdentifier(args[0]); err != nil {
			return err
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fqn, identifier, _ := parseFQNWithIdentifier(args[0])
		outputFile := args[1]

		c := getClient()
		r, err := c.GetArtifactUpload(
			cmd.Context(),
			&client.GetArtifactUploadParams{
				Source:     fqn.Source,
				Author:     fqn.Author,
				Name:       fqn.Name,
				Identifier: identifier,
			},
		)

		handleResponse(
			&GenericResponseWithBody{Response: r},
			err,
			"Downloading Artifact...",
		)

		//nolint:gosec // File creation from user input is intended here
		outFile, err := os.Create(outputFile)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create output file")
		}

		//nolint:errcheck // Ignore close error
		defer outFile.Close()

		written, err := io.Copy(outFile, r.Body)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to write to output file")
		}

		log.Info().
			Int64("size", written).
			Str("file", outputFile).
			Msg("Artifact downloaded successfully")
	},
}

var artifactMetadataCmd = &cobra.Command{
	Use:     "metadata <fqn>",
	Aliases: []string{"meta", "info"},
	Short:   "Get Artifact Metadata",
	Long:    "Get artifact metadata from a provided FQN with an identifier (hash or tag)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("provide an FQN (Fully Qualified Name")
		}

		_, _, err := parseFQNWithIdentifier(args[0])

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		fqn, identifier, _ := parseFQNWithIdentifier(args[0])

		c := getClient()
		artifact, err := c.GetArtifactWithResponse(
			cmd.Context(),
			&client.GetArtifactParams{
				Source:     fqn.Source,
				Author:     fqn.Author,
				Name:       fqn.Name,
				Identifier: identifier,
			},
		)

		handleResponse(artifact, err, "")

		printArtifact(artifact.JSON200)
	},
}

var artifactDeleteCmd = &cobra.Command{
	Use:     "delete <fqn>",
	Aliases: []string{"remove"},
	Short:   "Delete an artifact",
	Long:    "Delete an artifact by its FQN (Fully Qualified Name) with its identifier (hash or tag)",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("provide an FQN (Fully Qualified Name)")
		}

		_, _, err := parseFQNWithIdentifier(args[0])

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		fqn, identifier, _ := parseFQNWithIdentifier(args[0])

		c := getClient()
		resp, err := c.DeleteArtifactWithResponse(
			cmd.Context(),
			client.DeleteArtifactJSONRequestBody{
				Fqn:        fqn,
				Identifier: identifier,
			},
		)

		handleResponse(resp, err, "Artifact deleted successfully")
	},
}

var artifactTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage artifact tags",
	Long:  "Commands for managing artifact tags (add, remove)",
}

//nolint:dupl // Similar code as in artifactTagRemoveCmd
var artifactTagAddCmd = &cobra.Command{
	Use:     "add <fqn> <tag>",
	Aliases: []string{"create"},
	Short:   "Add tag to artifact",
	Long:    "Add a tag to an existing artifact",
	Args: func(cmd *cobra.Command, args []string) error {
		//nolint:mnd // Two arguments are expected
		if len(args) != 2 {
			return errors.New("provide an FQN (Fully Qualified Name) and a tag")
		}

		_, identifier, err := parseFQNWithIdentifier(args[0])
		if err != nil {
			return err
		}

		if !strings.HasPrefix(identifier, "hash:") {
			return fmt.Errorf(
				"can only add tags to artifacts identified by hash, got identifier: %s",
				identifier,
			)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fqn, identifier, _ := parseFQNWithIdentifier(args[0])
		tag := args[1]

		identifier = strings.TrimPrefix(identifier, "hash:")

		c := getClient()
		resp, err := c.PostArtifactTagWithResponse(
			cmd.Context(),
			client.PostArtifactTagJSONRequestBody{
				Fqn:         fqn,
				VersionHash: identifier,
				NewTag:      tag,
			},
		)

		handleResponse(resp, err, "Tag added successfully")
	},
}

//nolint:dupl // Similar code as in artifactTagAddCmd
var artifactTagRemoveCmd = &cobra.Command{
	Use:     "remove <fqn> <tag>",
	Aliases: []string{"delete"},
	Short:   "Remove tag from artifact",
	Long:    "Remove a tag from an existing artifact",
	Args: func(cmd *cobra.Command, args []string) error {
		//nolint:mnd // Two arguments are expected
		if len(args) != 2 {
			return errors.New("provide an FQN (Fully Qualified Name) and a tag")
		}

		_, identifier, err := parseFQNWithIdentifier(args[0])
		if err != nil {
			return err
		}

		if !strings.HasPrefix(identifier, "hash:") {
			return fmt.Errorf(
				"can only remove tags from artifacts identified by hash, got identifier: %s",
				identifier,
			)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fqn, identifier, _ := parseFQNWithIdentifier(args[0])
		tag := args[1]

		identifier = strings.TrimPrefix(identifier, "hash:")

		c := getClient()
		resp, err := c.DeleteArtifactTagWithResponse(
			cmd.Context(),
			client.DeleteArtifactTagJSONRequestBody{
				Fqn:         fqn,
				VersionHash: identifier,
				Tag:         tag,
			},
		)

		handleResponse(resp, err, "Tag removed successfully")
	},
}

func init() {
	rootCmd.AddCommand(artifactCmd)

	// List/Query command
	artifactCmd.AddCommand(artifactListCmd)
	artifactListCmd.Flags().
		StringP("source", "s", "", "Only list matching source")
	artifactListCmd.Flags().
		StringP("author", "a", "", "Only list matching author")
	artifactListCmd.Flags().
		StringP("name", "n", "", "Only list matching name")

	// Upload command
	artifactCmd.AddCommand(artifactUploadCmd)
	artifactUploadCmd.Flags().
		StringP("tags", "t", "", "Space separated list of tags to add to the upload")

	// Download command
	artifactCmd.AddCommand(artifactDownloadCmd)

	// Get metadata command
	artifactCmd.AddCommand(artifactMetadataCmd)

	// Delete artifact command
	artifactCmd.AddCommand(artifactDeleteCmd)

	// Tag management command
	artifactCmd.AddCommand(artifactTagCmd)

	// Tag add command
	artifactTagCmd.AddCommand(artifactTagAddCmd)

	// Tag remove command
	artifactTagCmd.AddCommand(artifactTagRemoveCmd)
}

func parseFQN(fqn string) (client.FQN, error) {
	parts := strings.Split(fqn, "/")
	//nolint:mnd // FQN consists of 3 parts
	if len(parts) != 3 {
		return client.FQN{}, errors.New("provided FQN invalid")
	}

	return client.FQN{
		Source: parts[0],
		Author: parts[1],
		Name:   parts[2],
	}, nil
}

func parseFQNWithIdentifier(fqn string) (client.FQN, string, error) {
	parts := strings.Split(fqn, "/")
	//nolint:mnd // FQN consists of 3 parts separated by "/"
	if len(parts) != 3 {
		return client.FQN{}, "", errors.New(
			"provided FQN with identifier invalid",
		)
	}

	//nolint:mnd // Last part consists of <name>:<identifier>
	nameAndIdentifier := strings.SplitN(parts[2], ":", 2)
	//nolint:mnd // Check that last part consists of two parts
	if len(nameAndIdentifier) != 2 {
		return client.FQN{}, "", errors.New(
			"provided FQN with identifier invalid",
		)
	}

	return client.FQN{
		Source: parts[0],
		Author: parts[1],
		Name:   nameAndIdentifier[0],
	}, nameAndIdentifier[1], nil
}
