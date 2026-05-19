package cli

import (
	"bufio"
	"databaseHandler/catalog"
	"databaseHandler/storage"
	"databaseHandler/tree"
	"fmt"
	"os"
	"strings"
)

func StartLoop(s *storage.Storage, t *tree.Tree) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("cli ready to use")

	for {
		fmt.Print(">	")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(input)

		if len(parts) == 0 {
			continue
		}

		switch parts[0]{
		case "exit":
			return
		case "create":
			if parts[1] == "table" {
				 CatalogEntry(parts)
			}
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
		case "insert":
			InsertEntry(parts, s, t)	
		case "select":
			SelectEntry(parts, s, t)
		default:
			fmt.Println("unknown, command:", parts[0])
		}
	}

}
