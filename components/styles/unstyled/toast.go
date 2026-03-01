package unstyled

import "github.com/livetemplate/lvt/components/styles"

func toastStyles() styles.ToastStyles {
	return styles.ToastStyles{
		Container:      "lvt-toast-container",
		ContainerWidth: "lvt-toast-container__width",
		Toast:          "lvt-toast",
		ToastInner:     "lvt-toast__inner",
		ToastLayout:    "lvt-toast__layout",
		IconWrap:       "lvt-toast__icon-wrap",
		ContentWrap:    "lvt-toast__content",
		ContentIcon:    "lvt-toast__content--with-icon",
		Title:          "lvt-toast__title",
		Body:           "lvt-toast__body",
		BodyWithTitle:  "lvt-toast__body--with-title",
		DismissWrap:    "lvt-toast__dismiss-wrap",
		DismissBtn:     "lvt-toast__dismiss-btn",
		DismissIcon:    "lvt-toast__dismiss-icon",
		// Position classes
		PosTopRight:     "lvt-toast-container--top-right",
		PosTopLeft:      "lvt-toast-container--top-left",
		PosTopCenter:    "lvt-toast-container--top-center",
		PosBottomRight:  "lvt-toast-container--bottom-right",
		PosBottomLeft:   "lvt-toast-container--bottom-left",
		PosBottomCenter: "lvt-toast-container--bottom-center",
		// Type classes
		TypeSuccess: "lvt-toast--success",
		TypeWarning: "lvt-toast--warning",
		TypeError:   "lvt-toast--error",
		TypeInfo:    "lvt-toast--info",
	}
}
