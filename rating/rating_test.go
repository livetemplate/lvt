package rating

import (
	"testing"
)

func TestNew(t *testing.T) {
	r := New("test-rating")

	if r.ID() != "test-rating" {
		t.Errorf("Expected ID 'test-rating', got '%s'", r.ID())
	}
	if r.Namespace() != "rating" {
		t.Errorf("Expected namespace 'rating', got '%s'", r.Namespace())
	}
	if r.MaxStars != 5 {
		t.Errorf("Expected default MaxStars 5, got %d", r.MaxStars)
	}
	if r.Value != 0 {
		t.Errorf("Expected default Value 0, got %f", r.Value)
	}
	if r.Size != "md" {
		t.Errorf("Expected default Size 'md', got '%s'", r.Size)
	}
	if r.Character != "★" {
		t.Errorf("Expected default Character '★', got '%s'", r.Character)
	}
}

func TestNewReadonly(t *testing.T) {
	r := NewReadonly("test", 4.5)

	if r.Value != 4.5 {
		t.Errorf("Expected Value 4.5, got %f", r.Value)
	}
	if !r.Readonly {
		t.Error("Expected Readonly to be true")
	}
}

func TestWithValue(t *testing.T) {
	r := New("test", WithValue(3.5))
	if r.Value != 4 { // Rounds without AllowHalf
		t.Errorf("Expected Value 4 (rounded), got %f", r.Value)
	}
}

func TestWithValueAndHalf(t *testing.T) {
	r := New("test", WithAllowHalf(true), WithValue(3.5))
	if r.Value != 3.5 {
		t.Errorf("Expected Value 3.5, got %f", r.Value)
	}
}

func TestWithMaxStars(t *testing.T) {
	r := New("test", WithMaxStars(10))
	if r.MaxStars != 10 {
		t.Errorf("Expected MaxStars 10, got %d", r.MaxStars)
	}
}

func TestWithAllowHalf(t *testing.T) {
	r := New("test", WithAllowHalf(true))
	if !r.AllowHalf {
		t.Error("Expected AllowHalf to be true")
	}
}

func TestWithAllowClear(t *testing.T) {
	r := New("test", WithAllowClear(true))
	if !r.AllowClear {
		t.Error("Expected AllowClear to be true")
	}
}

func TestWithReadonly(t *testing.T) {
	r := New("test", WithReadonly(true))
	if !r.Readonly {
		t.Error("Expected Readonly to be true")
	}
}

func TestWithSize(t *testing.T) {
	r := New("test", WithSize("lg"))
	if r.Size != "lg" {
		t.Errorf("Expected Size 'lg', got '%s'", r.Size)
	}
}

func TestWithColor(t *testing.T) {
	r := New("test", WithColor("red"))
	if r.Color != "red" {
		t.Errorf("Expected Color 'red', got '%s'", r.Color)
	}
}

func TestWithShowValue(t *testing.T) {
	r := New("test", WithShowValue(true))
	if !r.ShowValue {
		t.Error("Expected ShowValue to be true")
	}
}

func TestWithShowCount(t *testing.T) {
	r := New("test", WithShowCount(true), WithCount(100))
	if !r.ShowCount {
		t.Error("Expected ShowCount to be true")
	}
	if r.Count != 100 {
		t.Errorf("Expected Count 100, got %d", r.Count)
	}
}

func TestWithLabel(t *testing.T) {
	r := New("test", WithLabel("Product Rating"))
	if r.Label != "Product Rating" {
		t.Errorf("Expected Label 'Product Rating', got '%s'", r.Label)
	}
}

func TestWithCharacter(t *testing.T) {
	r := New("test", WithCharacter("♥"))
	if r.Character != "♥" {
		t.Errorf("Expected Character '♥', got '%s'", r.Character)
	}
}

func TestWithStyled(t *testing.T) {
	r := New("test", WithStyled(false))
	if r.IsStyled() {
		t.Error("Expected IsStyled to be false")
	}
}

