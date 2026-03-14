package generator

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// InjectEmbeddedChild modifies the parent handler (.go file) to integrate an embedded child resource.
//
// It performs the following modifications:
//  1. Add import for the child package
//  2. Add *child.EmbeddedState field to the parent state struct
//  3. Add *child.EmbeddedController field to the parent controller struct
//  4. Initialize the child controller and state in the Handler function
//  5. Add child template files to ParseFiles
//  6. Append forwarding methods (ChildAdd, ChildEdit, ChildUpdate, ChildDelete, ChildCancelEdit)
//  7. Hook into View/Mount to load child data
//
// This function is idempotent — it checks for existing markers before injecting.
func InjectEmbeddedChild(parentGoPath string, childData ResourceData) error {
	content, err := os.ReadFile(parentGoPath)
	if err != nil {
		return fmt.Errorf("failed to read parent handler: %w", err)
	}

	src := string(content)

	// Idempotency: check if already injected
	marker := fmt.Sprintf("%s.EmbeddedState", childData.PackageName)
	if strings.Contains(src, marker) {
		return nil // already injected
	}

	titleCaser := cases.Title(language.English)
	childPkg := childData.PackageName
	childSingular := childData.ResourceNameSingular // e.g. "Comment"
	childPlural := childData.ResourceNamePlural     // e.g. "Comments"
	parentSingular := childData.ParentResourceSingular
	parentPkg := childData.ParentPackageName
	parentNameCap := titleCaser.String(parentPkg)

	// 1. Add import for child package
	childImport := fmt.Sprintf("\t\"%s/app/%s\"", childData.ModuleName, childPkg)
	src, err = injectImport(src, childImport)
	if err != nil {
		return fmt.Errorf("failed to inject import: %w", err)
	}

	// 2. Add EmbeddedState field to parent state struct
	stateField := fmt.Sprintf("\t%s *%s.EmbeddedState `json:\"%s\"`", childPlural, childPkg, childPkg)
	src, err = injectStructField(src, parentNameCap+"State", stateField)
	if err != nil {
		return fmt.Errorf("failed to inject state field: %w", err)
	}

	// 3. Add EmbeddedController field to parent controller struct
	ctrlField := fmt.Sprintf("\t%sCtrl *%s.EmbeddedController", childPlural, childPkg)
	src, err = injectStructField(src, parentNameCap+"Controller", ctrlField)
	if err != nil {
		return fmt.Errorf("failed to inject controller field: %w", err)
	}

	// 4. Initialize controller in the Handler function
	// Find: controller := &ParentController{
	//            Queries: queries,
	//        }
	// Add: controller.CommentsCtrl = comments.NewEmbeddedController(queries)
	ctrlInitPattern := fmt.Sprintf(`controller := &%sController{`, parentNameCap)
	ctrlInitIdx := strings.Index(src, ctrlInitPattern)
	if ctrlInitIdx == -1 {
		return fmt.Errorf("could not find controller initialization in parent handler")
	}
	// Find the closing brace of the struct literal
	closingBrace := findClosingBrace(src, ctrlInitIdx)
	if closingBrace == -1 {
		return fmt.Errorf("could not find closing brace of controller initialization")
	}
	ctrlInit := fmt.Sprintf("\n\tcontroller.%sCtrl = %s.NewEmbeddedController(queries)", childPlural, childPkg)
	src = src[:closingBrace+1] + ctrlInit + src[closingBrace+1:]

	// 5. Initialize child state in initial state
	// Find: initialState := &ParentState{
	// Add: Comments: &comments.EmbeddedState{},
	stateInitPattern := fmt.Sprintf(`initialState := &%sState{`, parentNameCap)
	stateInitIdx := strings.Index(src, stateInitPattern)
	if stateInitIdx == -1 {
		return fmt.Errorf("could not find initial state in parent handler")
	}
	stateClosing := findClosingBrace(src, stateInitIdx)
	if stateClosing == -1 {
		return fmt.Errorf("could not find closing brace of initial state")
	}
	// Insert before closing brace
	stateInit := fmt.Sprintf("\n\t\t%s: &%s.EmbeddedState{},", childPlural, childPkg)
	src = src[:stateClosing] + stateInit + "\n\t" + src[stateClosing:]

	// 6. Initialize child state in resourceState (page-mode detail view)
	// Note: The current handler template uses a single shared handler (no resourceState).
	// These steps are retained for backward compatibility with older generated code.
	resourceStatePattern := fmt.Sprintf(`resourceState := &%sState{`, parentNameCap)
	resourceStateIdx := strings.Index(src, resourceStatePattern)
	if resourceStateIdx != -1 {
		resourceStateClosing := findClosingBrace(src, resourceStateIdx)
		if resourceStateClosing != -1 {
			// Detect indentation of closing brace to match field indentation
			indent := detectIndent(src, resourceStateClosing)
			resourceStateInit := fmt.Sprintf("\n%s\t%s: &%s.EmbeddedState{},", indent, childPlural, childPkg)
			src = src[:resourceStateClosing] + resourceStateInit + "\n" + indent + src[resourceStateClosing:]
		}
	}

	// 6b. Load child data in page-mode after the parent item is found
	// Find the comment "// Mount with custom state for this URL" which follows the item-finding loop,
	// and inject child loading before it.
	mountComment := "// Mount with custom state for this URL"
	mountCommentIdx := strings.Index(src, mountComment)
	if mountCommentIdx != -1 && resourceStateIdx != -1 {
		// Detect indentation of the mount comment line to match
		indent := detectIndent(src, mountCommentIdx)
		childLoad := fmt.Sprintf("%sresourceState.%s, _ = controller.%sCtrl.Load(resourceState.%s, context.Background(), resourceID)\n\n%s",
			indent, childPlural, childPlural, childPlural, indent)
		// Find the start of the mount comment line (after previous newline)
		lineStart := strings.LastIndex(src[:mountCommentIdx], "\n")
		if lineStart != -1 {
			src = src[:lineStart] + "\n" + childLoad + src[mountCommentIdx:]
		}
	}

	// 7. Add child template to WithParseFiles inside livetemplate.New()
	// (renumbered from step 6 after adding resourceState injections above)
	// The child template must be loaded during New() so that template flattening
	// can resolve cross-template references (e.g., {{template "comments:section" .Comments}}).
	// A separate ParseFiles() after Must() is too late — flattening happens inside New().
	childTmplPath := fmt.Sprintf(`"app/%s/%s.tmpl"`, childPkg, childPkg)
	parentTmplPath := fmt.Sprintf(`"app/%s/%s.tmpl"`, parentPkg, parentPkg)

	// Check if WithParseFiles already exists (e.g., from a previous child injection)
	if strings.Contains(src, "livetemplate.WithParseFiles(") {
		// Append child template to existing WithParseFiles if not already present
		if !strings.Contains(src, childTmplPath) {
			withParseFilesEnd := strings.Index(src, "livetemplate.WithParseFiles(")
			if withParseFilesEnd != -1 {
				// Find the closing paren of WithParseFiles(...)
				afterWPF := src[withParseFilesEnd:]
				closeParen := strings.Index(afterWPF, ")")
				if closeParen != -1 {
					insertPos := withParseFilesEnd + closeParen
					src = src[:insertPos] + ", " + childTmplPath + src[insertPos:]
				}
			}
		}
	} else {
		// No WithParseFiles yet — add it as an option inside livetemplate.New()
		// Insert before the closing "))" of livetemplate.Must(livetemplate.New(...))
		withParseFiles := fmt.Sprintf("\n\t\tlivetemplate.WithParseFiles(%s, %s),", parentTmplPath, childTmplPath)

		// Find the component templates closing and insert after it
		componentEnd := strings.Index(src, "livetemplate.WithComponentTemplates(")
		if componentEnd != -1 {
			// Find the matching closing ")," for WithComponentTemplates
			afterComp := src[componentEnd:]
			// Find closing pattern: "),\n" after the component templates block
			depth := 0
			insertAfter := -1
			for i := 0; i < len(afterComp); i++ {
				if afterComp[i] == '(' {
					depth++
				} else if afterComp[i] == ')' {
					depth--
					if depth == 0 {
						insertAfter = componentEnd + i + 1
						// Skip the trailing comma if present
						if insertAfter < len(src) && src[insertAfter] == ',' {
							insertAfter++
						}
						break
					}
				}
			}
			if insertAfter != -1 {
				src = src[:insertAfter] + withParseFiles + src[insertAfter:]
			}
		} else {
			// Fallback: insert before the closing "))" of Must(New(...))
			mustClose := strings.Index(src, "\t))")
			if mustClose != -1 {
				src = src[:mustClose] + withParseFiles + "\n" + src[mustClose:]
			}
		}
	}

	// Remove any standalone ParseFiles call that was previously added
	// (it would fail because flattening happens inside New())
	standaloneParse := fmt.Sprintf("\tif _, err := baseTmpl.ParseFiles(%s, %s); err != nil {\n\t\tlog.Fatalf(\"Failed to parse template: %%v\", err)\n\t}\n", parentTmplPath, childTmplPath)
	src = strings.Replace(src, standaloneParse, "", 1)
	// Also remove the original single-file ParseFiles if present
	singleParse := fmt.Sprintf("\tif _, err := baseTmpl.ParseFiles(%s); err != nil {\n\t\tlog.Fatalf(\"Failed to parse template: %%v\", err)\n\t}\n", parentTmplPath)
	src = strings.Replace(src, singleParse, "", 1)

	// Remove unused "log" import if no longer referenced
	if !strings.Contains(src, "log.") {
		src = strings.Replace(src, "\t\"log\"\n", "", 1)
	}

	// 7. Hook into View method to load child data
	src = injectChildLoadIntoView(src, parentNameCap, parentSingular, childPlural, childPkg)

	// 8. Hook into Mount method to load child data
	src = injectChildLoadIntoMount(src, parentNameCap, childPlural, childPkg)

	// 9. Append forwarding methods
	methods := generateForwardingMethods(parentNameCap, parentSingular, childSingular, childPlural, childPkg)
	src += "\n" + methods

	return os.WriteFile(parentGoPath, []byte(src), 0644)
}

