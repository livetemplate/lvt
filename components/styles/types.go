package styles

// AccordionStyles defines CSS classes for the accordion component.
type AccordionStyles struct {
	Root         string // outer container: border, divide, rounded
	Item         string // each accordion item wrapper
	ItemDisabled string // disabled item opacity
	Header       string // button: flex, padding, text, hover
	HeaderIcon   string // icon span in header
	ChevronIcon  string // chevron SVG: size, color, transition
	ChevronOpen  string // rotated state for open items
	Content      string // content panel: padding, text
}

// AutocompleteStyles defines CSS classes for the autocomplete component.
type AutocompleteStyles struct {
	Root             string // outer container: relative
	InputWrapper     string // wraps input + loading/clear icons
	Input            string // text input: border, padding, ring
	InputLoading     string // extra padding when loading spinner shown
	LoadingWrapper   string // loading icon container: absolute positioning
	LoadingIcon      string // spinner SVG: size, color, animation
	ClearBtn         string // clear-all button: absolute positioning
	ClearIcon        string // clear icon SVG: size, color, hover
	SelectedBadge    string // selected item pill wrapper: bg, text, padding, rounded
	SelectedBadgeBtn string // remove button inside badge: color, hover
	Dropdown         string // suggestions list: absolute, z-index, border, shadow
	Option           string // suggestion item: padding, cursor
	OptionDisabled   string // disabled option: opacity, cursor
	OptionActive     string // highlighted option: bg, text color
	OptionLabel      string // label text: font-medium
	OptionDesc       string // description: text-sm, color
	OptionDescAlt    string // description when highlighted
	OptionIcon       string // icon in option: margin
	OptionLayout     string // flex layout for option content
	Empty            string // no suggestions message: padding, text
}

// BreadcrumbsStyles defines CSS classes for the breadcrumbs component.
type BreadcrumbsStyles struct {
	Nav          string // nav element: flex
	List         string // ol: flex, items-center, spacing
	ListItem     string // li: flex items-center
	HomeLink     string // home icon link: color, hover
	HomeIcon     string // home SVG: size
	SrOnly       string // screen reader text
	SepText      string // text separator: color, margin
	SepIcon      string // chevron SVG separator: size, color
	Ellipsis     string // collapsed ellipsis: color
	ItemIcon     string // icon before label: margin
	SizeSm       string // text-sm
	SizeMd       string // text-base
	SizeLg       string // text-lg
}

// BreadcrumbItemStyles defines CSS classes for individual breadcrumb items.
type BreadcrumbItemStyles struct {
	Current  string // current page: color, font-weight
	Disabled string // disabled: color, cursor
	Link     string // clickable link: color, hover
}

// DatatableStyles defines CSS classes for the datatable component.
type DatatableStyles struct {
	Root            string // outer container: overflow
	FilterInput     string // search input: padding, border, ring
	FilterWrapper   string // filter container: margin
	LoadingWrapper  string // loading overlay: flex, padding
	LoadingIcon     string // spinner SVG: size, color, animation
	TableWrapper    string // table scroll container: overflow, border, rounded
	Table           string // table element: min-width, divide
	Thead           string // thead: bg-gray
	Th              string // th: padding, text, uppercase
	ThSortable      string // sortable th: cursor, hover
	ThAlignCenter   string // centered th
	ThAlignRight    string // right-aligned th
	ThContent       string // th inner flex: items-center, gap
	SortIcon        string // sort SVG: size
	SortIconIdle    string // unsorted: muted color
	Tbody           string // tbody: bg, divide
	Td              string // td: padding, text-sm
	TdCompact       string // compact row padding
	TdAlignCenter   string // centered td
	TdAlignRight    string // right-aligned td
	RowStriped      string // striped even rows: bg
	RowHover        string // hoverable rows: hover bg
	RowSelected     string // selected row: bg
	Checkbox        string // row/header checkbox: size, border, color, ring
	EmptyCell       string // empty state cell: padding, text, center
	Pagination      string // pagination wrapper: flex, padding, border
	PageInfo        string // page info text: size, color
	PageInfoBold    string // bold counts in page info
	PageBtn         string // pagination button: padding, border, rounded
	PageBtnDisabled string // disabled page button: opacity, cursor
}

