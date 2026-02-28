package toast

import (
	"html/template"
	"strings"
	"testing"
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
}

func TestContainer_AddSuccess(t *testing.T) {
	c := New("test")
	c.AddSuccess("Success", "Success message")

	if c.Messages[0].Type != Success {
		t.Error("expected Success type")
	}
}

func TestContainer_AddWarning(t *testing.T) {
	c := New("test")
	c.AddWarning("Warning", "Warning message")

	if c.Messages[0].Type != Warning {
		t.Error("expected Warning type")
	}
}

func TestContainer_AddError(t *testing.T) {
	c := New("test")
	c.AddError("Error", "Error message")

	if c.Messages[0].Type != Error {
		t.Error("expected Error type")
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
		c := New("test", WithPosition(tc.pos))
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

	for _, tc := range tests {
		classes := GetTypeClasses(tc.typ)
		if !strings.Contains(classes, tc.contains) {
			t.Errorf("Type %s: expected classes to contain '%s', got '%s'", tc.typ, tc.contains, classes)
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
	if !strings.Contains(html, `data-toast-container="notifications"`) {
		t.Error("expected data-toast-container attribute")
	}
	if !strings.Contains(html, `aria-live="polite"`) {
		t.Error("expected aria-live attribute")
	}
	if !strings.Contains(html, "Success!") {
		t.Error("expected toast title in output")
	}
	if !strings.Contains(html, "Your changes have been saved.") {
		t.Error("expected toast body in output")
	}
	if !strings.Contains(html, `lvt-click="dismiss_toast_notifications"`) {
		t.Error("expected dismiss button with lvt-click")
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
	// Unstyled version should not have Tailwind classes
	if strings.Contains(html, "fixed z-50") {
		t.Error("unstyled template should not have Tailwind classes")
	}
	// But should still have functional attributes
	if !strings.Contains(html, `role="alert"`) {
		t.Error("unstyled template should have role=alert")
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
	// Should render container even if empty
	if !strings.Contains(html, `data-toast-container="test"`) {
		t.Error("expected empty container to render")
	}
}
