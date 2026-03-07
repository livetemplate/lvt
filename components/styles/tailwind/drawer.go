package tailwind

import "github.com/livetemplate/lvt/components/styles"

func drawerStyles() styles.DrawerStyles {
	return styles.DrawerStyles{
		Root:      "fixed inset-0 z-40",
		Overlay:   "fixed inset-0 bg-black bg-opacity-50 transition-opacity",
		PanelBase: "fixed bg-white shadow-xl transform transition-transform duration-300 ease-in-out",
		Header:    "flex items-center justify-between px-4 py-3 border-b border-gray-200",
		Title:     "text-lg font-semibold text-gray-900",
		CloseBtn:  "p-2 text-gray-400 hover:text-gray-600 rounded-md hover:bg-gray-100",
		CloseIcon: "w-5 h-5",
		Content:   "flex-1 overflow-y-auto p-4",
		// Position classes
		PositionLeft:   "left-0 top-0 h-full",
		PositionRight:  "right-0 top-0 h-full",
		PositionTop:    "top-0 left-0 w-full",
		PositionBottom: "bottom-0 left-0 w-full",
		// Horizontal size classes
		SizeSmH:   "w-64",
		SizeMdH:   "w-80",
		SizeLgH:   "w-96",
		SizeXlH:   "w-[32rem]",
		SizeFullH: "w-full",
		// Vertical size classes
		SizeSmV:   "h-48",
		SizeMdV:   "h-64",
		SizeLgV:   "h-96",
		SizeXlV:   "h-[32rem]",
		SizeFullV: "h-full",
		// Transform classes
		TransformOpen:   "translate-x-0 translate-y-0",
		TransformLeft:   "-translate-x-full",
		TransformRight:  "translate-x-full",
		TransformTop:    "-translate-y-full",
		TransformBottom: "translate-y-full",
	}
}
