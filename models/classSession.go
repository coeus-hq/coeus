package models

import "fmt"

type ClassSession struct {
	ID         int
	SectionID  int
	ScheduleID int
	InProgress bool
	Data       string
	CreatedAt  string
	UpdatedAt  string
}

type Participant struct {
	ID        int
	SessionID int
	UserID    int
	JoinedAt  string
}

// ** CREATE **
// Join adds a participant to a class session.
// It returns any error encountered and the participant table id.
func (s *ClassSession) Join(sessionID int, studentID int) (int, error) {
	db := NewDB()

	// Check if the student is already in the class session
	var count int
	err := db.QueryRow(`
	SELECT
		COUNT(*)
	FROM
		participants
	WHERE
		session_id = $1
	AND
		user_id = $2`,
		sessionID, studentID).Scan(&count)
	if err != nil {
		return 0, err
	}

	// If the student is already in the class session, return
	if count > 0 {
		return 0, nil
	}

	// Insert participant into the database and return the participant ID
	sqlStatement := `
        INSERT INTO
            participants
        VALUES 
            (NULL,
            $1,
            $2,
            datetime('now')
            )
        RETURNING id`

	var participantID int
	err = db.QueryRow(sqlStatement, sessionID, studentID).Scan(&participantID)
	if err != nil {
		return 0, err
	}

	return participantID, nil
}

// Start sets the in_progress value for a given class session id to true.
// It returns the class session id and any error encountered.
func (s *ClassSession) Start(sectionID int) (int, error) {
	db := NewDB()

	// Update the in_progress field
	res, err := db.Exec(`
	UPDATE
		class_session 
	SET
		in_progress = 1
	WHERE
		section_id = $1`,
		sectionID)
	if err != nil {
		return 0, err
	}

	// Check if a row was affected
	rowsAffected, err := res.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return 0, fmt.Errorf("No rows updated")
	}

	// Select the class session id for the updated row
	var classSessionID int
	err = db.QueryRow("SELECT id FROM class_session WHERE section_id = $1", sectionID).Scan(&classSessionID)
	if err != nil {
		return 0, err
	}

	return classSessionID, nil
}

// ** READ **
// GetParticipantCount returns a list of participants for a given class session id.
// It returns the count of participants and any error encountered.
func (s ClassSession) GetParticipantCount(classSessionID int) (int, error) {
	db := NewDB()

	var count int
	err := db.QueryRow(`
	SELECT
		COUNT(*)
	FROM 
		participants 
	WHERE
		session_id = $1`,
		classSessionID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetParticipants given class session id.
// It returns a slice of Participants structs and any error encountered.
func (s Participant) GetParticipants(classSessionID int) ([]Participant, error) {
	db := NewDB()

	var participants []Participant
	rows, err := db.Query(`
	SELECT
		id,
		session_id,
		user_id,
		joined_at
	FROM
		participants
	WHERE

		session_id = $1`,
		classSessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Participant{}
		err := rows.Scan(&p.ID, &p.SessionID, &p.UserID, &p.JoinedAt)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	return participants, nil
}

// GetID returns the class session id for a given section id.
// It returns the class session id and any error encountered.
func (s ClassSession) GetID(sectionID int) (int, error) {
	db := NewDB()

	var classSessionID int
	err := db.QueryRow(`
	SELECT
		id
	FROM 
		class_session 
	WHERE
		section_id = $1`,
		sectionID).Scan(&classSessionID)
	if err != nil {
		return 0, err
	}

	return classSessionID, nil
}

// GetInProgress returns the in_progress value for a given class session id.
// It returns the in_progress value and any error encountered.
func (s ClassSession) GetInProgress(classSessionID int) (bool, error) {
	db := NewDB()

	var inProgress bool
	err := db.QueryRow(`
	SELECT
		in_progress
	FROM 
		class_session 
	WHERE
		id = $1`,
		classSessionID).Scan(&inProgress)
	if err != nil {
		return false, err
	}

	return inProgress, nil
}

// GetSectionID returns the section id for a given class session id.
// It returns the section id and any error encountered.
func (s ClassSession) GetSectionID(classSessionID int) (int, error) {
	db := NewDB()

	var sectionID int
	err := db.QueryRow(`
	SELECT
		section_id
	FROM 
		class_session 
	WHERE
		id = $1`,
		classSessionID).Scan(&sectionID)
	if err != nil {
		return 0, err
	}

	return sectionID, nil
}

// ** UPDATE **
// End sets the in_progress value for a given class session id to false, removes all participants, and removes all questions.
// It returns any error encountered.
func (s *ClassSession) End(classSessionID int) error {
	db := NewDB()

	// Set the class session to not in progress
	_, err := db.Exec(`
	UPDATE
		class_session 
	SET
		in_progress = 0
	WHERE
		id = $1`,
		classSessionID)
	if err != nil {
		return err
	}

	// Remove all participants from the class session
	_, err = db.Exec(`
	DELETE FROM
		participants
	WHERE
		session_id = $1`,
		classSessionID)
	if err != nil {
		return err
	}

	// Remove all questions from the class session
	_, err = db.Exec(`
	DELETE FROM
		question
	WHERE
		session_id = $1`,
		classSessionID)
	if err != nil {
		return err
	}

	return nil
}