// DatepickerStyles defines CSS classes for the datepicker component.
type DatepickerStyles struct {
	Root         string // outer container: relative, inline-block
	TriggerBtn   string // trigger button: padding, border, shadow, hover, ring
	TriggerText  string // flex layout for trigger content
	Placeholder  string // placeholder text color
	TriggerIcon  string // calendar SVG: size, color
	Dropdown     string // calendar dropdown: absolute, z, border, shadow, padding
	Calendar     string // calendar container: width
	NavWrapper   string // month nav: flex, justify-between, margin
	NavBtn       string // prev/next buttons: padding, hover, rounded
	NavIcon      string // nav arrow SVG: size
	MonthLabel   string // month/year text: font, color
	WeekdayGrid  string // weekday header grid: grid-cols-7
	Weekday      string // weekday label: center, text, font, color
	WeekGrid     string // day row grid: grid-cols-7, gap
	DayBtn       string // day button: size, text, rounded, flex, center
	DayOutMonth  string // day outside current month: muted
	DayDisabled  string // disabled day: muted, cursor
	DaySelected  string // selected day: bg, text-white
	DayToday     string // today: border
	DayDefault   string // normal day: text, hover
	HelperText   string // helper text (e.g. "Select start/end date"): text size, color, margin
	Footer       string // footer area: flex, justify-between, margin, border
	TodayBtn     string // "Today" link: text, color, hover
	ClearBtn     string // "Clear" link: text, color, hover
}

// DrawerStyles defines CSS classes for the drawer component.
type DrawerStyles struct {
	Root      string // fixed overlay container: inset, z-index
	Overlay   string // backdrop: fixed, inset, bg-black/50, transition
	PanelBase string // panel shared: fixed, bg-white, shadow, transform, transition
	Header    string // flex header: items-center, justify-between, padding, border
	Title     string // h2: text-lg, font-semibold, color
	CloseBtn  string // close button: padding, color, hover, rounded
	CloseIcon string // close SVG: size
	Content   string // content area: flex-1, overflow-y-auto, padding
	// Position classes (from PositionClass)
	PositionLeft   string
	PositionRight  string
	PositionTop    string
	PositionBottom string
	// Size classes (from SizeClass)
	SizeSmH string // horizontal sm: w-64
	SizeMdH string // horizontal md: w-80
	SizeLgH string // horizontal lg: w-96
	SizeXlH string // horizontal xl: w-[32rem]
	SizeFullH string // horizontal full: w-full
	SizeSmV string // vertical sm: h-48
	SizeMdV string // vertical md: h-64
	SizeLgV string // vertical lg: h-96
	SizeXlV string // vertical xl: h-[32rem]
	SizeFullV string // vertical full: h-full
	// Transform classes (from TransformClass)
	TransformOpen        string // visible: translate-x-0 translate-y-0
	TransformLeft        string // hidden left: -translate-x-full
	TransformRight       string // hidden right: translate-x-full
	TransformTop         string // hidden top: -translate-y-full
	TransformBottom      string // hidden bottom: translate-y-full
}

// DropdownStyles defines CSS classes for the dropdown component.
type DropdownStyles struct {
	Root            string // outer container: relative, inline-block
	TriggerBtn      string // trigger button: width, padding, border, shadow, hover, ring
	SelectedText    string // selected label: truncate
	SearchWrapper   string // inner wrapper for search input + icons: relative positioning
	TriggerIconWrap string // decorative icon wrapper: absolute, inset-y, flex, padding, pointer-events-none
	TriggerIcon     string // chevron SVG: size, color
	ClearBtnWrap    string // clear button wrapper: absolute, inset-y, flex, padding (pointer-events enabled)
	Dropdown        string // options panel: absolute, z, border, shadow, max-h, overflow
	Option          string // option: padding, cursor, hover
	OptionDisabled  string // disabled option: opacity, cursor
	OptionSelected  string // selected option: bg highlight
}

