package main

import (
	"flag"
	"fmt"
	"os"
	"io"
	"encoding/csv"
)

type Cell struct {
	Key string
	Type string
	Value string
}

type Row struct {
	Cells []Cell
}

type Obj struct {
	Name string
	Rows []Row
}

func usage() {
	msg := `USAGE:
	csv2jsobj [OBJNAME PATH]
`
	fmt.Println(msg)
}

func parse(objname string, path_ string) Obj {
	f, err := os.Open(path_)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	var cols []string

	obj := Obj{}
	obj.Name = objname

	for {
		cols, err = reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		row := Row{}
		const nMember = 3
		nItems := len(cols) / nMember
		for i := 0; i < nItems; i++ {
			cell := Cell{}
			cell.Key = cols[0 + i * nMember]
			cell.Type = cols[1 + i * nMember]
			cell.Value = cols[2 + i * nMember]
			row.Cells = append(row.Cells, cell)
		}
		obj.Rows = append(obj.Rows, row)

	}
	return obj
}

func output(obj Obj) {
	fmt.Printf("const %s = [\n", obj.Name)

	for _, row := range obj.Rows {
	fmt.Println("  {")
		for _, cell := range row.Cells {
			if cell.Type == "string" {
				fmt.Printf("    %s: \"%s\",\n", cell.Key, cell.Value)

			} else if cell.Type == "number" {
				fmt.Printf("    %s: %s,\n", cell.Key, cell.Value)

			} else if cell.Type == "null" {
				fmt.Printf("    %s: null,\n", cell.Key)

			}
		}
	fmt.Println("  },")
	}

	fmt.Println("]")
}

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		usage()
		return
	}

	obj := parse(flag.Arg(0), flag.Arg(1))
	output(obj)
}

