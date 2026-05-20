package seeder

import (
	"databaseHandler/catalog"
	"databaseHandler/storage"
	"databaseHandler/tree"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type SeederConfig struct {
    TableName string
    Columns   []catalog.Column
    RowCount  int
}

type BenchmarkConfig struct {
    Degrees   []int
    RowCounts []int
    Columns   []catalog.Column
    TableName string
}

func RunBenchmark(config BenchmarkConfig) {

		file, _ := os.Create("benchmark_results.txt")
		defer file.Close()

		fmt.Fprintln(file, "degree,rows,insert_ms,search_target,search_ns")

    for _, degree := range config.Degrees {
        for _, rows := range config.RowCounts {

						os.Remove("data.db")
						os.Remove("catalog.json")
            s, _ := storage.NewStorage()
            tr := tree.NewTree(degree)

            seedConfig := SeederConfig{
                TableName: config.TableName,
                Columns:   config.Columns,
                RowCount:  rows,
            }

            insertDuration := Seed(seedConfig, s, &tr)

            searchTargets := []int{
                0,
                rows - 1,
                rand.Intn(rows),
            }

            for _, target := range searchTargets {
                start := time.Now()
                tree.Search(target, tr.Root)
                searchDuration := time.Since(start)

                fmt.Fprintf(file, "%d,%d,%d,%d,%d\n",
                    degree,
                    rows,
                    insertDuration.Milliseconds(),
                    target,
                    searchDuration.Nanoseconds(),
                )
            }

            s.Close()
        }
    }
		fmt.Println("benchmark done :D")
}

func Seed(config SeederConfig, s *storage.Storage, t *tree.Tree) time.Duration {
    // create table if it doesn't exist
    catalog.CreateTable(config.TableName, config.Columns)

    start := time.Now()

    for i := 0; i < config.RowCount; i++ {
        record := map[string]string{}
        for _, col := range config.Columns {
            record[col.Name] = generateValue(col.Type)
        }

        id, _ := catalog.NextID(config.TableName)
        record["id"] = strconv.Itoa(id)

        offset, size, _ := s.Write(record)
        tree.AddKey(id, offset, size, t)
    }

    return time.Since(start)
}

func generateValue(colType string) string {
    switch colType {
    case "int":
        return strconv.Itoa(rand.Intn(10000))
    case "string":
        words := []string{"alice", "bob", "charlie", "dave", "eve", "frank"}
        return words[rand.Intn(len(words))]
    default:
        return "unknown"
    }
}
