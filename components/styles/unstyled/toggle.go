package unstyled

import "github.com/livetemplate/lvt/components/styles"

func toggleStyles() styles.ToggleStyles {
	return styles.ToggleStyles{
		UseCustomTrack: false,
		Label:          "lvt-toggle__label",
		LabelDisabled:  "lvt-toggle__label--disabled",
		Input:          "lvt-toggle__input",
		LabelText:      "lvt-toggle__label-text",
		Description:    "lvt-toggle__description",
		Track:          "lvt-toggle__track",
		Knob:           "lvt-toggle__knob",
		// Track size classes
		TrackSm: "lvt-toggle--sm",
		TrackMd: "lvt-toggle--md",
		TrackLg: "lvt-toggle--lg",
		// Knob size classes
		KnobSm: "lvt-toggle__knob--sm",
		KnobMd: "lvt-toggle__knob--md",
		KnobLg: "lvt-toggle__knob--lg",
		// Track color classes
		TrackChecked:           "lvt-toggle__track--checked",
		TrackUnchecked:         "lvt-toggle__track--unchecked",
		TrackCheckedDisabled:   "lvt-toggle__track--checked-disabled",
		TrackUncheckedDisabled: "lvt-toggle__track--unchecked-disabled",
		// Knob translate classes
		KnobUnchecked: "lvt-toggle__knob--unchecked",
		KnobSmChecked: "lvt-toggle__knob--sm-checked",
		KnobMdChecked: "lvt-toggle__knob--md-checked",
		KnobLgChecked: "lvt-toggle__knob--lg-checked",
	}
}

func checkboxStyles() styles.CheckboxStyles {
	return styles.CheckboxStyles{
		UseCustomCheckbox: false,
		Label:             "lvt-checkbox__label",
		LabelDisabled:     "lvt-checkbox__label--disabled",
		Input:             "lvt-checkbox__input",
		CheckboxWrap:      "lvt-checkbox__wrap",
		CheckboxBox:       "lvt-checkbox__box",
		CheckIcon:         "lvt-checkbox__check-icon",
		LabelTextWrap:     "lvt-checkbox__label-text-wrap",
		LabelText:         "lvt-checkbox__label-text",
		Description:       "lvt-checkbox__description",
		// State classes
		StateChecked:           "lvt-checkbox--checked",
		StateUnchecked:         "lvt-checkbox--unchecked",
		StateCheckedDisabled:   "lvt-checkbox--checked-disabled",
		StateUncheckedDisabled: "lvt-checkbox--unchecked-disabled",
	}
}
