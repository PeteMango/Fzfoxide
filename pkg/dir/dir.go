package dir

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/petemango/fzfoxide/pkg/fuzzy"
)

type Directory struct {
	Path         string    `json:"path"`
	Count        int       `json:"count"`
	LastAccessed time.Time `json:"last_accessed"`
}

var getDatabasePath = func() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	dbPath := filepath.Join(usr.HomeDir, ".fzfoxide")
	return dbPath, nil
}

func readDatabase() ([]Directory, error) {
	dbPath, err := getDatabasePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return []Directory{}, nil
	}

	data, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, err
	}

	var dirs []Directory
	err = json.Unmarshal(data, &dirs)
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func writeDatabase(dirs []Directory) error {
	dbPath, err := getDatabasePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(dirs, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(dbPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func RecordDirectory(dirPath string) error {
	dirs, err := readDatabase()
	if err != nil {
		return err
	}

	found := false
	now := time.Now()
	for i, dir := range dirs {
		if dir.Path == dirPath {
			dirs[i].Count++
			dirs[i].LastAccessed = now
			found = true
			break
		}
	}
	if !found {
		newEntry := Directory{
			Path:         dirPath,
			Count:        1,
			LastAccessed: now,
		}
		dirs = append(dirs, newEntry)
	}

	err = writeDatabase(dirs)
	if err != nil {
		return err
	}
	return nil
}

func Run(input string) (string, error) {
	if filepath.IsAbs(input) {

		if info, err := os.Stat(input); err == nil && info.IsDir() {
			err = RecordDirectory(input)
			if err != nil {
				return "", err
			}
			return input, nil
		}
	} else {
		absPath, err := filepath.Abs(input)
		if err == nil {
			if info, err := os.Stat(absPath); err == nil && info.IsDir() {
				err = RecordDirectory(absPath)
				if err != nil {
					return "", err
				}
				return absPath, nil
			}
		}
	}

	dirs, err := readDatabase()
	if err != nil {
		return "", err
	}

	if len(dirs) == 0 {
		return "", fmt.Errorf("no directories in the database")
	}

	type Match struct {
		Dir        Directory
		Similarity float64
	}

	matches := []Match{}
	for _, dir := range dirs {
		components := strings.Split(dir.Path, string(filepath.Separator))
		maxSimilarity := 0.0
		for _, component := range components {
			similarity := fuzzy.SimilarityPercentage(strings.ToLower(component), strings.ToLower(input))
			if similarity > maxSimilarity {
				maxSimilarity = similarity
			}
		}
		if maxSimilarity >= 50.0 {
			matches = append(matches, Match{Dir: dir, Similarity: maxSimilarity})
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no matching directory found with at least 50%% similarity")
	}

	bestMatch := matches[0]
	for _, match := range matches[1:] {
		if match.Similarity > bestMatch.Similarity {
			bestMatch = match
		} else if match.Similarity == bestMatch.Similarity && match.Dir.Count > bestMatch.Dir.Count {
			bestMatch = match
		}
	}

	err = RecordDirectory(bestMatch.Dir.Path)
	if err != nil {
		return "", err
	}
	return bestMatch.Dir.Path, nil
}
