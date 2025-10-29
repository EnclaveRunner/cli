package cmd

import (
	"cli/client"
	"cli/config"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/list"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rs/zerolog/log"
)

const (
	colorPrimaryGreen = "#6f7f37"
)

func getClient() *client.ClientWithResponses {
	if config.Cfg.APIServerURL == "" {
		log.Fatal().Msg("API server URL not configured")
	}

	if config.Cfg.Auth == nil {
		log.Fatal().Msg("Authentication not configured")
	}

	c, err := client.NewClientWithResponses(
		config.Cfg.APIServerURL,
		client.WithRequestEditorFn(
			func(ctx context.Context, req *http.Request) error {
				req.Header.Set("Authorization", config.Cfg.Auth.GetAuthHeader())

				return nil
			},
		),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create API client")
	}

	return c
}

func handleResponse(resp *http.Response, successMsg string) bool {
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode < 300:
		if successMsg != "" {
			log.Info().Msg(successMsg)
		}

		return true

	case resp.StatusCode == http.StatusUnauthorized:
		log.Error().Msg("Unauthorized: Invalid credentials")

		return false

	case resp.StatusCode == http.StatusForbidden:
		log.Error().
			Msg("Forbidden: You do not have permission to perform this action")

		return false

	case resp.StatusCode == http.StatusNotFound:
		log.Error().Msg("Not Found: The requested resource does not exist")

		return false

	case resp.StatusCode == http.StatusInternalServerError:
		log.Error().
			Msg("Internal Server Error: An error occurred on the server. Look at the server logs for more details.")

		return false

	default:
		log.Error().Msgf("Request failed with status code %d", resp.StatusCode)
		if resp.Body != nil {
			defer resp.Body.Close()
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error().Err(err).Msg("Failed to read response body")

				return false
			}

			// Check if body is json with error field
			var dest client.ErrGeneric
			if err := json.Unmarshal(bodyBytes, &dest); err != nil {
				log.Error().Msgf("Error: %s", string(bodyBytes))

				return false
			}

			log.Error().Msgf("Error: %s", dest.Error)
		}

		return false
	}
}

func printSlice(arr []string) {
	conv := make([]any, len(arr))
	for i, v := range arr {
		conv[i] = v
	}

	l := list.New(conv...).Enumerator(list.Arabic).
		ItemStyle(lipgloss.NewStyle().Bold(true).PaddingLeft(1)).
		EnumeratorStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(colorPrimaryGreen)))

	fmt.Println(l)
}

func printUsers(users []*client.UserResponse) {
	data := make([][]string, len(users))
	headers := []string{"ID", "Username", "Display Name"}

	for i, user := range users {
		data[i] = []string{
			user.Id,
			user.Name,
			user.DisplayName,
		}
	}

	printTable(data, headers)
}

func printTable(data [][]string, headers []string) {
	baseStyle := lipgloss.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Bold(true).Foreground(lipgloss.Color(colorPrimaryGreen))

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("244"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			default:
				return baseStyle
			}
		})

	t.Headers(headers...)
	t.Rows(data...)
	fmt.Println(t)
}

func getUserById(ctx context.Context, userId string) (*client.UserResponse, error) {
	c := getClient()
	params := &client.GetUsersUserParams{
		UserId: userId,
	}

	resp, err := c.GetUsersUserWithResponse(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user")

		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: status %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}

// getUsersByIds fetches multiple users concurrently by their IDs
func getUsersByIds(ctx context.Context, userIds []string) ([]*client.UserResponse, error) {
	if len(userIds) == 0 {
		return []*client.UserResponse{}, nil
	}

	type userResult struct {
		user  *client.UserResponse
		err   error
		index int
	}

	results := make(chan userResult, len(userIds))
	var wg sync.WaitGroup

	// Launch concurrent requests
	for i, userId := range userIds {
		wg.Add(1)
		go func(id string, idx int) {
			defer wg.Done()
			user, err := getUserById(ctx, id)
			results <- userResult{user: user, err: err, index: idx}
		}(userId, i)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results in order
	users := make([]*client.UserResponse, len(userIds))
	var errs []error

	for result := range results {
		if result.err != nil {
			errs = append(errs, result.err)

			continue
		}
		users[result.index] = result.user
	}

	// Return error if any requests failed
	if len(errs) > 0 {
		errMsg := fmt.Sprintf("failed to fetch %d user(s):", len(errs))
		for i, err := range errs {
			errMsg += fmt.Sprintf("\n  [%d] %v", i+1, err)
		}
		return nil, fmt.Errorf("%s", errMsg)
	}

	return users, nil
}
