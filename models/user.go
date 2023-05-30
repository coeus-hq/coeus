package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Email     string
	LastName  string
	FirstName string
	CreatedAt string
	UpdatedAt string
}

// ** CREATE **
// Add adds a new user to the database.
// It returns the ID of the newly inserted user and any error encountered.
func (u *User) Add(Email string, password string, LastName string, FirstName string) (int64, error) {
	var result sql.Result
	var id int64

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// Open a new database connection
	db := NewDB()

	// Check if user already exists
	var count int
	sqlStatement := `
		SELECT 
			COUNT(*) 
		FROM 
			user 
		WHERE 
			email = $1`

	err = db.QueryRow(sqlStatement, Email).Scan(&count)
	if err != nil {
		// Handle the error
		return 0, err
	}
	if count > 0 {
		// User already exists, return err
		return 0, errors.New("User already exists")
	}

	// Add user to database
	sqlStatement = `
		INSERT INTO 
			user 
		VALUES 
			(NULL,
			$1,
			$2, 
			$3, 
			$4, 
			datetime('now'), 
			datetime('now'))`

	result, err = db.Exec(sqlStatement, Email, hash, LastName, FirstName)
	if err != nil {
		return 0, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// Return true on success
	return id, err
}

// AddUserToOrganization adds a user to an organization.
// It returns an int and any error encountered.
func (u User) AddUserToOrganization(userID int, organizationID int) error {

	db := NewDB()
	sqlStatement := `
	INSERT INTO 
		user_organization 
		(user_id, organization_id) 
	VALUES 
		($1, $2);`

	_, err := db.Exec(sqlStatement, userID, organizationID)
	if err != nil {
		return err
	}

	return nil
}

// ** READ **
// Get retrieves a user from the database.
// It returns a User object and any error encountered.
func (u User) Get(id int64) (User, error) {
	var Email string
	var LastName string
	var FirstName string
	var CreatedAt string
	var UpdatedAt string

	db := NewDB()
	sqlStatement := `
	SELECT 
		id, 
		email, 
		last_name, 
		first_name, 
		created_at, 
		updated_at 
	FROM 
		user 
	WHERE 
		id=$1;`

	row := db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&id, &Email, &LastName, &FirstName, &CreatedAt, &UpdatedAt); err {
	case sql.ErrNoRows:
		return User{}, err
	case nil:
		user := User{id, Email, LastName, FirstName, CreatedAt, UpdatedAt}
		return user, err
	default:
		return User{}, err
	}
}

// GetAll retrieves all users from the database, including their highest moderator type.
// It returns a slice of User objects and any error encountered.
func (u User) GetAll() ([]struct {
	ID                   int64
	Email                string
	LastName             string
	FirstName            string
	CreatedAt            string
	UpdatedAt            string
	HighestModeratorType string
}, error) {
	db := NewDB()
	sqlStatement := `
	SELECT
	    user.id,
	    user.email,
	    user.last_name,
	    user.first_name,
	    user.created_at,
	    user.updated_at,
	    CASE MAX(CASE moderator.type
	        WHEN 'instructor' THEN 4
	        WHEN 'teacher assistant' THEN 3
	        WHEN 'moderator' THEN 2
	        WHEN 'student' THEN 1
	        ELSE 0
	    END)
	    WHEN 4 THEN 'instructor'
	    WHEN 3 THEN 'teacher assistant'
	    WHEN 2 THEN 'moderator'
	    WHEN 1 THEN 'student'
	    ELSE 'none'
	    END AS highest_moderator_type
	FROM
	    user
	    LEFT JOIN moderator ON user.id = moderator.user_id
	GROUP BY
	    user.id,
	    user.email,
	    user.last_name,
	    user.first_name,
	    user.created_at,
	    user.updated_at`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []struct {
		ID                   int64
		Email                string
		LastName             string
		FirstName            string
		CreatedAt            string
		UpdatedAt            string
		HighestModeratorType string
	}{}
	for rows.Next() {
		var id int64
		var Email string
		var LastName string
		var FirstName string
		var CreatedAt string
		var UpdatedAt string
		var HighestModeratorType string

		err = rows.Scan(&id, &Email, &LastName, &FirstName, &CreatedAt, &UpdatedAt, &HighestModeratorType)
		if err != nil {
			return nil, err
		}

		users = append(users, struct {
			ID                   int64
			Email                string
			LastName             string
			FirstName            string
			CreatedAt            string
			UpdatedAt            string
			HighestModeratorType string
		}{
			ID:                   id,
			Email:                Email,
			LastName:             LastName,
			FirstName:            FirstName,
			CreatedAt:            CreatedAt,
			UpdatedAt:            UpdatedAt,
			HighestModeratorType: HighestModeratorType,
		})
	}

	return users, nil
}

// Authenticate returns user_id from the user table on success or 0 on failure
func (u User) Authenticate(Email string, password string) int {
	var id int
	var hash string

	db := NewDB()
	sqlStatement := `
		SELECT 
			id, 
			Email, 
			hash 
		FROM 
			user 
		WHERE 
			Email=$1;`

	row := db.QueryRow(sqlStatement, Email)
	switch err := row.Scan(&id, &Email, &hash); err {
	case sql.ErrNoRows:
		return 0
	case nil:
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err == nil {
			return id
		}
	default:
		return 0
	}
	return 0
}