// InjectEmbeddedChildTemplate modifies the parent template (.tmpl file) to include the child section.
//
// It finds the detail page section and inserts {{template "child:section" .Child}} before the closing {{end}}.
// This function is idempotent.
func InjectEmbeddedChildTemplate(parentTmplPath string, childData ResourceData) error {
	content, err := os.ReadFile(parentTmplPath)
	if err != nil {
		return fmt.Errorf("failed to read parent template: %w", err)
	}

	src := string(content)

	// Idempotency check
	childSection := fmt.Sprintf(`{{template "%s:section"`, childData.PackageName)
	if strings.Contains(src, childSection) {
		return nil // already injected
	}

	childPlural := childData.ResourceNamePlural
	childPkg := childData.PackageName

	// Find "<!-- Detail Content -->" marker in the detail page section
	detailMarker := "<!-- Detail Content -->"
	detailIdx := strings.Index(src, detailMarker)
	if detailIdx == -1 {
		// Fallback: find the detailPage define block's last {{end}} before closing
		// Look for {{define "detailPage"}} and find a good insertion point
		detailDefine := `{{define "detailPage"}}`
		detailDefIdx := strings.Index(src, detailDefine)
		if detailDefIdx == -1 {
			return fmt.Errorf("could not find detail page section in parent template; expected '<!-- Detail Content -->' marker or {{define \"detailPage\"}}")
		}

		// Find the second-to-last {{end}} in this block (before the final closing {{end}})
		// Insert before the view-mode's closing section
		// As a simple heuristic, find "{{end}}" that closes the detailPage define
		// and insert before it
		afterDefine := src[detailDefIdx:]
		lastEndIdx := strings.LastIndex(afterDefine, "{{end}}")
		if lastEndIdx == -1 {
			return fmt.Errorf("could not find closing {{end}} in detailPage block")
		}

		insertPos := detailDefIdx + lastEndIdx
		injection := fmt.Sprintf("\n  {{template \"%s:section\" .%s}}\n", childPkg, childPlural)
		src = src[:insertPos] + injection + src[insertPos:]
	} else {
		// Insert after the detail content section — find the end of the detail fields div
		// Look for the closing </div> after the detail marker that ends the max-width div
		afterMarker := src[detailIdx:]
		// Find the pattern: </div> followed by whitespace/newlines then {{end}}
		// This is the closing of the detail view's field display section
		closingDivPattern := regexp.MustCompile(`</div>\s*\n\s*\{\{end\}\}`)
		loc := closingDivPattern.FindStringIndex(afterMarker)
		if loc == nil {
			// Fallback: insert before the last {{end}} in the block
			lastEndIdx := strings.LastIndex(afterMarker, "{{end}}")
			if lastEndIdx == -1 {
				return fmt.Errorf("could not find insertion point for child section in parent template")
			}
			insertPos := detailIdx + lastEndIdx
			injection := fmt.Sprintf("\n  {{template \"%s:section\" .%s}}\n", childPkg, childPlural)
			src = src[:insertPos] + injection + src[insertPos:]
		} else {
			// Insert after the </div> but before the {{end}}
			insertPos := detailIdx + loc[0] + len("</div>")
			injection := fmt.Sprintf("\n\n  {{template \"%s:section\" .%s}}", childPkg, childPlural)
			src = src[:insertPos] + injection + src[insertPos:]
		}
	}

	return os.WriteFile(parentTmplPath, []byte(src), 0644)
}

