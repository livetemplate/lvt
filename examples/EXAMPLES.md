# Component Usage Examples

This document shows usage examples for all 21 LiveTemplate components.

## Form Controls

### Dropdown

```go
import "github.com/livetemplate/components/dropdown"

// Basic dropdown
countrySelect := dropdown.New("country",
    dropdown.WithPlaceholder("Select a country"),
    dropdown.WithOptions([]dropdown.Option{
        {Value: "us", Label: "United States"},
        {Value: "ca", Label: "Canada"},
        {Value: "uk", Label: "United Kingdom"},
    }),
)

// Searchable dropdown
searchSelect := dropdown.New("search",
    dropdown.WithPlaceholder("Search..."),
    dropdown.WithSearchable(true),
    dropdown.WithOptions(countries),
)

// Multi-select dropdown
tagsSelect := dropdown.New("tags",
    dropdown.WithPlaceholder("Select tags..."),
    dropdown.WithMultiple(true),
    dropdown.WithOptions(tags),
)
```

Template:
```html
{{template "lvt:dropdown:default:v1" .CountrySelect}}
{{template "lvt:dropdown:searchable:v1" .SearchSelect}}
{{template "lvt:dropdown:multi:v1" .TagsSelect}}
```

### Autocomplete

```go
import "github.com/livetemplate/components/autocomplete"

userSearch := autocomplete.New("user-search",
    autocomplete.WithPlaceholder("Search users..."),
    autocomplete.WithMinChars(2),
    autocomplete.WithDebounceMs(300),
)
```

Template:
```html
{{template "lvt:autocomplete:default:v1" .UserSearch}}
```

### Date Picker

```go
import "github.com/livetemplate/components/datepicker"

// Single date
birthDate := datepicker.New("birth-date",
    datepicker.WithPlaceholder("Select date"),
    datepicker.WithFormat("2006-01-02"),
)

// Date range
dateRange := datepicker.New("date-range",
    datepicker.WithPlaceholder("Select range"),
    datepicker.WithRange(true),
)

// Inline calendar
inlineCalendar := datepicker.New("calendar",
    datepicker.WithInline(true),
)
```

Template:
```html
{{template "lvt:datepicker:single:v1" .BirthDate}}
{{template "lvt:datepicker:range:v1" .DateRange}}
{{template "lvt:datepicker:inline:v1" .InlineCalendar}}
```

### Time Picker

```go
import "github.com/livetemplate/components/timepicker"

meetingTime := timepicker.New("meeting-time",
    timepicker.WithPlaceholder("Select time"),
    timepicker.WithFormat24Hour(false),
    timepicker.WithMinuteStep(15),
)
```

Template:
```html
{{template "lvt:timepicker:default:v1" .MeetingTime}}
```

### Tags Input

```go
import "github.com/livetemplate/components/tagsinput"

skills := tagsinput.New("skills",
    tagsinput.WithPlaceholder("Add skills..."),
    tagsinput.WithTags([]string{"Go", "TypeScript"}),
    tagsinput.WithMaxTags(10),
)
```

Template:
```html
{{template "lvt:tagsinput:default:v1" .Skills}}
```

### Toggle

```go
import "github.com/livetemplate/components/toggle"

darkMode := toggle.New("dark-mode",
    toggle.WithLabel("Dark Mode"),
    toggle.WithChecked(false),
)
```

Template:
```html
{{template "lvt:toggle:default:v1" .DarkMode}}
{{template "lvt:toggle:checkbox:v1" .AcceptTerms}}
```

### Rating

```go
import "github.com/livetemplate/components/rating"

productRating := rating.New("product-rating",
    rating.WithMaxStars(5),
    rating.WithValue(4),
    rating.WithReadOnly(false),
)
```

Template:
```html
{{template "lvt:rating:default:v1" .ProductRating}}
```

## Layout Components

### Tabs

