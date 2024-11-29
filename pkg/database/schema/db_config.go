package schema

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/bioform/go-web-app-template/pkg/database/find"
)

type DBConfig struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Output   string
}

// a. Postgres DSN
// Postgres DSNs can have the format:
//
//   - Connection URI: postgres://user:password@host:port/database
//   - Key-value format: user=username password=pass host=hostname port=5432 dbname=database
//
// You can use the net/url package to parse the connection URI or split the key-value format manually.

// b. MySQL DSN
// MySQL DSNs usually look like this:
//
//   - user:password@tcp(host:port)/database
//
// You can use regular expressions or string manipulation to extract the components.

// c. SQLite DSN
// SQLite DSNs are typically just the database file path (e.g., file::memory: or /path/to/database.db).

var dbSchemaFile = "db/schema.sql"

func ReadSchemaFile() ([]byte, error) {
	path, err := find.File(dbSchemaFile)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

func ParseDSN(dbType, dsn string) (DBConfig, error) {
	switch dbType {
	case "postgres":
		// Parse DSN as URL
		u, err := url.Parse(dsn)
		if err != nil {
			return DBConfig{}, fmt.Errorf("invalid Postgres DSN: %w", err)
		}

		user := u.User.Username()
		password, _ := u.User.Password()

		return DBConfig{
			Type:     "postgres",
			Host:     u.Hostname(),
			Port:     u.Port(),
			User:     user,
			Password: password,
			Database: strings.TrimPrefix(u.Path, "/"),
			Output:   dbSchemaFile,
		}, nil

	case "mysql":
		// MySQL DSN pattern: user:password@tcp(host:port)/database
		re := regexp.MustCompile(`^(?P<user>[^:]+):(?P<password>[^@]*)@tcp\((?P<host>[^:]+):(?P<port>\d+)\)/(?P<database>[^?]+)`)
		match := re.FindStringSubmatch(dsn)
		if match == nil {
			return DBConfig{}, errors.New("invalid MySQL DSN")
		}

		return DBConfig{
			Type:     "mysql",
			Host:     match[3],
			Port:     match[4],
			User:     match[1],
			Password: match[2],
			Database: match[5],
			Output:   dbSchemaFile,
		}, nil

	case "sqlite":
		// SQLite DSN is just the file path
		return DBConfig{
			Type:     "sqlite",
			Database: dsn,
			Output:   dbSchemaFile,
		}, nil

	default:
		return DBConfig{}, fmt.Errorf("unsupported database type: %s", dbType)
	}
}
