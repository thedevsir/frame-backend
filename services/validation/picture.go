package validation

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"github.com/thedevsir/frame-backend/services/errors"
)

func ValidatePicture(size int64, r io.Reader, exFormats string, exSize int64) error {

	_, format, err := image.DecodeConfig(r)
	if err != nil {
		return errors.ErrPictureNotValid
	}

	if strings.Index(exFormats, format) == -1 {
		return errors.ErrPictureNotValid
	}

	if exSize < size {
		return errors.ErrPictureTooLarge
	}

	return nil
}