// MenuStyles defines CSS classes for the menu component.
type MenuStyles struct {
	Root           string // outer container: relative, inline-block
	TriggerBtn     string // trigger button: flex, padding, text, border, shadow, hover, ring
	TriggerIcon    string // chevron: size
	TriggerIconOpen string // rotated chevron
	Panel          string // dropdown panel: absolute, z, w, rounded, bg, shadow, ring
	PanelInner     string // inner padding: py
	Divider        string // hr: margin, border
	SectionHeader  string // group label: padding, text, uppercase
	Item           string // menu item: flex, padding, text
	ItemDisabled   string // disabled item: color, cursor
	ItemActive     string // active/highlighted item: bg, color
	ItemDefault    string // default item: color, hover
	ItemIcon       string // item icon: size
	Badge          string // item badge: padding, text, rounded
	BadgeRed       string // red badge variant
	BadgeBlue      string // blue badge variant
	BadgeGreen     string // green badge variant
	BadgeDefault   string // default badge variant
	Shortcut       string // keyboard shortcut: text, color
	// Position variants
	PositionBottomLeft  string
	PositionBottomRight string
	PositionTopLeft     string
	PositionTopRight    string
}

// ModalStyles defines CSS classes for the modal component.
type ModalStyles struct {
	Root             string // fixed overlay container: inset, z-50, overflow
	Overlay          string // backdrop: fixed, inset, bg-black/50, transition
	ContainerCentered string // centered wrapper: flex, min-h-full, items-center, justify-center, p-4
	ContainerTop     string // top-aligned wrapper: flex, min-h-full, items-start, pt-16, justify-center, p-4
	Panel            string // dialog panel: relative, w-full, bg-white, rounded, shadow, transition
	Header           string // header bar: flex, items-center, justify-between, padding, border
	Title            string // h3: text-lg, font-semibold, color
	CloseBtn         string // close button: padding, color, hover, rounded
	CloseIcon        string // close SVG: size
	Body             string // body area: padding
	BodyScrollable   string // scrollable body: padding, max-h, overflow-y-auto
	// Size classes (from SizeClass)
	SizeSm   string // max-w-sm
	SizeMd   string // max-w-lg
	SizeLg   string // max-w-2xl
	SizeXl   string // max-w-4xl
	SizeFull string // max-w-full mx-4
}

// ConfirmModalStyles defines CSS classes for the confirm modal component.
type ConfirmModalStyles struct {
	ShowIconCircle       bool   // structural: whether to render the icon circle
	Root                 string // fixed overlay container
	Overlay              string // backdrop
	Container            string // centering wrapper
	Panel                string // dialog panel: relative, w-full, max-w-md, bg, rounded, shadow, padding
	IconCircle           string // icon circle base: mx-auto, flex, items-center, justify-center, size, rounded-full
	IconCircleDestructive string // destructive bg: bg-red-100
	IconCircleDefault    string // default bg: bg-yellow-100
	IconSvg              string // icon SVG: size
	Content              string // text content wrapper: text-center
	Title                string // h3: text-lg, font-semibold, color, margin
	Message              string // p: text-sm, color
	Actions              string // button container: margin, flex, gap, justify
	CancelBtn            string // cancel button: padding, text, border, rounded, hover
	ConfirmBtnBase       string // confirm button base: padding, text, rounded
	ConfirmDestructive   string // destructive confirm: bg-red, hover, text-white
	ConfirmDefault       string // default confirm: bg-blue, hover, text-white
	// Icon color classes (from IconClass)
	IconWarning          string // text-yellow-500
	IconWarningDestructive string // text-red-500
	IconInfo             string // text-blue-500
	IconSuccess          string // text-green-500
	IconError            string // text-red-500
	IconDefault          string // text-gray-500
}

// SheetStyles defines CSS classes for the sheet modal component.
type SheetStyles struct {
	Root      string // fixed overlay container: inset, z-50
	Overlay   string // backdrop: fixed, inset, bg-black/50, transition
	PanelBase string // panel shared: fixed, bg-white, shadow, transform, transition
	Header    string // header: flex, items-center, justify-between, padding, border
	Title     string // h3: text-lg, font-semibold, color
	CloseBtn  string // close button: padding, color, hover, rounded
	CloseIcon string // close SVG: size
	Content   string // content area: flex-1, overflow-y-auto, padding
	// Position classes (from PositionClass)
	PositionLeft   string
	PositionRight  string
	PositionTop    string
	PositionBottom string
	// Horizontal size classes
	SizeSmH  string // w-64
	SizeMdH  string // w-80
	SizeLgH  string // w-96
	SizeXlH  string // w-[32rem]
	SizeFullH string // w-full
	// Vertical size classes
	SizeSmV  string // h-48
	SizeMdV  string // h-64
	SizeLgV  string // h-96
	SizeXlV  string // h-[32rem]
	SizeFullV string // h-full
	// Transform classes
	TransformOpen   string // translate-x-0 translate-y-0
	TransformLeft   string // -translate-x-full
	TransformRight  string // translate-x-full
	TransformTop    string // -translate-y-full
	TransformBottom string // translate-y-full
}

