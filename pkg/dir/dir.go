package dir

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"
)

type Directory struct {
	Path         string    `json:"path"`
	Count        int       `json:"count"`
	LastAccessed time.Time `json:"last_accessed"`
}

func getDatabasePath() (string, error) {
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

	/* empty database */
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

func RecordDir(dirPath string) error {
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
	/* check if its already an absolute path */
	if filepath.IsAbs(input) {

		/* if absolute, just go into that directory */
		if info, err := os.Stat(input); err == nil && info.IsDir() {
			err = RecordDir(input)
			if err != nil {
				return "", err
			}
			return input, nil
		}
	} else {
		/* check relative path */
		absPath, err := filepath.Abs(input)
		if err == nil {
			if info, err := os.Stat(absPath); err == nil && info.IsDir() {
				// found absolute path relative to this path
				err = RecordDir(absPath)
				if err != nil {
					return "", err
				}
				return absPath, nil
			}
		}
	}

	/* if not a regular cd command, look up in database */
	dirs, err := readDatabase()
	if err != nil {
		return "", err
	}
	matches := []Directory{}
	for _, dir := range dirs {
		if strings.Contains(dir.Path, input) {
			matches = append(matches, dir)
		}
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no matching directory found")
	}

	/* find the most popular match */
	bestMatch := matches[0]
	for _, match := range matches[1:] {
		if match.Count > bestMatch.Count {
			bestMatch = match
		}
	}

	/* record this path */
	err = RecordDir(bestMatch.Path)
	if err != nil {
		return "", err
	}
	return bestMatch.Path, nil
}
