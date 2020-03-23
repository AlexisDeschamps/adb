package sendy_sync

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	SyncStatus_Synced = 1
	SyncStatus_Error  = 2
)

type supporterBasic struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string `db:"email"`
}

func StartSupportersSendySync(db *sqlx.DB, sendyAPIKey, sendyAllADBList string) {
	for {
		log.Println("Starting supporters to sendy sync")
		syncSupportersToSendy(db, sendyAPIKey, sendyAllADBList)
		log.Println("Finished supporters to sendy sync")
		time.Sleep(6 * time.Minute)
	}
}

func syncSupportersToSendy(db *sqlx.DB, sendyAPIKey, sendyList string) {
	// First, get everyone in ADB that isn't already in sendy.

	query := `
SELECT
  id,
  first_name,
  last_name,
  s.email as email
FROM
  supporters s
LEFT JOIN
  supporters_sendy_sync ss
ON s.id = ss.supporter_id
WHERE
  ss.supporter_id IS NULL
  AND s.email != ''
LIMIT 1000
`
	var supporters []supporterBasic
	if err := db.Select(&supporters, query); err != nil {
		log.Printf("Could not query supporters: %s\n", err)
		return
	}

	for _, supporter := range supporters {
		log.Printf("Attempting to sync supporter to sendy: (%d, %s)\n", supporter.ID, supporter.Email)
		var name string
		trimFirstName := strings.TrimSpace(supporter.FirstName)
		trimLastName := strings.TrimSpace(supporter.LastName)
		if trimFirstName != "" && trimLastName != "" {
			name = trimFirstName + " " + trimLastName
		} else if trimFirstName != "" {
			name = trimFirstName
		} else if trimLastName != "" {
			name = trimLastName
		}
		formData := url.Values{
			"api_key": {sendyAPIKey},
			"name":    {name},
			"email":   {supporter.Email},
			"list":    {sendyList},
			"boolean": {"true"},
		}
		resp, err := http.PostForm("https://sendy.wayneformayor.com/subscribe", formData)
		if err != nil {
			log.Printf("Error subscribing email %s for supporter %d: %s\n", supporter.Email, supporter.ID, err)
			continue
		}
		if resp.StatusCode != 200 {
			log.Printf("Error: Bad response code from sendy during subscriber sync for (%d, %s): %d\n", supporter.ID, supporter.Email, resp.StatusCode)
			continue
		}

		var emailExistsInSendy bool
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading sendy response: %v\n", err)
			return
		}
		if bytes.Equal(body, []byte("1")) ||
			bytes.Equal(body, []byte("Already subscribed.")) {
			log.Printf("Received successfuly response from sendy\n")
			emailExistsInSendy = true
		} else {
			log.Printf("Error adding email to sendy: (%d, %s)\n", supporter.ID, supporter.Email)
		}

		err = updateSupportersSendySync(db, supporter, sendyList, emailExistsInSendy)
		if err != nil {
			log.Printf("Error updating supporters_sendy_sync db: %s", err)
			return
		}

	}

}

func updateSupportersSendySync(db *sqlx.DB, supporter supporterBasic, sendyListID string, emailExistsInSendy bool) error {
	query := `
INSERT INTO supporters_sendy_sync (
  supporter_id,
  sendy_list_id,
  email,
  sync_status,
  sync_timestamp
) VALUES (
  ?,
  ?,
  ?,
  ?,
  UTC_TIMESTAMP()
)
`
	syncStatus := SyncStatus_Synced
	if !emailExistsInSendy {
		syncStatus = SyncStatus_Error
	}
	_, err := db.Exec(query,
		supporter.ID,
		sendyListID,
		supporter.Email,
		syncStatus)
	if err != nil {
		return errors.Wrapf(err, "Error adding supporter to supporters_sendy_sync: %d", supporter.ID)
	}
	return nil
}
