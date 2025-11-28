package cmd

import (
	"cli/client"
	"cli/config"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rs/zerolog/log"
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

// ResponseWithBody is an interface that matches all generated API response
// types
// All response types have StatusCode() method and Body field
type ResponseWithBody interface {
	StatusCode() int
}

type GenericResponseWithBody struct {
	Response *http.Response
}

func (r *GenericResponseWithBody) StatusCode() int {
	return r.Response.StatusCode
}

// handleResponse safely handles API responses, checking for nil before
// accessing fields
// This prevents segfaults when an error occurs and resp is nil
// It uses reflection to safely access the Body field which exists on all
// response types
func handleResponse(resp ResponseWithBody, err error, successMsg string) {
	if err != nil {
		log.Fatal().Err(err).Msg("Request failed")
	}

	if resp == nil {
		log.Fatal().Msg("Request failed: no response received")
	}

	// Extract body using reflection to safely access the Body field
	// All generated response types have: Body []byte field
	var body []byte

	// Use reflection to access the Body field - all response types have this
	// field
	respValue := reflect.ValueOf(resp)
	if respValue.Kind() == reflect.Pointer {
		respValue = respValue.Elem()
	}

	// Get the Body field
	bodyField := respValue.FieldByName("Body")
	if bodyField.IsValid() && bodyField.Kind() == reflect.Slice {
		body = bodyField.Bytes()
	} else {
		// Fallback: no body available
		body = []byte{}
	}

	switch {
	case resp.StatusCode() >= 200 && resp.StatusCode() < 300:
		if successMsg != "" {
			log.Info().Msg(TextPrimary.Render(successMsg))
		}

	case resp.StatusCode() == http.StatusUnauthorized:
		log.Fatal().Msg("Unauthorized: Invalid credentials")

	case resp.StatusCode() == http.StatusForbidden:
		log.Fatal().
			Msg("Forbidden: You do not have permission to perform this action")

	case resp.StatusCode() == http.StatusNotFound:
		log.Fatal().Msg("Not Found: The requested resource does not exist")

	case resp.StatusCode() == http.StatusInternalServerError:
		log.Fatal().
			Msg("Internal Server Error: An error occurred on the server. Look at the server logs for more details.")

	default:
		if len(body) > 0 {
			// Check if body is json with error field
			var dest client.ErrGeneric
			if err := json.Unmarshal(body, &dest); err != nil {
				log.Fatal().
					Msgf("Request failed with status code %d: %s", resp.StatusCode(), string(body))
			} else {
				log.Fatal().Msgf("Request failed with status code %d: %s", resp.StatusCode(), dest.Error)
			}
		} else {
			log.Fatal().Msgf("Request failed with status code %d", resp.StatusCode())
		}
	}
}

func printStringTable(arr []string, header string) {
	data := make([][]string, len(arr))
	for i, v := range arr {
		data[i] = []string{v}
	}

	headers := []string{header}
	printTable(data, headers)
}

// RoleInfo contains role name with additional metadata
type RoleInfo struct {
	Role        string
	UserCount   int
	PolicyCount int
}

func printRoles(roles []RoleInfo) {
	data := make([][]string, len(roles))
	headers := []string{"ROLE", "USERS", "POLICIES"}

	for i, role := range roles {
		data[i] = []string{
			role.Role,
			strconv.Itoa(role.UserCount),
			strconv.Itoa(role.PolicyCount),
		}
	}

	printTable(data, headers)
}

// ResourceGroupInfo contains resource group name with additional metadata
type ResourceGroupInfo struct {
	ResourceGroup string
	EndpointCount int
	PolicyCount   int
}

func printResourceGroups(groups []ResourceGroupInfo) {
	data := make([][]string, len(groups))
	headers := []string{"RESOURCE GROUP", "ENDPOINTS", "POLICIES"}

	for i, group := range groups {
		data[i] = []string{
			group.ResourceGroup,
			strconv.Itoa(group.EndpointCount),
			strconv.Itoa(group.PolicyCount),
		}
	}

	printTable(data, headers)
}

// getRoleInfo fetches role metadata (user count and policy count)
//
//nolint:dupl // Similar code exists for resource groups
func getRoleInfo(ctx context.Context, roles []string) []RoleInfo {
	c := getClient()

	// Get all policies to count per role
	policiesResp, err := c.GetRbacPolicyWithResponse(ctx)
	if err != nil || policiesResp.JSON200 == nil {
		// If we can't get policies, return roles with zero counts
		result := make([]RoleInfo, len(roles))
		for i, role := range roles {
			result[i] = RoleInfo{Role: role, UserCount: 0, PolicyCount: 0}
		}

		return result
	}

	policies := *policiesResp.JSON200

	// Count policies and users per role
	roleCounts := make(map[string]*RoleInfo)
	for _, role := range roles {
		roleCounts[role] = &RoleInfo{Role: role, UserCount: 0, PolicyCount: 0}
	}

	// Count policies
	for _, policy := range policies {
		if info, exists := roleCounts[policy.Role]; exists {
			info.PolicyCount++
		}
	}

	// Count users per role concurrently
	type roleUserCount struct {
		role  string
		count int
	}

	resultChan := make(chan roleUserCount, len(roles))
	var wg sync.WaitGroup

	for _, role := range roles {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			params := &client.GetRbacRoleParams{Role: r}
			resp, err := c.GetRbacRoleWithResponse(ctx, params)
			if err == nil && resp.JSON200 != nil {
				resultChan <- roleUserCount{role: r, count: len(*resp.JSON200)}
			} else {
				resultChan <- roleUserCount{role: r, count: 0}
			}
		}(role)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if info, exists := roleCounts[result.role]; exists {
			info.UserCount = result.count
		}
	}

	// Convert map to slice in original order
	result := make([]RoleInfo, len(roles))
	for i, role := range roles {
		result[i] = *roleCounts[role]
	}

	return result
}

