package mysqldb

import (
    "database/sql"
    "errors"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "github.com/harishb2k/dbsupport"
)

var ErrorConnectionNotCreated error = errors.New("failed to create connection")

// Settings contains all info to connection to a MySQL.
type Settings struct {
    Url string
}

// mysqlDatabase is internal structure to keep DB info.
type mysqlDatabase struct {
    dbSession *sql.DB
}

// Initialize a MySQL connection
func (db *mysqlDatabase) Initialize(settings interface{}) (error) {
    var err error
    if mysqlSettings, ok := settings.(Settings); ok {
        if db.dbSession, err = sql.Open("mysql", mysqlSettings.Url); err != nil {
            fmt.Printf("Failed to create mysql connection: url=%s\n", mysqlSettings.Url)
            return ErrorConnectionNotCreated
        }
    }
    return nil
}

// NewQueryInterface givens a query interface to query
func (db *mysqlDatabase) NewQueryInterface() (dbsupport.QueryInterface) {
    return mysqlQueryInterface{session: db.dbSession}
}

// New will give a new object
func New() dbsupport.Database {
    return &mysqlDatabase{}
}

// mysqlQueryInterface provides all Sql query
type mysqlQueryInterface struct {
    session *sql.DB
}

// QueryOne gives a single row
func (qi mysqlQueryInterface) QueryOne(query string, mapper dbsupport.RowMapper, params ... interface{}) (interface{}, error) {

    // Query DB
    rows, err := qi.session.Query(query, params...)
    if err != nil {
        return nil, err
    }

    // Check if we have data in result set, if we have rows then map row using mapper and return if no error
    if rows.Next() {
        row, err := mapper.MapRow(rows)
        if err != nil {
            return nil, err
        }
        return row, nil
    }

    return nil, errors.New("no rows found")
}

// QueryAll gets all rows
func (qi mysqlQueryInterface) QueryAll(query string, mapper dbsupport.RowMapper, recordConsumer dbsupport.RowConsumer, params ... interface{}) (error) {

    // Fetch data from DB
    rows, err := qi.session.Query(query, params...)
    if err != nil {
        return errors.New("failed to create connection")
    }

    // Process all rows and pass it to consumer
    for rows.Next() {
        if data, err := mapper.MapRow(rows); err != nil {
            return errors.New("failed to map result with mapper")
        } else {
            recordConsumer(data)
        }
    }
    return nil
}
