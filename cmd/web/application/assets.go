package application

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (app *Config) EnsureAssetsDir() error {
	root := app.GetAssetsRoot()
	if _, err := os.Stat(root); os.IsNotExist(err) {
		return os.Mkdir(root, 0750)
	}
	return nil
}

func GenerateAssetPath() string {
	base := make([]byte, 32)
	_, err := rand.Read(base)
	if err != nil {
		panic("failed to generate random bytes, you might be on an old version of Go or Linux incompatible with rand.Read")
	}
	id := base64.RawURLEncoding.EncodeToString(base)
	// Output file will always be png because of bg-removal
	return fmt.Sprintf("%s.png", id)
}

func (app *Config) GetAssetDiskPath(assetPath string) (string, error) {
	rootAbs, err := filepath.Abs(app.GetAssetsRoot())
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

func (app *Config) GetAssetURL(assetPath string) string {
	return fmt.Sprintf("http://localhost:%s/assets/%s", app.GetPort(), assetPath)
}

func MediaTypeToExtension(mediaType string) string {
	parts := strings.Split(mediaType, "/")
	if len(parts) != 2 {
		return ".bin"
	}
	return "." + parts[1]
}

func (app *Config) placeholderImage() string {
	return "/assets/placeholder.jpg"
}
