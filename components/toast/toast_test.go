package toast

import (
	"html/template"
	"strings"
	"testing"

	_ "github.com/livetemplate/lvt/components/styles/tailwind"
	_ "github.com/livetemplate/lvt/components/styles/unstyled"
)

func TestNew(t *testing.T) {
	c := New("notifications")

	if c.ID() != "notifications" {
		t.Errorf("expected ID 'notifications', got '%s'", c.ID())
	}
	if c.Position != TopRight {
		t.Errorf("expected Position TopRight, got %s", c.Position)
	}
	if c.Count() != 0 {
		t.Error("expected 0 messages initially")
	}
}

func TestNewWithOptions(t *testing.T) {
	c := New("test",
		WithPosition(BottomLeft),
		WithMaxVisible(5),
		WithStyled(false),
	)

	if c.Position != BottomLeft {
		t.Errorf("expected Position BottomLeft, got %s", c.Position)
	}
	if c.MaxVisible != 5 {
		t.Errorf("expected MaxVisible 5, got %d", c.MaxVisible)
	}
	if c.IsStyled() {
		t.Error("expected IsStyled to be false")
	}
}

func TestContainer_Add(t *testing.T) {
	c := New("test")

	c.Add(Message{Title: "Test", Body: "Message"})

	if c.Count() != 1 {
		t.Error("expected 1 message after Add")
	}
	if c.Messages[0].ID == "" {
		t.Error("expected auto-generated ID")
	}
}

func TestContainer_AddWithMaxVisible(t *testing.T) {
	c := New("test", WithMaxVisible(2))

	c.Add(Message{Title: "1"})
	c.Add(Message{Title: "2"})
	c.Add(Message{Title: "3"})

	if c.Count() != 2 {
		t.Errorf("expected 2 messages (MaxVisible), got %d", c.Count())
	}
	// Oldest should be removed
	if c.Messages[0].Title != "2" {
		t.Error("expected oldest message to be removed")
	}
}

func TestContainer_AddInfo(t *testing.T) {
	c := New("test")
	c.AddInfo("Info", "Info message")

	if c.Count() != 1 {
		t.Error("expected 1 message")
	}
	if c.Messages[0].Type != Info {
		t.Error("expected Info type")
	}
	if c.Messages[0].AutoDismissMS != DefaultAutoDismissMS {
		t.Errorf("expected AutoDismissMS=%d for info, got %d", DefaultAutoDismissMS, c.Messages[0].AutoDismissMS)
	}
}

func TestContainer_AddSuccess(t *testing.T) {
	c := New("test")
	c.AddSuccess("Success", "Success message")

	if c.Messages[0].Type != Success {
		t.Error("expected Success type")
	}
	if c.Messages[0].AutoDismissMS != DefaultAutoDismissMS {
		t.Errorf("expected AutoDismissMS=%d for success, got %d", DefaultAutoDismissMS, c.Messages[0].AutoDismissMS)
	}
}

func TestContainer_AddWarning(t *testing.T) {
	c := New("test")
	c.AddWarning("Warning", "Warning message")

	if c.Messages[0].Type != Warning {
		t.Error("expected Warning type")
	}
	if c.Messages[0].AutoDismissMS != 0 {
		t.Errorf("expected AutoDismissMS=0 for warning, got %d", c.Messages[0].AutoDismissMS)
	}
}

func TestContainer_AddError(t *testing.T) {
	c := New("test")
	c.AddError("Error", "Error message")

	if c.Messages[0].Type != Error {
		t.Error("expected Error type")
	}
	if c.Messages[0].AutoDismissMS != 0 {
		t.Errorf("expected AutoDismissMS=0 for error, got %d", c.Messages[0].AutoDismissMS)
	}
}

func TestContainer_Dismiss(t *testing.T) {
	c := New("test")
	c.Add(Message{ID: "msg1", Title: "1"})
	c.Add(Message{ID: "msg2", Title: "2"})

	c.Dismiss("msg1")

	if c.Count() != 1 {
		t.Error("expected 1 message after dismiss")
	}
	if c.Messages[0].ID != "msg2" {
		t.Error("expected msg2 to remain")
	}
}

