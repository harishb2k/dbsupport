Step 1)  Make a main.go file
====
package main
import "github.com/harishb2k/dbsupport/mysqldb"
func main() {
    mysqldb.SchemaBuilderMain()
}
====

Step 2) Build   schema_builder
export GOPATH=${PWD}:$HOME/Go
export GOROOT=/usr/local/opt/go/libexec
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:$GOROOT/bin
go build -o schema_builder src/main.go