package checker

import (
	"github.com/yanpin-dev/propeller/pkg/health"
	"database/sql"
)

// DBChecker is a checker that check a database connection
type DBChecker struct {
	name       string
	checkSQL   string
	versionSQL string
	db         *sql.DB
}

// NewRedisChecker returns a new db.DBChecker with the given URL
func NewDBChecker(name, checkSQL, versionSQL string, db *sql.DB) health.Checker {
	return DBChecker{name: name, checkSQL: checkSQL, versionSQL: versionSQL, db: db}
}

// NewMySQLChecker returns a new db.DBChecker configured for use in MySQL
func NewMySQLChecker(db *sql.DB) health.Checker {
	return NewDBChecker("mysql", "SELECT 1", "SELECT VERSION()", db)
}

// NewPostgreSQLChecker returns a new db.DBChecker configured for use in PostgreSQL
func NewPostgreSQLChecker(db *sql.DB) health.Checker {
	return NewDBChecker("postgres", "SELECT 1", "SELECT VERSION()", db)
}

// Check execute two queries in the database
// The first is a simple one used to verify if the database is up
// If is Up then another query is executed, querying for the database version
func (c DBChecker) Check() health.Health {
	var (
		version string
		ok      string
	)

	h := health.NewHealth()

	if c.db == nil {
		h.Down().AddInfo("error", "Empty resource")
		return h
	}

	err := c.db.QueryRow(c.checkSQL).Scan(&ok)

	if err != nil {
		h.Down().AddInfo("error", err.Error())
		return h
	}

	// We are gonna make it versionSQL optional, as I cannot change the API,
	// we decided to ignore if the versionSQL is empty.
	if c.versionSQL != "" {
		err = c.db.QueryRow(c.versionSQL).Scan(&version)

		if err != nil {
			h.Down().AddInfo("error", err.Error())
			return h
		}

		h.Up().AddInfo("version", version)
	}

	return h
}

func (c DBChecker) Name() string {
	return c.name
}
