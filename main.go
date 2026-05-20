package main

import (
	"databaseHandler/cli"
	"databaseHandler/storage"
	"databaseHandler/tree"
	"databaseHandler/seeder"
	"databaseHandler/catalog"
	"fmt"
)

func main() {
	mode := "benchmark"	

	if mode == "benchmark" {
		seeder.RunBenchmark(seeder.BenchmarkConfig{
				Degrees:   []int{10, 25, 100, 200},
				RowCounts: []int{100, 10000, 1000000},
				TableName: "users",
				Columns: []catalog.Column{
						{Name: "name", Type: "string"},
						{Name: "age", Type: "int"},
				},
		})
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
