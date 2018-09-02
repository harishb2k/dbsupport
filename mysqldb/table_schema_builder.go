package mysqldb

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/iancoleman/strcase"
    "os"
    "strings"
)

type columnMetadata struct {
    ColumnName    string
    IsNullAllowed bool
    IsUnsigned    bool
    DataType      string

    internalColumnSchema string
    externalColumnSchema string
    mapperString         string
}

type tableMetadata struct {
    columns []columnMetadata

    ExternalTypeName string
    InternalTypeName string
}

func (tableMetadata *tableMetadata) MapRow(row interface{}) (interface{}, error) {
    if row, ok := row.(*sql.Rows); ok {
        if data, err := tableMetadata.Map(row); err == nil {
            return data, nil
        }
    }
    return nil, errors.New("some error.")
}

func (tableMetadata *tableMetadata) Map(row *sql.Rows) (interface{}, error) {
    if tableMetadata.columns == nil {
        tableMetadata.columns = [] columnMetadata{}
    }

    var columnName string
    var isNullabel string
    var columnType string
    var dataType string
    row.Scan(&columnName, &isNullabel, &columnType, &dataType)

    var columnMetadata = columnMetadata{}
    columnMetadata.ColumnName = strcase.ToCamel(columnName)
    columnMetadata.DataType = dataType
    columnMetadata.IsNullAllowed = isNullabel == "YES"
    columnMetadata.IsUnsigned = strings.Contains(columnType, "unsigned")

    fmt.Printf("ColumnName=%40s, DataType=%20s, IsNullAllowed=%10t, IsUnsigned=%10t \n", columnMetadata.ColumnName, columnMetadata.DataType, columnMetadata.IsNullAllowed, columnMetadata.IsUnsigned)

    switch columnMetadata.DataType {
    case "enum":
        columnMetadata.handleIntColumn()
    case "int":
        columnMetadata.handleIntColumn()
    case "bigint":
        columnMetadata.handleBigintColumn()
    case "tinyint":
        columnMetadata.handleTinyintColumn()
    case "varchar":
        columnMetadata.handleStringColumn()
    case "mediumtext":
        columnMetadata.handleStringColumn()
    case "datetime":
        columnMetadata.handleTimestampColumn()
    case "timestamp":
        columnMetadata.handleDatetimeColumn()

    }
    tableMetadata.columns = append(tableMetadata.columns, columnMetadata)

    return nil, nil
}

func (db mysqlQueryInterface) PrintSchema(database string, table string, packageName string) {
    var tableMetadata = &tableMetadata{}
    tableMetadata.ExternalTypeName = strcase.ToCamel(table)
    tableMetadata.InternalTypeName = "Db" + strcase.ToCamel(table)

    var query = "SELECT column_name, is_nullable, column_type, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '" + database + "' AND table_name = '" + table + "'"
    fmt.Println(query)

    err := db.QueryAll(query, tableMetadata, func(data interface{}) {
    })

    if err != nil {
        fmt.Println("Failed to run query", err)
        return
    }

    f, _ := os.Create(table + "_auto.go")
    defer f.Close()

    f.WriteString("package " + packageName + " \n import \"database/sql \n")

    f.WriteString("type " + tableMetadata.InternalTypeName + " struct {\n")
    f.WriteString("\n")
    for _, k := range tableMetadata.columns {
        f.WriteString(k.internalColumnSchema)
        f.WriteString("\n")
    }
    f.WriteString("}\n")

    f.WriteString("type " + tableMetadata.ExternalTypeName + " struct {")
    for _, k := range tableMetadata.columns {
        f.WriteString(k.externalColumnSchema)
        f.WriteString("\n")
    }
    f.WriteString("}\n")

    f.WriteString("func " + tableMetadata.InternalTypeName + "To" + tableMetadata.ExternalTypeName + "Mapper(source *" + tableMetadata.InternalTypeName + ", destination *" + tableMetadata.ExternalTypeName + " ) {")
    for _, k := range tableMetadata.columns {
        f.WriteString(k.mapperString)
        f.WriteString("\n")
    }
    f.WriteString("}\n")

    f.WriteString("func Convert" + tableMetadata.InternalTypeName + "To" + tableMetadata.ExternalTypeName + "(source *" + tableMetadata.InternalTypeName + ") " + tableMetadata.ExternalTypeName + "{ \n")
    f.WriteString("destination := " + tableMetadata.ExternalTypeName + "{} \n")
    f.WriteString("DbCustomersToCustomersMapper(source, &destination) \n")
    f.WriteString("return destination \n")
    f.WriteString("}\n")
}

