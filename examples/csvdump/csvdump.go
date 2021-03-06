/*
   Package main in csvdump represents a cursor->csv dumper

   Copyright 2013 Tamás Gulácsi

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/errgo.v1"
	"gopkg.in/goracle.v1/examples/connect"
	"gopkg.in/goracle.v1/oracle"
)

func getQuery(table, where string, columns []string) string {
	if strings.HasPrefix(table, "SELECT ") {
		return table
	}
	cols := "*"
	if len(columns) > 0 {
		cols = strings.Join(columns, ", ")
	}
	if where == "" {
		return "SELECT " + cols + " FROM " + table
	}
	return "SELECT " + cols + " FROM " + table + " WHERE " + where
}

func dump(w io.Writer, qry string) error {
	//db, err := connect.GetConnection("")
	cx, err := connect.GetRawConnection("")
	if err != nil {
		return errgo.Newf("error connecting to database: %s", err)
	}
	//defer db.Close()
	defer cx.Close()
	//rows, err := db.Query(qry)
	cu := cx.NewCursor()
	defer cu.Close()
	err = cu.Execute(qry, nil, nil)
	if err != nil {
		return errgo.Newf("error executing %q: %s", qry, err)
	}
	//defer rows.Close()
	//columns, err := rows.Columns()
	columns, err := GetColumns(cu)
	if err != nil {
		return errgo.Newf("error getting column converters: %s", err)
	}
	log.Printf("columns: %#v", columns)

	bw := bufio.NewWriterSize(w, 65536)
	defer bw.Flush()
	for i, col := range columns {
		if i > 0 {
			bw.Write([]byte{';'})
		}
		bw.Write([]byte{'"'})
		bw.WriteString(col.Name)
		bw.Write([]byte{'"'})
	}
	bw.Write([]byte{'\n'})
	n := 0
	rows, err := cu.FetchMany(1000)
	for err == nil && len(rows) > 0 {
		for _, row := range rows {
			for i, data := range row {
				if i > 0 {
					bw.Write([]byte{';'})
				}
				if data == nil {
					continue
				}
				bw.WriteString(columns[i].String(data))
			}
			bw.Write([]byte{'\n'})
			n++
		}
		rows, err = cu.FetchMany(1000)
	}
	log.Printf("written %d rows.", n)
	if err != nil && err != io.EOF {
		return errgo.Newf("error fetching rows from %s: %s", cu, err)
	}
	return nil
}

type ColConverter func(interface{}) string

type Column struct {
	Name   string
	String ColConverter
}

func GetColumns(cu *oracle.Cursor) (cols []Column, err error) {
	desc, err := cu.GetDescription()
	if err != nil {
		return nil, errgo.Newf("error getting description for %s: %s", cu, err)
	}
	//log.Printf("columns: %s", columns)
	log.Printf("desc: %#v", desc)
	//cols := godrv.ColumnDescriber(rows).DescribeColumns()
	//log.Printf("cols: %s", cols)
	//for rows.Next() {
	var ok bool
	cols = make([]Column, len(desc))
	for i, col := range desc {
		cols[i].Name = col.Name
		if cols[i].String, ok = converters[col.Type]; !ok {
			log.Fatalf("no converter for type %d (column name: %s)", col.Type, col.Name)
		}
	}
	return cols, nil
}

var converters = map[int]ColConverter{
	1: func(data interface{}) string { //VARCHAR2
		return fmt.Sprintf("%q", data.(string))
	},
	6: func(data interface{}) string { //NUMBER
		return fmt.Sprintf("%v", data)
	},
	96: func(data interface{}) string { //CHAR
		return fmt.Sprintf("%q", data.(string))
	},
	156: func(data interface{}) string { //DATE
		return `"` + data.(time.Time).Format(time.RFC3339) + `"`
	},
}

func main() {
	var (
		where   string
		columns []string
	)

	flag.Parse()
	if flag.NArg() > 1 {
		where = flag.Arg(1)
		if flag.NArg() > 2 {
			columns = flag.Args()[2:]
		}
	}
	qry := getQuery(flag.Arg(0), where, columns)
	if err := dump(os.Stdout, qry); err != nil {
		log.Printf("error dumping: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}
