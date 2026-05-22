package main

import (
	"databaseHandler/cli"
	"databaseHandler/storage"
	"databaseHandler/tree"
	"databaseHandler/seeder"
	"databaseHandler/catalog"
	"databaseHandler/lineardb"
	"fmt"
)

func main() {
	mode := "benchmark"	

	if mode == "benchmark" {
		seeder.RunBenchmark(seeder.BenchmarkConfig{
				Degrees:   []int{10, 100, 200, 300, 400, 500, 600, 700, 800, 900, 1000},
				RowCounts: []int{10, 100000, 200000, 300000, 400000, 500000, 600000, 700000, 800000, 900000, 1000000},
				TableName: "users",
				Columns: []catalog.Column{
						{Name: "name", Type: "string"},
						{Name: "age", Type: "int"},
				},
		})
		return
	}

 if mode == "linearbenchmark" {
    lineardb.LinearRunBenchmark(lineardb.LinearBenchmarkConfig{
        RowCounts: []int{10, 100000, 200000, 300000, 400000, 500000, 600000, 700000, 800000, 900000, 1000000},
        TableName: "users",
        Columns: []catalog.Column{
            {Name: "name", Type: "string"},
            {Name: "age", Type: "int"},
        },
    })
		return
	}

	if mode == "linear" {
    s, _ := storage.NewStorage()
    defer s.Close()
    list := lineardb.NewList()
    lineardb.LinearStartLoop(s, &list)
    return
	}
	

	tr := tree.NewTree(10)
	s, error := storage.NewStorage()
	if error != nil {
		fmt.Println("falied to open storage:", error)
		return
	}
	defer s.Close()

	cli.StartLoop(s, &tr)
}
