package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (c *APIConfig) EnsureAssetsDir() error {
	if _, err := os.Stat(c.assetsRoot); os.IsNotExist(err) {
		return os.Mkdir(c.assetsRoot, 0750)
	}
	return nil
}

func getAssetPath(mediaType string) string {
	base := make([]byte, 32)
	_, err := rand.Read(base)
	if err != nil {
		panic("failed to generate random bytes")
	}
	id := base64.RawURLEncoding.EncodeToString(base)
	ext := mediaTypeToExtension(mediaType)
	return fmt.Sprintf("%s%s", id, ext)
}

func (c *APIConfig) getAssetDiskPath(assetPath string) (string, error) {
	rootAbs, err := filepath.Abs(c.assetsRoot)
	if err != nil {
		return "", err
	}

	joined := filepath.Join(rootAbs, assetPath)
	cleaned := filepath.Clean(joined)

	if strings.HasPrefix(cleaned, rootAbs) {
		return cleaned, nil
	}

	return "", errors.New("asset path is outside of assets root")
}

func (c *APIConfig) getAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", c.port, assetPath)
}

func mediaTypeToExtension(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}
