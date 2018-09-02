package mysqldb

import (
    "flag"
    "fmt"
    "strconv"
)

func SchemaBuilderMain() {
    var user, password, host, database, table, packageName string
    var port = 3306
    flag.IntVar(&port, "port", 3306, "MySQL Port")
    flag.StringVar(&user, "user", "<some user name>", "Database user name")
    flag.StringVar(&password, "password", "", "Database password")
    flag.StringVar(&host, "host", "localhost", "MySQL Host")
    flag.StringVar(&database, "database", "<some db>", "Database name")
    flag.StringVar(&table, "table", "<some table>", "Table name")
    flag.StringVar(&packageName, "package", "main", "Package name")
    flag.Parse()

    var url = user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + database + ""
    fmt.Printf("Input Arguments -> host=%s, port=%d, user=%s, password=%s, database=%s, table=%s, url=%s\n", host, port, user, password, database, table, url)

    var t = New()
    if err := t.Initialize(Settings{Url: url}); err != nil {
        fmt.Println("Error in DB init", err)
        return
    }

    var queryInterface = t.NewQueryInterface()
    queryInterface.PrintSchema(database, table, packageName)
}
