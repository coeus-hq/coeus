package models

import (
	"database/sql"
	"errors"
)

type Moderator struct {
	ID        int
	UserID    int
	SectionID sql.NullInt64
	Type      string
}

type ModeratorInfo struct {
	ID        int
	Email     string
	LastName  string
	FirstName string
	UserType  string
	SectionID int
}

// ** CREATE **
// Add a new moderator to the database for a given section.
// It returns the ID of the new moderator on success and any encountered error.
func (m *Moderator) Add(userID int, SectionID int, Type string) (int64, error) {
	db := NewDB()
	var result sql.Result
	var ID int64

	// Check if user already exists as a moderator for this section
	var count int
	sqlStatement := `
		SELECT 
			COUNT(*) 
		FROM 
			moderator 
		WHERE 
			section_id = $1
			AND user_id = $2
		`
	err := db.QueryRow(sqlStatement, SectionID, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("User already exists as a moderator for this section")
	}

	// Add user to moderator table
	sqlStatement = `
		INSERT INTO 
			moderator 
		VALUES 
			(NULL,
			$1,
			$2, 
			$3)`
	result, err = db.Exec(sqlStatement, userID, SectionID, Type)
	if err != nil {
		return 0, err
	}
	ID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return ID, err
}

// Add a new moderator to the database without a given section.
// It returns the ID of the new moderator on success and any encountered error.
func (m *Moderator) AdminAdd(userID int, Type string) (int, error) {
	db := NewDB()
	var result sql.Result
	var ID int64

	// Check if user already exists as a moderator for this section
	var count int
	sqlStatement := `
		SELECT 
			COUNT(*) 
		FROM 
			moderator 
		WHERE 
			section_id IS NULL
			AND user_id = $1
		`
	err := db.QueryRow(sqlStatement, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, errors.New("User already exists as a moderator.")
	}

	// Add user to moderator table
	sqlStatement = `
		INSERT INTO 
			moderator 
		VALUES 
			(NULL,
			$1,
			NULL, 
			$2)`
	result, err = db.Exec(sqlStatement, userID, Type)
	if err != nil {
		return 0, err
	}
	ID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Parse the ID into an int
	return int(ID), err
}

// ** READ **
// Get a moderator from the database by ID.
// It returns the moderator struct on success and any encountered error.
func (m Moderator) Get(userID int) (Moderator, error) {
	var moderator Moderator
	db := NewDB()
	sqlStatement := `
		SELECT 
			*
		FROM 
			moderator 
		WHERE 
			user_id = $1;`
	row := db.QueryRow(sqlStatement, userID)
	switch err := row.Scan(&moderator.ID, &moderator.UserID, &moderator.SectionID, &moderator.Type); err {
	case sql.ErrNoRows:
		return Moderator{}, err
	case nil:
		return moderator, err
	default:
		return Moderator{}, err
	}
}