```go
import "github.com/livetemplate/components/tabs"

// Horizontal tabs
horizontalTabs := tabs.New("demo-tabs",
    tabs.WithOrientation(tabs.OrientationHorizontal),
    tabs.WithItems([]*tabs.TabItem{
        tabs.NewTabItem("overview", "Overview", "Overview content..."),
        tabs.NewTabItem("features", "Features", "Features content..."),
        tabs.NewTabItem("pricing", "Pricing", "Pricing content..."),
    }),
)

// Vertical tabs
verticalTabs := tabs.New("settings-tabs",
    tabs.WithOrientation(tabs.OrientationVertical),
    tabs.WithItems(settingsItems),
)

// Pills style
pillTabs := tabs.New("pill-tabs",
    tabs.WithVariant(tabs.VariantPills),
    tabs.WithItems(items),
)
```

Template:
```html
{{template "lvt:tabs:horizontal:v1" .HorizontalTabs}}
{{template "lvt:tabs:vertical:v1" .VerticalTabs}}
{{template "lvt:tabs:pills:v1" .PillTabs}}
```

### Accordion

```go
import "github.com/livetemplate/components/accordion"

faq := accordion.New("faq",
    accordion.WithItems([]*accordion.AccordionItem{
        accordion.NewItem("q1", "What is LiveTemplate?", "LiveTemplate is..."),
        accordion.NewItem("q2", "How do I get started?", "Run go get..."),
    }),
    accordion.WithAllowMultiple(true),
)

// Single-open accordion
singleAccordion := accordion.New("single",
    accordion.WithItems(items),
    accordion.WithAllowMultiple(false),
)
```

Template:
```html
{{template "lvt:accordion:default:v1" .FAQ}}
{{template "lvt:accordion:single:v1" .SingleAccordion}}
```

### Modal

```go
import "github.com/livetemplate/components/modal"

confirmModal := modal.New("confirm-delete",
    modal.WithTitle("Confirm Delete"),
    modal.WithSize(modal.SizeSmall),
)

sheetModal := modal.New("settings",
    modal.WithTitle("Settings"),
    modal.WithVariant(modal.VariantSheet),
)
```

Template:
```html
{{template "lvt:modal:default:v1" .ConfirmModal}}
{{template "lvt:modal:confirm:v1" .DeleteConfirm}}
{{template "lvt:modal:sheet:v1" .SettingsSheet}}
```

### Drawer

```go
import "github.com/livetemplate/components/drawer"

settingsDrawer := drawer.New("settings",
    drawer.WithTitle("Settings"),
    drawer.WithPosition(drawer.PositionRight),
    drawer.WithWidth("400px"),
)
```

Template:
```html
{{template "lvt:drawer:default:v1" .SettingsDrawer}}
```

## Feedback Components

### Toast

```go
import "github.com/livetemplate/components/toast"

toastContainer := toast.NewContainer("notifications",
    toast.WithPosition(toast.PositionTopRight),
    toast.WithMaxToasts(5),
)

// Add toasts programmatically
toastContainer.Add(toast.NewToast("success", "Saved!", toast.TypeSuccess))
toastContainer.Add(toast.NewToast("error", "Failed!", toast.TypeError))
```

Template:
```html
{{template "lvt:toast:container:v1" .ToastContainer}}
```

### Tooltip

```go
import "github.com/livetemplate/components/tooltip"

helpTip := tooltip.New("help-tip",
    tooltip.WithContent("Click here for help"),
    tooltip.WithPosition(tooltip.PositionTop),
)
```

Template:
```html
{{template "lvt:tooltip:default:v1" .HelpTip}}
```

### Popover

```go
import "github.com/livetemplate/components/popover"

userPopover := popover.New("user-info",
    popover.WithTitle("User Details"),
    popover.WithPosition(popover.PositionBottom),
)
```

Template:
```html
{{template "lvt:popover:default:v1" .UserPopover}}
```

### Progress

```go
import "github.com/livetemplate/components/progress"

// Linear progress bar
uploadProgress := progress.New("upload",
    progress.WithValue(65),
    progress.WithMax(100),
    progress.WithShowLabel(true),
)

// Circular progress
circularProgress := progress.New("loading",
    progress.WithValue(75),
    progress.WithVariant(progress.VariantCircular),
)

// Spinner
spinner := progress.New("spinner",
    progress.WithVariant(progress.VariantSpinner),
)
```

