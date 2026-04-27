package views

// ModalConfirmedMsg is sent when the user confirms a modal dialog.
type ModalConfirmedMsg struct{}

// ModalCancelledMsg is sent when the user cancels a modal dialog.
type ModalCancelledMsg struct{}

// FormSubmittedMsg carries the field values when the user submits a form.
type FormSubmittedMsg struct {
	Values []string
}

// FormCancelledMsg is sent when the user cancels a form.
type FormCancelledMsg struct{}

// FormDeleteUserMsg is sent from the modal to trigger an async user delete.
type FormDeleteUserMsg struct{ Name string }

// UserDeletedMsg is sent after a user delete operation completes.
type UserDeletedMsg struct{ Err error }

// UserCreatedMsg is sent after a user create operation completes.
type UserCreatedMsg struct{ Err error }

// RoleDeletedMsg is sent after a role delete operation completes.
type RoleDeletedMsg struct{ Err error }

// RoleCreatedMsg is sent after a role create operation completes.
type RoleCreatedMsg struct{ Err error }

// ResourceGroupDeletedMsg is sent after a resource group delete operation completes.
type ResourceGroupDeletedMsg struct{ Err error }

// ResourceGroupCreatedMsg is sent after a resource group create operation completes.
type ResourceGroupCreatedMsg struct{ Err error }

// PolicyDeletedMsg is sent after a policy delete operation completes.
type PolicyDeletedMsg struct{ Err error }

// PolicyCreatedMsg is sent after a policy create operation completes.
type PolicyCreatedMsg struct{ Err error }
