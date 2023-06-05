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
	APIKey               string
	Email                string
	Onboarding           string
	CreatedAt            string
	UpdatedAt            string
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

	// Insert new organization
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
			datetime('now')
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
	var APIKey string
	var Email string
	var Onboarding string

	db := NewDB()
	sqlStatement := `
		SELECT 
			* 
		FROM 
			organization 
		WHERE 
			id = $1;`

	row := db.QueryRow(sqlStatement, orgID)
	switch err := row.Scan(&id, &Name, &OrganizationTimezone, &LogoPath, &APIKey, &Email, &Onboarding, &CreatedAt, &UpdatedAt); err {
	case sql.ErrNoRows:
		return Organization{}, err
	case nil:
		organization := Organization{id, Name, OrganizationTimezone, LogoPath.String, APIKey, Email, Onboarding, CreatedAt, UpdatedAt}
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
	var apiKey string

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
	if apiKey != "" {
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

// DatabaseReset resets the database.
// It returns any error encountered.
func (o Organization) DatabaseReset() error {
	db := NewDB()

	sqlStatement := `
	
		DROP TABLE IF EXISTS user;
		DROP TABLE IF EXISTS site;
		DROP TABLE IF EXISTS setting;
		DROP TABLE IF EXISTS course;
		DROP TABLE IF EXISTS section;
		DROP TABLE IF EXISTS enrollment;
		DROP TABLE IF EXISTS class_session;
		DROP TABLE IF EXISTS question;
		DROP TABLE IF EXISTS participants;
		DROP TABLE IF EXISTS schedule;
		DROP TABLE IF EXISTS vote;
		DROP TABLE IF EXISTS moderator;
		DROP TABLE IF EXISTS user_organization;
		DROP TABLE IF EXISTS organization;
		DROP TABLE IF EXISTS verify_user;
		DROP TABLE IF EXISTS attendance;
		DROP TABLE IF EXISTS user_attendance;
		DROP TABLE IF EXISTS is_admin;
		DROP TABLE IF EXISTS schema_version;
	
		CREATE TABLE schema_version (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	
		CREATE TABLE is_admin (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL
		);
	
		CREATE TABLE user(
		   id INTEGER PRIMARY KEY AUTOINCREMENT,
		   email VARCHAR(128) NOT NULL,
		   hash TEXT NOT NULL,
		   last_name VARCHAR(40) NOT NULL,
		   first_name VARCHAR(40) NOT NULL,
		   created_at TEXT NOT NULL,
		   updated_at TEXT NOT NULL
		);
	
		CREATE TABLE user_organization (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			organization_id INTEGER NOT NULL
		);
	
		CREATE TABLE setting(
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  user_id INTEGER NOT NULL,
		  theme BOOLEAN NOT NULL,
		  timezone_offset INTEGER NOT NULL
		);
	
		CREATE TABLE organization(
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  name VARCHAR(256) NOT NULL,
		  organization_timezone VARCHAR(5),
		  logo_path TEXT,
		  api_key TEXT,
		  email TEXT,
		  onboarding_complete BOOLEAN NOT NULL,
		  created_at TEXT NOT NULL,
		  updated_at TEXT NOT NULL
		);
	
		CREATE TABLE course(
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  number TEXT NOT NULL,
		  title TEXT NOT NULL,
		  start_date TEXT NOT NULL,
		  end_date TEXT NOT NULL,
		  semester TEXT NOT NULL,
		  year INTEGER NOT NULL,
		  created_at TEXT NOT NULL,
		  updated_at TEXT NOT NULL
		);
	
		CREATE TABLE section(
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  course_id INTEGER NOT NULL REFERENCES course(id),
		  name VARCHAR(256) NOT NULL, 
		  created_at TEXT NOT NULL,
		  updated_at TEXT NOT NULL
		);
	
		CREATE TABLE enrollment(
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  section_id INTEGER NOT NULL REFERENCES section(id),
		  user_id INTEGER NOT NULL REFERENCES user(id)
		);
	
		CREATE TABLE schedule(
		  id INTEGER PRIMARY KEY AUTOINCREMENT,
		  section_id INTEGER NOT NULL REFERENCES section(id),
		  day STRING NOT NULL,
		  timeslot STRING NOT NULL
		);
	
		CREATE TABLE class_session(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			section_id INTEGER NOT NULL,
			schedule_id INTEGER NOT NULL,
			in_progress BOOLEAN NOT NULL,
			data TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
	
		CREATE TABLE attendance(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			section_id INTEGER NOT NULL,
			instructor_id INTEGER NOT NULL,
			class_session_id INTEGER NOT NULL,
			attended_at TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
	
		CREATE TABLE user_attendance(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			attendance_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			status TEXT CHECK(status IN ('present', 'absent', 'excused', 'late'))
		);
	
		CREATE TABLE question(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			text VARCHAR(140) NOT NULL,
			votes INTEGER NOT NULL,
			answered BOOLEAN NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		);
	
		CREATE TABLE participants(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			joined_at TEXT NOT NULL
		);
	
		CREATE TABLE vote(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			question_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL
		);
	
		CREATE TABLE moderator(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			section_id INTEGER,
			type TEXT CHECK(type IN ('student', 'moderator', 'teacher assistant', 'instructor'))
		);
	
		CREATE TABLE verify_user(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			email TEXT NOT NULL,
			token TEXT NOT NULL,
			expiration TEXT NOT NULL,
			created_at TEXT NOT NULL,
			status TEXT CHECK(status IN ('pending', 'verified', 'expired'))
		);

`

	_, err := db.Exec(sqlStatement)
	if err != nil {
		return fmt.Errorf("unable to delete: %v", err)
	}
	return nil
}

// DatabaseSeed seeds the database with an example data set.
// It returns any error encountered.
func (o Organization) DatabaseSeed() error {
	db := NewDB()

	sqlStatement := `
	DELETE FROM user WHERE id NOT IN (SELECT user_id FROM is_admin);
	UPDATE user SET id = 2000 WHERE id = 1;
	UPDATE is_admin SET user_id = 2000 WHERE user_id = 1;
	DELETE FROM participants;
	DELETE FROM attendance;
	DELETE FROM class_session;
	DELETE FROM course;
	DELETE FROM enrollment;
	DELETE FROM moderator;
	DELETE FROM question;
	DELETE FROM schedule;
	DELETE FROM schema_version;
	DELETE FROM section;
	DELETE FROM setting;
	DELETE FROM user_attendance;
	DELETE FROM user_organization;
	DELETE FROM verify_user;
	DELETE FROM vote;
	INSERT INTO user VALUES (1, 'student@coeus.education', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', '1', 'Student', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (NULL, 'ta@coeus.education', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'A', 'T', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (3, 'instructor@coeus.education', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'I', 'I', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (NULL, 'whalencollin@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Whalen', 'Collin', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (999, 'testuser@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'U', 'T', datetime('now'), datetime('now'));
	INSERT INTO setting VALUES (1, 1, false, 0);
	INSERT INTO setting VALUES (2, 2, true, -360);
	INSERT INTO setting VALUES (3, 3, false, 0);
	INSERT INTO setting VALUES (4, 4, false, 0);
	INSERT INTO setting VALUES (5, 999, false, 0);
	INSERT INTO course VALUES (1, 'CSI 1113', 'Intro to C/C++', '2023-05-14', '2023-01-14', 'Spring', 2023, datetime('now'), datetime('now'));
	INSERT INTO course VALUES (2, 'CSCI 1115','Exploring CSCI: C++', '2023-05-14', '2023-01-14', 'Spring', 2023, datetime('now'), datetime('now'));
	INSERT INTO course VALUES (3, 'CSCI 1133','Intro to Programming Concepts', '2023-05-14', '2023-01-14', 'Spring', 2023, datetime('now'), datetime('now'));
	INSERT INTO course VALUES (4, 'CSCI 1135', 'Exploring CSCI: Python', '2023-05-14', '2023-01-14', 'Spring', 2023, datetime('now'), datetime('now'));
	INSERT INTO course VALUES (5, 'CSCI 1913','Intro to Algs. & Program Dev.', '2023-05-14', '2023-01-14', 'Spring', 2023, datetime('now'), datetime('now'));
	INSERT INTO course VALUES (6, 'CSCI 1933','Intro Algs & Data Str.', '2023-05-14', '2023-01-14', 'Spring', 2023, datetime('now'), datetime('now'));
	INSERT INTO section VALUES (1, 1, '1', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (1, 1, 'M W F', '11:15 AM-12:05 PM');
	INSERT INTO section VALUES (2, 1, '10', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (2, 2, 'M W F', '03:35 PM-04:25 PM');
	INSERT INTO section VALUES (3, 2, '1 (model tests ID)', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (3, 3, 'M W F', '10:10 AM-11:00 AM ');
	INSERT INTO section VALUES (4, 3, '1', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (4, 4, 'M W F', '10:10 AM-11:00 AM ');
	INSERT INTO section VALUES (5, 3, '10', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (5, 5, 'W', '06:30 PM-09:00 PM');
	INSERT INTO section VALUES (6, 3, '20', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (6, 6, 'M W F', '10:10 AM-11:00 AM');
	INSERT INTO section VALUES (7, 3, '30', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (7, 7, 'M W F', '03:35 PM-04:25 PM');
	INSERT INTO section VALUES (8, 3, '40', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (8, 8, 'M W F', '12:20 PM-01:10 PM');
	INSERT INTO section VALUES (9, 3, '50', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (9, 9, 'M W F', '12:20 PM-01:10 PM');
	INSERT INTO section VALUES (10, 4, '1', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (10, 10, 'Tu', '06:00 PM-08:00 PM');
	INSERT INTO section VALUES (11, 5, '1', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (11, 11, 'M W F', '02:30 PM-03:20 PM');
	INSERT INTO section VALUES (12, 5, '10', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (12, 12, 'W', '12:20 PM-02:15 PM');
	INSERT INTO section VALUES (13, 6, '10', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (13, 13, 'Th', '06:30 PM-09:00 PM');
	INSERT INTO section VALUES (14, 6, '20', datetime('now'), datetime('now'));
	INSERT INTO schedule VALUES (14, 14, 'M W F', '11:15 AM-12:05 PM');
	INSERT INTO enrollment VALUES (NULL, 1, 1);
	INSERT INTO enrollment VALUES (NULL, 3, 1);
	INSERT INTO enrollment VALUES (NULL, 3, 2);
	INSERT INTO enrollment VALUES (NULL, 3, 3);
	INSERT INTO enrollment VALUES (NULL, 3, 4);
	INSERT INTO enrollment VALUES (NULL, 10, 1);
	INSERT INTO enrollment VALUES (NULL, 4, 1);
	INSERT INTO enrollment VALUES (NULL, 3, 1);
	INSERT INTO enrollment VALUES (NULL, 1, 2);
	INSERT INTO enrollment VALUES (NULL, 10, 2);
	INSERT INTO enrollment VALUES (NULL, 1, 2);
	INSERT INTO class_session(id, section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (1, 10, 1, true, 'Session data 1 (model tests ID)', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (2, 2, false, 'Session data 2', '2022-01-02 11:00:00', '2022-01-02 11:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (3, 3, true, 'Session data 3', '2022-01-03 12:00:00', '2022-01-03 12:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (1, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (4, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (5, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (6, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (7, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (8, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (9, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (11, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (12, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (13, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO class_session(section_id, schedule_id, in_progress, data, created_at, updated_at) VALUES (14, 1, false, 'Session data 4', '2022-01-01 10:00:00', '2022-01-01 10:00:00');
	INSERT INTO question(id, session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (1, 10, 1, 'Question 1 for session 1 (model tests ID)', 5, true, '2022-01-01 10:05:00', '2022-01-01 11:00:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (10, 1, 'Question 2 for session 1', 2, false, '2022-01-01 10:10:00', '2022-01-01 11:00:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (10, 1, 'Question 1 for session 2', 0, false, '2022-01-02 11:05:00', '2022-01-02 11:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (11, 4, 'Question 1 for session 3', 3, false, '2022-01-03 12:05:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (11, 5, 'Question 2 for session 3', 1, false, '2022-01-03 12:10:00', '2022-01-03 12:30:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (1, 1, '2022-01-01 10:00:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (1, 2, '2022-01-01 10:02:00');
	INSERT INTO moderator(user_id, section_id, type) VALUES (1, 3, 'student');
	INSERT INTO moderator(user_id, section_id, type) VALUES (3, NULL, 'instructor');
	INSERT INTO moderator(user_id, section_id, type) VALUES (999, 3, 'instructor');
	INSERT INTO moderator(user_id, section_id, type) VALUES (3, 3, 'instructor');
	INSERT INTO moderator(user_id, section_id, type) VALUES (1, 4, 'student');
	INSERT INTO moderator(user_id, section_id, type) VALUES (1, 1, 'student');
	INSERT INTO moderator(user_id, section_id, type) VALUES (1, 10, 'student');
	INSERT INTO moderator(user_id, section_id, type) VALUES (2, 10, 'student');
	INSERT INTO moderator(user_id, section_id, type) VALUES (2, 1, 'student');
	INSERT INTO user_organization (user_id, organization_id) VALUES (1, 1);
	INSERT INTO user_organization (user_id, organization_id) VALUES (2, 1);
	INSERT INTO user_organization (user_id, organization_id) VALUES (3, 1);
	INSERT INTO user_organization (user_id, organization_id) VALUES (4, 1);
	INSERT INTO user_organization (user_id, organization_id) VALUES (999, 1);
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'Course 1115 is a test course.', 15, true, '2022-01-02 11:05:00', '2022-01-02 11:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 4, 'The session ID is 3 and section ID is 3.', 3, false, '2022-01-03 12:05:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'Are there specific prerequisites needed to enroll in this course, or can anyone join?', 36, false, '2022-01-03 12:15:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'How are the grades in this class calculated? Are there any specific factors considered?', 29, true, '2022-01-03 12:20:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'Are there any group projects assigned in this course? If so, how will they be structured?', 25, false, '2022-01-03 12:25:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'What type of assignments can students expect throughout the course? Are there examples available?', 38, true, '2022-01-03 12:30:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'Will there be any guest speakers visiting during the semester? If so, who might we expect? Any days set aside for guests and open house?', 17, false, '2022-01-03 12:35:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'How can students best prepare for exams in this course? Are there any recommended strategies? Is there a large amount of reading in 1115?', 40, true, '2022-01-03 13:00:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'What is the most effective approach to succeed in this class? Are there any tips to follow?', 45, false, '2022-01-03 13:05:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'Are there any online resources or platforms that can help with this courses material?', 28, true, '2022-01-03 13:10:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'Will the course teach about data structures and algorithms in-depth? What can we expect? What classes do we need to have already taken?', 35, false, '2022-01-03 13:15:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'Is there a grading curve applied to this course? If so, how does it affect the final grades? How much of our grade are final exams worth?', 32, true, '2022-01-03 13:20:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'What specific software applications or tools will we be using throughout this class? Are there any alternatives?', 27, false, '2022-01-03 13:25:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'Are late assignments accepted in this course? If so, are there penalties or restrictions to consider?', 41, true, '2022-01-03 13:30:00', '2022-01-03 12:30:00');
	INSERT INTO question(session_id, user_id, text, votes, answered, created_at, updated_at) VALUES (3, 1, 'Do we need to have prior knowledge of any specific programming languages for this course? Which ones?', 30, false, '2022-01-03 13:35:00', '2022-01-03 12:30:00');
	INSERT INTO enrollment VALUES (NULL, 3, 1);
	INSERT INTO enrollment VALUES (NULL, 3, 2);
	INSERT INTO enrollment VALUES (NULL, 3, 3);
	INSERT INTO enrollment VALUES (NULL, 3, 4);
	INSERT INTO enrollment VALUES (NULL, 3, 5);
	INSERT INTO enrollment VALUES (NULL, 3, 6);
	INSERT INTO enrollment VALUES (NULL, 3, 7);
	INSERT INTO enrollment VALUES (NULL, 3, 8);
	INSERT INTO enrollment VALUES (NULL, 3, 9);
	INSERT INTO enrollment VALUES (NULL, 3, 10);
	INSERT INTO enrollment VALUES (NULL, 3, 11);
	INSERT INTO enrollment VALUES (NULL, 3, 12);
	INSERT INTO enrollment VALUES (NULL, 3, 13);
	INSERT INTO enrollment VALUES (NULL, 3, 14);
	INSERT INTO enrollment VALUES (NULL, 3, 15);
	INSERT INTO enrollment VALUES (NULL, 3, 16);
	INSERT INTO enrollment VALUES (NULL, 3, 17);
	INSERT INTO enrollment VALUES (NULL, 3, 18);
	INSERT INTO enrollment VALUES (NULL, 3, 19);
	INSERT INTO enrollment VALUES (NULL, 3, 20);
	INSERT INTO enrollment VALUES (NULL, 3, 21);
	INSERT INTO enrollment VALUES (NULL, 3, 22);
	INSERT INTO enrollment VALUES (NULL, 3, 23);
	INSERT INTO enrollment VALUES (NULL, 3, 24);
	INSERT INTO enrollment VALUES (NULL, 3, 25);
	INSERT INTO enrollment VALUES (NULL, 3, 26);
	INSERT INTO enrollment VALUES (NULL, 3, 27);
	INSERT INTO enrollment VALUES (NULL, 3, 28);
	INSERT INTO enrollment VALUES (NULL, 3, 29);
	INSERT INTO enrollment VALUES (NULL, 3, 30);
	INSERT INTO enrollment VALUES (NULL, 3, 31);
	INSERT INTO enrollment VALUES (NULL, 3, 32);
	INSERT INTO enrollment VALUES (NULL, 3, 33);
	INSERT INTO enrollment VALUES (NULL, 3, 34);
	INSERT INTO enrollment VALUES (NULL, 3, 35);
	INSERT INTO enrollment VALUES (NULL, 3, 36);
	INSERT INTO enrollment VALUES (NULL, 3, 37);
	INSERT INTO enrollment VALUES (NULL, 3, 38);
	INSERT INTO enrollment VALUES (NULL, 3, 39);
	INSERT INTO enrollment VALUES (NULL, 3, 40);
	INSERT INTO enrollment VALUES (NULL, 3, 41);
	INSERT INTO enrollment VALUES (NULL, 3, 42);
	INSERT INTO enrollment VALUES (NULL, 3, 43);
	INSERT INTO enrollment VALUES (NULL, 3, 44);
	INSERT INTO enrollment VALUES (NULL, 3, 45);
	INSERT INTO enrollment VALUES (NULL, 3, 46);
	INSERT INTO enrollment VALUES (NULL, 3, 47);
	INSERT INTO enrollment VALUES (NULL, 3, 48);
	INSERT INTO enrollment VALUES (NULL, 3, 49);
	INSERT INTO enrollment VALUES (NULL, 3, 50);
	INSERT INTO enrollment VALUES (NULL, 3, 51);
	INSERT INTO enrollment VALUES (NULL, 3, 52);
	INSERT INTO enrollment VALUES (NULL, 3, 53);
	INSERT INTO enrollment VALUES (NULL, 3, 54);
	INSERT INTO enrollment VALUES (NULL, 3, 55);
	INSERT INTO enrollment VALUES (NULL, 3, 56);
	INSERT INTO enrollment VALUES (NULL, 3, 57);
	INSERT INTO enrollment VALUES (NULL, 3, 58);
	INSERT INTO enrollment VALUES (NULL, 3, 59);
	INSERT INTO enrollment VALUES (NULL, 3, 60);
	INSERT INTO enrollment VALUES (NULL, 3, 61);
	INSERT INTO enrollment VALUES (NULL, 3, 62);
	INSERT INTO enrollment VALUES (NULL, 3, 63);
	INSERT INTO enrollment VALUES (NULL, 3, 64);
	INSERT INTO enrollment VALUES (NULL, 3, 65);
	INSERT INTO enrollment VALUES (NULL, 3, 66);
	INSERT INTO enrollment VALUES (NULL, 3, 67);
	INSERT INTO enrollment VALUES (NULL, 3, 68);
	INSERT INTO enrollment VALUES (NULL, 3, 69);
	INSERT INTO enrollment VALUES (NULL, 3, 70);
	INSERT INTO enrollment VALUES (NULL, 3, 71);
	INSERT INTO enrollment VALUES (NULL, 3, 72);
	INSERT INTO enrollment VALUES (NULL, 3, 73);
	INSERT INTO enrollment VALUES (NULL, 3, 74);
	INSERT INTO enrollment VALUES (NULL, 3, 75);
	INSERT INTO enrollment VALUES (NULL, 3, 76);
	INSERT INTO enrollment VALUES (NULL, 3, 77);
	INSERT INTO enrollment VALUES (NULL, 3, 78);
	INSERT INTO enrollment VALUES (NULL, 3, 79);
	INSERT INTO enrollment VALUES (NULL, 3, 80);
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 1, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 2, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 3, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 4, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 5, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 6, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 7, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 8, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 9, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 10, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 11, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 12, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 13, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 14, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 15, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 16, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 17, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 18, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 19, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 20, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 21, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 22, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 23, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 24, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 25, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 26, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 27, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 28, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 29, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 30, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 31, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 32, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 33, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 34, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 35, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 36, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 37, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 38, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 39, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 40, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 41, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 42, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 43, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 44, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 45, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 46, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 47, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 48, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 49, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 50, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 51, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 52, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 53, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 54, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 55, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 56, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 57, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 58, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 59, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 60, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 61, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 62, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 63, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 64, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 65, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 66, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 67, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 68, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 69, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 70, '2022-01-01 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 71, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 72, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 73, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 74, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 75, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 76, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 77, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 78, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 79, '2022-01-02 10:02:00');
	INSERT INTO participants(session_id, user_id, joined_at) VALUES (3, 80, '2022-01-01 10:02:00');
	INSERT INTO user VALUES (5, 'janesmith@yahoo.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Smith', 'Jane', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (6, 'johnsmith@hotmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Smith', 'John', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (7, 'sarahlee@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Lee', 'Sarah', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (8, 'michaelchang@hotmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Chang', 'Michael', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (9, 'davidnguyen@yahoo.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Nguyen', 'David', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (10, 'amandafranklin@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Franklin', 'Amanda', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (11, 'johndoe1@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Anderson1', 'Alice', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (12, 'johndoe2@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Brown2', 'Benjamin', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (13, 'johndoe3@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Campbell3', 'Caroline', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (14, 'johndoe4@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Davis4', 'Daniel', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (15, 'johndoe5@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Edwards5', 'Emily', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (16, 'johndoe6@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Fisher6', 'Frank', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (17, 'johndoe7@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Gibson7', 'Gabriella', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (18, 'johndoe8@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Harris8', 'Henry', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (19, 'johndoe9@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Ingram9', 'Isabella', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (20, 'johndoe10@gmail.com', '$2y$10$GHDgqrCPap7Ej9vK25aZugZSh.ipt3D5d7GI1vFdycXERZpBdyri%', 'Ingram9', 'Isabella', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (21, 'janedoe1@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Adams1', 'Anna', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (22, 'janedoe2@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Brown2', 'Brandon', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (23, 'janedoe3@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Clark3', 'Cameron', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (24, 'janedoe4@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Davis4', 'Dylan', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (25, 'janedoe5@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Evans5', 'Ellie', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (26, 'janedoe6@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Fisher6', 'Finn', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (27, 'janedoe7@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Garcia7', 'Gabriel', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (28, 'janedoe8@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Hernandez8', 'Haley', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (29, 'janedoe9@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Irwin9', 'Isaac', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (30, 'janedoe10@gmail.com', '$2y$10$GHDgqrCPap7Ej9vK25aZugZSh.ipt3D5d7GI1vFdycXERZpBdyri%', 'Irwin9', 'Isaac', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (31, 'janedoe11@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Jones11', 'Jasmine', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (32, 'janedoe12@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Kim12', 'Kaitlyn', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (33, 'janedoe13@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Lee13', 'Landon', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (34, 'janedoe14@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Martinez14', 'Megan', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (35, 'janedoe15@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Nelson15', 'Noah', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (36, 'janedoe16@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Owens16', 'Oliver', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (37, 'janedoe17@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Perez17', 'Paige', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (38, 'janedoe18@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Quinn18', 'Quinn', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (39, 'janedoe19@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Rodriguez19', 'Riley', datetime('now'), datetime('now'));
	INSERT INTO user VALUES (40, 'janedoe20@gmail.com', '$2a$10$6rF4ewi/ZealdOt9ghvYJeyA4Oh/VKME/kzbd7Yw3MdL5.frlKNae', 'Rodriguez19', 'Riley', datetime('now'), datetime('now'));
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (44, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (45, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (46, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (47, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (48, 1, 'absent');
	INSERT INTO attendance (section_id, instructor_id, class_session_id, attended_at, created_at, updated_at) VALUES (3, 3, 3, '2023-05-03 09:05:00', '2023-05-01 09:05:00', '2023-05-01 09:05:00');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (4, 1, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (5, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (6, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (7, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (8, 1, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (9, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (10, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (11, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (12, 1, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (13, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (14, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (15, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (16, 1, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (17, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (18, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (19, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (20, 1, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (21, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (22, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (23, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (24, 1, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (25, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (26, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (27, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (28, 1, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (29, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (30, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (31, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (32, 1, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (33, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (34, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (35, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (36, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (37, 1, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (38, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (39, 1, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (40, 1, 'present');
	INSERT INTO attendance (section_id, instructor_id, class_session_id, attended_at, created_at, updated_at) VALUES (1, 3, 1, '2023-05-01 09:05:00', '2023-05-01 09:05:00', '2023-05-01 09:05:00');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (5, 2, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (6, 2, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (7, 2, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (8, 2, 'present');
	INSERT INTO attendance (section_id, instructor_id, class_session_id, attended_at, created_at, updated_at) VALUES (2, 3, 2, '2023-05-02 09:05:00', '2023-05-01 09:05:00', '2023-05-01 09:05:00');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (4, 3, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (6, 3, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (10, 3, 'absent');
	INSERT INTO attendance (section_id, instructor_id, class_session_id, attended_at, created_at, updated_at) VALUES (4, 3, 4, '2023-05-04 09:05:00', '2023-05-01 09:05:00', '2023-05-01 09:05:00');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (4, 4, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (5, 4, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (6, 4, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (7, 4, 'absent');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (8, 4, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (9, 4, 'present');
	INSERT INTO user_attendance(user_id, attendance_id, status) VALUES (10, 4, 'absent');

		`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		return fmt.Errorf("unable to insert: %v", err)
	}

	return nil
}