// PopoverStyles defines CSS classes for the popover component.
type PopoverStyles struct {
	Root      string // outer container: relative, inline-block
	Panel     string // popover panel: absolute, z, bg, rounded, shadow, border
	Header    string // header area: flex, items-center, justify-between, padding, border
	Title     string // h3: text-sm, font-semibold, color
	CloseBtn  string // close button: padding, color, hover, rounded
	CloseIcon string // close icon SVG: size
	Body      string // body area: padding, text-sm, color
	// Position classes (from PositionClasses) — 12 positions
	PosTop         string
	PosTopStart    string
	PosTopEnd      string
	PosBottom      string
	PosBottomStart string
	PosBottomEnd   string
	PosLeft        string
	PosLeftStart   string
	PosLeftEnd     string
	PosRight       string
	PosRightStart  string
	PosRightEnd    string
	// Arrow classes (from ArrowClasses)
	ArrowTop    string
	ArrowBottom string
	ArrowLeft   string
	ArrowRight  string
}

// ProgressStyles defines CSS classes for the linear progress bar component.
type ProgressStyles struct {
	Root          string // outer container: w-full
	Header        string // label/value wrapper: flex, justify-between, margin
	LabelText     string // label text: text-sm, font-medium, color
	ValueText     string // percentage text: text-sm, font-medium, color
	Track         string // track container: w-full, bg, rounded-full, overflow
	Bar           string // progress bar: rounded-full, transition
	BarIndeterminate string // indeterminate bar: animation
	Striped       string // striped pattern class
	Animated      string // animated stripes class
	// Size classes (from SizeClass)
	SizeXs string // h-1
	SizeSm string // h-2
	SizeMd string // h-4
	SizeLg string // h-6
	// Color classes (from ColorClass)
	ColorPrimary string // bg-blue-500
	ColorSuccess string // bg-green-500
	ColorWarning string // bg-yellow-500
	ColorDanger  string // bg-red-500
	ColorInfo    string // bg-cyan-500
}

// CircularProgressStyles defines CSS classes for the circular progress component.
type CircularProgressStyles struct {
	Root             string // container: relative, inline-flex, items-center, justify-center
	SvgIndeterminate string // spinning SVG: animate-spin
	TrackCircle      string // background circle: color
	BarCircle        string // progress circle: transition
	Label            string // center label: absolute, text-sm, font-medium, color
	// Color classes (from ColorClass — text-* variant)
	ColorPrimary string // text-blue-500
	ColorSuccess string // text-green-500
	ColorWarning string // text-yellow-500
	ColorDanger  string // text-red-500
	ColorInfo    string // text-cyan-500
}

// SpinnerStyles defines CSS classes for the spinner component.
type SpinnerStyles struct {
	Root       string // container: inline-flex, items-center
	Svg        string // spinner SVG: animate-spin
	TrackClass string // background circle: opacity
	BarClass   string // foreground path: opacity
	SrOnly     string // screen reader label: sr-only
	// Size classes (from SizeClass)
	SizeSm string // w-4 h-4
	SizeMd string // w-6 h-6
	SizeLg string // w-8 h-8
	SizeXl string // w-12 h-12
	// Color classes (from ColorClass — text-* variant)
	ColorPrimary string
	ColorSuccess string
	ColorWarning string
	ColorDanger  string
	ColorInfo    string
}

