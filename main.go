package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type alias struct {
	name         string
	destinations []string
}

// writeAliases writes aliases in alias: destination1, destination2, ... format to w
// Aliases are sorted first.
func writeAliases(aliases []alias, w io.Writer) error {
	for _, alias := range aliases {
		_, err := fmt.Fprintf(w, "%s: %s\n", alias.name, strings.Join(alias.destinations, ", "))
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s postgres://user:pass@dbserver/dbname", os.Args[0])
	}
	sqlx.MustConnect("postgres", os.Args[1])

	var aliases []alias

	err := writeAliases(aliases, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
