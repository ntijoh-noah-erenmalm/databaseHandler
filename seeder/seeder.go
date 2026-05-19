package seeder

import (
    "databaseHandler/catalog"
    "databaseHandler/storage"
    "databaseHandler/tree"
    "math/rand"
    "strconv"
    "time"
		"fmt"
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
    for _, degree := range config.Degrees {
        for _, rows := range config.RowCounts {
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

            fmt.Printf("\ndegree=%d rows=%d insert=%v\n", degree, rows, insertDuration)
            for _, target := range searchTargets {
                start := time.Now()
                tree.Search(target, tr.Root)
                searchDuration := time.Since(start)
                fmt.Printf("  search target=%d duration=%v\n", target, searchDuration)
            }

            s.Close()
        }
    }
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
