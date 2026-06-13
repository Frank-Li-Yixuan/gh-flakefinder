package githubapi

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadFixture(path string) (DataSet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DataSet{}, fmt.Errorf("read fixture %q: %w", path, err)
	}
	var fixture DataSet
	if err := json.Unmarshal(data, &fixture); err != nil {
		return DataSet{}, fmt.Errorf("parse fixture %q: %w", path, err)
	}
	normalizeDataSet(&fixture)
	return fixture, nil
}

func normalizeDataSet(data *DataSet) {
	for i := range data.Runs {
		if data.Runs[i].Attempt == 0 {
			data.Runs[i].Attempt = 1
		}
	}
	for i := range data.Jobs {
		if data.Jobs[i].Attempt == 0 {
			data.Jobs[i].Attempt = 1
		}
	}
}
