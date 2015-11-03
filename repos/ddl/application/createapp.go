package main

import (
	"github.com/xtraclabs/roll/repos/ddl"
)

func main() {
	ddl.DeleteTable("Application")
	ddl.CreateAppTable()
}
