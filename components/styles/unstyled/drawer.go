package unstyled

import "github.com/livetemplate/lvt/components/styles"

func drawerStyles() styles.DrawerStyles {
	return styles.DrawerStyles{
		Root:      "lvt-drawer",
		Overlay:   "lvt-drawer__overlay",
		PanelBase: "lvt-drawer__panel",
		Header:    "lvt-drawer__header",
		Title:     "lvt-drawer__title",
		CloseBtn:  "lvt-drawer__close-btn",
		CloseIcon: "lvt-drawer__close-icon",
		Content:   "lvt-drawer__content",
		// Position classes
		PositionLeft:   "lvt-drawer--left",
		PositionRight:  "lvt-drawer--right",
		PositionTop:    "lvt-drawer--top",
		PositionBottom: "lvt-drawer--bottom",
		// Horizontal size classes
		SizeSmH:   "lvt-drawer--sm-h",
		SizeMdH:   "lvt-drawer--md-h",
		SizeLgH:   "lvt-drawer--lg-h",
		SizeXlH:   "lvt-drawer--xl-h",
		SizeFullH: "lvt-drawer--full-h",
		// Vertical size classes
		SizeSmV:   "lvt-drawer--sm-v",
		SizeMdV:   "lvt-drawer--md-v",
		SizeLgV:   "lvt-drawer--lg-v",
		SizeXlV:   "lvt-drawer--xl-v",
		SizeFullV: "lvt-drawer--full-v",
		// Transform classes
		TransformOpen:   "lvt-drawer--open",
		TransformLeft:   "lvt-drawer--closed-left",
		TransformRight:  "lvt-drawer--closed-right",
		TransformTop:    "lvt-drawer--closed-top",
		TransformBottom: "lvt-drawer--closed-bottom",
	}
}
