// Package imaging provides image processing utilities for generated applications.
package imaging

import (
	"bytes"
	"fmt"
	"io"

	"github.com/disintegration/imaging"
)

// GenerateThumbnail reads an image from src, resizes it to fit within
// maxWidth x maxHeight (preserving aspect ratio), and returns the result
// as a JPEG-encoded reader.
func GenerateThumbnail(src io.Reader, maxWidth, maxHeight int) (io.Reader, error) {
	img, err := imaging.Decode(src, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("imaging: decode: %w", err)
	}

	thumb := imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
	img = nil // allow GC of full-resolution image before encoding

	var buf bytes.Buffer
	if err := imaging.Encode(&buf, thumb, imaging.JPEG, imaging.JPEGQuality(85)); err != nil {
		return nil, fmt.Errorf("imaging: encode: %w", err)
	}

	return &buf, nil
}