func TestContainer_DismissNonExistent(t *testing.T) {
	c := New("test")
	c.Add(Message{ID: "msg1", Title: "1"})

	c.Dismiss("nonexistent")

	if c.Count() != 1 {
		t.Error("expected messages to remain unchanged")
	}
}

func TestContainer_DismissAll(t *testing.T) {
	c := New("test")
	c.Add(Message{Title: "1"})
	c.Add(Message{Title: "2"})

	c.DismissAll()

	if c.Count() != 0 {
		t.Error("expected 0 messages after DismissAll")
	}
}

func TestContainer_HasMessages(t *testing.T) {
	c := New("test")

	if c.HasMessages() {
		t.Error("expected HasMessages to be false initially")
	}

	c.Add(Message{Title: "Test"})

	if !c.HasMessages() {
		t.Error("expected HasMessages to be true after Add")
	}
}

func TestContainer_VisibleMessages(t *testing.T) {
	c := New("test", WithMaxVisible(2))

	c.Add(Message{ID: "1", Title: "1"})
	c.Add(Message{ID: "2", Title: "2"})
	c.Add(Message{ID: "3", Title: "3"})

	visible := c.VisibleMessages()
	if len(visible) != 2 {
		t.Errorf("expected 2 visible messages, got %d", len(visible))
	}
}

func TestContainer_VisibleMessagesNoLimit(t *testing.T) {
	c := New("test") // MaxVisible = 0 (unlimited)

	c.Add(Message{Title: "1"})
	c.Add(Message{Title: "2"})
	c.Add(Message{Title: "3"})

	visible := c.VisibleMessages()
	if len(visible) != 3 {
		t.Errorf("expected 3 visible messages with no limit, got %d", len(visible))
	}
}

func TestContainer_GetPositionClasses(t *testing.T) {
	tests := []struct {
		pos      Position
		expected string
	}{
		{TopRight, "top-4 right-4"},
		{TopLeft, "top-4 left-4"},
		{TopCenter, "top-4 left-1/2 -translate-x-1/2"},
		{BottomRight, "bottom-4 right-4"},
		{BottomLeft, "bottom-4 left-4"},
		{BottomCenter, "bottom-4 left-1/2 -translate-x-1/2"},
	}

	for _, tc := range tests {
		c := New("test", WithPosition(tc.pos), WithStyled(true))
		if c.GetPositionClasses() != tc.expected {
			t.Errorf("Position %s: expected '%s', got '%s'", tc.pos, tc.expected, c.GetPositionClasses())
		}
	}
}

func TestGetTypeClasses(t *testing.T) {
	tests := []struct {
		typ      Type
		contains string
	}{
		{Info, "blue"},
		{Success, "green"},
		{Warning, "yellow"},
		{Error, "red"},
	}

	c := New("test", WithStyled(true))
	for _, tc := range tests {
		classes := c.GetTypeClasses(tc.typ)
		if !strings.Contains(classes, tc.contains) {
			t.Errorf("Type %s: expected classes to contain '%s', got '%s'", tc.typ, tc.contains, classes)
		}
	}
}

func TestGetTypeClasses_StringInput(t *testing.T) {
	// After JSON round-trip, Type values may arrive as plain strings.
	// GetTypeClasses must handle both Type and string inputs.
	tests := []struct {
		input    interface{}
		contains string
	}{
		{"info", "blue"},
		{"success", "green"},
		{"warning", "yellow"},
		{"error", "red"},
		{42, "blue"}, // unknown type defaults to Info
	}

	c := New("test", WithStyled(true))
	for _, tc := range tests {
		classes := c.GetTypeClasses(tc.input)
		if !strings.Contains(classes, tc.contains) {
			t.Errorf("Input %v: expected classes to contain '%s', got '%s'", tc.input, tc.contains, classes)
		}
	}
}

