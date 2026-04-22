package output

import (
	"cli/internal/styles"
	"strconv"
	"strings"
	"time"

	"github.com/EnclaveRunner/sdk-go/enclave"
)

// UserColumns defines table columns for enclave.User.
var UserColumns = []Column{
	{
		Header: "NAME",
		Extract: func(r any) string {
			u, _ := r.(enclave.User)

			return u.Name
		},
	},
	{
		Header: "DISPLAY NAME",
		Extract: func(r any) string {
			u, _ := r.(enclave.User)

			return u.DisplayName
		},
	},
	{
		Header: "ROLES",
		Extract: func(r any) string {
			u, _ := r.(enclave.User)

			return strings.Join(u.Roles, ", ")
		},
	},
}

// RoleColumns defines table columns for enclave.Role.
var RoleColumns = []Column{
	{
		Header: "NAME",
		Extract: func(r any) string {
			role, _ := r.(enclave.Role)

			return role.Name
		},
	},
	{Header: "USERS", Extract: func(r any) string {
		role, _ := r.(enclave.Role)

		return strconv.Itoa(len(role.Users))
	}},
	{
		Header: "USER LIST",
		Extract: func(r any) string {
			role, _ := r.(enclave.Role)

			return strings.Join(role.Users, ", ")
		},
	},
}

// ResourceGroupColumns defines table columns for enclave.ResourceGroup.
var ResourceGroupColumns = []Column{
	{
		Header: "NAME",
		Extract: func(r any) string {
			rg, _ := r.(enclave.ResourceGroup)

			return rg.Name
		},
	},
	{Header: "ENDPOINTS", Extract: func(r any) string {
		rg, _ := r.(enclave.ResourceGroup)

		return strconv.Itoa(len(rg.Endpoints))
	}},
	{Header: "ENDPOINT LIST", Extract: func(r any) string {
		rg, _ := r.(enclave.ResourceGroup)

		return strings.Join(rg.Endpoints, ", ")
	}},
}

// PolicyColumns defines table columns for enclave.Policy.
var PolicyColumns = []Column{
	{
		Header: "ROLE",
		Extract: func(r any) string {
			p, _ := r.(enclave.Policy)

			return p.Role
		},
	},
	{
		Header: "RESOURCE GROUP",
		Extract: func(r any) string {
			p, _ := r.(enclave.Policy)

			return p.ResourceGroup
		},
	},
	{
		Header: "METHOD",
		Extract: func(r any) string {
			p, _ := r.(enclave.Policy)

			return string(p.Method)
		},
	},
}

// TaskColumns defines table columns for enclave.Task.
var TaskColumns = []Column{
	{Header: "ID", Extract: func(r any) string {
		t, _ := r.(enclave.Task)

		return t.ID
	}},
	{
		Header: "SOURCE",
		Extract: func(r any) string {
			t, _ := r.(enclave.Task)

			return t.Source
		},
	},
	{Header: "STATE", MinWidth: 14, Extract: func(r any) string {
		t, _ := r.(enclave.Task)

		return styles.TaskStateBadge(t.Status.State)
	}},
	{Header: "RETRIES", Extract: func(r any) string {
		t, _ := r.(enclave.Task)

		return strconv.Itoa(t.Status.Retries)
	}},
	{Header: "LAST ERROR", MinWidth: 20, Extract: func(r any) string {
		t, _ := r.(enclave.Task)
		e := t.Status.LastError
		if len(e) > 40 {
			return e[:40] + "…"
		}

		return e
	}},
	{Header: "NEXT PROCESS", Extract: func(r any) string {
		t, _ := r.(enclave.Task)
		if t.Status.NextProcessAt.IsZero() {
			return "-"
		}

		return t.Status.NextProcessAt.Format(time.RFC3339)
	}},
}

// TaskLogColumns defines table columns for enclave.TaskLog.
var TaskLogColumns = []Column{
	{Header: "TIME", Extract: func(r any) string {
		l, _ := r.(enclave.TaskLog)

		return l.Timestamp.Format("15:04:05.000")
	}},
	{
		Header:   "LEVEL",
		MinWidth: 7,
		Extract: func(r any) string {
			l, _ := r.(enclave.TaskLog)

			return l.Level
		},
	},
	{
		Header: "ISSUER",
		Extract: func(r any) string {
			l, _ := r.(enclave.TaskLog)

			return l.Issuer
		},
	},
	{
		Header:   "MESSAGE",
		MinWidth: 30,
		Extract: func(r any) string {
			l, _ := r.(enclave.TaskLog)

			return l.Message
		},
	},
}

// ArtifactColumns defines table columns for enclave.Artifact.
var ArtifactColumns = []Column{
	{
		Header: "NAMESPACE",
		Extract: func(r any) string {
			a, _ := r.(enclave.Artifact)

			return a.Namespace
		},
	},
	{
		Header: "NAME",
		Extract: func(r any) string {
			a, _ := r.(enclave.Artifact)

			return a.Name
		},
	},
	{Header: "HASH", MinWidth: 16, Extract: func(r any) string {
		a, _ := r.(enclave.Artifact)
		h := a.VersionHash
		if len(h) > 16 {
			return h[:16]
		}

		return h
	}},
	{
		Header: "TAGS",
		Extract: func(r any) string {
			a, _ := r.(enclave.Artifact)

			return strings.Join(a.Tags, ", ")
		},
	},
	{Header: "CREATED", Extract: func(r any) string {
		a, _ := r.(enclave.Artifact)

		return a.CreatedAt.Format("2006-01-02 15:04")
	}},
	{Header: "PULLS", Extract: func(r any) string {
		a, _ := r.(enclave.Artifact)

		return strconv.Itoa(a.Pulls)
	}},
}
