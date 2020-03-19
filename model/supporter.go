package model

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Supporter struct {
	ID                int    `db:"id"`
	FirstName         string `db:"first_name"`
	LastName          string `db:"last_name"`
	Email             string `db:"email"`
	Phone             string `db:"phone"`
	LocationAddress1  string `db:"location_address1"`
	LocationAddress2  string `db:"location_address2"`
	LocationCity      string `db:"location_city"`
	LocationState     string `db:"location_state"`
	LocationZIP       string `db:"location_zip"`
	Source            string `db:"source"`
	DateSourced       string `db:"date_sourced"`
	RequestedLawnSign bool   `db:"requested_lawn_sign"`
	RequestedPoster   bool   `db:"requested_poster"`
	Voter             bool   `db:"voter"`

	IssueHousing              bool `db:"issue_housing"`
	IssueHomelessness         bool `db:"issue_homelessness"`
	IssueClimate              bool `db:"issue_climate"`
	IssuePublicSafety         bool `db:"issue_public_safety"`
	IssuePoliceAccountability bool `db:"issue_police_accountability"`
	IssueTransit              bool `db:"issue_transit"`
	IssueEconomicEquality     bool `db:"issue_economic_equality"`
	IssuePublicHealth         bool `db:"issue_public_health"`
	IssueAnimalRights         bool `db:"issue_animal_rights"`

	InterestDonate      bool `db:"interest_donate"`
	InterestAttendEvent bool `db:"interest_attend_event"`
	InterestVolunteer   bool `db:"interest_volunteer"`
	InterestHostEvent   bool `db:"interest_host_event"`

	Notes string `db:"notes"`

	RequiresFollowup bool `db:"requires_followup"`

	// Canvasser string `db:"canvasser"`
	// CanvassLeader string `db:"canvass_leader"`
}

type SupporterJSON struct {
	ID                int    `json:"id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	Phone             string `json:"phone"`
	LocationAddress1  string `json:"location_address1"`
	LocationAddress2  string `json:"location_address2"`
	LocationCity      string `json:"location_city"`
	LocationState     string `json:"location_state"`
	LocationZIP       string `json:"location_zip"`
	Source            string `json:"source"`
	DateSourced       string `json:"date_sourced"`
	RequestedLawnSign bool   `json:"requested_lawn_sign"`
	RequestedPoster   bool   `json:"requested_poster"`
	Voter             bool   `json:"voter"`

	IssueHousing              bool `json:"issue_housing"`
	IssueHomelessness         bool `json:"issue_homelessness"`
	IssueClimate              bool `json:"issue_climate"`
	IssuePublicSafety         bool `json:"issue_public_safety"`
	IssuePoliceAccountability bool `json:"issue_police_accountability"`
	IssueTransit              bool `json:"issue_transit"`
	IssueEconomicEquality     bool `json:"issue_economic_equality"`
	IssuePublicHealth         bool `json:"issue_public_health"`
	IssueAnimalRights         bool `json:"issue_animal_rights"`

	InterestDonate      bool `json:"interest_donate"`
	InterestAttendEvent bool `json:"interest_attend_event"`
	InterestVolunteer   bool `json:"interest_volunteer"`
	InterestHostEvent   bool `json:"interest_host_event"`

	Notes string `json:"notes"`

	RequiresFollowup bool `json:"requires_followup"`

	// Canvasser string `json:"canvasser"`
	// CanvassLeader string `json:"canvass_leader"`
}

func CleanSupporterData(body io.Reader) (Supporter, error) {
	var supporterJSON SupporterJSON
	err := json.NewDecoder(body).Decode(&supporterJSON)
	if err != nil {
		return Supporter{}, err
	}

	supporter := Supporter{
		ID:                        supporterJSON.ID,
		FirstName:                 strings.TrimSpace(supporterJSON.FirstName),
		LastName:                  strings.TrimSpace(supporterJSON.LastName),
		Email:                     strings.TrimSpace(supporterJSON.Email),
		Phone:                     strings.TrimSpace(supporterJSON.Phone),
		LocationAddress1:          strings.TrimSpace(supporterJSON.LocationAddress1),
		LocationAddress2:          strings.TrimSpace(supporterJSON.LocationAddress2),
		LocationCity:              strings.TrimSpace(supporterJSON.LocationCity),
		LocationState:             strings.TrimSpace(supporterJSON.LocationState),
		LocationZIP:               strings.TrimSpace(supporterJSON.LocationZIP),
		Source:                    strings.TrimSpace(supporterJSON.Source),
		DateSourced:               strings.TrimSpace(supporterJSON.DateSourced),
		RequestedLawnSign:         supporterJSON.RequestedLawnSign,
		RequestedPoster:           supporterJSON.RequestedPoster,
		Voter:                     supporterJSON.Voter,
		IssueHousing:              supporterJSON.IssueHousing,
		IssueHomelessness:         supporterJSON.IssueHomelessness,
		IssueClimate:              supporterJSON.IssueClimate,
		IssuePublicSafety:         supporterJSON.IssuePublicSafety,
		IssuePoliceAccountability: supporterJSON.IssuePoliceAccountability,
		IssueTransit:              supporterJSON.IssueTransit,
		IssueEconomicEquality:     supporterJSON.IssueEconomicEquality,
		IssuePublicHealth:         supporterJSON.IssuePublicHealth,
		IssueAnimalRights:         supporterJSON.IssueAnimalRights,
		InterestDonate:            supporterJSON.InterestDonate,
		InterestAttendEvent:       supporterJSON.InterestAttendEvent,
		InterestVolunteer:         supporterJSON.InterestVolunteer,
		InterestHostEvent:         supporterJSON.InterestHostEvent,
		Notes:                     strings.TrimSpace(supporterJSON.Notes),
		RequiresFollowup:          supporterJSON.RequiresFollowup,
	}
	return supporter, nil
}

func CreateSupporter(db *sqlx.DB, supporter Supporter) (int, error) {
	if supporter.ID != 0 {
		return 0, errors.New("Cannot create supporter when ID != 0")
	}
	if supporter.Email == "" && supporter.Phone == "" {
		return 0, errors.New("Cannot create supporter if either email or phone isn't set")
	}

	result, err := db.NamedExec(`
INSERT INTO supporters (
  first_name,
  last_name,
  email,
  phone,
  location_address1,
  location_address2,
  location_city,
  location_state,
  location_zip,
  source,
  date_sourced,
  requested_lawn_sign,
  requested_poster,
  voter,
  issue_housing,
  issue_homelessness,
  issue_climate,
  issue_public_safety,
  issue_police_accountability,
  issue_transit,
  issue_economic_equality,
  issue_public_health,
  issue_animal_rights,
  interest_donate,
  interest_attend_event,
  interest_volunteer,
  interest_host_event,
  notes,
  requires_followup
) VALUES (

  :first_name,
  :last_name,
  :email,
  :phone,
  :location_address1,
  :location_address2,
  :location_city,
  :location_state,
  :location_zip,
  :source,
  :date_sourced,
  :requested_lawn_sign,
  :requested_poster,
  :voter,
  :issue_housing,
  :issue_homelessness,
  :issue_climate,
  :issue_public_safety,
  :issue_police_accountability,
  :issue_transit,
  :issue_economic_equality,
  :issue_public_health,
  :issue_animal_rights,
  :interest_donate,
  :interest_attend_event,
  :interest_volunteer,
  :interest_host_event,
  :notes,
  :requires_followup

)
`, supporter)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not create supporter: %s", supporter.Email)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, errors.Wrapf(err, "Could not get LastInsertId for supporter %s", supporter.Email)
	}
	return int(id), nil
}
