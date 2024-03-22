package img

import "github.com/disintegration/imaging"

// IsImage by filename.
func IsImage(filename string) bool {
	if _, err := imaging.FormatFromFilename(filename); err != nil {
		return false
	}
	return true
}
