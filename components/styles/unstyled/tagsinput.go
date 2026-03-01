package unstyled

import "github.com/livetemplate/lvt/components/styles"

func tagsInputStyles() styles.TagsInputStyles {
	return styles.TagsInputStyles{
		Root:          "lvt-tags-input",
		Wrapper:       "lvt-tags-input__wrapper",
		Tag:           "lvt-tags-input__tag",
		TagRemoveBtn:  "lvt-tags-input__tag-remove",
		TagRemoveIcon: "lvt-tags-input__tag-remove-icon",
		Input:         "lvt-tags-input__input",
		Dropdown:      "lvt-tags-input__dropdown",
		Suggestion:    "lvt-tags-input__suggestion",
		Counter:       "lvt-tags-input__counter",
	}
}
