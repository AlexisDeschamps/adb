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

type GetSupporterOptions struct {
	ID                 int    `json:"id"`
	Order              int    `json:"order"`
	OrderField         string `json:"order_field"`
	Filter             string `json:"filter"`
	RestrictToBerkeley bool   `json:"restrict_to_berkeley"`
}

var validSupporterOrderFields = map[string]struct{}{
	"first_name":   struct{}{},
	"email":        struct{}{},
	"phone":        struct{}{},
	"date_sourced": struct{}{},
}

func (s Supporter) ToJSON() SupporterJSON {
	return buildSupporterJSONArray([]Supporter{s})[0]
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

func CreateUpdateSupporter(db *sqlx.DB, supporter Supporter) (int, error) {
	if supporter.ID == 0 {
		return CreateSupporter(db, supporter)
	}
	return UpdateSupporter(db, supporter)
}

func UpdateSupporter(db *sqlx.DB, supporter Supporter) (int, error) {
	if supporter.ID == 0 {
		return 0, errors.New("Cannot update supporter when ID == 0")
	}
	if supporter.Email == "" && supporter.Phone == "" {
		return 0, errors.New("Cannot update supporter if either email or phone isn't set")
	}

	_, err := db.NamedExec(`
UPDATE supporters
SET
  first_name = :first_name,
  last_name = :last_name,
  email = :email,
  phone = :phone,
  location_address1 = :location_address1,
  location_address2 = :location_address2,
  location_city = :location_city,
  location_state = :location_state,
  location_zip = :location_zip,
  source = :source,
  date_sourced = :date_sourced,
  requested_lawn_sign = :requested_lawn_sign,
  requested_poster = :requested_poster,
  voter = :voter,
  issue_housing = :issue_housing,
  issue_homelessness = :issue_homelessness,
  issue_climate = :issue_climate,
  issue_public_safety = :issue_public_safety,
  issue_police_accountability = :issue_police_accountability,
  issue_transit = :issue_transit,
  issue_economic_equality = :issue_economic_equality,
  issue_public_health = :issue_public_health,
  issue_animal_rights = :issue_animal_rights,
  interest_donate = :interest_donate,
  interest_attend_event = :interest_attend_event,
  interest_volunteer = :interest_volunteer,
  interest_host_event = :interest_host_event,
  notes = :notes,
  requires_followup = :requires_followup
WHERE
  id = :id`, supporter)

	if err != nil {
		return 0, errors.Wrap(err, "failed to update supporter")
	}
	return supporter.ID, nil
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

func validateGetSupporterOptions(o GetSupporterOptions) (GetSupporterOptions, error) {
	// Set defaults
	if o.Order == 0 {
		o.Order = DescOrder
	}
	if o.OrderField == "" {
		o.OrderField = "date_sourced"
	}

	if o.Order != DescOrder && o.Order != AscOrder {
		return GetSupporterOptions{}, errors.New("Supporters Range order must be ascending or descending")
	}
	if _, ok := validSupporterOrderFields[o.OrderField]; !ok {
		return GetSupporterOptions{}, errors.New("Supporter OrderField is not valid")
	}
	return o, nil
}

func CleanGetSupporterOptions(body io.Reader) (GetSupporterOptions, error) {
	var getSupporterOptions GetSupporterOptions
	err := json.NewDecoder(body).Decode(&getSupporterOptions)
	if err != nil {
		return GetSupporterOptions{}, err
	}
	getSupporterOptions, err = validateGetSupporterOptions(getSupporterOptions)
	if err != nil {
		return GetSupporterOptions{}, err
	}
	return getSupporterOptions, nil
}

func GetSupporter(db *sqlx.DB, options GetSupporterOptions) (Supporter, error) {
	if options.ID == 0 {
		return Supporter{}, errors.New("GetSupporter: Must pass in ID")
	}
	supporters, err := GetSupporters(db, options)
	if err != nil {
		return Supporter{}, errors.Wrap(err, "GetSupporter: error getting supporter from db")
	}
	if len(supporters) == 0 {
		return Supporter{}, errors.New("GetSupporter: supporter not found")
	}
	if len(supporters) > 1 {
		return Supporter{}, errors.New("GetSupporter: too many supporters for ID")
	}
	return supporters[0], nil
}

func GetSupportersJSON(db *sqlx.DB, options GetSupporterOptions) ([]SupporterJSON, error) {
	if options.ID != 0 {
		return nil, errors.New("GetSupportersJSON: Cannot include ID in options")
	}
	return getSupportersJSON(db, options)
}

func getSupportersJSON(db *sqlx.DB, options GetSupporterOptions) ([]SupporterJSON, error) {
	supporters, err := GetSupporters(db, options)
	if err != nil {
		return nil, err
	}
	return buildSupporterJSONArray(supporters), nil
}

func GetSupporters(db *sqlx.DB, options GetSupporterOptions) ([]Supporter, error) {
	// Redundant options validation
	var err error
	options, err = validateGetSupporterOptions(options)
	if err != nil {
		return nil, err
	}

	query := `
SELECT
  id,
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
FROM supporters s
`

	var queryArgs []interface{}

	if options.ID != 0 {
		// retrieve specific supporter
		query += " WHERE s.id = ? "
		queryArgs = append(queryArgs, options.ID)
	} else {
		// get multiple supporters
		query += ` WHERE NOT s.deleted `

		if options.RestrictToBerkeley {
			query += ` AND
  location_zip IN ('94701', '94702', '94703', '94704', '94705', '94706', '94707', '94708', '94709', '94710', '94712', '94720')
`
		}

		orderField := options.OrderField
		// Default to date_sourced if orderField isn't specified
		if orderField == "" {
			orderField = "date_sourced"
		}
		// Paranoid check b/c this could be a sql injection.
		if _, ok := validSupporterOrderFields[orderField]; !ok {
			return nil, errors.New("Invalid OrderField")
		}

		query += " ORDER BY " + options.OrderField
		if options.Order == DescOrder {
			query += " DESC "
		}
		// Add ID as second thing to order by so the order is the same
		// every time.
		query += ", s.id DESC "
	}

	var supporters []Supporter
	if err := db.Select(&supporters, query, queryArgs...); err != nil {
		return nil, errors.Wrapf(err, "fail to get supporters for uid %d", options.ID)
	}

	return supporters, nil
}

func buildSupporterJSONArray(supporters []Supporter) []SupporterJSON {
	var supportersJSON []SupporterJSON

	for _, s := range supporters {
		supportersJSON = append(supportersJSON, SupporterJSON{
			ID:                        s.ID,
			FirstName:                 s.FirstName,
			LastName:                  s.LastName,
			Email:                     s.Email,
			Phone:                     s.Phone,
			LocationAddress1:          s.LocationAddress1,
			LocationAddress2:          s.LocationAddress2,
			LocationCity:              s.LocationCity,
			LocationState:             s.LocationState,
			LocationZIP:               s.LocationZIP,
			Source:                    s.Source,
			DateSourced:               s.DateSourced,
			RequestedLawnSign:         s.RequestedLawnSign,
			RequestedPoster:           s.RequestedPoster,
			Voter:                     s.Voter,
			IssueHousing:              s.IssueHousing,
			IssueHomelessness:         s.IssueHomelessness,
			IssueClimate:              s.IssueClimate,
			IssuePublicSafety:         s.IssuePublicSafety,
			IssuePoliceAccountability: s.IssuePoliceAccountability,
			IssueTransit:              s.IssueTransit,
			IssueEconomicEquality:     s.IssueEconomicEquality,
			IssuePublicHealth:         s.IssuePublicHealth,
			IssueAnimalRights:         s.IssueAnimalRights,
			InterestDonate:            s.InterestDonate,
			InterestAttendEvent:       s.InterestAttendEvent,
			InterestVolunteer:         s.InterestVolunteer,
			InterestHostEvent:         s.InterestHostEvent,
			Notes:                     s.Notes,
			RequiresFollowup:          s.RequiresFollowup,
		})
	}
	return supportersJSON
}

func DeleteSupporter(db *sqlx.DB, supporterID int) error {
	if supporterID == 0 {
		return errors.New("supporterID cannot be 0")
	}
	query := `
UPDATE supporters
SET
  deleted = true
WHERE
  id = ?
LIMIT 1 -- id is a primary key, so there should only ever be one, but do this just in case.
`
	result, err := db.Exec(query, supporterID)
	if err != nil {
		return errors.Wrap(err, "error deleting supporter")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "error getting rowsaffected when deleting supporter")
	}
	if rowsAffected == 0 {
		return errors.Errorf("supporter for id does not exist: %d", supporterID)
	}
	return nil
}