// RatingStyles defines CSS classes for the rating component.
type RatingStyles struct {
	Root          string // outer container: inline-flex, items-center, gap
	Label         string // label text: text-sm, color
	StarsWrapper  string // stars container: inline-flex
	StarBtn       string // star button: cursor, transition
	StarReadonly  string // readonly star: no cursor
	HalfStarOuter string // half star relative wrapper
	HalfStarInner string // half star overflow clip
	ValueText     string // value display: margin, text-sm, font-medium, color
	CountText     string // count display: margin, text-sm, color
	// Size classes (from SizeClass)
	SizeSm string // text-lg
	SizeMd string // text-2xl
	SizeLg string // text-3xl
	SizeXl string // text-4xl
	// Color classes (from ColorClass)
	ColorYellow    string // text-yellow-400
	ColorRed       string // text-red-500
	ColorBlue      string // text-blue-500
	ColorGreen     string // text-green-500
	// Empty color classes (from EmptyColorClass)
	EmptyDefault   string // text-gray-300
	EmptyYellow    string // text-yellow-200
	EmptyRed       string // text-red-200
	EmptyBlue      string // text-blue-200
	EmptyGreen     string // text-green-200
}

// SkeletonStyles defines CSS classes for the skeleton loading component.
type SkeletonStyles struct {
	Base           string // skeleton element: bg-gray-200
	MultiLineWrap  string // multi-line container: space-y
	// Shape classes (from ShapeClass)
	ShapeCircle    string // rounded-full
	ShapeRounded   string // rounded-md
	// Animation classes (from AnimationClass)
	AnimationPulse string // animate-pulse
	AnimationWave  string // animate-shimmer
}

// AvatarSkeletonStyles defines CSS classes for the avatar skeleton component.
type AvatarSkeletonStyles struct {
	Root      string // container: relative, inline-block
	Avatar    string // avatar circle: bg, rounded-full, animate
	Badge     string // status badge: absolute, size, bg, rounded-full, border, animate
	// Size classes (from SizeClass)
	SizeSm string // w-8 h-8
	SizeMd string // w-12 h-12
	SizeLg string // w-16 h-16
	SizeXl string // w-24 h-24
}

// CardSkeletonStyles defines CSS classes for the card skeleton component.
type CardSkeletonStyles struct {
	Root          string // card container: bg, rounded, shadow, overflow
	Image         string // image placeholder: bg, animate
	Body          string // body area: padding, spacing
	TitleLine     string // title skeleton: height, bg, rounded, animate, width
	DescLine      string // description line: height, bg, rounded, animate
	DescLineLast  string // last description line: narrower width
	DescWrapper   string // description lines container: spacing
	Footer        string // footer area: border, padding
	FooterAvatar  string // footer avatar: size, bg, rounded-full, animate
	FooterContent string // footer text area: flex, spacing
	FooterLine1   string // footer line 1: height, bg, rounded, animate
	FooterLine2   string // footer line 2: height, bg, rounded, animate
}

// TabsStyles defines CSS classes for the tabs component.
type TabsStyles struct {
	Root       string // outer container: w-full
	TabList    string // tab list border: border-b
	Nav        string // nav container: flex, negative-margin, spacing
	Tab        string // tab button base: padding, border-b, font, text, whitespace
	TabActive  string // active tab: border-color, text-color
	TabDefault string // inactive tab: border-transparent, text, hover
	TabDisabled string // disabled tab: opacity, cursor
	TabIcon    string // icon in tab: margin
	Badge      string // badge: margin, padding, text, rounded, bg, color
	// Vertical-specific
	VerticalRoot    string // vertical container: flex
	VerticalTabList string // vertical tab list: border-r instead of border-b
	VerticalNav     string // vertical nav: flex-col
	VerticalTab     string // vertical tab: border-r instead of border-b
	// Pills-specific
	PillsNav     string // pills nav: flex, gap
	PillTab      string // pill button base: padding, rounded, font, text
	PillActive   string // active pill: bg, text
	PillDefault  string // inactive pill: text, hover
	PillDisabled string // disabled pill: opacity, cursor
}

// TagsInputStyles defines CSS classes for the tags input component.
type TagsInputStyles struct {
	Root          string // outer container: relative
	Wrapper       string // tags + input wrapper: flex-wrap, padding, border, rounded, ring
	Tag           string // tag: inline-flex, items-center, padding, text, bg, color, rounded
	TagRemoveBtn  string // remove button: flex, size, color, hover, rounded
	TagRemoveIcon string // remove icon SVG: size
	Input         string // text input: flex-1, min-w, outline-none, bg-transparent, text
	Dropdown      string // suggestions: absolute, z, border, shadow, max-h, overflow
	Suggestion    string // suggestion item: padding, text, cursor, hover
	Counter       string // tag counter: margin, text-xs, color
}

