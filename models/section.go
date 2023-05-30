package models

import (
	"database/sql"
	"strings"
)

type Section struct {
	ID        int
	CourseId  int
	Name      string
	CreatedAt string
	UpdatedAt string
	Enrolled  bool
}

type Schedual struct {
	ID        int
	SectionId int
	Day       string
	Time      string
}

type Enrollment struct {
	ID        int
	SectionId int
	UserId    int
}

// ** CREATE **
// AddEnrollment takes a section id and a user id.
// It returns any encountered errors.
func (s Section) AddEnrollment(SectionId int, UserId int) error {
	db := NewDB()

	sqlStatement := `
	INSERT INTO
		enrollment
		(section_id,
		user_id)
	VALUES
		($1,
		$2)
	`
	_, err := db.Exec(sqlStatement, SectionId, UserId)
	if err != nil {
		return err
	}

	return nil
}

// AddClassSession adds a course section and schedual to a class session table by section id
func (s *Section) AddClassSession(sectionID, scheduleID int) error {
	db := NewDB()
	defer db.Close()

	// add the course section to the class session table
	sqlStatement := `
	INSERT INTO
		class_session
	VALUES
		(NULL,
		$1,
		$2,
		0,
		'Session data',
		datetime('now'), 
		datetime('now')
		)`
	_, err := db.Exec(sqlStatement, sectionID, scheduleID)
	if err != nil {
		return err
	}

	return nil
}

// CreateSchedule given a section id and a schedule, it creates a schedule for a section and returns the schedule id
func (s *Section) CreateSchedule(sectionID int, schedule string) (int, error) {
	db := NewDB()
	defer db.Close()

	// Split the schedule string into an array of days and times.
	scheduleArray := strings.Split(schedule, "|")
	scheduleDays := scheduleArray[0]
	scheduleTime := scheduleArray[1]

	// Add the course section to the class session table and return the inserted ID.
	sqlStatement := `
	INSERT INTO
		schedule
	VALUES
		(NULL,
		$1,
		$2,
		$3
		)
	RETURNING id
	`
	var scheduleID int
	err := db.QueryRow(sqlStatement, sectionID, scheduleDays, scheduleTime).Scan(&scheduleID)
	if err != nil {
		return 0, err
	}

	return scheduleID, nil
}

