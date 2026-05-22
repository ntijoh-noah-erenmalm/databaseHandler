package lineardb

import (
	"bufio"
	"databaseHandler/catalog"
	"databaseHandler/storage"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type List struct {
	Root *Node
	Last *Node
}

type Node struct {
	Value int
	Offset int64
	Size int
	Next *Node
}

func LinearStartLoop(s *storage.Storage, list *List) {
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Println("linear db ready to use")

    for {
        fmt.Print("> ")
        scanner.Scan()
        input := strings.TrimSpace(scanner.Text())
        parts := strings.Fields(input)

        if len(parts) == 0 {
            continue
        }

        switch parts[0] {
        case "exit":
            fmt.Println("bye")
            return
        case "create":
            if len(parts) > 1 && parts[1] == "table" {
                LinearCatalogEntry(parts)
            }
        case "insert":
            if len(parts) > 1 && parts[1] == "into" {
                LinearInsertEntry(parts, s, list)
            }
        case "select":
            LinearSelectEntry(parts, s, list)
        case "show":
            if len(parts) > 1 && parts[1] == "tables" {
                tables := catalog.GetTables()
                if len(tables) == 0 {
                    fmt.Println("no tables found")
                    continue
                }
                for _, t := range tables {
                    fmt.Println(t.Name)
                    for _, col := range t.Columns {
                        fmt.Println(" ", col.Name, col.Type)
                    }
                }
            }
        default:
            fmt.Println("unknown command:", parts[0])
        }
    }
}

func LinearCatalogEntry(command []string) {
    tableName := command[2]
    raw := strings.Join(command[3:], " ")
    raw = strings.Trim(raw, "()")
    columnDefs := strings.Split(raw, ",")
    var columns []catalog.Column
    for _, def := range columnDefs {
        def = strings.TrimSpace(def)
        colParts := strings.Fields(def)
        if len(colParts) >= 2 {
            columns = append(columns, catalog.Column{
                Name: colParts[0],
                Type: colParts[1],
            })
        }
    }
    table, err := catalog.CreateTable(tableName, columns)
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    fmt.Println("created table", table.Name)
}

func LinearInsertEntry(command []string, s *storage.Storage, list *List) {
    tableName := command[2]
    raw := strings.Join(command[4:], " ")
    raw = strings.Trim(raw, "()")
    valueParts := strings.Split(raw, ",")

    table, err := catalog.GetTable(tableName)
    if err != nil {
        fmt.Println("error:", err)
        return
    }

    if len(valueParts) != len(table.Columns)-1 {
        fmt.Println("error: wrong number of values")
        return
    }

    record := map[string]string{}
    valueIndex := 0
    for _, col := range table.Columns {
        if col.Name == "id" {
            continue
        }
        record[col.Name] = strings.TrimSpace(valueParts[valueIndex])
        valueIndex++
    }

    id, err := catalog.NextID(tableName)
    if err != nil {
        fmt.Println("error:", err)
        return
    }
    record["id"] = strconv.Itoa(id)

    offset, size, err := s.Write(record)
    if err != nil {
        fmt.Println("error writing to storage:", err)
        return
    }

    ListAdd(id, offset, size, list)
    fmt.Println("inserted into", tableName, "with id", id)
}

func LinearSelectEntry(command []string, s *storage.Storage, list *List) {
    // select * from users where id = 1
    tableName := command[3]
    _, err := catalog.GetTable(tableName)
    if err != nil {
        fmt.Println("error:", err)
        return
    }

    id, err := strconv.Atoi(command[7])
    if err != nil {
        fmt.Println("invalid id:", command[7])
        return
    }

    node := ListSearch(id, list)
    if node == nil {
        fmt.Println("record not found")
        return
    }

    record, err := s.Read(node.Offset, node.Size)
    if err != nil {
        fmt.Println("error reading from storage:", err)
        return
    }

    fmt.Println(record)
}



func NewList() List {
	return List{Root: nil, Last: nil}
}

func ListAdd(value int, offset int64, size int, list *List){
	newNode := &Node{			
		Value: value,
		Offset: offset,
		Size: size,
		Next: nil,
	}

	if list.Root == nil {
		list.Root = newNode
		list.Last = newNode
		return
	}

	

	list.Last.Next = newNode
	list.Last = newNode
}

func ListSearch(value int, list *List) *Node {
	currentNode := list.Root
	for currentNode != nil {
		if currentNode.Value == value {
			return currentNode
		}
		currentNode = currentNode.Next
	}	
	return nil
}

func LinearRunBenchmark(config LinearBenchmarkConfig) {
    file, _ := os.Create("linear_benchmark_results.txt")
    defer file.Close()
    fmt.Fprintln(file, "rows,insert_ms,search_target,search_ns")

    for _, rows := range config.RowCounts {
        os.Remove("data.db")
        os.Remove("catalog.json")

        s, _ := storage.NewStorage()
        list := NewList()

        catalog.CreateTable(config.TableName, config.Columns)

        start := time.Now()
        for i := 0; i < rows; i++ {
            record := map[string]string{}
            for _, col := range config.Columns {
                record[col.Name] = generateValue(col.Type)
            }
            id, _ := catalog.NextID(config.TableName)
            record["id"] = strconv.Itoa(id)
            offset, size, _ := s.Write(record)
            ListAdd(id, offset, size, &list)
        }
        insertDuration := time.Since(start)

        searchTargets := []int{
            0,
            rows - 1,
            rand.Intn(rows),
        }

        for _, target := range searchTargets {
            start := time.Now()
            ListSearch(target, &list)
            searchDuration := time.Since(start)
            fmt.Fprintf(file, "%d,%d,%d,%d\n",
                rows,
                insertDuration.Milliseconds(),
                target,
                searchDuration.Nanoseconds(),
            )
        }

        s.Close()
    }
    fmt.Println("linear benchmark done :D")
}

type LinearBenchmarkConfig struct {
    RowCounts  []int
    Columns    []catalog.Column
    TableName  string
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
