package sendy_sync

import (
	"bytes"
	"fmt"
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

type SendyLists struct {
	AllADB                                 string
	PublicHealthOnly                       string
	PublicHealthClimate                    string
	PublicHealthHousingHomelessness        string
	PublicHealthClimateHousingHomelessness string
	ClimateOnly                            string
	ClimateHousingHomelessness             string
	HousingHomelessnessOnly                string
}

// Choose which lists to sync to.
type syncOptions struct {
	// if allADB is true, then the other options are ignored.
	allADB bool

	// send to whatever combination of issues that are set (e.g.
	// if only publicHealth is set, add them to the public
	// health-only list).
	publicHealth        bool
	climate             bool
	housingHomelessness bool
}

func StartSupportersSendySync(db *sqlx.DB, sendyAPIKey string, lists SendyLists) {
	// Iterate through (publicHealth, climate,
	// housingHomelessness) subsets using a bitfield. E.g. 001 is
	// publicHealth, 010 is climate, 011 is publicHealth and
	// climate, etc.
	publicHealthMask := 1
	climateMask := 1 << 1
	housingHomelessnessMask := 1 << 2

	maxSet := 1 << 3
	subsetToList := map[int]string{
		publicHealthMask:                                         lists.PublicHealthOnly,
		publicHealthMask | climateMask:                           lists.PublicHealthClimate,
		publicHealthMask | housingHomelessnessMask:               lists.PublicHealthHousingHomelessness,
		publicHealthMask | climateMask | housingHomelessnessMask: lists.PublicHealthClimateHousingHomelessness,
		climateMask:                           lists.ClimateOnly,
		climateMask | housingHomelessnessMask: lists.ClimateHousingHomelessness,
		housingHomelessnessMask:               lists.HousingHomelessnessOnly,
	}
	for {
		if lists.AllADB != "" {
			log.Println("Starting supporters to sendy sync")
			syncSupportersToSendy(db, sendyAPIKey, lists.AllADB, syncOptions{allADB: true})
			log.Println("Finished supporters to sendy sync")
		} else {
			log.Println("Not syncing supporters to all ADB list")
		}
		for subset := 1; subset < maxSet; subset++ {
			list, ok := subsetToList[subset]
			if !ok {
				// Fatal because this is a programmer
				// error, every subset should exist in
				// subsetToList.
				log.Fatalf("Missing sendy sync list: %d", subset)
			}
			options := syncOptions{}
			log.Println("Supporters sync: Processing the following subset:")
			if publicHealthMask&subset > 0 {
				log.Println("Public Health")
				options.publicHealth = true
			}
			if climateMask&subset > 0 {
				log.Println("Climate")
				options.climate = true
			}
			if housingHomelessnessMask&subset > 0 {
				log.Println("Housing & Homelessness")
				options.housingHomelessness = true
			}
			if list == "" {
				log.Println("Not syncing this subset because the list id is empty.")
				continue
			}
			syncSupportersToSendy(db, sendyAPIKey, list, options)
			log.Println("Finished syncing subset to sendy")
		}
		time.Sleep(6 * time.Minute)
	}
}

func syncSupportersToSendy(db *sqlx.DB, sendyAPIKey string, sendyList string, options syncOptions) {
	// First, get everyone in ADB that isn't already in sendy.

	var whereStr string
	// Only add to the "all adb" list if the user hasn't
	// been added to any lists, or the user hasn't already
	// been added to the "all adb" list.
	var isIssueSublist bool
	if !options.allADB {
		isIssueSublist = true

		whereStr += fmt.Sprintf(`
  AND issue_public_health = %v
  AND issue_climate = %v `, options.publicHealth, options.climate)
		if options.housingHomelessness {
			whereStr += " AND (issue_housing = true OR issue_homelessness = true) "
		} else {
			whereStr += " AND (issue_housing = false AND issue_homelessness = false) "
		}
	}

	query := fmt.Sprintf(`
SELECT
  id,
  first_name,
  last_name,
  email
FROM
  supporters
WHERE
  email != ''
  AND id NOT IN (
    SELECT
      supporter_id as id
    FROM supporters_sendy_sync
    WHERE
      is_issue_sublist = %v

  )
  %s
LIMIT 1000
`, isIssueSublist, whereStr)

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

		err = updateSupportersSendySync(db, supporter, sendyList, emailExistsInSendy, isIssueSublist)
		if err != nil {
			log.Printf("Error updating supporters_sendy_sync db: %s", err)
			return
		}

	}

}

func updateSupportersSendySync(db *sqlx.DB, supporter supporterBasic, sendyListID string, emailExistsInSendy bool, isIssueSublist bool) error {
	query := `
INSERT INTO supporters_sendy_sync (
  supporter_id,
  sendy_list_id,
  email,
  is_issue_sublist,
  sync_status,
  sync_timestamp
) VALUES (
  ?,
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
		isIssueSublist,
		syncStatus)
	if err != nil {
		return errors.Wrapf(err, "Error adding supporter to supporters_sendy_sync: %d", supporter.ID)
	}
	return nil
}