// injectImport adds an import line to the import block if not already present.
func injectImport(src, importLine string) (string, error) {
	if strings.Contains(src, importLine) {
		return src, nil
	}

	// Find import block closing paren
	importStart := strings.Index(src, "import (")
	if importStart == -1 {
		return src, fmt.Errorf("could not find import block")
	}
	importEndRel := strings.Index(src[importStart:], "\n)")
	if importEndRel == -1 {
		return src, fmt.Errorf("could not find end of import block")
	}
	insertPos := importStart + importEndRel
	src = src[:insertPos] + "\n" + importLine + src[insertPos:]
	return src, nil
}

// injectStructField adds a field to a named struct if not already present.
func injectStructField(src, structName, fieldLine string) (string, error) {
	if strings.Contains(src, fieldLine) {
		return src, nil
	}

	// Find: type StructName struct {
	pattern := fmt.Sprintf("type %s struct {", structName)
	idx := strings.Index(src, pattern)
	if idx == -1 {
		return src, fmt.Errorf("could not find struct %s", structName)
	}

	// Find the closing brace
	closing := findClosingBrace(src, idx)
	if closing == -1 {
		return src, fmt.Errorf("could not find closing brace of struct %s", structName)
	}

	// Insert field before closing brace
	src = src[:closing] + fieldLine + "\n" + src[closing:]
	return src, nil
}

