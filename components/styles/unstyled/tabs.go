package unstyled

import "github.com/livetemplate/lvt/components/styles"

func tabsStyles() styles.TabsStyles {
	return styles.TabsStyles{
		Root:        "lvt-tabs",
		TabList:     "lvt-tabs__list",
		Nav:         "lvt-tabs__nav",
		Tab:         "lvt-tabs__tab",
		TabActive:   "lvt-tabs__tab--active",
		TabDefault:  "lvt-tabs__tab--default",
		TabDisabled: "lvt-tabs__tab--disabled",
		TabIcon:     "lvt-tabs__tab-icon",
		Badge:       "lvt-tabs__badge",
		// Vertical-specific
		VerticalRoot:    "lvt-tabs--vertical",
		VerticalTabList: "lvt-tabs__list--vertical",
		VerticalNav:     "lvt-tabs__nav--vertical",
		VerticalTab:     "lvt-tabs__tab--vertical",
		// Pills-specific
		PillsNav:     "lvt-tabs__nav--pills",
		PillTab:      "lvt-tabs__pill",
		PillActive:   "lvt-tabs__pill--active",
		PillDefault:  "lvt-tabs__pill--default",
		PillDisabled: "lvt-tabs__pill--disabled",
	}
}