// ** READ **
// For admin course table to get all
// Count returns the number of sections in the database.
// It returns the number of sections and any encountered errors.
func (s Section) Count() (int, error) {
	var count int
	db := NewDB()
	sqlStatement := `
	SELECT
	 COUNT(*)
	FROM
	 section
	;`
	row := db.QueryRow(sqlStatement)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountByInstructor takes an instructor id and returns the number of sections.
// It returns the number of sections and any encountered errors.
func (s Section) CountByInstructor(instructorID int) (int, error) {
	var count int
	db := NewDB()
	sqlStatement := `
	SELECT
	 COUNT(*)
	FROM
	 section
	JOIN
	 moderator
	ON
	 section.id = moderator.section_id
	AND
	 moderator.user_id = ?
	AND
	 moderator.type = 'instructor'
	;`
	row := db.QueryRow(sqlStatement, instructorID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Get takes a course section id.
// Returns a section struct and any encountered errors.
func (s Section) Get(sectionID int) (Section, error) {
	db := NewDB()

	var ID int
	var Name string
	var CreatedAt string
	var UpdatedAt string
	var CourseId int
	sqlStatement := `
	SELECT
		id,
	 	course_id,
	 	name,
	 	created_at, 
		updated_at 
	FROM
		section
	WHERE
	  	id = $1;
	  `
	row := db.QueryRow(sqlStatement, CourseId)

	switch err := row.Scan(&ID, &CourseId, &Name, &CreatedAt, &UpdatedAt); err {
	case sql.ErrNoRows:
		return Section{}, err
	case nil:
		section := Section{ID, CourseId, Name, CreatedAt, UpdatedAt, false}
		return section, err
	default:
		return Section{}, err
	}
}

// GetByCourse takes a course id.
// It returns a slice of section structs and any encountered errors.
func (s Section) GetByCourse(CourseId int) ([]Section, error) {
	var sections []Section
	db := NewDB()
	defer db.Close()
	sqlStatement := `
	SELECT
		id,
	 	course_id,
		name,
		created_at,
		updated_at
	FROM
		section
	WHERE
		course_id = $1;
	`
	rows, err := db.Query(sqlStatement, CourseId)
	if err != nil {
		return sections, err
	}
	defer rows.Close()
	for rows.Next() {
		var ID int
		var CourseId int
		var Name string
		var CreatedAt string
		var UpdatedAt string
		if err := rows.Scan(&ID, &CourseId, &Name, &CreatedAt, &UpdatedAt); err != nil {
			return sections, err
		}
		section := Section{ID, CourseId, Name, CreatedAt, UpdatedAt, false}
		sections = append(sections, section)
	}
	if err := rows.Err(); err != nil {
		return sections, err
	}

	return sections, nil
}

// GetEnrolledSections takes a user id.
// It returns a slice of section ids and any encountered errors.
func (s Section) GetEnrolledSections(UserId int) ([]int, error) {
	db := NewDB()

	sqlStatement := `
	SELECT
		section_id
	FROM
		enrollment
	WHERE
		user_id = $1
	`
	rows, err := db.Query(sqlStatement, UserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sectionIDs []int
	for rows.Next() {
		var sectionID int
		if err := rows.Scan(&sectionID); err != nil {
			return nil, err
		}
		sectionIDs = append(sectionIDs, sectionID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sectionIDs, nil
}

// CheckSectionEnrollment takes a user id and a course number.
// It returns a section id and any encountered errors.
func (s Section) CheckSectionEnrollment(UserId int, CourseNumber string) (int, error) {
	var SectionId int
	db := NewDB()
	sqlStatement := `
	SELECT
		section_id
	FROM
		enrollment
		INNER JOIN section ON enrollment.section_id = section.id
		INNER JOIN course ON section.course_id = course.id
	WHERE
		enrollment.user_id = $1
		AND course.number = $2;
	`
	row := db.QueryRow(sqlStatement, UserId, CourseNumber)

	switch err := row.Scan(&SectionId); err {
	case sql.ErrNoRows:
		return SectionId, err
	case nil:
		return SectionId, err
	default:
		return SectionId, err
	}
}

// GetEnrolledUsersBySectionID gets the list of enrollment structs for users who are enrolled in the section.
// It returns a slice of Enrollment structs and any encountered errors.
func (s *Enrollment) GetEnrolledUsersBySectionID(sectionID int) ([]Enrollment, error) {
	var enrollments []Enrollment
	db := NewDB()
	defer db.Close()
	sqlStatement := `
        SELECT
            enrollment.id,
            enrollment.section_id,
            enrollment.user_id
        FROM
            enrollment
        WHERE
            enrollment.section_id = $1
    `
	rows, err := db.Query(sqlStatement, sectionID)
	if err != nil {
		return enrollments, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, sectionId, userId int
		if err := rows.Scan(&id, &sectionId, &userId); err != nil {
			return enrollments, err
		}
		enrollment := Enrollment{ID: id, SectionId: sectionId, UserId: userId}
		enrollments = append(enrollments, enrollment)
	}
	return enrollments, nil
}

// GetSchedualBySectionID takes a section id.
// It returns a Schedual struct and any encountered errors.
func (s *Schedual) GetSchedualBySectionID(sectionID int) (Schedual, error) {
	db := NewDB()
	defer db.Close()

	var schedual Schedual

	sqlStatement := `
	SELECT
		id,
		section_id,
		day,
		timeslot
	FROM
		schedule
	WHERE
		section_id = $1
	LIMIT 1
	`
	err := db.QueryRow(sqlStatement, sectionID).Scan(&schedual.ID, &schedual.SectionId, &schedual.Day, &schedual.Time)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle no rows returned error here if needed
			return schedual, nil
		} else {
			return schedual, err
		}
	}

	return schedual, nil
}

// ** DELETE **
// DeleteByID takes a section id.
// It returns any encountered errors.
func (s Section) DeleteByID(id int) error {
	db := NewDB()

	sqlStatement := `
	DELETE FROM
		enrollment
	WHERE
		section_id = $1
	`
	result, err := db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteBySectionId takes a section id and a user id.
// It returns any encountered errors.
func (s Section) DeleteBySectionId(SectionId int, UserId int) error {
	db := NewDB()
	sqlStatement := `
	DELETE FROM
		enrollment
    WHERE
		section_id
	IN (
    SELECT
		id
    FROM
		section
    WHERE 
		course_id = (
    SELECT
		course_id
    FROM
		section
    WHERE
		id = $1))
	AND
		user_id = $2
	`
	result, err := db.Exec(sqlStatement, SectionId, UserId)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

// DeleteByCourseIdAndName takes a course id and a section number.
// It returns any encountered errors.
func (s Section) DeleteByCourseIdAndSection(courseID int, sectionNumber int) error {
	db := NewDB()

	sqlStatement := `
	DELETE FROM
		section
	WHERE
		course_id = $1
	AND
	 	name = $2
	`
	_, err := db.Exec(sqlStatement, courseID, sectionNumber)
	if err != nil {
		return err
	}
	return nil
}

// ** TESTING **
// THIS IS ONLY FOR TESTING PURPOSES IN MODEL TESTS
// GetSectionIDByCourseIDAndSectionNumber takes a course id and a section number.
// It returns a section id and any encountered errors.
func (s Section) GetSectionIDByCourseIDAndSectionNumber(courseID int, sectionNumber int) (int, error) {
	var sectionID int
	db := NewDB()
	sqlStatement := `
	SELECT
		id
	FROM
		section
	WHERE
		course_id = $1
	AND
	 	name = $2
	`
	row := db.QueryRow(sqlStatement, courseID, sectionNumber)

	switch err := row.Scan(&sectionID); err {
	case sql.ErrNoRows:
		return sectionID, err
	case nil:
		return sectionID, err
	default:
		return sectionID, err
	}
}
