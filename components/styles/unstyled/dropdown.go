package unstyled

import "github.com/livetemplate/lvt/components/styles"

func dropdownStyles() styles.DropdownStyles {
	return styles.DropdownStyles{
		Root:            "lvt-dropdown",
		TriggerBtn:      "lvt-dropdown__trigger",
		SelectedText:    "lvt-dropdown__selected-text",
		TriggerIconWrap: "lvt-dropdown__trigger-icon-wrap",
		TriggerIcon:     "lvt-dropdown__trigger-icon",
		Dropdown:        "lvt-dropdown__panel",
		Option:          "lvt-dropdown__option",
		OptionDisabled:  "lvt-dropdown__option--disabled",
		OptionSelected:  "lvt-dropdown__option--selected",
	}
}
