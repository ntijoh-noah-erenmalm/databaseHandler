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
	mode := "cli"	

	if mode == "benchmark" {
		seeder.RunBenchmark(seeder.BenchmarkConfig{
				Degrees:   []int{3, 5, 10, 20},
				RowCounts: []int{100, 1000, 10000},
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
