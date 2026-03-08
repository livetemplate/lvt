package generator

// ComponentUsage tracks which UI components are needed by a generated resource.
type ComponentUsage struct {
	UseModal    bool // delete confirmation modal
	UseToast    bool // CRUD feedback notifications
	UseDropdown bool // select field dropdowns
}

// ComputeComponentUsage determines which components a resource needs
// based on its field types.
func ComputeComponentUsage(data ResourceData) ComponentUsage {
	usage := ComponentUsage{
		UseModal: true, // always: delete confirmation
		UseToast: true, // always: CRUD feedback
	}

	for _, f := range data.Fields {
		if f.IsSelect {
			usage.UseDropdown = true
			break
		}
	}

	return usage
}
