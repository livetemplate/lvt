package unstyled

import "github.com/livetemplate/lvt/components/styles"

func skeletonStyles() styles.SkeletonStyles {
	return styles.SkeletonStyles{
		Base:           "lvt-skeleton",
		MultiLineWrap:  "lvt-skeleton__multi-line",
		ShapeCircle:    "lvt-skeleton--circle",
		ShapeRounded:   "lvt-skeleton--rounded",
		AnimationPulse: "lvt-skeleton--pulse",
		AnimationWave:  "lvt-skeleton--wave",
	}
}

func avatarSkeletonStyles() styles.AvatarSkeletonStyles {
	return styles.AvatarSkeletonStyles{
		Root:   "lvt-avatar-skeleton",
		Avatar: "lvt-avatar-skeleton__avatar",
		Badge:  "lvt-avatar-skeleton__badge",
		SizeSm: "lvt-avatar-skeleton--sm",
		SizeMd: "lvt-avatar-skeleton--md",
		SizeLg: "lvt-avatar-skeleton--lg",
		SizeXl: "lvt-avatar-skeleton--xl",
	}
}

func cardSkeletonStyles() styles.CardSkeletonStyles {
	return styles.CardSkeletonStyles{
		Root:          "lvt-card-skeleton",
		Image:         "lvt-card-skeleton__image",
		Body:          "lvt-card-skeleton__body",
		TitleLine:     "lvt-card-skeleton__title",
		DescLine:      "lvt-card-skeleton__desc-line",
		DescLineLast:  "lvt-card-skeleton__desc-line--last",
		DescWrapper:   "lvt-card-skeleton__desc",
		Footer:        "lvt-card-skeleton__footer",
		FooterAvatar:  "lvt-card-skeleton__footer-avatar",
		FooterContent: "lvt-card-skeleton__footer-content",
		FooterLine1:   "lvt-card-skeleton__footer-line1",
		FooterLine2:   "lvt-card-skeleton__footer-line2",
	}
}
