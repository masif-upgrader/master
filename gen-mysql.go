// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if errGMS := genMysqlSchema(); errGMS != nil {
		fmt.Println(errGMS.Error())
		os.Exit(1)
	}
}

func genMysqlSchema() error {
	rawSchema, errRF := ioutil.ReadFile("schema/mysql/schema.sql")
	if errRF != nil {
		return errRF
	}

	ddls := []string{}

	for _, rawDdl := range bytes.Split(rawSchema, []byte(";")) {
		ddl := bytes.Trim(rawDdl, " \n")
		if len(ddl) > 0 {
			ddls = append(ddls, string(ddl))
		}
	}

	return ioutil.WriteFile("mysql.go", []byte(fmt.Sprintf("package main\nvar mysqlDdls = %#v", ddls)), 0666)
}
