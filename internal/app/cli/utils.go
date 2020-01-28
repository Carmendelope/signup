/*
 * Copyright 2020 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"encoding/base64"
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetPath resolves a given path by adding support for relative paths.
func GetPath(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, _ := user.Current()
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	if strings.HasPrefix(path, "../") {
		abs, _ := filepath.Abs("../")
		return strings.Replace(path, "..", abs, 1)
	}
	if strings.HasPrefix(path, ".") {
		abs, _ := filepath.Abs("./")
		return strings.Replace(path, ".", abs, 1)
	}
	return path
}

// PhotoPathToBase64 reads a image an convert the content to a base64 string
func PhotoPathToBase64(path string) (string, derrors.Error) {
	// if there is no path -> empty image
	if path == "" {
		return "", nil
	}

	viErr := ValidateImage(path)
	if viErr != nil {
		return "", viErr
	}

	convertedPath := GetPath(path)
	content, err := ioutil.ReadFile(convertedPath)
	if err != nil {
		return "", derrors.AsError(err, "cannot read image")
	}

	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(content)

	return imgBase64Str, nil
}

// ValidateImage validates that the image is jpg or png and wights under 1 MB
func ValidateImage(photoPath string) derrors.Error {
	// Check extension
	photoExt := filepath.Ext(photoPath)
	log.Debug().Str("extension", photoExt).Msg("image extension")
	if photoExt != ".jpg" && photoExt != ".JPG" && photoExt != ".jpeg" && photoExt != ".JPEG" && photoExt != ".png" && photoExt != ".PNG" {
		log.Error().Msg("invalid image format, please use jpg or png")
	}

	// Check size
	photoFile, err := os.Stat(photoPath)
	if err != nil {
		log.Error().Err(err).Msg("cannot read photo")
		return derrors.NewGenericError("cannot read photo")
	} else {
		if photoFile.Size() > 1024*1024 {
			log.Error().Msg("image too big, should weight less than 1 MB")
			return derrors.NewGenericError("image too big, should weight less than 1 MB")
		}
	}

	// Valid image
	return nil
}
