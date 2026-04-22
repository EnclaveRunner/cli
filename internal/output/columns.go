package output

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"cli/internal/styles"

	"github.com/EnclaveRunner/sdk-go/enclave"
)

// UserColumns defines table columns for enclave.User.
var UserColumns = []Column{
	{Header: "NAME", Extract: func(r any) string { return r.(enclave.User).Name }},
	{Header: "DISPLAY NAME", Extract: func(r any) string { return r.(enclave.User).DisplayName }},
	{Header: "ROLES", Extract: func(r any) string { return strings.Join(r.(enclave.User).Roles, ", ") }},
}

// RoleColumns defines table columns for enclave.Role.
var RoleColumns = []Column{
	{Header: "NAME", Extract: func(r any) string { return r.(enclave.Role).Name }},
	{Header: "USERS", Extract: func(r any) string {
		users := r.(enclave.Role).Users
		return fmt.Sprintf("%d", len(users))
	}},
	{Header: "USER LIST", Extract: func(r any) string { return strings.Join(r.(enclave.Role).Users, ", ") }},
}

// ResourceGroupColumns defines table columns for enclave.ResourceGroup.
var ResourceGroupColumns = []Column{
	{Header: "NAME", Extract: func(r any) string { return r.(enclave.ResourceGroup).Name }},
	{Header: "ENDPOINTS", Extract: func(r any) string {
		return fmt.Sprintf("%d", len(r.(enclave.ResourceGroup).Endpoints))
	}},
	{Header: "ENDPOINT LIST", Extract: func(r any) string {
		return strings.Join(r.(enclave.ResourceGroup).Endpoints, ", ")
	}},
}

// PolicyColumns defines table columns for enclave.Policy.
var PolicyColumns = []Column{
	{Header: "ROLE", Extract: func(r any) string { return r.(enclave.Policy).Role }},
	{Header: "RESOURCE GROUP", Extract: func(r any) string { return r.(enclave.Policy).ResourceGroup }},
	{Header: "METHOD", Extract: func(r any) string { return string(r.(enclave.Policy).Method) }},
}

// TaskColumns defines table columns for enclave.Task.
var TaskColumns = []Column{
	{Header: "ID", Extract: func(r any) string { return r.(enclave.Task).ID }},
	{Header: "SOURCE", Extract: func(r any) string { return r.(enclave.Task).Source }},
	{Header: "STATE", MinWidth: 14, Extract: func(r any) string {
		return styles.TaskStateBadge(r.(enclave.Task).Status.State)
	}},
	{Header: "RETRIES", Extract: func(r any) string {
		return strconv.Itoa(r.(enclave.Task).Status.Retries)
	}},
	{Header: "LAST ERROR", MinWidth: 20, Extract: func(r any) string {
		e := r.(enclave.Task).Status.LastError
		if len(e) > 40 {
			return e[:40] + "…"
		}
		return e
	}},
	{Header: "NEXT PROCESS", Extract: func(r any) string {
		t := r.(enclave.Task).Status.NextProcessAt
		if t.IsZero() {
			return "-"
		}
		return t.Format(time.RFC3339)
	}},
}

// TaskLogColumns defines table columns for enclave.TaskLog.
var TaskLogColumns = []Column{
	{Header: "TIME", Extract: func(r any) string {
		return r.(enclave.TaskLog).Timestamp.Format("15:04:05.000")
	}},
	{Header: "LEVEL", MinWidth: 7, Extract: func(r any) string { return r.(enclave.TaskLog).Level }},
	{Header: "ISSUER", Extract: func(r any) string { return r.(enclave.TaskLog).Issuer }},
	{Header: "MESSAGE", MinWidth: 30, Extract: func(r any) string { return r.(enclave.TaskLog).Message }},
}

// ArtifactColumns defines table columns for enclave.Artifact.
var ArtifactColumns = []Column{
	{Header: "NAMESPACE", Extract: func(r any) string { return r.(enclave.Artifact).Namespace }},
	{Header: "NAME", Extract: func(r any) string { return r.(enclave.Artifact).Name }},
	{Header: "HASH", MinWidth: 16, Extract: func(r any) string {
		h := r.(enclave.Artifact).VersionHash
		if len(h) > 16 {
			return h[:16]
		}
		return h
	}},
	{Header: "TAGS", Extract: func(r any) string { return strings.Join(r.(enclave.Artifact).Tags, ", ") }},
	{Header: "CREATED", Extract: func(r any) string {
		return r.(enclave.Artifact).CreatedAt.Format("2006-01-02 15:04")
	}},
	{Header: "PULLS", Extract: func(r any) string {
		return fmt.Sprintf("%d", r.(enclave.Artifact).Pulls)
	}},
}
