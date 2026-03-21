package imaging

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io"
	"testing"
)

// createTestImage creates a simple PNG image for testing.
func createTestImage(w, h int) io.Reader {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return &buf
}

func TestGenerateThumbnail(t *testing.T) {
	src := createTestImage(800, 600)

	result, err := GenerateThumbnail(src, 200, 200)
	if err != nil {
		t.Fatalf("GenerateThumbnail() error = %v", err)
	}

	// Decode the result to verify it's a valid image
	data, err := io.ReadAll(result)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	if len(data) == 0 {
		t.Error("GenerateThumbnail() returned empty result")
	}

	// Verify JPEG header (starts with 0xFF 0xD8)
	if data[0] != 0xFF || data[1] != 0xD8 {
		t.Error("GenerateThumbnail() output is not JPEG")
	}
}

func TestGenerateThumbnail_PreservesAspectRatio(t *testing.T) {
	src := createTestImage(800, 400) // 2:1 ratio

	result, err := GenerateThumbnail(src, 200, 200)
	if err != nil {
		t.Fatalf("GenerateThumbnail() error = %v", err)
	}

	data, err := io.ReadAll(result)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	// Decode to check dimensions
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// With 800x400 input and 200x200 fit, width should be 200, height should be 100
	if width != 200 {
		t.Errorf("width = %d, want 200", width)
	}
	if height != 100 {
		t.Errorf("height = %d, want 100", height)
	}
}

func TestGenerateThumbnail_SmallImage(t *testing.T) {
	src := createTestImage(50, 50)

	result, err := GenerateThumbnail(src, 200, 200)
	if err != nil {
		t.Fatalf("GenerateThumbnail() error = %v", err)
	}

	data, err := io.ReadAll(result)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	bounds := img.Bounds()
	// imaging.Fit does not upscale
	if bounds.Dx() != 50 || bounds.Dy() != 50 {
		t.Errorf("small image should not be upscaled: got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestGenerateThumbnail_InvalidInput(t *testing.T) {
	src := bytes.NewReader([]byte("not an image"))

	_, err := GenerateThumbnail(src, 200, 200)
	if err == nil {
		t.Error("GenerateThumbnail() expected error for invalid input")
	}
}