func TestGetTypeIcon(t *testing.T) {
	tests := []Type{Info, Success, Warning, Error}

	for _, typ := range tests {
		icon := GetTypeIcon(typ)
		if icon == "" {
			t.Errorf("Type %s: expected non-empty icon", typ)
		}
		if !strings.Contains(icon, "svg") {
			t.Errorf("Type %s: expected SVG icon", typ)
		}
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage(
		WithTitle("Test Title"),
		WithBody("Test Body"),
		WithType(Success),
		WithDismissible(false),
		WithIcon("icon"),
	)

	if msg.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got '%s'", msg.Title)
	}
	if msg.Body != "Test Body" {
		t.Errorf("expected body 'Test Body', got '%s'", msg.Body)
	}
	if msg.Type != Success {
		t.Errorf("expected type Success, got %s", msg.Type)
	}
	if msg.Dismissible {
		t.Error("expected Dismissible to be false")
	}
	if msg.Icon != "icon" {
		t.Errorf("expected icon 'icon', got '%s'", msg.Icon)
	}
}

func TestNewMessageDefaults(t *testing.T) {
	msg := NewMessage()

	if msg.Type != Info {
		t.Errorf("expected default type Info, got %s", msg.Type)
	}
	if !msg.Dismissible {
		t.Error("expected default Dismissible to be true")
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()

	if ts == nil {
		t.Fatal("expected non-nil TemplateSet")
	}
	if ts.Namespace != "toast" {
		t.Errorf("expected namespace 'toast', got '%s'", ts.Namespace)
	}
}

func TestTemplateRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	c := New("notifications", WithStyled(true))
	c.Add(Message{ID: "msg1", Title: "Success!", Body: "Your changes have been saved.", Type: Success, Dismissible: true})

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:toast:container:v1", c)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	if !strings.Contains(html, `data-toast-trigger="notifications"`) {
		t.Error("expected data-toast-trigger attribute")
	}
	if !strings.Contains(html, `hidden`) {
		t.Error("expected hidden attribute on trigger span")
	}
	if !strings.Contains(html, `aria-hidden="true"`) {
		t.Error("expected aria-hidden attribute on trigger span")
	}
	// Pending messages should be serialized in data-pending attribute
	if !strings.Contains(html, `data-pending='[`) {
		t.Error("expected data-pending attribute with JSON messages")
	}
	// html/template escapes quotes inside attributes as &#34;
	if !strings.Contains(html, `&#34;title&#34;:&#34;Success!&#34;`) {
		t.Errorf("expected toast title in data-pending JSON, got:\n%s", html)
	}
}

func TestUnstyledTemplateRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	c := New("test", WithStyled(false))
	c.Add(Message{Title: "Test", Body: "Body"})

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:toast:container:v1", c)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	// Trigger-only template is style-agnostic; both styled and unstyled
	// render the same hidden span with data-toast-trigger.
	if !strings.Contains(html, `data-toast-trigger="test"`) {
		t.Error("expected data-toast-trigger attribute")
	}
	if !strings.Contains(html, `hidden`) {
		t.Error("expected hidden attribute on trigger span")
	}
	if !strings.Contains(html, `data-pending='[`) {
		t.Error("expected data-pending attribute with JSON messages")
	}
}

func TestAutoDismissRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	c := New("test", WithStyled(true))
	c.AddSuccess("Done", "Item created")

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:toast:container:v1", c)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	// Auto-dismiss metadata is now embedded in the JSON payload inside data-pending.
	// html/template escapes quotes as &#34; inside attribute values.
	// The success toast should carry dismissMS in the JSON.
	if !strings.Contains(html, `&#34;dismissMS&#34;:5000`) {
		t.Errorf("expected dismissMS field in data-pending JSON for success toast, got:\n%s", html)
	}

	// Error toast should have dismissMS=0 (no auto-dismiss)
	c2 := New("test2", WithStyled(true))
	c2.AddError("Oops", "Something went wrong")

	buf.Reset()
	err = tmpl.ExecuteTemplate(&buf, "lvt:toast:container:v1", c2)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html2 := buf.String()
	if !strings.Contains(html2, `&#34;dismissMS&#34;:0`) {
		t.Errorf("expected dismissMS:0 in data-pending JSON for error toast (no auto-dismiss), got:\n%s", html2)
	}
}

