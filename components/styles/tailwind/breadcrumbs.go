package tailwind

import "github.com/livetemplate/lvt/components/styles"

func breadcrumbsStyles() styles.BreadcrumbsStyles {
	return styles.BreadcrumbsStyles{
		Nav:      "flex",
		List:     "flex items-center space-x-2",
		ListItem: "flex items-center",
		HomeLink: "text-gray-500 hover:text-gray-700",
		HomeIcon: "w-5 h-5",
		SrOnly:   "sr-only",
		SepText:  "text-gray-400 mx-1",
		SepIcon:  "w-4 h-4 text-gray-400",
		Ellipsis: "text-gray-400",
		ItemIcon: "mr-1",
		SizeSm:   "text-sm",
		SizeMd:   "text-base",
		SizeLg:   "text-lg",
	}
}

func breadcrumbItemStyles() styles.BreadcrumbItemStyles {
	return styles.BreadcrumbItemStyles{
		Current:  "text-gray-700 font-medium",
		Disabled: "text-gray-400 cursor-not-allowed",
		Link:     "text-gray-500 hover:text-gray-700",
	}
}
