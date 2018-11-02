package main

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type aliases map[string][]string

// writeAliases writes aliases in alias: destination1, destination2, ... format to w
func writeAliases(aliases aliases, w io.Writer) error {
	var sortedAliases []string
	for alias := range aliases {
		sortedAliases = append(sortedAliases, alias)
	}
	sort.Strings(sortedAliases)

	for _, alias := range sortedAliases {
		_, err := fmt.Fprintf(w, "%s: %s\n", alias, strings.Join(aliases[alias], ", "))
		if err != nil {
			return err
		}
	}
	return nil
}

// Result of a join between officerships, officership_members and members.
type officership struct {
	Name        string         // The name of the officership
	Alias       string         `db:"email_alias"` // The email alias for this officership
	StartDate   pq.NullTime    `db:"start_date"`
	EndDate     pq.NullTime    `db:"end_date"`      // Start and end of incumbent officership, if it exists
	ServerName  sql.NullString `db:"server_name"`   // The email of the incumbent officer, if they exist
	VacantAlias sql.NullString `db:"vacant_alias"`  // The Alias of the officership to which mail should be sent if this one is vacant.
	Email       sql.NullString `db:"email_address"` // The email set on the user's profile
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s postgres://user:pass@dbserver/dbname", os.Args[0])
	}
	db := sqlx.MustConnect("postgres", os.Args[1])
	var officerships []officership
	err := db.Select(
		&officerships,
		`SELECT officerships.name, officerships.email_alias, ifvacantofficership.email_alias AS vacant_alias, officerships.email_alias, member_officerships.start_date, member_officerships.end_date, members.server_name, members.email_address
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

	theAliases := make(aliases)

	// Loop through officerships. If filled, use member's server name (local mailbox) or profile email. If unfilled, use if_unfilled value from database.
	for _, offShip := range officerships {
		var email string
		if offShip.ServerName.Valid && offShip.ServerName.String != "" {
			email = offShip.ServerName.String
		} else if offShip.Email.Valid && offShip.Email.String != "" {
			log.Warnf("Officer %s has no server name, defaulting to using email", offShip.Name)
			email = offShip.Email.String
		} else {
			log.Warnf("Officer %s is unfilled or has no email or server name set", offShip.Name)
                        email = offShip.VacantAlias.String
		}
		theAliases[offShip.Alias] = append(theAliases[offShip.Alias], email)
	}

	err = writeAliases(theAliases, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