// GetUserInitials returns the user's initials
func (u User) GetUserInitials(id int) (string, error) {
	var FirstName string
	var LastName string

	db := NewDB()
	sqlStatement := `
	SELECT 
		first_name, 
		last_name 
	FROM 
		user 
	WHERE 
		id=$1;`

	row := db.QueryRow(sqlStatement, id)
	switch err := row.Scan(&FirstName, &LastName); err {
	case sql.ErrNoRows:
		return "", err
	case nil:
		return FirstName[0:1] + LastName[0:1], err
	default:
		return "", err
	}
}

// GetUserId by email.
// Returns the user's id.
func (u User) GetUserId(Email string) (int, error) {
	var id int

	db := NewDB()
	sqlStatement := `
	SELECT 
		id 
	FROM 
		user 
	WHERE 
		email=$1;`

	row := db.QueryRow(sqlStatement, Email)
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		return 0, err
	case nil:
		return id, err
	default:
		return 0, err
	}
}

// Count returns the number of users in the database.
// It returns an int and any error encountered.
func (u User) Count() (int, error) {
	var count int

	db := NewDB()
	sqlStatement := `
	SELECT 
		COUNT(*) 
	FROM 
		user;`

	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&count); err {
	case sql.ErrNoRows:
		return 0, err
	case nil:
		return count, err
	default:
		return 0, err
	}
}

// OrganizationID returns the organization id for a user.
// It returns an int and any error encountered.
func (u User) OrganizationID(userID int) (int, error) {
	var organization_id int

	db := NewDB()
	sqlStatement := `
	SELECT 
		organization_id 
	FROM 
	user_organization 
	WHERE 
		user_id = $1;`

	row := db.QueryRow(sqlStatement, userID)
	switch err := row.Scan(&organization_id); err {
	case sql.ErrNoRows:
		return 0, err
	case nil:
		return organization_id, err
	default:
		return 0, err
	}
}

// GetOrganizationID returns the organization id for the first organization in the database.
// It returns an int and any error encountered.
func (u User) GetOrganizationID() (int, error) {
	var organization_id int

	db := NewDB()
	sqlStatement := `
	SELECT 
		id 
	FROM 
		organization;`

	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&organization_id); err {
	case sql.ErrNoRows:
		return 0, err
	case nil:
		return organization_id, err
	default:
		return 0, err
	}
}

// ** UPDATE **
// Update updates a user in the database.
// It returns a User object and any error encountered.
func (u User) Update(id int64, Email string, LastName string, FirstName string, password ...string) error {
	db := NewDB()

	// Check if user exists
	var count int
	sqlStatement := `
		SELECT 
			COUNT(*) 
		FROM 
			user 
		WHERE 
			id = $1`

	err := db.QueryRow(sqlStatement, id).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		// User doesn't exist
		return errors.New("User doesn't exist")
	}

	// Update password if provided
	if len(password) > 0 && password[0] != "" {
		// Hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(password[0]), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Update password in the database
		sqlStatement = `
			UPDATE 
				user 
			SET 
				hash = $1, 
				updated_at = datetime('now') 
			WHERE 
				id = $2`

		_, err = db.Exec(sqlStatement, hash, id)
		if err != nil {
			return err
		}
	}

	// Update user in database
	sqlStatement = `
		UPDATE 
			user 
		SET 
			email = $1, 
			last_name = $2, 
			first_name = $3, 
			updated_at = datetime('now') 
		WHERE 
			id = $4`

	_, err = db.Exec(sqlStatement, Email, LastName, FirstName, id)
	if err != nil {
		return err
	}

	return nil
}

// UpdatePassword updates a user's password in the database.
// It returns a User object and any error encountered.
func (u User) UpdatePassword(id int, password string) error {
	db := NewDB()

	// Check if user exists
	var count int
	sqlStatement := `

		SELECT
			COUNT(*)
		FROM
			user
		WHERE
			id = $1`

	err := db.QueryRow(sqlStatement, id).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		// User doesn't exist
		return errors.New("User doesn't exist")
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password in the database
	sqlStatement = `
		UPDATE
			user
		SET
			hash = $1,
			updated_at = datetime('now')
		WHERE
			id = $2`

	_, err = db.Exec(sqlStatement, hash, id)
	if err != nil {
		return err
	}

	return nil
}

// ** DELETE **
// Delete deletes a user from the database.
// It returns a User object and any error encountered.
func (u User) Delete(id int64) error {
	db := NewDB()

	// Check if user exists
	var count int
	sqlStatement := `
		SELECT 
			COUNT(*) 
		FROM 
			user 
		WHERE 
			id = $1`

	err := db.QueryRow(sqlStatement, id).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		// User doesn't exist
		return errors.New("User doesn't exist")
	}

	// Delete user from database
	sqlStatement = `
		DELETE FROM 
			user 
		WHERE 
			id = $1`

	_, err = db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}

	return nil
}