func TestSetValue(t *testing.T) {
	r := New("test")
	r.SetValue(3)

	if r.Value != 3 {
		t.Errorf("Expected Value 3, got %f", r.Value)
	}
}

func TestSetValueClampMin(t *testing.T) {
	r := New("test")
	r.SetValue(-1)

	if r.Value != 0 {
		t.Errorf("Expected Value 0 (clamped), got %f", r.Value)
	}
}

func TestSetValueClampMax(t *testing.T) {
	r := New("test", WithMaxStars(5))
	r.SetValue(10)

	if r.Value != 5 {
		t.Errorf("Expected Value 5 (clamped), got %f", r.Value)
	}
}

func TestSetValueRoundsToHalf(t *testing.T) {
	r := New("test", WithAllowHalf(true))
	r.SetValue(3.7)

	if r.Value != 3.5 {
		t.Errorf("Expected Value 3.5 (rounded to half), got %f", r.Value)
	}

	r.SetValue(3.8)
	if r.Value != 4 {
		t.Errorf("Expected Value 4 (rounded up), got %f", r.Value)
	}
}

func TestClear(t *testing.T) {
	r := New("test", WithValue(4))
	r.Clear()

	if r.Value != 0 {
		t.Errorf("Expected Value 0 after Clear, got %f", r.Value)
	}
}

func TestClick(t *testing.T) {
	r := New("test")
	r.Click(3)

	if r.Value != 3 {
		t.Errorf("Expected Value 3 after click, got %f", r.Value)
	}
}

func TestClickReadonly(t *testing.T) {
	r := New("test", WithReadonly(true), WithValue(2))
	r.Click(4)

	if r.Value != 2 {
		t.Error("Expected Value unchanged for readonly")
	}
}

func TestClickClear(t *testing.T) {
	r := New("test", WithAllowClear(true), WithValue(3))
	r.Click(3) // Click same value

	if r.Value != 0 {
		t.Errorf("Expected Value 0 (cleared), got %f", r.Value)
	}
}

func TestClickHalf(t *testing.T) {
	r := New("test", WithAllowHalf(true))
	r.ClickHalf(3, true) // First half

	if r.Value != 2.5 {
		t.Errorf("Expected Value 2.5 for half click, got %f", r.Value)
	}

	r.ClickHalf(3, false) // Second half
	if r.Value != 3 {
		t.Errorf("Expected Value 3 for full click, got %f", r.Value)
	}
}

func TestHover(t *testing.T) {
	r := New("test")
	r.Hover(4)

	if r.HoverValue != 4 {
		t.Errorf("Expected HoverValue 4, got %f", r.HoverValue)
	}
}

func TestHoverReadonly(t *testing.T) {
	r := New("test", WithReadonly(true))
	r.Hover(4)

	if r.HoverValue != -1 {
		t.Error("Expected HoverValue unchanged for readonly")
	}
}

func TestHoverHalf(t *testing.T) {
	r := New("test", WithAllowHalf(true))
	r.HoverHalf(3, true)

	if r.HoverValue != 2.5 {
		t.Errorf("Expected HoverValue 2.5, got %f", r.HoverValue)
	}
}

func TestLeave(t *testing.T) {
	r := New("test")
	r.Hover(4)
	r.Leave()

	if r.HoverValue != -1 {
		t.Errorf("Expected HoverValue -1, got %f", r.HoverValue)
	}
}

func TestDisplayValue(t *testing.T) {
	r := New("test", WithValue(3))

	// Without hover, returns actual value
	if r.DisplayValue() != 3 {
		t.Errorf("Expected DisplayValue 3, got %f", r.DisplayValue())
	}

	// With hover, returns hover value
	r.Hover(5)
	if r.DisplayValue() != 5 {
		t.Errorf("Expected DisplayValue 5 (hover), got %f", r.DisplayValue())
	}
}

