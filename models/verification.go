package models

import (
	"coeus/email"
)

type VerifyUser struct {
	ID       int64
	Email    string
	Password string
	Token    string
	Status   string
}

// CreateToken resets the password for a user by finding the user by email and creating a VerifyUser record than can be used to reset the password
// Returns an error if the user is not found and if the user is found the VerifyUser ID is returned
func (v VerifyUser) CreateToken(userEmail string) (int, error) {

	// Open a new database connection
	db := NewDB()

	userID, err := new(User).GetUserId(userEmail)
	if err != nil {
		return 0, err
	}

	// Convert userID to int64
	userIDInt64 := int64(userID)

	user, err := new(User).Get(userIDInt64)

	token := email.SendForgotPasswordEmail(userEmail, user.FirstName)
	_, err = db.Exec("INSERT INTO verify_user (user_id, email, token, expiration, created_at, status) VALUES (?, ?, ?, datetime('now'), datetime('now'), ?)", userID, userEmail, token, "pending")
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// MatchToken verifies the token for a user and returns the user ID if the token is valid
// Returns an error if the token is not valid
func (v VerifyUser) MatchToken(token string) (bool, int, error) {
	// Open a new database connection
	db := NewDB()
	var userID int
	err := db.QueryRow("SELECT user_id FROM verify_user WHERE token = ?", token).Scan(&userID)
	if err != nil {
		return false, 0, err
	}
	return true, userID, nil
}

// GetToken returns the token for a user by ID
// Returns an error if the token is not found
func (v VerifyUser) GetToken(userID int) (string, error) {
	// Open a new database connection
	db := NewDB()
	var token string
	err := db.QueryRow("SELECT token FROM verify_user WHERE user_id = ?", userID).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}

// RemoveToken removes the all tokens for a user by ID if the token(s) are found
// Return no errors
func (v VerifyUser) DeleteToken(userID int) error {
	// Open a new database connection
	db := NewDB()
	_, err := db.Exec("DELETE FROM verify_user WHERE user_id = ?", userID)
	if err != nil {
		return nil
	}
	return nil
}
