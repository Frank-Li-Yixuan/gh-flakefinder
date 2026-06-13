package report

import (
	"bytes"
	"encoding/json"

	"github.com/Frank-Li-Yixuan/gh-flakefinder/internal/detect"
)

func JSON(repository string, suspects []detect.Suspect) (string, error) {
	if suspects == nil {
		suspects = []detect.Suspect{}
	}
	envelope := struct {
		Repository   string           `json:"repository"`
		GeneratedBy  string           `json:"generated_by"`
		SuspectCount int              `json:"suspect_count"`
		Suspects     []detect.Suspect `json:"suspects"`
	}{
		Repository:   repository,
		GeneratedBy:  "gh-flakefinder",
		SuspectCount: len(suspects),
		Suspects:     suspects,
	}
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(envelope); err != nil {
		return "", err
	}
	return buf.String(), nil
}
