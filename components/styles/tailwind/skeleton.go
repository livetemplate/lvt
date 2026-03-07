package tailwind

import "github.com/livetemplate/lvt/components/styles"

func skeletonStyles() styles.SkeletonStyles {
	return styles.SkeletonStyles{
		Base:          "bg-gray-200",
		MultiLineWrap: "space-y-2",
		// Shape classes
		ShapeCircle:  "rounded-full",
		ShapeRounded: "rounded-md",
		// Animation classes
		AnimationPulse: "animate-pulse",
		AnimationWave:  "animate-shimmer",
	}
}

func avatarSkeletonStyles() styles.AvatarSkeletonStyles {
	return styles.AvatarSkeletonStyles{
		Root:   "relative inline-block",
		Avatar: "bg-gray-200 rounded-full animate-pulse",
		Badge:  "absolute bottom-0 right-0 w-3 h-3 bg-gray-300 rounded-full border-2 border-white animate-pulse",
		// Size classes
		SizeSm: "w-8 h-8",
		SizeMd: "w-12 h-12",
		SizeLg: "w-16 h-16",
		SizeXl: "w-24 h-24",
	}
}

func cardSkeletonStyles() styles.CardSkeletonStyles {
	return styles.CardSkeletonStyles{
		Root:          "bg-white rounded-lg shadow overflow-hidden",
		Image:         "bg-gray-200 animate-pulse",
		Body:          "p-4 space-y-3",
		TitleLine:     "h-5 bg-gray-200 rounded animate-pulse w-3/4",
		DescLine:      "h-3 bg-gray-200 rounded animate-pulse",
		DescLineLast:  "h-3 bg-gray-200 rounded animate-pulse w-1/2",
		DescWrapper:   "space-y-2",
		Footer:        "flex items-center gap-4 pt-2",
		FooterAvatar:  "h-8 w-8 bg-gray-200 rounded-full animate-pulse",
		FooterContent: "flex-1 space-y-2",
		FooterLine1:   "h-3 bg-gray-200 rounded animate-pulse w-1/3",
		FooterLine2:   "h-2 bg-gray-200 rounded animate-pulse w-1/4",
	}
}