// detectIndent returns the whitespace indentation for the line containing position idx.
func detectIndent(src string, idx int) string {
	// Walk backwards to find the start of the line
	lineStart := strings.LastIndex(src[:idx], "\n")
	if lineStart == -1 {
		lineStart = 0
	} else {
		lineStart++ // skip the newline itself
	}
	// Extract leading whitespace
	var indent strings.Builder
	for i := lineStart; i < idx; i++ {
		if src[i] == '\t' || src[i] == ' ' {
			indent.WriteByte(src[i])
		} else {
			break
		}
	}
	return indent.String()
}

// findClosingBrace finds the matching } for the first { at or after startIdx.
func findClosingBrace(src string, startIdx int) int {
	depth := 0
	started := false
	for i := startIdx; i < len(src); i++ {
		if src[i] == '{' {
			depth++
			started = true
		} else if src[i] == '}' {
			depth--
			if started && depth == 0 {
				return i
			}
		}
	}
	return -1
}

// injectChildLoadIntoView adds child data loading to the View method.
func injectChildLoadIntoView(src, parentNameCap, parentSingular, childPlural, childPkg string) string {
	// Find the View method and add child loading after the existing item loading
	// Pattern: state.EditingParent = &itemCopy ... break
	// Add: state.ChildPlural, _ = c.ChildPluralCtrl.Load(state.ChildPlural, context.Background(), state.EditingID)
	viewMethod := fmt.Sprintf("func (c *%sController) View(", parentNameCap)
	viewIdx := strings.Index(src, viewMethod)
	if viewIdx == -1 {
		return src
	}

	// Find "state.LastUpdated = formatTime()" within the View method
	afterView := src[viewIdx:]
	lastUpdatedIdx := strings.Index(afterView, "state.LastUpdated = formatTime()")
	if lastUpdatedIdx == -1 {
		return src
	}

	insertPos := viewIdx + lastUpdatedIdx
	loadCode := fmt.Sprintf("\n\tif state.EditingID != \"\" {\n\t\tstate.%s, _ = c.%sCtrl.Load(state.%s, context.Background(), state.EditingID)\n\t}\n\n\t",
		childPlural, childPlural, childPlural)
	src = src[:insertPos] + loadCode + src[insertPos:]

	return src
}

