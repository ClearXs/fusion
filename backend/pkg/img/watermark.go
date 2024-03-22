package img

import (
	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"io"
	"path/filepath"
)

type Location = byte

// watermark location
const (
	TopLeft     Location = 0
	TopRight    Location = 1
	BottomLeft  Location = 2
	BottomRight Location = 3
	Center      Location = 4
)

type Watermark struct {
	src image.Image
}

// From filepath create Watermark instance
func From(filepath string) (*Watermark, error) {
	img, err := imaging.Open(filepath)
	if err != nil {
		return nil, err
	}
	return &Watermark{src: img}, nil
}

// NewWatermark from io.Reader assign to input create Watermark
func NewWatermark(input io.Reader) (*Watermark, error) {
	img, err := imaging.Decode(input)
	if err != nil {
		return nil, err
	}
	return &Watermark{src: img}, nil
}

// Encode text mask to src image as watermark and through location adjust watermark place
func (watermark *Watermark) Encode(w io.Writer, text string, l Location, f imaging.Format) error {
	inputBounds := watermark.src.Bounds()

	// Create a new RGBA image with the same size as the input image
	output := imaging.New(inputBounds.Dx(), inputBounds.Dy(), color.NRGBA{})

	// Draw the input image onto the result image
	draw.Draw(output, inputBounds, watermark.src, image.Point{}, draw.Src)

	// Create a Mask
	mask := image.NewRGBA(image.Rect(0, 0, inputBounds.Dx(), inputBounds.Dy()))
	// Draw the input image onto the result image
	draw.Draw(mask, inputBounds, watermark.src, image.Point{}, draw.Src)
	point := watermark.calculateLocation(l, inputBounds)
	d := &font.Drawer{
		Dst:  mask,
		Src:  image.Black,
		Face: basicfont.Face7x13,
		Dot:  fixed.P(point.X, point.Y),
	}
	d.DrawString(text)

	// Draw the text image onto the result image
	draw.Draw(output, mask.Bounds(), mask, image.Point{}, draw.Over)
	return imaging.Encode(w, output, f)
}

// EncodeFromFilepath from path gain image then invoke EncodeFromImage
func (watermark *Watermark) EncodeFromFilepath(w io.Writer, path string, l Location) error {
	img, err := imaging.Open(path)
	if err != nil {
		return err
	}
	ext := filepath.Ext(path)
	format, err := imaging.FormatFromExtension(ext)
	if err != nil {
		return err
	}
	return watermark.EncodeFromImage(w, img, l, format)
}

// EncodeFromImage base on mask image to src input image, and through location adjust watermark place
func (watermark *Watermark) EncodeFromImage(w io.Writer, mask image.Image, l Location, format imaging.Format) error {
	inputBounds := watermark.src.Bounds()

	// Create a new RGBA image with the same size as the input image
	output := imaging.New(inputBounds.Dx(), inputBounds.Dy(), color.NRGBA{})
	point := watermark.calculateLocation(l, inputBounds)

	// Draw maks to output
	draw.Draw(output, mask.Bounds(), mask, point, draw.Over)

	return imaging.Encode(w, output, format)
}

// calculateLocation calculate watermark location base on Location
func (watermark *Watermark) calculateLocation(l Location, inputBounds image.Rectangle) image.Point {
	switch l {
	case TopLeft:
		return image.Point{}
	case TopRight:
		return image.Point{Y: inputBounds.Dy()}
	case BottomLeft:
		return image.Point{X: inputBounds.Dx()}
	case BottomRight:
		return image.Point{X: inputBounds.Dx(), Y: inputBounds.Dy()}
	case Center:
		return image.Point{X: inputBounds.Dx() / 2, Y: inputBounds.Dy() / 2}
	default:
		return image.Point{}
	}
}