// getResourceGroupInfo fetches resource group metadata (endpoint count and
// policy count)
//
//nolint:dupl // Similar code exists for roles
func getResourceGroupInfo(
	ctx context.Context,
	resourceGroups []string,
) []ResourceGroupInfo {
	c := getClient()

	// Get all policies to count per resource group
	policiesResp, err := c.GetRbacPolicyWithResponse(ctx)
	if err != nil || policiesResp.JSON200 == nil {
		// If we can't get policies, return groups with zero counts
		result := make([]ResourceGroupInfo, len(resourceGroups))
		for i, rg := range resourceGroups {
			result[i] = ResourceGroupInfo{
				ResourceGroup: rg,
				EndpointCount: 0,
				PolicyCount:   0,
			}
		}

		return result
	}

	policies := *policiesResp.JSON200

	// Count policies per resource group
	rgCounts := make(map[string]*ResourceGroupInfo)
	for _, rg := range resourceGroups {
		rgCounts[rg] = &ResourceGroupInfo{
			ResourceGroup: rg,
			EndpointCount: 0,
			PolicyCount:   0,
		}
	}

	for _, policy := range policies {
		if info, exists := rgCounts[policy.ResourceGroup]; exists {
			info.PolicyCount++
		}
	}

	// Count endpoints per resource group concurrently
	type rgEndpointCount struct {
		resourceGroup string
		count         int
	}

	resultChan := make(chan rgEndpointCount, len(resourceGroups))
	var wg sync.WaitGroup

	for _, rg := range resourceGroups {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			params := &client.GetRbacResourceGroupParams{ResourceGroup: r}
			resp, err := c.GetRbacResourceGroupWithResponse(ctx, params)
			if err == nil && resp.JSON200 != nil {
				resultChan <- rgEndpointCount{resourceGroup: r, count: len(*resp.JSON200)}
			} else {
				resultChan <- rgEndpointCount{resourceGroup: r, count: 0}
			}
		}(rg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if info, exists := rgCounts[result.resourceGroup]; exists {
			info.EndpointCount = result.count
		}
	}

	// Convert map to slice in original order
	result := make([]ResourceGroupInfo, len(resourceGroups))
	for i, rg := range resourceGroups {
		result[i] = *rgCounts[rg]
	}

	return result
}

