package cli 

import (
	"databaseHandler/catalog"
	"strings"
	"fmt"
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

