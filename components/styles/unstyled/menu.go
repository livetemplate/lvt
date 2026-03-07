package unstyled

import "github.com/livetemplate/lvt/components/styles"

func menuStyles() styles.MenuStyles {
	return styles.MenuStyles{
		Root:            "lvt-menu",
		TriggerBtn:      "lvt-menu__trigger",
		TriggerIcon:     "lvt-menu__trigger-icon",
		TriggerIconOpen: "lvt-menu__trigger-icon--open",
		Panel:           "lvt-menu__panel",
		PanelInner:      "lvt-menu__panel-inner",
		Divider:         "lvt-menu__divider",
		SectionHeader:   "lvt-menu__section-header",
		Item:            "lvt-menu__item",
		ItemDisabled:    "lvt-menu__item--disabled",
		ItemActive:      "lvt-menu__item--active",
		ItemDefault:     "lvt-menu__item--default",
		ItemIcon:        "lvt-menu__item-icon",
		Badge:           "lvt-menu__badge",
		BadgeRed:        "lvt-menu__badge--red",
		BadgeBlue:       "lvt-menu__badge--blue",
		BadgeGreen:      "lvt-menu__badge--green",
		BadgeDefault:    "lvt-menu__badge--default",
		Shortcut:        "lvt-menu__shortcut",
		// Position variants
		PositionBottomLeft:  "lvt-menu--bottom-left",
		PositionBottomRight: "lvt-menu--bottom-right",
		PositionTopLeft:     "lvt-menu--top-left",
		PositionTopRight:    "lvt-menu--top-right",
	}
}