// TimelineStyles defines CSS classes for the timeline container.
type TimelineStyles struct {
	VerticalRoot   string // vertical container: relative
	HorizontalRoot string // horizontal container: flex, flex-row, spacing, overflow
	Connector      string // vertical connector line: absolute, left, top, bottom, width, bg
}

// TimelineItemStyles defines CSS classes for timeline items.
type TimelineItemStyles struct {
	VerticalItem      string // vertical item wrapper: relative, padding-left, padding-bottom
	HorizontalItem    string // horizontal item wrapper: flex-shrink-0, flex-col, items-center
	IndicatorVertical string // vertical indicator: absolute, left, size, rounded-full, flex, items-center, justify-center
	IndicatorHoriz    string // horizontal indicator: size, rounded-full, flex, items-center, justify-center, margin-bottom
	IndicatorIcon     string // icon in indicator: text-white, text-sm
	IndicatorDot      string // default dot: size, rounded-full, bg-white
	CheckIcon         string // completed check SVG: size, text-white
	ContentVertical   string // vertical content: flex-1
	ContentHoriz      string // horizontal content: text-center, max-w
	Time              string // timestamp: text, color
	TimeVertical      string // vertical-specific time styling
	TimeHoriz         string // horizontal-specific time styling
	Title             string // h4: text, font-semibold, color
	TitleHoriz        string // horizontal title: smaller text
	Description       string // p: margin, text-sm, color
	HorizConnector    string // horizontal connector line: absolute, top, left, size, bg
	// Indicator color classes (from IndicatorClass)
	ColorGray   string // bg-gray-400
	ColorBlue   string // bg-blue-500
	ColorGreen  string // bg-green-500
	ColorYellow string // bg-yellow-500
	ColorRed    string // bg-red-500
	ColorPurple string // bg-purple-500
	// Status classes (from StatusClass)
	StatusDefault  string // bg-gray-400 text-white
	StatusPending  string // bg-gray-200 text-gray-500
	StatusActive   string // bg-blue-500 text-white ring
	StatusComplete string // bg-green-500 text-white
	StatusError    string // bg-red-500 text-white
	// Ring classes (from RingClass)
	RingBlue   string // ring-4 ring-blue-100
	RingGreen  string // ring-4 ring-green-100
	RingYellow string // ring-4 ring-yellow-100
	RingRed    string // ring-4 ring-red-100
	RingPurple string // ring-4 ring-purple-100
	RingGray   string // ring-4 ring-gray-100
}

// TimepickerStyles defines CSS classes for the timepicker component.
type TimepickerStyles struct {
	Root          string // outer container: relative, inline-block
	TriggerBtn    string // trigger button: padding, border, shadow, hover, ring
	TriggerText   string // trigger content: flex, items-center, justify-between
	Placeholder   string // placeholder text color
	TriggerIcon   string // clock SVG: size, color
	Dropdown      string // picker panel: absolute, z, padding, border, rounded, shadow
	SpinnerLayout string // spinners container: flex, items-center, gap
	SpinnerCol    string // column for each spinner: flex-col, items-center
	SpinnerBtn    string // up/down buttons: padding, hover, rounded
	SpinnerIcon   string // spinner arrow SVG: size
	SpinnerInput  string // numeric display: width, text-center, text-lg, font, border, rounded
	Separator     string // colon separator: text-lg, font-semibold
	PeriodCol     string // AM/PM column: flex-col, items-center, margin
	PeriodActive  string // active period button: bg-blue, text-white
	PeriodDefault string // inactive period button: bg-gray, text, hover
	Footer        string // footer: flex, justify-between, margin, border
	NowBtn        string // "Now" button: text-sm, color, hover
	ClearBtn      string // "Clear" button: text-sm, color, hover
}

