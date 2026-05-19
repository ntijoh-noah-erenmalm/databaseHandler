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
	NextID int
}

type Catalog struct {
	Tables []Table `json:"tables"`
}

func CreateTable(tableName string, columns []Column) (Table, error) {
	// Create table tableName (name type, name2 type2...)
		
	// 1. check if catalog.bin exists, if not create it
	catalog := loadCatalog()
	// 2. check if table name already exists in catalog.bin
	// 3. if it does, return error
	for _, t := range catalog.Tables {
		if t.Name == tableName {
			return Table{}, fmt.Errorf("table %s already exists", tableName)
		}
	}

	for _, col := range columns {
		if col.Name == "id" {
			return Table{}, fmt.Errorf("id is a reserved column name")
		}
	}

	columns = append([]Column{{Name: "id", Type: "int"}}, columns...)

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

func GetTable(name string) (Table, error) {
    catalog := loadCatalog()
    for _, t := range catalog.Tables {
        if t.Name == name {
            return t, nil
        }
    }
    return Table{}, fmt.Errorf("table %s not found", name)
}

func NextID(tableName string) (int, error) {
    cat := loadCatalog()
    for i, t := range cat.Tables {
        if t.Name == tableName {
            id := t.NextID
            cat.Tables[i].NextID++
            saveCatalog(cat)
            return id, nil
        }
    }
    return 0, fmt.Errorf("table %s not found", tableName)
}
