package models

import (
	"database/sql"
	"fmt"
)

type Organization struct {
	id                   int
	Name                 string
	OrganizationTimezone string
	LogoPath             string
	APIKey               sql.NullString
	Email                sql.NullString
	Onboarding           string
	CreatedAt            string
	UpdatedAt            string
	IsDemo               bool
}

// ** CREATE **
// Add adds a new organization to the database there can only be one organization in the database.
// It returns any error encountered.
func (o Organization) Add(name, organizationTimezone, logoPath, apiKey, email string) (int, error) {
	db := NewDB()

	// Delete any existing organization
	deleteSqlStatement := `
		DELETE FROM
			organization
	`
	_, err := db.Exec(deleteSqlStatement)
	if err != nil {
		return 0, fmt.Errorf("unable to delete existing organization: %v", err)
	}

	// Insert new organization with the is_demo flag set to false
	sqlStatement := `
		INSERT INTO
			organization
		VALUES
			(
			NULL,
			$1,
			$2,
			$3,
			$4,
			$5,
			true,
			datetime('now'),
			datetime('now'),
			false
			)
		RETURNING id`
	var ID int
	err = db.QueryRow(sqlStatement, name, organizationTimezone, logoPath, apiKey, email).Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("unable to insert: %v", err)
	}

	return ID, nil
}

// SetAdmin sets the admin for the organization given the user id.
// It returns any error encountered.
func (o Organization) SetAdmin(userID int) error {
	db := NewDB()

	sqlStatement := `
		INSERT INTO
			is_admin
		VALUES
			(
			NULL,
			$1
			)`
	_, err := db.Exec(sqlStatement, userID)
	if err != nil {
		return fmt.Errorf("unable to insert: %v", err)
	}

	return nil
}

// ** READ **
// GetStatus retrieves the status of the organization.
// It returns a string and any error encountered.
func (o Organization) GetStatus() (string, error) {
	db := NewDB()
	var status string
	sqlStatement := `
	SELECT
		is_demo
	FROM
		organization
	LIMIT
		1;
	`

	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&status); err {
	case sql.ErrNoRows:
		return "", err
	case nil:
		return status, nil
	default:
		return "", err
	}
}

// Get retrieves an organization from the database.
// It returns an Organization object and any error encountered.
func (o Organization) Get(orgID int) (Organization, error) {
	var id int
	var Name string
	var OrganizationTimezone string
	var LogoPath sql.NullString
	var CreatedAt string
	var UpdatedAt string
	var APIKey sql.NullString
	var Email sql.NullString
	var Onboarding string
	var IsDemo bool

	db := NewDB()
	sqlStatement := `
		SELECT 
			* 
		FROM 
			organization 
		WHERE 
			id = $1;`

	row := db.QueryRow(sqlStatement, orgID)
	switch err := row.Scan(&id, &Name, &OrganizationTimezone, &LogoPath, &APIKey, &Email, &Onboarding, &CreatedAt, &UpdatedAt, &IsDemo); err {
	case sql.ErrNoRows:
		return Organization{}, err
	case nil:
		organization := Organization{id, Name, OrganizationTimezone, LogoPath.String, APIKey, Email, Onboarding, CreatedAt, UpdatedAt, IsDemo}
		return organization, err
	default:
		return Organization{}, err
	}
}

// OrganizationExists checks if any organization exists in the database.
// It returns true if the organization exists and false if it does not.
func (o Organization) OrganizationExists() bool {
	db := NewDB()
	var count int

	sqlStatement := `
	SELECT
	COUNT(*)
	FROM
		 organization
	`
	err := db.QueryRow(sqlStatement).Scan(&count)
	if err != nil {
		return false
	}
	if count > 0 {
		return true
	}
	return false
}

// GetOrganizationID returns the id of the only organization in the database.

// It returns the organization id and any error encountered.
func (o Organization) GetOrganizationID() (int, error) {
	db := NewDB()
	var orgID int

	sqlStatement := `
	SELECT
		id
	FROM
		 organization
	`
	err := db.QueryRow(sqlStatement).Scan(&orgID)
	if err != nil {
		return 0, err
	}
	return orgID, nil
}

// GetAdminID finds the user id of the admin for the organization.
// It returns the user id of the admin and any error encountered.
func (o Organization) GetAdminID() (int, error) {
	db := NewDB()
	var userID int

	sqlStatement := `
	SELECT
		user_id
	FROM
		is_admin
	`
	err := db.QueryRow(sqlStatement).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// CheckAPIKey checks if the API key is not null for the organization.
// It returns true if the API key is not null and an error if one occurs during the process.

func (o Organization) CheckAPIKey() (bool, error) {
	db := NewDB()
	var apiKey sql.NullString // Change the type to sql.NullString

	sqlStatement := `
	SELECT
		api_key
	FROM
		organization
	`
	err := db.QueryRow(sqlStatement).Scan(&apiKey)
	if err != nil {
		return false, fmt.Errorf("failed to execute query: %w", err)
	}
	if apiKey.Valid { // Check if the value is not null
		return true, nil
	}
	return false, nil
}

// ** UPDATE **
// StoreLogo stores the path to a logo in the database for an organization.
// It returns any error encountered.
func (o Organization) StoreLogo(orgID int, path string) error {
	db := NewDB()

	sqlStatement := fmt.Sprintf(`
		UPDATE 
			organization 
		SET 
			logo_path = '%s' 
		WHERE 
			id = %d`, path, orgID)

	_, err := db.Exec(sqlStatement)
	if err != nil {
		return fmt.Errorf("unable to update: %v", err)
	}

	return nil
}

// ** DELETE **
// Delete deletes an organization from the database.
// It returns any error encountered.
func (o Organization) Delete(orgID int) error {

	db := NewDB()

	sqlStatement := `
		DELETE FROM
			organization
		WHERE
			id = $1`
	_, err := db.Exec(sqlStatement, orgID)
	if err != nil {
		return fmt.Errorf("unable to delete: %v", err)
	}
	return nil
}

// DeleteAdmin deletes an admin from the database.
// It returns any error encountered.
func (o Organization) DeleteAdmin(userID int) error {

	db := NewDB()

	sqlStatement := `
		DELETE FROM
			is_admin
		WHERE
			user_id = $1`
	_, err := db.Exec(sqlStatement, userID)
	if err != nil {
		return fmt.Errorf("unable to delete: %v", err)
	}
	return nil
}
