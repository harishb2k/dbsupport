package dbsupport

// RowConsumer is called for each row as it is processed by QueryInterface.
type RowConsumer = func(data interface{})

// Database define all methods which should be supported.
type Database interface {
    Initialize(settings interface{}) (error)
    NewQueryInterface() QueryInterface
}

// RowMapper takes a input and maps to a row.
type RowMapper interface {
    MapRow(row interface{}) (interface{}, error)
}

// QueryInterface provides all method supported by database.
type QueryInterface interface {
    // QueryOne results a single result
    QueryOne(query string, mapper RowMapper, params ... interface{}) (interface{}, error)

    // QueryAll result all rows. It calls recordConsumer function for each row. The consumer should store
    // all the rows in some list if needed
    QueryAll(query string, mapper RowMapper, recordConsumer RowConsumer, params ... interface{}) (error)

    // Helper method to generate a file with Golang struct of table
    PrintSchema(database string, table string, packageName string)
}
