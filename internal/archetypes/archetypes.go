package archetypes

import (
	"embed"
	"os"
	"path/filepath"
	"strings"
)

//go:embed *.md
var archetypesFS embed.FS

var overridesDir = "archetypes"

func SetOverridesDir(dir string) {
	overridesDir = dir
}

func GetOverridesDir() string {
	return overridesDir
}

func Read(slug string) ([]byte, error) {
	filename := slug + ".md"
	// Try override first
	if data, err := os.ReadFile(filepath.Join(overridesDir, filename)); err == nil {
		return data, nil
	}
	return archetypesFS.ReadFile(filename)
}

func List() ([]string, error) {
	slugMap := make(map[string]bool)

	// List embedded
	entries, err := archetypesFS.ReadDir(".")
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				slugMap[strings.TrimSuffix(e.Name(), ".md")] = true
			}
		}
	}

	// List overrides
	if entries, err := os.ReadDir(overridesDir); err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
				slugMap[strings.TrimSuffix(e.Name(), ".md")] = true
			}
		}
	}

	var slugs []string
	for s := range slugMap {
		slugs = append(slugs, s)
	}
	return slugs, nil
}

func Exists(slug string) bool {
	filename := slug + ".md"
	if _, err := os.Stat(filepath.Join(overridesDir, filename)); err == nil {
		return true
	}
	_, err := archetypesFS.Open(filename)
	return err == nil
}

// Get returns the archetype content for the given slug, or an error if not found.
func Get(slug string) (string, error) {
	data, err := Read(slug)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteToTemp writes the archetype content to a temporary file and returns its path.
// This is useful for passing the file to external commands.
func WriteToTemp(slug string) (string, func(), error) {
	data, err := Read(slug)
	if err != nil {
		return "", nil, err
	}
	tmpFile, err := os.CreateTemp("", slug+"-*.md")
	if err != nil {
		return "", nil, err
	}
	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", nil, err
	}
	tmpFile.Close()
	cleanup := func() { os.Remove(tmpFile.Name()) }
	return tmpFile.Name(), cleanup, nil
}
