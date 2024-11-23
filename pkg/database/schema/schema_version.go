package schema

import (
	"bufio"
	"fmt"
	"io/fs"
	"strconv"
	"strings"

	"github.com/bioform/go-web-app-template/pkg/util"
)

const (
	// Path to your migrations folder
	SchemaFilePath      = "db/schema.sql"
	SchemaVersionPrefix = "-- Latest Migration: "
)

func SchemaVersion(fsys fs.FS) (int64, error) {
	if fsys == nil {
		fsys = util.OsFS{}
	}

	// Open the schema file for reading.
	file, err := fsys.Open(SchemaFilePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open schema file: %w", err)
	}
	defer file.Close()

	// Read only the first few lines of the file.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " ")
		if strings.HasPrefix(line, SchemaVersionPrefix) {
			migrationID, err := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(line, SchemaVersionPrefix)), 10, 64)
			if err != nil {
				return 0, fmt.Errorf("invalid migration ID in schema file: %w", err)
			}
			return migrationID, nil
		}
		if len(line) != 0 && !strings.HasPrefix(line, "--") {
			return 0, fmt.Errorf("latest migration identifier not found in schema file")
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("failed to read schema file: %w", err)
	}

	return 0, fmt.Errorf("latest migration identifier not found in schema file")
}