Template:
```html
{{template "lvt:progress:default:v1" .UploadProgress}}
{{template "lvt:progress:circular:v1" .CircularProgress}}
{{template "lvt:progress:spinner:v1" .Spinner}}
```

### Skeleton

```go
import "github.com/livetemplate/components/skeleton"

// Card skeleton
cardSkeleton := skeleton.New("card-loading",
    skeleton.WithVariant(skeleton.VariantCard),
)

// Avatar skeleton
avatarSkeleton := skeleton.New("avatar-loading",
    skeleton.WithVariant(skeleton.VariantAvatar),
)
```

Template:
```html
{{template "lvt:skeleton:default:v1" .TextSkeleton}}
{{template "lvt:skeleton:avatar:v1" .AvatarSkeleton}}
{{template "lvt:skeleton:card:v1" .CardSkeleton}}
```

## Data Display Components

### Data Table

```go
import "github.com/livetemplate/components/datatable"

usersTable := datatable.New("users",
    datatable.WithColumns([]datatable.Column{
        {Key: "name", Label: "Name", Sortable: true},
        {Key: "email", Label: "Email", Sortable: true},
        {Key: "role", Label: "Role"},
        {Key: "actions", Label: "Actions"},
    }),
    datatable.WithPageSize(10),
    datatable.WithSearchable(true),
)
```

Template:
```html
{{template "lvt:datatable:default:v1" .UsersTable}}
```

### Timeline

```go
import "github.com/livetemplate/components/timeline"

projectTimeline := timeline.New("project",
    timeline.WithItems([]*timeline.TimelineItem{
        timeline.NewItem("1", "Project Started", "Initial commit", time.Now().AddDate(0, -3, 0)),
        timeline.NewItem("2", "Alpha Release", "First alpha", time.Now().AddDate(0, -2, 0)),
        timeline.NewItem("3", "v1.0 Release", "Production ready!", time.Now()),
    }),
    timeline.WithOrientation(timeline.OrientationVertical),
)
```

Template:
```html
{{template "lvt:timeline:default:v1" .ProjectTimeline}}
```

### Breadcrumbs

```go
import "github.com/livetemplate/components/breadcrumbs"

navBreadcrumbs := breadcrumbs.New("nav",
    breadcrumbs.WithItems([]*breadcrumbs.BreadcrumbItem{
        breadcrumbs.NewItem("home", "Home", "/"),
        breadcrumbs.NewItem("products", "Products", "/products"),
        breadcrumbs.NewItem("details", "Product Details", ""),
    }),
    breadcrumbs.WithSeparator(breadcrumbs.SeparatorChevron),
)
```

Template:
```html
{{template "lvt:breadcrumbs:default:v1" .NavBreadcrumbs}}
```

## Navigation Components

### Menu

```go
import "github.com/livetemplate/components/menu"

sideMenu := menu.New("sidebar",
    menu.WithItems([]*menu.MenuItem{
        menu.NewItem("dashboard", "Dashboard", "/dashboard"),
        menu.NewItem("users", "Users", "/users"),
        menu.NewItem("settings", "Settings", "/settings"),
    }),
)

// Nested menu
nestedMenu := menu.New("nav",
    menu.WithItems([]*menu.MenuItem{
        menu.NewItem("products", "Products", "").WithChildren([]*menu.MenuItem{
            menu.NewItem("all", "All Products", "/products"),
            menu.NewItem("new", "New Product", "/products/new"),
        }),
    }),
)
```

Template:
```html
{{template "lvt:menu:default:v1" .SideMenu}}
{{template "lvt:menu:nested:v1" .NestedMenu}}
```

## Unstyled Variants

All components support an unstyled mode for custom CSS:

```go
// Use WithStyled(false) for any component
dropdown.New("id", dropdown.WithStyled(false))
tabs.New("id", tabs.WithStyled(false))
modal.New("id", modal.WithStyled(false))
```

This renders semantic HTML without Tailwind classes, allowing you to apply your own CSS.
