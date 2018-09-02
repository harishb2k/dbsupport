Step 1)  Make a main.go file
====
package main

import (
    "flag"
    "fmt"
    "github.com/harishb2k/dbsupport/mysqldb"
    "strconv"
)

func main()  {
    var user, password, host, database, table, packageName string
    var port = 3306
    flag.IntVar(&port, "port", 3306, "MySQL Port")
    flag.StringVar(&user, "user", "<some user name>", "Database user name")
    flag.StringVar(&password, "password", "<some password>", "Database password")
    flag.StringVar(&host, "host", "localhost", "MySQL Host")
    flag.StringVar(&database, "database", "<some db>", "Database name")
    flag.StringVar(&table, "table", "<some table>", "Table name")
    flag.StringVar(&packageName, "package", "main", "Package name")
    flag.Parse()

    var url = user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + database + ""
    fmt.Printf("Input Arguments -> host=%s, port=%d, user=%s, password=%s, database=%s, table=%s, url=%s\n", host, port, user, password, database, table, url)

    var db = mysqldb.Db{Url: url}
    db.Initialize()
    db.PrintSchema(database, table, packageName)
}
====

Step 2) Build   schema_builder
export GOPATH=$HOME/Go:${PWD}
export GOROOT=/usr/local/opt/go/libexec
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin
go build -o schema_builder src/main.go