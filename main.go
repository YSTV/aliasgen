package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type alias struct {
	name         string
	destinations []string
}

// writeAliases writes aliases in alias: destination1, destination2, ... format to w
func writeAliases(aliases []alias, w io.Writer) error {
	for _, alias := range aliases {
		_, err := fmt.Fprintf(w, "%s: %s\n", alias.name, strings.Join(alias.destinations, ", "))
		if err != nil {
			return err
		}
	}
	return nil
}

// Result of a join between officerships, officership_members and members.
type officership struct {
	Name               string      // The name of the officership
	Alias              string      `db:"email_alias"` // The email alias for this officership
	StartDate, EndDate pq.NullTime // Start and end of incumbant officership, if it exists
	Email              string      // The email of the incumbant officer, if they exist
	VacantAlias        string      // The Alias of the officership to which mail should be sent if this one is vacant.
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s postgres://user:pass@dbserver/dbname", os.Args[0])
	}
	db := sqlx.MustConnect("postgres", os.Args[1])
	// Loop through officerships. If filled, use member's email. If vacant use team head's email. If vacant and team head,
	// send to station director.
	var officerships []officership
	err := db.Select(
		&officerships,
		`SELECT officerships.name, officerships.email_alias, ifvacantofficership.email_alias AS vacant_alias, officerships.email_alias, member_officerships.start_date, member_officerships.end_date, members.server_name
			FROM officerships
			LEFT JOIN member_officerships ON (
				officerships.name=member_officerships.officership_name AND member_officerships.start_date < NOW() AND (
					member_officerships.end_date > NOW() OR member_officerships.end_date IS NULL
				)
			) LEFT JOIN officerships AS ifvacantofficership ON (
				officerships.if_unfilled=ifvacantofficership.name
			) LEFT JOIN members ON (
				member_officerships.member_id=members.id
			)
			WHERE officerships.is_current=true;`,
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, offShip := range officerships {
		fmt.Println(offShip)
	}

	var aliases []alias

	err = writeAliases(aliases, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
