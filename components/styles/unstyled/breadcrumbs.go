package unstyled

import "github.com/livetemplate/lvt/components/styles"

func breadcrumbsStyles() styles.BreadcrumbsStyles {
	return styles.BreadcrumbsStyles{
		Nav:      "lvt-breadcrumbs",
		List:     "lvt-breadcrumbs__list",
		ListItem: "lvt-breadcrumbs__list-item",
		HomeLink: "lvt-breadcrumbs__home-link",
		HomeIcon: "lvt-breadcrumbs__home-icon",
		SrOnly:   "lvt-sr-only",
		SepText:  "lvt-breadcrumbs__separator",
		SepIcon:  "lvt-breadcrumbs__separator-icon",
		Ellipsis: "lvt-breadcrumbs__ellipsis",
		ItemIcon: "lvt-breadcrumbs__item-icon",
		SizeSm:   "lvt-breadcrumbs--sm",
		SizeMd:   "lvt-breadcrumbs--md",
		SizeLg:   "lvt-breadcrumbs--lg",
	}
}

func breadcrumbItemStyles() styles.BreadcrumbItemStyles {
	return styles.BreadcrumbItemStyles{
		Current:  "lvt-breadcrumbs__item--current",
		Disabled: "lvt-breadcrumbs__item--disabled",
		Link:     "lvt-breadcrumbs__item-link",
	}
}