/*
func (db *Db) PrintSchema(database string, table string, packageName string) {
    var tableMetadata = &tableMetadata{}
    tableMetadata.ExternalTypeName = strcase.ToCamel(table)
    tableMetadata.InternalTypeName = "Db" + strcase.ToCamel(table)

    var query = "SELECT column_name, is_nullable, column_type, data_type FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = '" + database + "' AND table_name = '" + table + "'"
    fmt.Println(query)
    if _, err := db.QueryAll(query, tableMetadata); err != nil {
        fmt.Println("Failed to run query", err)
        return
    }

    f, _ := os.Create(table + "_auto.go")
    defer f.Close()

    f.WriteString("package " + packageName + " \n import \"database/sql \n")

    f.WriteString("type " + tableMetadata.InternalTypeName + " struct {\n")
    f.WriteString("\n")
    for _, k := range tableMetadata.columns {
        f.WriteString(k.internalColumnSchema)
        f.WriteString("\n")
    }
    f.WriteString("}\n")

    f.WriteString("type " + tableMetadata.ExternalTypeName + " struct {")
    for _, k := range tableMetadata.columns {
        f.WriteString(k.externalColumnSchema)
        f.WriteString("\n")
    }
    f.WriteString("}\n")

    f.WriteString("func " + tableMetadata.InternalTypeName + "To" + tableMetadata.ExternalTypeName + "Mapper(source *" + tableMetadata.InternalTypeName + ", destination *" + tableMetadata.ExternalTypeName + " ) {")
    for _, k := range tableMetadata.columns {
        f.WriteString(k.mapperString)
        f.WriteString("\n")
    }
    f.WriteString("}\n")

    f.WriteString("func Convert" + tableMetadata.InternalTypeName + "To" + tableMetadata.ExternalTypeName + "(source *" + tableMetadata.InternalTypeName + ") " + tableMetadata.ExternalTypeName + "{ \n")
    f.WriteString("destination := " + tableMetadata.ExternalTypeName + "{} \n")
    f.WriteString("DbCustomersToCustomersMapper(source, &destination) \n")
    f.WriteString("return destination \n")
    f.WriteString("}\n")
}
*/
func (metadata *columnMetadata) handleStringColumn() {
    metadata.internalColumnSchema = metadata.ColumnName + " sql.NullString"
    metadata.externalColumnSchema = metadata.ColumnName + " string"

    metadata.mapperString = "if source." + metadata.ColumnName + ".Valid {"
    metadata.mapperString += "  destination." + metadata.ColumnName + " = " + "source." + metadata.ColumnName + ".String"
    metadata.mapperString += "}"
}

func (metadata *columnMetadata) handleTimestampColumn() {
    metadata.internalColumnSchema = metadata.ColumnName + " sql.NullInt64"
    metadata.externalColumnSchema = metadata.ColumnName + " int"

    metadata.mapperString = "if source." + metadata.ColumnName + ".Valid {"
    metadata.mapperString += "  destination." + metadata.ColumnName + " = " + "int(source." + metadata.ColumnName + ".Int64)"
    metadata.mapperString += "}"
}

func (metadata *columnMetadata) handleDatetimeColumn() {
    metadata.internalColumnSchema = metadata.ColumnName + " sql.NullInt64"
    metadata.externalColumnSchema = metadata.ColumnName + " int"

    metadata.mapperString = "if source." + metadata.ColumnName + ".Valid {"
    metadata.mapperString += "  destination." + metadata.ColumnName + " = " + "int(source." + metadata.ColumnName + ".Int64)"
    metadata.mapperString += "}"
}

func (metadata *columnMetadata) handleIntColumn() {
    metadata.internalColumnSchema = metadata.ColumnName + " sql.NullInt64"

    if metadata.IsUnsigned {
        metadata.externalColumnSchema = metadata.ColumnName + " int"
    } else {
        metadata.externalColumnSchema = metadata.ColumnName + " uint"
    }

    if metadata.IsUnsigned {
        metadata.mapperString = "if source." + metadata.ColumnName + ".Valid {"
        metadata.mapperString += "  destination." + metadata.ColumnName + " = " + "int(source." + metadata.ColumnName + ".Int64)"
        metadata.mapperString += "}"
    } else {
        metadata.mapperString = "if source." + metadata.ColumnName + ".Valid {"
        metadata.mapperString += "  destination." + metadata.ColumnName + " = " + "uint(source." + metadata.ColumnName + ".Int64)"
        metadata.mapperString += "}"
    }
}

func (metadata *columnMetadata) handleBigintColumn() {
    metadata.internalColumnSchema = metadata.ColumnName + " sql.NullInt64"

    if metadata.IsUnsigned {
        metadata.externalColumnSchema = metadata.ColumnName + " int64"
    } else {
        metadata.externalColumnSchema = metadata.ColumnName + " uint64"
    }

    if metadata.IsUnsigned {
        metadata.mapperString = "if source." + metadata.ColumnName + ".Valid {"
        metadata.mapperString += "  destination." + metadata.ColumnName + " = " + "int64(source." + metadata.ColumnName + ".Int64)"
        metadata.mapperString += "}"
    } else {
        metadata.mapperString = "if source." + metadata.ColumnName + ".Valid {"
        metadata.mapperString += "  destination." + metadata.ColumnName + " = " + "uint64(source." + metadata.ColumnName + ".Int64)"
        metadata.mapperString += "}"
    }
}

func (metadata *columnMetadata) handleTinyintColumn() {
    metadata.internalColumnSchema = metadata.ColumnName + " sql.NullInt64"

    if metadata.IsUnsigned {
        metadata.externalColumnSchema = metadata.ColumnName + " int8"
    } else {
        metadata.externalColumnSchema = metadata.ColumnName + " uint8"
    }

    if metadata.IsUnsigned {
        metadata.mapperString = "if source." + metadata.ColumnName + ".Valid {"
        metadata.mapperString += "  destination." + metadata.ColumnName + " = " + "int8(source." + metadata.ColumnName + ".Int64)"
        metadata.mapperString += "}"
    } else {
        metadata.mapperString = "if source." + metadata.ColumnName + ".Valid {"
        metadata.mapperString += "  destination." + metadata.ColumnName + " = " + "uint8(source." + metadata.ColumnName + ".Int64)"
        metadata.mapperString += "}"
    }
}
