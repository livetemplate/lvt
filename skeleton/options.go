package skeleton

// Option is a functional option for configuring skeletons.
type Option func(*Skeleton)

// WithWidth sets the skeleton width.
func WithWidth(width string) Option {
	return func(s *Skeleton) {
		s.Width = width
	}
}

// WithHeight sets the skeleton height.
func WithHeight(height string) Option {
	return func(s *Skeleton) {
		s.Height = height
	}
}

// WithShape sets the skeleton shape.
func WithShape(shape Shape) Option {
	return func(s *Skeleton) {
		s.Shape = shape
	}
}

// WithAnimation sets the animation type.
func WithAnimation(animation Animation) Option {
	return func(s *Skeleton) {
		s.Animation = animation
	}
}

// WithLines sets the number of lines for multi-line skeleton.
func WithLines(lines int) Option {
	return func(s *Skeleton) {
		s.Lines = lines
	}
}

// WithLineHeight sets the line height for multi-line spacing.
func WithLineHeight(height string) Option {
	return func(s *Skeleton) {
		s.LineHeight = height
	}
}

// WithStyled enables Tailwind CSS styling.
func WithStyled(styled bool) Option {
	return func(s *Skeleton) {
		s.SetStyled(styled)
	}
}

// AvatarOption is a functional option for avatar skeletons.
type AvatarOption func(*AvatarSkeleton)

// WithAvatarSize sets the avatar size.
func WithAvatarSize(size string) AvatarOption {
	return func(a *AvatarSkeleton) {
		a.Size = size
	}
}

// WithAvatarBadge shows or hides the status badge.
func WithAvatarBadge(show bool) AvatarOption {
	return func(a *AvatarSkeleton) {
		a.ShowBadge = show
	}
}

// WithAvatarStyled enables Tailwind CSS styling.
func WithAvatarStyled(styled bool) AvatarOption {
	return func(a *AvatarSkeleton) {
		a.SetStyled(styled)
	}
}

// CardOption is a functional option for card skeletons.
type CardOption func(*CardSkeleton)

// WithCardImage shows or hides the image placeholder.
func WithCardImage(show bool) CardOption {
	return func(c *CardSkeleton) {
		c.ShowImage = show
	}
}

// WithCardImageHeight sets the image height.
func WithCardImageHeight(height string) CardOption {
	return func(c *CardSkeleton) {
		c.ImageHeight = height
	}
}

// WithCardTitle shows or hides the title placeholder.
func WithCardTitle(show bool) CardOption {
	return func(c *CardSkeleton) {
		c.ShowTitle = show
	}
}

// WithCardDescription shows or hides description lines.
func WithCardDescription(show bool, lines int) CardOption {
	return func(c *CardSkeleton) {
		c.ShowDescription = show
		c.DescriptionLines = lines
	}
}

// WithCardFooter shows or hides the footer area.
func WithCardFooter(show bool) CardOption {
	return func(c *CardSkeleton) {
		c.ShowFooter = show
	}
}

// WithCardStyled enables Tailwind CSS styling.
func WithCardStyled(styled bool) CardOption {
	return func(c *CardSkeleton) {
		c.SetStyled(styled)
	}
}
