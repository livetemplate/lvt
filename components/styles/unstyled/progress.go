package unstyled

import "github.com/livetemplate/lvt/components/styles"

func progressStyles() styles.ProgressStyles {
	return styles.ProgressStyles{
		Root:             "lvt-progress",
		Header:           "lvt-progress__header",
		LabelText:        "lvt-progress__label",
		ValueText:        "lvt-progress__value",
		Track:            "lvt-progress__track",
		Bar:              "lvt-progress__bar",
		BarIndeterminate: "lvt-progress__bar--indeterminate",
		Striped:          "lvt-progress__bar--striped",
		Animated:         "lvt-progress__bar--animated",
		// Size classes
		SizeXs: "lvt-progress--xs",
		SizeSm: "lvt-progress--sm",
		SizeMd: "lvt-progress--md",
		SizeLg: "lvt-progress--lg",
		// Color classes
		ColorPrimary: "lvt-progress--primary",
		ColorSuccess: "lvt-progress--success",
		ColorWarning: "lvt-progress--warning",
		ColorDanger:  "lvt-progress--danger",
		ColorInfo:    "lvt-progress--info",
	}
}

func circularProgressStyles() styles.CircularProgressStyles {
	return styles.CircularProgressStyles{
		Root:             "lvt-circular-progress",
		SvgIndeterminate: "lvt-circular-progress__svg--indeterminate",
		TrackCircle:      "lvt-circular-progress__track",
		BarCircle:        "lvt-circular-progress__bar",
		Label:            "lvt-circular-progress__label",
		// Color classes
		ColorPrimary: "lvt-circular-progress--primary",
		ColorSuccess: "lvt-circular-progress--success",
		ColorWarning: "lvt-circular-progress--warning",
		ColorDanger:  "lvt-circular-progress--danger",
		ColorInfo:    "lvt-circular-progress--info",
	}
}

func spinnerStyles() styles.SpinnerStyles {
	return styles.SpinnerStyles{
		Root:       "lvt-spinner",
		Svg:        "lvt-spinner__svg",
		TrackClass: "lvt-spinner__track",
		BarClass:   "lvt-spinner__bar",
		SrOnly:     "lvt-sr-only",
		// Size classes
		SizeSm: "lvt-spinner--sm",
		SizeMd: "lvt-spinner--md",
		SizeLg: "lvt-spinner--lg",
		SizeXl: "lvt-spinner--xl",
		// Color classes
		ColorPrimary: "lvt-spinner--primary",
		ColorSuccess: "lvt-spinner--success",
		ColorWarning: "lvt-spinner--warning",
		ColorDanger:  "lvt-spinner--danger",
		ColorInfo:    "lvt-spinner--info",
	}
}
