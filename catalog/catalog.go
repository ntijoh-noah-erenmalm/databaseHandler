package catalog

import (
	"encoding/json"
	"fmt"
	"os"
)


type Column struct {
	Name string
	Type string
}

type Table struct {
	Name string
	Columns []Column
}

type Catalog struct {
	Tables []Table `json:"tables"`
}

func CreateTable(tableName string, columns []Column) (Table, error) {


	// 1. check if catalog.bin exists, if not create it
	catalog := loadCatalog()
	// 2. check if table name already exists in catalog.bin
	// 3. if it does, return error
	for _, t := range catalog.Tables {
		if t.Name == tableName {
			return Table{}, fmt.Errorf("table %s already exists", tableName)
		}
	}
	// 4. if it doesn't, write table to catalog.bin
	table := Table{Name: tableName, Columns: columns}
	catalog.Tables = append(catalog.Tables, table)
	saveCatalog(catalog)
	// 5. return the table
	return table, nil
}

func loadCatalog() Catalog {
	
	data, err := os.ReadFile("catalog.json")
	if err != nil {
		return Catalog{Tables: []Table{}}
	}

    var catalog Catalog
    json.Unmarshal(data, &catalog)
    return catalog
}

func saveCatalog(catalog Catalog) {
    data, _ := json.Marshal(catalog)
    os.WriteFile("catalog.json", data, 0666)
}

func GetTables() []Table {
    catalog := loadCatalog()
    return catalog.Tables
}
