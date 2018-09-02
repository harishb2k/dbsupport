Step 1)  Make a main.go file
====
package main
import "github.com/harishb2k/dbsupport/mysqldb"
func main() {
    mysqldb.SchemaBuilderMain()
}
====

Step 2)
Build   schema_builder
    export GOPATH=${PWD}:$HOME/Go
    export GOROOT=/usr/local/opt/go/libexec
    export PATH=$PATH:$GOPATH/bin
    export PATH=$PATH:$GOROOT/bin
    go build -o schema_builder src/main.go


Step 3)
 ./schema_builder --user <user> --host <localhost | your host> --database <db name --table <table name> --package anyname
 ./schema_builder --user abcd --password xyz --host localhost --database my_database --table my_table --package domainobjects




========================================== Usage ================================================

By using schema builder you can get following objects and helper method to work with this lib:


// your domain objects mapping to DB (it has sql.NullXXX value to work with Null values)
type DbXYZ struct {
    Id      sql.NullInt64
    Name sql.NullString
}

// your domain objects mapping to DB
type XYZ struct {
    Id      int
    Name    string
}

// Helper
func ConvertDbXYZToXYZ(source *DbXYZ) XYZ {
    ...
}




// Mapper to read sql.Row into your domain object (XYZ)
type Mapper struct{}

func (Mapper) MapRow(row interface{}) (interface{}, error) {
    dbObj := DbXYZ{}
    if row, ok := row.(*sql.Rows); ok {
        if err := row.Scan(&dbObj.Id, &dbObj.Name); err != nil {
            return nil, err
        }
        return ConvertDbXYZToXYZ(&dbObj), nil
    }
    return nil, errors.New("error in mapping row")
}



// Create a Database
db := mysqldb.New()
if err := db.Initialize(mysqldb.Settings{Url: url}); err != nil {
    fmt.Println("Error in database open")
    return
}

// Get a new query interface to make DB query
qi := db.NewQueryInterface()

// Find a object from DB
objectFound, err := qi.QueryOne("SELECT id, name from some_table limit 1", &Mapper{})
if err != nil {
    fmt.Println("Error in reading customer", err)
}
fmt.Println(objectFound)

// Find 10 objects from DB
qi.QueryAll("SELECT id, name from some_table limit 1", &Mapper{}, func(data interface{}) {
    if objectFound, ok := data.(XYZ); ok {
        fmt.Println(objectFound.Id, objectFound.Name)
    }
});