func TestIsStarFull(t *testing.T) {
	r := New("test", WithValue(3))

	if !r.IsStarFull(1) || !r.IsStarFull(2) || !r.IsStarFull(3) {
		t.Error("Expected stars 1-3 to be full")
	}
	if r.IsStarFull(4) || r.IsStarFull(5) {
		t.Error("Expected stars 4-5 to not be full")
	}
}

func TestIsStarHalf(t *testing.T) {
	r := New("test", WithAllowHalf(true), WithValue(3.5))

	if r.IsStarHalf(3) {
		t.Error("Expected star 3 to not be half (it's full)")
	}
	if !r.IsStarHalf(4) {
		t.Error("Expected star 4 to be half")
	}
	if r.IsStarHalf(5) {
		t.Error("Expected star 5 to not be half (it's empty)")
	}
}

func TestIsStarEmpty(t *testing.T) {
	r := New("test", WithValue(3))

	if r.IsStarEmpty(3) {
		t.Error("Expected star 3 to not be empty")
	}
	if !r.IsStarEmpty(4) || !r.IsStarEmpty(5) {
		t.Error("Expected stars 4-5 to be empty")
	}
}

func TestStarState(t *testing.T) {
	r := New("test", WithAllowHalf(true), WithValue(2.5))

	if r.StarState(1) != "full" || r.StarState(2) != "full" {
		t.Error("Expected stars 1-2 to be full")
	}
	if r.StarState(3) != "half" {
		t.Errorf("Expected star 3 to be half, got '%s'", r.StarState(3))
	}
	if r.StarState(4) != "empty" || r.StarState(5) != "empty" {
		t.Error("Expected stars 4-5 to be empty")
	}
}

func TestStars(t *testing.T) {
	r := New("test", WithMaxStars(5))
	stars := r.Stars()

	if len(stars) != 5 {
		t.Errorf("Expected 5 stars, got %d", len(stars))
	}
	if stars[0] != 1 || stars[4] != 5 {
		t.Error("Expected stars 1-5")
	}
}

func TestPercentage(t *testing.T) {
	r := New("test", WithMaxStars(5), WithValue(4))

	if r.Percentage() != 80 {
		t.Errorf("Expected Percentage 80, got %f", r.Percentage())
	}
}

func TestFormatValue(t *testing.T) {
	r := New("test", WithValue(4))

	if r.FormatValue() != "4" {
		t.Errorf("Expected '4', got '%s'", r.FormatValue())
	}
}

func TestFormatValueHalf(t *testing.T) {
	r := New("test", WithAllowHalf(true), WithValue(4.5))

	if r.FormatValue() != "4.5" {
		t.Errorf("Expected '4.5', got '%s'", r.FormatValue())
	}
}

func TestFormatValueWithMax(t *testing.T) {
	r := New("test", WithMaxStars(5), WithValue(4))

	if r.FormatValueWithMax() != "4/5" {
		t.Errorf("Expected '4/5', got '%s'", r.FormatValueWithMax())
	}
}

func TestSizeClass(t *testing.T) {
	r := New("test", WithSize("sm"))
	if r.SizeClass() != "text-lg" {
		t.Errorf("Expected 'text-lg' for sm, got '%s'", r.SizeClass())
	}

	r.Size = "lg"
	if r.SizeClass() != "text-3xl" {
		t.Errorf("Expected 'text-3xl' for lg, got '%s'", r.SizeClass())
	}
}

func TestColorClass(t *testing.T) {
	r := New("test", WithColor("red"))
	if r.ColorClass() != "text-red-500" {
		t.Errorf("Expected 'text-red-500', got '%s'", r.ColorClass())
	}

	r.Color = "yellow"
	if r.ColorClass() != "text-yellow-400" {
		t.Errorf("Expected 'text-yellow-400', got '%s'", r.ColorClass())
	}
}

func TestTemplates(t *testing.T) {
	ts := Templates()
	if ts == nil {
		t.Fatal("Expected Templates() to return a TemplateSet")
	}
}
