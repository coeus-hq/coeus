package models

import (
	"database/sql"
	"fmt"
)

type Setting struct {
	id             int
	UserId         int
	DarkTheme      bool
	TimezoneOffset int
}

// ** CREATE **
// Add adds a new setting to the database.
// It returns any error encountered.
func (s Setting) Add(UserId int) (int, error) {
	db := NewDB()

	// Get the Organization ID
	organizationID, err := new(Organization).GetOrganizationID()
	if err != nil {
		return 0, err
	}

	// Get the Organization timezone
	organization, err := new(Organization).Get(organizationID)
	if err != nil {
		return 0, err
	}

	sqlStatement := `
		INSERT INTO
			setting 
		VALUES
			(
			NULL,
			$1,
			0,
			$2
			)
		RETURNING id`
	var ID int
	err = db.QueryRow(sqlStatement, UserId, organization.OrganizationTimezone).Scan(&ID)
	if err != nil {
		return 0, fmt.Errorf("unable to insert: %v", err)
	}

	return ID, nil
}

// ** READ **
// Get retrieves a user's settings from the database.
// It returns a Setting object and any error encountered.
func (s Setting) Get(UserId int) (Setting, error) {
	var id int
	var DarkTheme bool
	var TimezoneOffset int

	db := NewDB()
	sqlStatement := `
		SELECT 
			* 
		FROM 
			setting 
		WHERE 
			user_id = $1;`

	row := db.QueryRow(sqlStatement, UserId)
	switch err := row.Scan(&id, &UserId, &DarkTheme, &TimezoneOffset); err {
	case sql.ErrNoRows:
		return Setting{}, err
	case nil:
		setting := Setting{id, UserId, DarkTheme, TimezoneOffset}
		return setting, err
	default:
		return Setting{}, err
	}
}

// ** UPDATE **
// Update setting for a user in the database.
// It returns any error encountered.
func (s Setting) Update(userID int, name string, value string, timezoneOffset int) error {
	db := NewDB()

	sqlStatement := fmt.Sprintf(`
        UPDATE 
            setting 
        SET 
            %s = %s,
            timezone_offset = %d
        WHERE 
            user_id = %d`, name, value, timezoneOffset, userID)

	_, err := db.Exec(sqlStatement)
	if err != nil {
		return fmt.Errorf("unable to update: %v", err)
	}

	return nil
}

// UpdateTimezone updates the timezone offset for a user in the database.
// It returns any error encountered.
func (s Setting) UpdateTimezone(userID int, timezoneOffset int) error {
	db := NewDB()

	sqlStatement := `
		UPDATE
			setting
		SET
			timezone_offset = $1
		WHERE
			user_id = $2`

	_, err := db.Exec(sqlStatement, timezoneOffset, userID)
	if err != nil {
		return fmt.Errorf("unable to update: %v", err)
	}

	return nil
}

// ToggleDarkTheme toggles the dark theme setting for a user in the database.
func (s Setting) ToggleDarkTheme(userID int) (bool, error) {
	db := NewDB()

	updateStatement := `
		UPDATE 
			setting 
		SET 
			theme = NOT theme
		WHERE 
			user_id = $1`

	_, err := db.Exec(updateStatement, userID)
	if err != nil {
		return false, fmt.Errorf("unable to update: %v", err)
	}

	// Retrieve the current value of theme for the user
	var newTheme bool
	selectStatement := `
		SELECT 
			theme 
		FROM 
			setting 
		WHERE 
			user_id = $1`
	err = db.QueryRow(selectStatement, userID).Scan(&newTheme)
	if err != nil {
		return false, fmt.Errorf("unable to retrieve new theme value: %v", err)
	}

	return newTheme, nil
}

// ** DELETE **
// Delete deletes a setting from the database.
// It returns any error encountered.
func (s Setting) Delete(UserId int) error {
	db := NewDB()

	sqlStatement := `
		DELETE FROM
			setting
		WHERE
			user_id = $1`
	_, err := db.Exec(sqlStatement, UserId)
	if err != nil {
		return fmt.Errorf("unable to delete: %v", err)
	}
	return nil
}