// ToastStyles defines CSS classes for the toast container component.
type ToastStyles struct {
	Container    string // fixed container: z, flex-col, gap, pointer-events-none
	ContainerWidth string // max width: w-full, max-w-sm
	Toast        string // individual toast: pointer-events-auto, w-full, bg, shadow, rounded, border, overflow
	ToastInner   string // toast padding: p-4
	ToastLayout  string // flex layout: flex, items-start
	IconWrap     string // icon container: flex-shrink-0
	ContentWrap  string // content area: flex-1
	ContentIcon  string // content when icon present: margin-left
	Title        string // title: text-sm, font-medium, color
	Body         string // body: text-sm, color
	BodyWithTitle string // body with title: margin-top
	DismissWrap  string // dismiss button container: margin, flex-shrink-0, flex
	DismissBtn   string // dismiss button: inline-flex, color, hover, ring, rounded
	DismissIcon  string // dismiss SVG: size
	// Position classes (from GetPositionClasses)
	PosTopRight     string
	PosTopLeft      string
	PosTopCenter    string
	PosBottomRight  string
	PosBottomLeft   string
	PosBottomCenter string
	// Type classes (from GetTypeClasses)
	TypeSuccess string // bg-green-50 border-green-200 text-green-800
	TypeWarning string // bg-yellow-50 border-yellow-200 text-yellow-800
	TypeError   string // bg-red-50 border-red-200 text-red-800
	TypeInfo    string // bg-blue-50 border-blue-200 text-blue-800
}

// ToggleStyles defines CSS classes for the toggle/switch component.
type ToggleStyles struct {
	UseCustomTrack bool   // structural: true = render track/knob, false = native checkbox only
	Label          string // label wrapper: inline-flex, items-center, gap, cursor-pointer
	LabelDisabled  string // disabled label: same + cursor-not-allowed, opacity
	Input          string // hidden input: sr-only peer
	LabelText      string // label text: text-sm, font-medium, color
	Description    string // description text: text-sm, color
	Track          string // track container: relative, rounded-full, transition
	Knob           string // knob circle: absolute, top, bg-white, rounded-full, shadow, transform, transition
	// Track size classes (from SizeClasses)
	TrackSm string // w-8 h-4
	TrackMd string // w-11 h-6
	TrackLg string // w-14 h-8
	// Knob size classes (from KnobSizeClasses)
	KnobSm string // w-3 h-3
	KnobMd string // w-5 h-5
	KnobLg string // w-6 h-6
	// Track color classes (from TrackColorClass)
	TrackChecked            string // bg-blue-600
	TrackUnchecked          string // bg-gray-200
	TrackCheckedDisabled    string // bg-blue-300
	TrackUncheckedDisabled  string // bg-gray-200
	// Knob translate classes (from KnobTranslateClass)
	KnobUnchecked  string // translate-x-0.5
	KnobSmChecked  string // translate-x-4
	KnobMdChecked  string // translate-x-5
	KnobLgChecked  string // translate-x-7
}

// CheckboxStyles defines CSS classes for the checkbox component.
type CheckboxStyles struct {
	UseCustomCheckbox bool   // structural: true = custom styled box, false = native checkbox
	Label             string // label wrapper: inline-flex, items-start, gap, cursor-pointer
	LabelDisabled     string // disabled label: same + cursor-not-allowed, opacity
	Input             string // hidden input: sr-only peer
	CheckboxWrap      string // checkbox container: relative, flex, items-center, justify-center
	CheckboxBox       string // checkbox visual: size, rounded, border, transition, flex, items-center, justify-center
	CheckIcon         string // check SVG: size, color
	LabelTextWrap     string // label text container: flex-1
	LabelText         string // label text: text-sm, font-medium, color
	Description       string // description text: text-sm, color
	// State classes (from CheckboxStateClass)
	StateChecked           string // bg-blue-600 border-blue-600
	StateUnchecked         string // bg-white border-gray-300
	StateCheckedDisabled   string // bg-blue-300 border-blue-300
	StateUncheckedDisabled string // bg-gray-100 border-gray-200
}

// TooltipStyles defines CSS classes for the tooltip component.
type TooltipStyles struct {
	Root  string // outer container: relative, inline-block
	Panel string // tooltip panel: absolute, z, padding, text, color, bg, rounded, shadow, whitespace
	// Position classes (from PositionClasses) — 12 positions
	PosTop         string
	PosTopStart    string
	PosTopEnd      string
	PosBottom      string
	PosBottomStart string
	PosBottomEnd   string
	PosLeft        string
	PosLeftStart   string
	PosLeftEnd     string
	PosRight       string
	PosRightStart  string
	PosRightEnd    string
	// Arrow classes (from ArrowClasses)
	ArrowTop    string
	ArrowBottom string
	ArrowLeft   string
	ArrowRight  string
}
