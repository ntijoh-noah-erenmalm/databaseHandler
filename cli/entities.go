package cli

import (
	"databaseHandler/catalog"
	"databaseHandler/storage"
	"databaseHandler/tree"
	"fmt"
	"strings"
	"strconv"
)

func CatalogEntry(command []string) {
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
	table, error := catalog.CreateTable(tableName, columns)
	if error != nil {
		fmt.Println("error:", error)
		return 
	}
	fmt.Println("created table", table.Name)
}

func InsertEntry(command []string, s *storage.Storage, t *tree.Tree) {
	// insert into tableName values (x, x...)
	tableName := command[2]
	raw := strings.Join(command[4:], " ")
	raw = strings.Trim(raw, "()")
	valueParts := strings.Split(raw, ",")

	table, error := catalog.GetTable(tableName)
	if error != nil {
		fmt.Println("error:", error)
		return
	}

	if len(valueParts) != len(table.Columns) - 1 {
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

	id, error := catalog.NextID(tableName)
	if error != nil {
		fmt.Println("error:",error)
		return
	}

	record["id"] = strconv.Itoa(id)

	offset, size, error := s.Write(record)
	if error != nil {
		fmt.Println("error writing to storage:", error)
		return
	}
	fmt.Println("wrote at offset", offset, "size", size)

	// insert into tree
	tree.AddKey(id, offset, size, t)
	fmt.Println("inserted into", tableName, record)
}

func SelectEntry(command []string, s *storage.Storage, t *tree.Tree) {
	// SELECT * FROM users WHERE id = 1
	// command = ["select", "*", "from", "users", "where", "id", "=", "1"]

	tableName := command[3]

	_, error := catalog.GetTable(tableName)
	if error != nil {
		fmt.Println("error:", error)
		return
	}

	id,_ := strconv.Atoi(command[7])

	key := tree.Search(id, t.Root)
	if key == nil {
		fmt.Println("record not ofund")
		return
	}

	record, error := s.Read(key.Offset, key.Size)
	if error != nil {
		fmt.Println("error reading from storage:", error)
		return
	}

	fmt.Println(record)
}