// EndpointInfo contains endpoint name with resource group assignment
type EndpointInfo struct {
	Endpoint      string
	ResourceGroup string
}

func printUser(user *client.UserResponse) {
	data := [][]string{
		{user.Id, user.Name, user.DisplayName},
	}
	headers := []string{"ID", "USERNAME", "DISPLAY NAME"}
	printTable(data, headers)
}

func printUsers(users []*client.UserResponse) {
	data := make([][]string, len(users))
	headers := []string{"ID", "USERNAME", "DISPLAY NAME"}

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
	headerStyle := baseStyle.Bold(true)
	rowStyle := baseStyle.Foreground(ColorPrimary)

	t := table.New().
		BorderBottom(false).
		BorderColumn(false).
		BorderHeader(false).
		BorderLeft(false).
		BorderRight(false).
		BorderRow(false).
		BorderTop(false).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case table.HeaderRow:
				return headerStyle
			default:
				return rowStyle
			}
		})

	t.Headers(headers...)
	t.Rows(data...)
	fmt.Println(t)
}

func printPolicies(policies []client.RBACPolicy) {
	data := make([][]string, len(policies))
	headers := []string{"ROLE", "RESOURCE GROUP", "PERMISSION"}

	for i, policy := range policies {
		data[i] = []string{
			policy.Role,
			policy.ResourceGroup,
			string(policy.Permission),
		}
	}

	printTable(data, headers)
}

func getUserById(ctx context.Context, userId string) *client.UserResponse {
	c := getClient()
	params := &client.GetUsersUserParams{
		UserId: &userId,
	}

	resp, err := c.GetUsersUserWithResponse(ctx, params)

	handleResponse(resp, err, "")

	return resp.JSON200
}

// getUsersByIds fetches multiple users concurrently by their IDs
func getUsersByIds(
	ctx context.Context,
	userIds []string,
) []*client.UserResponse {
	if len(userIds) == 0 {
		return []*client.UserResponse{}
	}

	type userResult struct {
		user  *client.UserResponse
		index int
	}

	results := make(chan userResult, len(userIds))
	var wg sync.WaitGroup

	// Launch concurrent requests
	for i, userId := range userIds {
		wg.Add(1)
		go func(id string, idx int) {
			defer wg.Done()
			user := getUserById(ctx, id)
			results <- userResult{user: user, index: idx}
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
		users[result.index] = result.user
	}

	// Return error if any requests failed
	if len(errs) > 0 {
		errMsg := fmt.Sprintf("failed to fetch %d user(s):", len(errs))
		for i, err := range errs {
			errMsg += fmt.Sprintf("\n  [%d] %v", i+1, err)
		}

		log.Fatal().Msg(errMsg)
	}

	return users
}

func getUserByName(ctx context.Context, username string) *client.UserResponse {
	c := getClient()
	params := &client.GetUsersUserParams{
		Name: &username,
	}

	resp, err := c.GetUsersUserWithResponse(ctx, params)

	handleResponse(resp, err, "")

	return resp.JSON200
}

func printArtifact(artifact *client.Artifact) {
	fqn := fmt.Sprintf(
		"%s/%s/%s",
		artifact.Fqn.Source,
		artifact.Fqn.Author,
		artifact.Fqn.Name,
	)
	tags := strings.Join(artifact.Tags, "\n")

	data := []string{
		fqn,
		artifact.VersionHash,
		tags,
		artifact.CreatedAt.Format("2006-01-02 15:04:05"),
		strconv.Itoa(artifact.Pulls),
	}

	printTable([][]string{data}, []string{
		"FQN",
		"HASH",
		"TAGS",
		"CREATED",
		"PULLS",
	})
}