func TestWithAutoDismissOption(t *testing.T) {
	msg := NewMessage(
		WithTitle("Custom"),
		WithBody("Custom dismiss"),
		WithAutoDismiss(3000),
	)
	if msg.AutoDismissMS != 3000 {
		t.Errorf("expected AutoDismissMS=3000, got %d", msg.AutoDismissMS)
	}

	// Negative values should be clamped to 0
	msg2 := NewMessage(WithAutoDismiss(-100))
	if msg2.AutoDismissMS != 0 {
		t.Errorf("expected negative duration clamped to 0, got %d", msg2.AutoDismissMS)
	}
}

func TestEmptyContainerRendering(t *testing.T) {
	ts := Templates()
	tmpl, err := template.New("test").ParseFS(ts.FS, ts.Pattern)
	if err != nil {
		t.Fatalf("failed to parse templates: %v", err)
	}

	c := New("test", WithStyled(true))

	var buf strings.Builder
	err = tmpl.ExecuteTemplate(&buf, "lvt:toast:container:v1", c)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	html := buf.String()
	// Should render trigger span even when empty
	if !strings.Contains(html, `data-toast-trigger="test"`) {
		t.Error("expected empty container to render trigger span")
	}
	// No data-pending when no messages
	if strings.Contains(html, `data-pending`) {
		t.Error("expected no data-pending attribute when container is empty")
	}
}

func TestTakePendingJSON_DrainAndCache(t *testing.T) {
	c := New("test")
	c.AddInfo("Hello", "World")

	// First call: drains messages and returns JSON.
	json1 := c.TakePendingJSON()
	if json1 == "" {
		t.Fatal("first call should return non-empty JSON")
	}
	if !strings.Contains(json1, `"title":"Hello"`) {
		t.Errorf("expected title in JSON, got %s", json1)
	}
	// Messages should be drained after first call.
	if c.Count() != 0 {
		t.Errorf("expected 0 messages after drain, got %d", c.Count())
	}

	// Second call (same render cycle): returns cached JSON (idempotent).
	json2 := c.TakePendingJSON()
	if json2 != json1 {
		t.Errorf("second call should return cached JSON;\n  got:  %s\n  want: %s", json2, json1)
	}

	// Third call: cache cleared, returns empty string.
	json3 := c.TakePendingJSON()
	if json3 != "" {
		t.Errorf("third call should return empty string, got %s", json3)
	}
}

func TestTakePendingJSON_NoMessages(t *testing.T) {
	c := New("test")

	result := c.TakePendingJSON()
	if result != "" {
		t.Errorf("expected empty string when no messages, got %s", result)
	}
}

func TestTakePendingJSON_AddAfterDrain(t *testing.T) {
	c := New("test")
	c.AddSuccess("First", "First message")

	// Drain
	json1 := c.TakePendingJSON()
	if json1 == "" {
		t.Fatal("expected non-empty JSON from first drain")
	}

	// Consume cache
	_ = c.TakePendingJSON()
	// Clear cache
	_ = c.TakePendingJSON()

	// Add new message after drain
	c.AddError("Second", "Second message")

	json4 := c.TakePendingJSON()
	if json4 == "" {
		t.Fatal("expected non-empty JSON after adding new message")
	}
	if !strings.Contains(json4, `"title":"Second"`) {
		t.Errorf("expected new message in JSON, got %s", json4)
	}
	// Should not contain the first message
	if strings.Contains(json4, `"title":"First"`) {
		t.Error("expected first message to be absent after drain")
	}
}
