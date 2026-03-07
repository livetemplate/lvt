package tailwind

import "github.com/livetemplate/lvt/components/styles"

func toastStyles() styles.ToastStyles {
	return styles.ToastStyles{
		Container:      "fixed z-50 flex flex-col gap-3 pointer-events-none",
		ContainerWidth: "w-full max-w-sm",
		Toast:          "pointer-events-auto w-full bg-white shadow-lg rounded-lg border overflow-hidden",
		ToastInner:     "p-4",
		ToastLayout:    "flex items-start",
		IconWrap:       "flex-shrink-0",
		ContentWrap:    "flex-1",
		ContentIcon:    "ml-3",
		Title:          "text-sm font-medium text-gray-900",
		Body:           "text-sm text-gray-500",
		BodyWithTitle:  "mt-1",
		DismissWrap:    "ml-4 flex-shrink-0 flex",
		DismissBtn:     "inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 rounded-md",
		DismissIcon:    "h-5 w-5",
		// Position classes
		PosTopRight:     "top-4 right-4",
		PosTopLeft:      "top-4 left-4",
		PosTopCenter:    "top-4 left-1/2 -translate-x-1/2",
		PosBottomRight:  "bottom-4 right-4",
		PosBottomLeft:   "bottom-4 left-4",
		PosBottomCenter: "bottom-4 left-1/2 -translate-x-1/2",
		// Type classes
		TypeSuccess: "bg-green-50 border-green-200 text-green-800",
		TypeWarning: "bg-yellow-50 border-yellow-200 text-yellow-800",
		TypeError:   "bg-red-50 border-red-200 text-red-800",
		TypeInfo:    "bg-blue-50 border-blue-200 text-blue-800",
	}
}
