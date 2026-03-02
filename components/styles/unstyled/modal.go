package unstyled

import "github.com/livetemplate/lvt/components/styles"

func modalStyles() styles.ModalStyles {
	return styles.ModalStyles{
		Root:              "lvt-modal",
		Overlay:           "lvt-modal__overlay",
		ContainerCentered: "lvt-modal__container--centered",
		ContainerTop:      "lvt-modal__container--top",
		Panel:             "lvt-modal__panel",
		Header:            "lvt-modal__header",
		Title:             "lvt-modal__title",
		CloseBtn:          "lvt-modal__close-btn",
		CloseIcon:         "lvt-modal__close-icon",
		Body:              "lvt-modal__body",
		BodyScrollable:    "lvt-modal__body--scrollable",
		// Size classes
		SizeSm:   "lvt-modal--sm",
		SizeMd:   "lvt-modal--md",
		SizeLg:   "lvt-modal--lg",
		SizeXl:   "lvt-modal--xl",
		SizeFull: "lvt-modal--full",
	}
}

func confirmModalStyles() styles.ConfirmModalStyles {
	return styles.ConfirmModalStyles{
		ShowIconCircle:         false,
		Root:                   "lvt-confirm-modal",
		Overlay:                "lvt-confirm-modal__overlay",
		Container:              "lvt-confirm-modal__container",
		Panel:                  "lvt-confirm-modal__panel",
		IconCircle:             "lvt-confirm-modal__icon-circle",
		IconCircleDestructive:  "lvt-confirm-modal__icon-circle--destructive",
		IconCircleDefault:      "lvt-confirm-modal__icon-circle--default",
		IconSvg:                "lvt-confirm-modal__icon-svg",
		Content:                "lvt-confirm-modal__content",
		Title:                  "lvt-confirm-modal__title",
		Message:                "lvt-confirm-modal__message",
		Actions:                "lvt-confirm-modal__actions",
		CancelBtn:              "lvt-confirm-modal__cancel-btn",
		ConfirmBtnBase:         "lvt-confirm-modal__confirm-btn",
		ConfirmDestructive:     "lvt-confirm-modal__confirm-btn--destructive",
		ConfirmDefault:         "lvt-confirm-modal__confirm-btn--default",
		IconWarning:            "lvt-confirm-modal__icon--warning",
		IconWarningDestructive: "lvt-confirm-modal__icon--warning-destructive",
		IconInfo:               "lvt-confirm-modal__icon--info",
		IconSuccess:            "lvt-confirm-modal__icon--success",
		IconError:              "lvt-confirm-modal__icon--error",
		IconDefault:            "lvt-confirm-modal__icon--default",
	}
}

func sheetStyles() styles.SheetStyles {
	return styles.SheetStyles{
		Root:      "lvt-sheet",
		Overlay:   "lvt-sheet__overlay",
		PanelBase: "lvt-sheet__panel",
		Header:    "lvt-sheet__header",
		Title:     "lvt-sheet__title",
		CloseBtn:  "lvt-sheet__close-btn",
		CloseIcon: "lvt-sheet__close-icon",
		Content:   "lvt-sheet__content",
		// Position classes
		PositionLeft:   "lvt-sheet--left",
		PositionRight:  "lvt-sheet--right",
		PositionTop:    "lvt-sheet--top",
		PositionBottom: "lvt-sheet--bottom",
		// Horizontal size classes
		SizeSmH:   "lvt-sheet--sm-h",
		SizeMdH:   "lvt-sheet--md-h",
		SizeLgH:   "lvt-sheet--lg-h",
		SizeXlH:   "lvt-sheet--xl-h",
		SizeFullH: "lvt-sheet--full-h",
		// Vertical size classes
		SizeSmV:   "lvt-sheet--sm-v",
		SizeMdV:   "lvt-sheet--md-v",
		SizeLgV:   "lvt-sheet--lg-v",
		SizeXlV:   "lvt-sheet--xl-v",
		SizeFullV: "lvt-sheet--full-v",
		// Transform classes
		TransformOpen:   "lvt-sheet--open",
		TransformLeft:   "lvt-sheet--closed-left",
		TransformRight:  "lvt-sheet--closed-right",
		TransformTop:    "lvt-sheet--closed-top",
		TransformBottom: "lvt-sheet--closed-bottom",
	}
}
