package main

import (
	"github.com/xtraclabs/roll/repos/ddl"
)

func main() {
	ddl.DeleteTable(ddl.ApplicationTableName)
	ddl.CreateAppTable()
}
