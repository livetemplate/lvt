package unstyled

import "github.com/livetemplate/lvt/components/styles"

func autocompleteStyles() styles.AutocompleteStyles {
	return styles.AutocompleteStyles{
		Root:           "lvt-autocomplete",
		InputWrapper:   "lvt-autocomplete__input-wrapper",
		Input:          "lvt-autocomplete__input",
		InputLoading:   "lvt-autocomplete__input--loading",
		LoadingWrapper: "lvt-autocomplete__loading-wrapper",
		LoadingIcon:    "lvt-autocomplete__loading-icon",
		ClearBtn:       "lvt-autocomplete__clear-btn",
		ClearIcon:      "lvt-autocomplete__clear-icon",
		Dropdown:       "lvt-autocomplete__dropdown",
		Option:         "lvt-autocomplete__option",
		OptionDisabled: "lvt-autocomplete__option--disabled",
		OptionActive:   "lvt-autocomplete__option--active",
		OptionLabel:    "lvt-autocomplete__option-label",
		OptionDesc:     "lvt-autocomplete__option-desc",
		OptionDescAlt:  "lvt-autocomplete__option-desc--alt",
		OptionIcon:     "lvt-autocomplete__option-icon",
		OptionLayout:   "lvt-autocomplete__option-layout",
		Empty:          "lvt-autocomplete__empty",
	}
}
