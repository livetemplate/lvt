package generator

import (
	"fmt"
	"text/template"
)

// CSSHelpers returns template functions for CSS framework selection
// Only Tailwind and None are supported - Bulma and Pico have been removed for simplification
func CSSHelpers() template.FuncMap {
	return template.FuncMap{
		// CDN link for framework
		"csscdn": func(framework string) string {
			switch framework {
			case "tailwind":
				return `<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>`
			case "none":
				return ""
			default:
				return `<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>` // default to tailwind
			}
		},

		// Container classes
		"containerClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "max-w-7xl mx-auto px-4 py-8"
			case "none":
				return ""
			default:
				return "max-w-7xl mx-auto px-4 py-8"
			}
		},

		// Section wrapper classes
		"sectionClass": func(framework string) string {
			return ""
		},

		// Box/Card classes
		"boxClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "bg-white shadow rounded-lg p-6 mb-6"
			case "none":
				return ""
			default:
				return "bg-white shadow rounded-lg p-6 mb-6"
			}
		},

		// Title classes
		"titleClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "text-3xl font-bold text-gray-900 mb-6"
			case "none":
				return ""
			default:
				return "text-3xl font-bold text-gray-900 mb-6"
			}
		},

		// Subtitle classes
		"subtitleClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "text-xl font-semibold text-gray-700 mb-4"
			case "none":
				return ""
			default:
				return "text-xl font-semibold text-gray-700 mb-4"
			}
		},

		// Field wrapper classes
		"fieldClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "mb-4"
			case "none":
				return ""
			default:
				return "mb-4"
			}
		},

		// Label classes
		"labelClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "block text-sm font-medium text-gray-700 mb-2"
			case "none":
				return ""
			default:
				return "block text-sm font-medium text-gray-700 mb-2"
			}
		},

		// Control wrapper (no longer needed without Bulma)
		"controlClass": func(framework string) string {
			return ""
		},

		// Input classes
		"inputClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
			case "none":
				return ""
			default:
				return "w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
			}
		},

		// Input error classes
		"inputErrorClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "border-red-500 focus:ring-red-500"
			case "none":
				return ""
			default:
				return "border-red-500 focus:ring-red-500"
			}
		},

		// Checkbox wrapper classes
		"checkboxClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "flex items-center"
			case "none":
				return ""
			default:
				return "flex items-center"
			}
		},

		// Button classes
		"buttonClass": func(framework string, variant string) string {
			switch framework {
			case "tailwind":
				if variant == "primary" {
					return "bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50"
				} else if variant == "secondary" {
					return "bg-gray-600 text-white px-2 py-1 text-sm rounded hover:bg-gray-700"
				}
				return "bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 disabled:opacity-50"
			case "none":
				return ""
			default:
				if variant == "primary" {
					return "bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 disabled:opacity-50"
				} else if variant == "secondary" {
					return "bg-gray-600 text-white px-2 py-1 text-sm rounded hover:bg-gray-700"
				}
				return "bg-red-600 text-white px-4 py-2 rounded-md hover:bg-red-700 disabled:opacity-50"
			}
		},

		// Table classes
		"tableClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "min-w-full divide-y divide-gray-200"
			case "none":
				return ""
			default:
				return "min-w-full divide-y divide-gray-200"
			}
		},

		// Table header classes
		"theadClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "bg-gray-50"
			default:
				return ""
			}
		},

		// Table header cell classes
		"thClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
			default:
				return ""
			}
		},

		// Table body classes
		"tbodyClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "bg-white divide-y divide-gray-200"
			default:
				return ""
			}
		},

		// Table row classes
		"trClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "hover:bg-gray-50"
			default:
				return ""
			}
		},

		// Table cell classes
		"tdClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "px-6 py-4 whitespace-nowrap text-sm text-gray-900"
			default:
				return ""
			}
		},

		// Select wrapper classes (no longer needed without Bulma)
		"selectWrapperClass": func(framework string) string {
			return ""
		},

		// Pagination wrapper classes
		"paginationClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "flex justify-between items-center mt-4"
			default:
				return ""
			}
		},

		// Help text classes
		"helpTextClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "text-sm text-red-600 mt-1"
			default:
				return ""
			}
		},

		// Check if framework needs semantic wrapper (no longer needed - only Pico used this)
		"needsWrapper": func(framework string) bool {
			return false
		},

		// Check if framework needs article tags (no longer needed - only Pico used this)
		"needsArticle": func(framework string) bool {
			return false
		},

		// Select element styling
		"selectClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
			default:
				return ""
			}
		},

		// Error state styling
		"errorClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "border-red-500"
			default:
				return ""
			}
		},

		// Table wrapper for overflow
		"needsTableWrapper": func(framework string) bool {
			return framework == "tailwind"
		},

		"tableWrapperClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "overflow-x-auto"
			default:
				return ""
			}
		},

		// Pagination button styling
		"paginationButtonClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "px-4 py-2 border border-gray-300 rounded hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
			default:
				return ""
			}
		},

		// Pagination info container
		"paginationInfoClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "flex items-center justify-center"
			default:
				return ""
			}
		},

		// Pagination current page indicator
		"paginationCurrentClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "px-4 py-2"
			default:
				return ""
			}
		},

		// Text/paragraph classes
		"textClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "text-gray-700"
			case "none":
				return ""
			default:
				return "text-gray-700"
			}
		},

		// Code block classes
		"codeClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "bg-gray-100 text-gray-800 rounded p-4"
			case "none":
				return ""
			default:
				return "bg-gray-100 text-gray-800 rounded p-4"
			}
		},

		// List classes
		"listClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "space-y-2"
			case "none":
				return ""
			default:
				return "space-y-2"
			}
		},

		// List item classes
		"listItemClass": func(framework string) string {
			return ""
		},

		// Link classes
		"linkClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "text-blue-600 hover:text-blue-800 underline"
			case "none":
				return ""
			default:
				return "text-blue-600 hover:text-blue-800 underline"
			}
		},

		// dict creates a map for passing multiple values to nested templates
		// Usage: {{template "formField" (dict "Field" . "CSS" $.CSSFramework)}}
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			if len(values)%2 != 0 {
				return nil, fmt.Errorf("dict requires even number of arguments (key-value pairs)")
			}
			m := make(map[string]interface{}, len(values)/2)
			for i := 0; i < len(values); i += 2 {
				key, ok := values[i].(string)
				if !ok {
					return nil, fmt.Errorf("dict keys must be strings")
				}
				m[key] = values[i+1]
			}
			return m, nil
		},

		// Loading indicator styling
		"loadingClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "text-gray-600 animate-pulse"
			default:
				return ""
			}
		},

		// Pagination active page indicator
		"paginationActiveClass": func(framework string) string {
			switch framework {
			case "tailwind":
				return "bg-blue-600 text-white px-3 py-1 rounded"
			default:
				return ""
			}
		},

		// Helper functions for numbered pagination
		"until": func(count int) []int {
			result := make([]int, count)
			for i := 0; i < count; i++ {
				result[i] = i
			}
			return result
		},

		"add": func(a, b int) int {
			return a + b
		},
	}
}