// GetAllModerators returns a list of all moderators for a given section.
func (m Moderator) GetAllModerators(sectionID int) ([]ModeratorInfo, error) {

	var moderatorInfos []ModeratorInfo
	db := NewDB()
	sqlStatement := `
        SELECT
            user.id,
            user.email,
            user.last_name,
            user.first_name,
            moderator.type,
			moderator.section_id
        FROM
            moderator
        JOIN
            user
        ON
            moderator.user_id = user.id
        WHERE
            moderator.section_id = ?
        AND
            (moderator.type = 'moderator'
            OR moderator.type = 'teacher assistant')
        `
	rows, err := db.Query(sqlStatement, sectionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var moderatorInfo ModeratorInfo
		err = rows.Scan(&moderatorInfo.ID, &moderatorInfo.Email, &moderatorInfo.LastName, &moderatorInfo.FirstName, &moderatorInfo.UserType, &moderatorInfo.SectionID)
		if err != nil {
			return nil, err
		}
		moderatorInfos = append(moderatorInfos, moderatorInfo)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return moderatorInfos, err
}

// GetStatus returns the status of a moderator given a section ID and user ID.
// It returns the moderator struct.
func (m Moderator) GetStatus(userID int, sectionID int) (Moderator, error) {

	var moderator Moderator
	db := NewDB()
	sqlStatement := `
		SELECT
			*
		FROM
			moderator
		WHERE
			section_id = $1
			AND user_id = $2`
	row := db.QueryRow(sqlStatement, sectionID, userID)
	switch err := row.Scan(&moderator.ID, &moderator.UserID, &moderator.SectionID, &moderator.Type); err {
	case sql.ErrNoRows:
		return Moderator{}, err
	case nil:
		return moderator, err
	default:
		return Moderator{}, err
	}
}

// IsInstructor returns true if the user is an instructor for any given section.
// It returns any encountered error.
func (m *Moderator) IsInstructor(userID int) (bool, error) {
	db := NewDB()
	var count int
	sqlStatement := `
		SELECT 
			COUNT(*) 
		FROM 
			moderator 
		WHERE 
			user_id = $1
			AND type = 'instructor'`
	err := db.QueryRow(sqlStatement, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// ** UPDATE **
// Update a moderator's type in the database for a given section or insert a new row if no existing moderator.
// It returns the ID of the updated/inserted moderator on success and any encountered error.
func (m *Moderator) Update(userID int, SectionID int, Type string) (int64, error) {
	db := NewDB()
	var result sql.Result
	var ID int64

	// Check if the user is already a moderator for the given section
	sqlCheck := `
		SELECT
			id
		FROM
			moderator
		WHERE
			section_id = $1
			AND user_id = $2
	`
	row := db.QueryRow(sqlCheck, SectionID, userID)
	err := row.Scan(&ID)

	if err == sql.ErrNoRows {
		// If the user is not a moderator, insert a new row
		sqlInsert := `
			INSERT INTO
				moderator (user_id, section_id, type)
			VALUES
				($1, $2, $3)
		`
		result, err = db.Exec(sqlInsert, userID, SectionID, Type)
		if err != nil {
			return 0, err
		}
		ID, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}
	} else if err == nil {
		// If the user is already a moderator, update the moderator type
		sqlUpdate := `
			UPDATE
				moderator
			SET
				type = $1
			WHERE
				section_id = $2
				AND user_id = $3
		`
		result, err = db.Exec(sqlUpdate, Type, SectionID, userID)
		if err != nil {
			return 0, err
		}
		// You can use ID from the SELECT query
	} else {
		// Handle any other error from the SELECT query
		return 0, err
	}

	return ID, err
}

// Update a user's moderator type in the database for a given NULL section.
// It returns the ID of the updated moderator on success and any encountered error.
func (m *Moderator) AdminUpdate(userID int, Type string) (int, error) {
	db := NewDB()
	var result sql.Result
	var ID int64

	// Check if user already exists as a moderator for this section
	var count int
	sqlStatement := `
		SELECT
			COUNT(*)
		FROM

			moderator
		WHERE
			section_id IS NULL
			AND user_id = $1
		`
	err := db.QueryRow(sqlStatement, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	if count == 0 {
		// Add user to moderator table
		sqlStatement = `
			INSERT INTO
				moderator
			VALUES
				(NULL,
				$1,
				NULL,
				$2)`
		result, err = db.Exec(sqlStatement, userID, Type)
		if err != nil {
			return 0, err
		}
		ID, err = result.LastInsertId()
		if err != nil {
			return 0, err
		}

		// Parse the ID into an int
		return int(ID), err

	} else {

		// Update the moderator type in the database
		sqlStatement = `
		UPDATE
    		moderator
		SET
   			type = $1
		WHERE
    		section_id IS NULL
    	AND user_id = $3;`
		result, err = db.Exec(sqlStatement, Type, userID)
		if err != nil {
			return 0, err
		}
		ID, err = result.LastInsertId()
		if err != nil {
		}
	}

	return int(ID), err
}

// ** DELETE **
// Delete a moderator from the database by section ID and user ID.
// It returns any encountered error.
func (m *Moderator) Delete(userID int, sectionID int) error {
	db := NewDB()
	sqlStatement := `
		DELETE FROM 
			moderator 
		WHERE 
			section_id = $1
			AND user_id = $2`
	_, err := db.Exec(sqlStatement, sectionID, userID)
	if err != nil {
		return err
	}
	return nil
}

// Delete a moderator from the database by user ID with a NULL section ID.
// It returns any encountered error.
func (m *Moderator) AdminDelete(userID int) error {
	db := NewDB()
	sqlStatement := `
		DELETE FROM 
			moderator 
		WHERE 
			section_id IS NULL
			AND user_id = $2`
	_, err := db.Exec(sqlStatement, userID)
	if err != nil {
		return err
	}
	return nil
}
