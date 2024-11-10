// dir/dir_test.go
package dir

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadWriteDatabase(t *testing.T) {
	tempDir := os.TempDir()
	dbPath := filepath.Join(tempDir, ".fzfoxide_test")
	originalGetDatabasePath := getDatabasePath
	getDatabasePath = func() (string, error) {
		return dbPath, nil
	}
	defer func() {
		getDatabasePath = originalGetDatabasePath
		os.Remove(dbPath)
	}()

	dirs := []Directory{
		{Path: "/home/user/Documents", Count: 2, LastAccessed: time.Now()},
		{Path: "/home/user/Downloads", Count: 1, LastAccessed: time.Now()},
	}
	err := writeDatabase(dirs)
	if err != nil {
		t.Fatalf("writeDatabase failed: %v", err)
	}

	readDirs, err := readDatabase()
	if err != nil {
		t.Fatalf("readDatabase failed: %v", err)
	}

	if len(readDirs) != len(dirs) {
		t.Fatalf("expected %d directories, got %d", len(dirs), len(readDirs))
	}

	for i, dir := range dirs {
		if readDirs[i].Path != dir.Path || readDirs[i].Count != dir.Count {
			t.Errorf("expected %v, got %v", dir, readDirs[i])
		}
	}
}

func TestRecordDirectory(t *testing.T) {
	tempDir := os.TempDir()
	dbPath := filepath.Join(tempDir, ".fzfoxide_test")
	originalGetDatabasePath := getDatabasePath
	getDatabasePath = func() (string, error) {
		return dbPath, nil
	}
	defer func() {
		getDatabasePath = originalGetDatabasePath
		os.Remove(dbPath)
	}()

	dirPath := "/home/user/Documents"
	err := RecordDirectory(dirPath)
	if err != nil {
		t.Fatalf("RecordDirectory failed: %v", err)
	}

	dirs, err := readDatabase()
	if err != nil {
		t.Fatalf("readDatabase failed: %v", err)
	}

	if len(dirs) != 1 {
		t.Fatalf("expected 1 directory, got %d", len(dirs))
	}

	if dirs[0].Path != dirPath || dirs[0].Count != 1 {
		t.Errorf("expected Path=%s, Count=1, got Path=%s, Count=%d", dirPath, dirs[0].Path, dirs[0].Count)
	}

	err = RecordDirectory(dirPath)
	if err != nil {
		t.Fatalf("RecordDirectory failed: %v", err)
	}

	dirs, err = readDatabase()
	if err != nil {
		t.Fatalf("readDatabase failed: %v", err)
	}

	if len(dirs) != 1 {
		t.Fatalf("expected 1 directory, got %d", len(dirs))
	}

	if dirs[0].Count != 2 {
		t.Errorf("expected Count=2, got %d", dirs[0].Count)
	}
}

func TestRunFuzzyMatch(t *testing.T) {
	tempDir := os.TempDir()
	dbPath := filepath.Join(tempDir, ".fzfoxide_test")
	originalGetDatabasePath := getDatabasePath
	getDatabasePath = func() (string, error) {
		return dbPath, nil
	}
	defer func() {
		getDatabasePath = originalGetDatabasePath
		os.Remove(dbPath)
	}()

	dirs := []Directory{
		{Path: "/home/user/Documents", Count: 2, LastAccessed: time.Now()},
		{Path: "/home/user/Downloads", Count: 1, LastAccessed: time.Now()},
		{Path: "/var/www/html", Count: 3, LastAccessed: time.Now()},
	}
	writeDatabase(dirs)

	input := "Documents"
	expected := "/home/user/Documents"
	result, err := Run(input)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}

	dirs, err = readDatabase()
	if err != nil {
		t.Fatalf("readDatabase failed: %v", err)
	}

	for _, dir := range dirs {
		if dir.Path == expected && dir.Count != 3 {
			t.Errorf("expected Count=3 for %s, got %d", expected, dir.Count)
		}
	}
}
