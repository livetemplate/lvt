package unstyled

import "github.com/livetemplate/lvt/components/styles"

func accordionStyles() styles.AccordionStyles {
	return styles.AccordionStyles{
		Root:         "lvt-accordion",
		Item:         "lvt-accordion__item",
		ItemDisabled: "lvt-accordion__item--disabled",
		Header:       "lvt-accordion__header",
		HeaderIcon:   "lvt-accordion__header-icon",
		ChevronIcon:  "lvt-accordion__chevron",
		ChevronOpen:  "lvt-accordion__chevron--open",
		Content:      "lvt-accordion__content",
	}
}
