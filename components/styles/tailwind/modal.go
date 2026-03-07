package tailwind

import "github.com/livetemplate/lvt/components/styles"

func modalStyles() styles.ModalStyles {
	return styles.ModalStyles{
		Root:              "fixed inset-0 z-50 overflow-y-auto",
		Overlay:           "fixed inset-0 bg-black bg-opacity-50 transition-opacity",
		ContainerCentered: "flex min-h-full items-center justify-center p-4",
		ContainerTop:      "flex min-h-full items-start pt-16 justify-center p-4",
		Panel:             "relative w-full bg-white rounded-lg shadow-xl transform transition-all",
		Header:            "flex items-center justify-between px-6 py-4 border-b border-gray-200",
		Title:             "text-lg font-semibold text-gray-900",
		CloseBtn:          "p-2 text-gray-400 hover:text-gray-600 rounded-md hover:bg-gray-100",
		CloseIcon:         "w-5 h-5",
		Body:              "px-6 py-4",
		BodyScrollable:    "px-6 py-4 max-h-[60vh] overflow-y-auto",
		// Size classes
		SizeSm:   "max-w-sm",
		SizeMd:   "max-w-lg",
		SizeLg:   "max-w-2xl",
		SizeXl:   "max-w-4xl",
		SizeFull: "max-w-full mx-4",
	}
}

func confirmModalStyles() styles.ConfirmModalStyles {
	return styles.ConfirmModalStyles{
		ShowIconCircle:         true,
		Root:                   "fixed inset-0 z-50 overflow-y-auto",
		Overlay:                "fixed inset-0 bg-black bg-opacity-50 transition-opacity",
		Container:              "flex min-h-full items-center justify-center p-4",
		Panel:                  "relative w-full max-w-md bg-white rounded-lg shadow-xl p-6",
		IconCircle:             "mx-auto flex items-center justify-center h-12 w-12 rounded-full mb-4",
		IconCircleDestructive:  "bg-red-100",
		IconCircleDefault:      "bg-yellow-100",
		IconSvg:                "h-6 w-6",
		Content:                "text-center",
		Title:                  "text-lg font-semibold text-gray-900 mb-2",
		Message:                "text-sm text-gray-500",
		Actions:                "mt-6 flex gap-3 justify-center",
		CancelBtn:              "px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50",
		ConfirmBtnBase:         "px-4 py-2 text-sm font-medium rounded-md",
		ConfirmDestructive:     "bg-red-600 hover:bg-red-700 text-white",
		ConfirmDefault:         "bg-blue-600 hover:bg-blue-700 text-white",
		IconWarning:            "text-yellow-500",
		IconWarningDestructive: "text-red-500",
		IconInfo:               "text-blue-500",
		IconSuccess:            "text-green-500",
		IconError:              "text-red-500",
		IconDefault:            "text-gray-500",
	}
}

func sheetStyles() styles.SheetStyles {
	return styles.SheetStyles{
		Root:      "fixed inset-0 z-50",
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
