package mysqldb

import (
    "database/sql"
    "errors"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
)

type RecordConsumer = func(data interface{})

type RowMapper interface {
    Map(row *sql.Rows) (interface{}, error)
}

type Db struct {
    Url       string
    dbSession *sql.DB
}

func (db *Db) Initialize() (error) {
    var err error
    if db.dbSession, err = sql.Open("mysql", db.Url); err != nil {
        fmt.Printf("Failed to create mysql connection: url=%s\n", db.Url)
        return errors.New("failed to create connection")
    }
    return nil
}

func (db *Db) QueryAll(query string, rowMapper RowMapper, param ... interface{}) ([]interface{}, error) {
    if rows, err := db.dbSession.Query(query, param...); err != nil {
        return nil, errors.New("failed to create connection")
    } else {
        var results []interface{}
        for rows.Next() {
            if data, err := rowMapper.Map(rows); err != nil {
                return nil, errors.New("failed to map result with mapper")
            } else {
                results = append(results, data)
            }
        }
        return results, nil
    }
}

func (db *Db) QueryOne(query string, rowMapper RowMapper, param ... interface{}) (interface{}, error) {
    if rows, err := db.dbSession.Query(query, param...); err != nil {
        return nil, errors.New("failed to create connection")
    } else {
        for rows.Next() {
            if data, err := rowMapper.Map(rows); err != nil {
                return nil, errors.New("failed to map result with mapper")
            } else {
                return data, nil
            }
        }
        return nil, errors.New("no result found")
    }
}

func (db *Db) QueryAllV2(query string, rowMapper RowMapper, recordConsumer RecordConsumer, param ... interface{}) (error) {

    // Fetch data from DB
    rows, err := db.dbSession.Query(query, param...)
    if err != nil {
        return errors.New("failed to create connection")
    }

    // Process all rows and pass it to consumer
    for rows.Next() {
        if data, err := rowMapper.Map(rows); err != nil {
            return errors.New("failed to map result with mapper")
        } else {
            recordConsumer(data)
        }
    }
    return nil
}
