package docs

import "embed"

//go:embed TUTORIAL-INICIANTE.md
var FS embed.FS

// ReadTutorial returns the content of the beginner tutorial.
func ReadTutorial() (string, error) {
	data, err := FS.ReadFile("TUTORIAL-INICIANTE.md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