// injectChildLoadIntoMount adds child data loading to the Mount method.
// For page-mode resources, Mount has a detail-view branch (when _resource_id is set)
// that needs to load child data alongside the parent item.
func injectChildLoadIntoMount(src, parentNameCap, childPlural, childPkg string) string {
	mountMethod := fmt.Sprintf("func (c *%sController) Mount(", parentNameCap)
	mountIdx := strings.Index(src, mountMethod)
	if mountIdx == -1 {
		return src
	}

	// Look for the page-mode detail branch: resourceID := ctx.GetString("_resource_id")
	// followed by "return state, nil"
	afterMount := src[mountIdx:]
	detailBranch := `resourceID := ctx.GetString("_resource_id")`
	detailIdx := strings.Index(afterMount, detailBranch)
	if detailIdx == -1 {
		// No page-mode detail branch in Mount — nothing to inject
		return src
	}

	// Find "return state, nil" within this detail branch
	afterDetail := afterMount[detailIdx:]
	returnIdx := strings.Index(afterDetail, "return state, nil")
	if returnIdx == -1 {
		return src
	}

	// Calculate absolute position and inject child loading before the return
	absPos := mountIdx + detailIdx + returnIdx
	indent := detectIndent(src, absPos)
	childLoad := fmt.Sprintf("state.%s, _ = c.%sCtrl.Load(state.%s, context.Background(), resourceID)\n%s",
		childPlural, childPlural, childPlural, indent)
	src = src[:absPos] + childLoad + src[absPos:]

	return src
}

// generateForwardingMethods creates the controller methods that forward to the embedded child controller.
func generateForwardingMethods(parentNameCap, parentSingular, childSingular, childPlural, childPkg string) string {
	parentState := parentNameCap + "State"

	var b strings.Builder

	// ChildAdd
	b.WriteString(fmt.Sprintf(`// %sAdd forwards to the embedded %s controller
func (c *%sController) %sAdd(state %s, ctx *livetemplate.Context) (%s, error) {
	var err error
	state.%s, err = c.%sCtrl.Add(state.%s, ctx, state.EditingID)
	return state, err
}

`, childSingular, childPkg, parentNameCap, childSingular, parentState, parentState,
		childPlural, childPlural, childPlural))

	// ChildEdit
	b.WriteString(fmt.Sprintf(`// %sEdit forwards to the embedded %s controller
func (c *%sController) %sEdit(state %s, ctx *livetemplate.Context) (%s, error) {
	var err error
	state.%s, err = c.%sCtrl.Edit(state.%s, ctx)
	return state, err
}

`, childSingular, childPkg, parentNameCap, childSingular, parentState, parentState,
		childPlural, childPlural, childPlural))

	// ChildUpdate
	b.WriteString(fmt.Sprintf(`// %sUpdate forwards to the embedded %s controller
func (c *%sController) %sUpdate(state %s, ctx *livetemplate.Context) (%s, error) {
	var err error
	state.%s, err = c.%sCtrl.Update(state.%s, ctx, state.EditingID)
	return state, err
}

`, childSingular, childPkg, parentNameCap, childSingular, parentState, parentState,
		childPlural, childPlural, childPlural))

	// ChildDelete
	b.WriteString(fmt.Sprintf(`// %sDelete forwards to the embedded %s controller
func (c *%sController) %sDelete(state %s, ctx *livetemplate.Context) (%s, error) {
	var err error
	state.%s, err = c.%sCtrl.Delete(state.%s, ctx, state.EditingID)
	return state, err
}

`, childSingular, childPkg, parentNameCap, childSingular, parentState, parentState,
		childPlural, childPlural, childPlural))

	// ChildCancelEdit
	b.WriteString(fmt.Sprintf(`// %sCancelEdit forwards to the embedded %s controller
func (c *%sController) %sCancelEdit(state %s, ctx *livetemplate.Context) (%s, error) {
	var err error
	state.%s, err = c.%sCtrl.CancelEdit(state.%s, ctx)
	return state, err
}
`, childSingular, childPkg, parentNameCap, childSingular, parentState, parentState,
		childPlural, childPlural, childPlural))

	return b.String()
}
