package models

import (
	"database/sql"
)

type Attendance struct {
	ID           int
	Section      int
	InstructorID int
	ClassSession int
	Date         string
	CreatedAt    string
	UpdatedAt    string
	CourseSection
}

type UserAttendance struct {
	ID         int
	Attendance int
	UserID     int
	Status     string
}

type CourseSection struct {
	SectionID    int
	CourseID     int
	CourseNumber string
	CourseTitle  string
	SectionName  string
	AttendanceID int
	Date         string
}

type AttendanceRecord struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Status    string
}

// ** CREATE **

// Add a new attendance record to the database for a given class session.
// It returns the ID of the new attendance record on success and any encountered error.
func (a *Attendance) Add(classSessionID, sectionIDInt, instructorID int) (int, error) {
	db := NewDB()
	var result sql.Result
	var ID int64
	// Add attendance record
	sqlStatement := `
		INSERT INTO
			attendance
		VALUES
			(NULL,
			$1,	
			$2,
			$3,
			datetime('now'),
			datetime('now'),
            datetime('now')
			)`
	result, err := db.Exec(sqlStatement, sectionIDInt, instructorID, classSessionID)
	if err != nil {
		return 0, err
	}

	// Get ID of newly created attendance record
	ID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ID), nil
}

// Add a new user_attendance record to the database for a given user and attendance record.
// It returns the ID of the new user_attendance record on success and any encountered error.
func (a *Attendance) AddUserAttendance(attendanceID int, userID int, status string) (int64, error) {
	db := NewDB()
	var result sql.Result
	var ID int64
	// Add user_attendance record
	sqlStatement := `
		INSERT INTO
			user_attendance
		VALUES
			(NULL,
			$1,
			$2,
			$3)`
	result, err := db.Exec(sqlStatement, attendanceID, userID, status)
	if err != nil {
		return 0, err
	}

	// Get ID of newly created user_attendance record
	ID, err = result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return ID, nil
}

// ** READ **

// GetCoursesByInstructor gets all courses associated with a given instructor.
// It returns a slice of CourseSection structs and any encountered error.
func (c *CourseSection) GetCoursesByInstructor(instructorID int) ([]CourseSection, error) {
	db := NewDB()
	var courses []CourseSection

	sqlStatement := `
		SELECT DISTINCT
			c.id as course_id,
			c.number as course_number,
			c.title as course_title
		FROM
			course c
		JOIN
			section s ON c.id = s.course_id
		JOIN
			attendance a ON s.id = a.section_id
		WHERE
			a.instructor_id = $1
		`
	rows, err := db.Query(sqlStatement, instructorID)
	if err != nil {
		return courses, err
	}
	defer rows.Close()

	for rows.Next() {
		var cs CourseSection
		err = rows.Scan(&cs.CourseID, &cs.CourseNumber, &cs.CourseTitle)
		if err != nil {
			return courses, err
		}
		courses = append(courses, cs)
	}

	return courses, nil
}

// GetSectionsByInstructorAndCourse gets all sections associated with a given course id and instructor id.
// It returns a slice of CourseSection structs and any encountered error.
func (c *CourseSection) GetSectionsByInstructorAndCourse(instructorID int, courseID int) ([]CourseSection, error) {
	db := NewDB()
	var sections []CourseSection

	sqlStatement := `
		SELECT DISTINCT
			s.id as section_id,
			c.number as course_number,
			c.title as course_title,
			s.name as section_name,
			a.id as attendance_id,
			strftime('%Y-%m-%d', a.attended_at) as date
		FROM
			section s
		JOIN
			course c ON s.course_id = c.id
		JOIN
			attendance a ON s.id = a.section_id
		WHERE
			a.instructor_id = $1
		AND
			c.id = $2
		`
	rows, err := db.Query(sqlStatement, instructorID, courseID)
	if err != nil {
		return sections, err
	}
	defer rows.Close()

	for rows.Next() {
		var cs CourseSection
		err = rows.Scan(&cs.SectionID, &cs.CourseNumber, &cs.CourseTitle, &cs.SectionName, &cs.AttendanceID, &cs.Date)
		if err != nil {
			return sections, err
		}
		sections = append(sections, cs)
	}

	return sections, nil
}

// GetByInstructor gets all attendance records for a given instructor along with the course title and section info for that record.
// It returns a slice of Attendance structs and any encountered error.
func (a *Attendance) GetByInstructor(instructorID int) ([]Attendance, error) {
	db := NewDB()
	var attendance []Attendance

	sqlStatement := `
		SELECT
			a.id,
			a.section_id,
			a.instructor_id,
			a.class_session_id,
			strftime('%Y-%m-%d', a.attended_at) as attended_at,
			a.created_at,
			a.updated_at,
			c.number as course_number,
			c.title as course_title,
			s.name as section_name
		FROM
			attendance a
		JOIN
			section s ON a.section_id = s.id
		JOIN
			course c ON s.course_id = c.id
		WHERE
			a.instructor_id = $1
		`
	rows, err := db.Query(sqlStatement, instructorID)
	if err != nil {
		return attendance, err
	}
	defer rows.Close()

	for rows.Next() {
		var a Attendance
		err = rows.Scan(&a.ID, &a.Section, &a.InstructorID, &a.ClassSession, &a.Date, &a.CreatedAt, &a.UpdatedAt, &a.CourseNumber, &a.CourseTitle, &a.SectionName)
		if err != nil {
			return attendance, err
		}
		attendance = append(attendance, a)
	}

	return attendance, nil
}

// Get all attendance records for a given class session.
// It returns a slice of Attendance structs and any encountered error.
func (a *Attendance) GetByClassSession(classSessionID int) ([]Attendance, error) {
	db := NewDB()
	var attendance []Attendance

	sqlStatement := `
		SELECT 
			* 
		FROM 
			attendance 
		WHERE 
			class_session_id = $1
		`
	rows, err := db.Query(sqlStatement, classSessionID)
	if err != nil {
		return attendance, err
	}
	defer rows.Close()

	for rows.Next() {
		var a Attendance
		err = rows.Scan(&a.ID, &a.Section, &a.InstructorID, &a.ClassSession, &a.Date, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return attendance, err
		}
		attendance = append(attendance, a)
	}

	return attendance, nil
}

// GetStudentsByAttendance gets all students for a given attendance ID.
// It returns a slice of User structs and any encountered error.
func (ar *AttendanceRecord) GetStudentsByAttendance(attendanceID int) ([]AttendanceRecord, error) {
	db := NewDB()
	var attendanceRecords []AttendanceRecord

	sqlStatement := `
	SELECT
		user.id,
		user.email,
		user.first_name,
		user.last_name,
		user_attendance.status
	FROM
		user
	JOIN
		user_attendance ON user.id = user_attendance.user_id
	JOIN
		attendance ON user_attendance.attendance_id = attendance.id
	WHERE
		attendance.id = $1
		`

	rows, err := db.Query(sqlStatement, attendanceID)
	if err != nil {
		return attendanceRecords, err
	}
	defer rows.Close()

	for rows.Next() {
		var ar AttendanceRecord
		err = rows.Scan(&ar.ID, &ar.Email, &ar.FirstName, &ar.LastName, &ar.Status)
		if err != nil {
			return attendanceRecords, err
		}
		attendanceRecords = append(attendanceRecords, ar)
	}

	return attendanceRecords, nil
}

// GetAttendanceIDBySectionID gets the attendance ID for a given section ID returning the most recent attendance record for that section.
// It returns the attendance ID and any encountered error.
func (a *Attendance) GetAttendanceIDBySectionID(sectionID int) (int, error) {
	db := NewDB()
	var attendanceID int

	sqlStatement := `
		SELECT
			id
		FROM
			attendance
		WHERE
			section_id = $1
		ORDER BY
			created_at DESC
		LIMIT 1
		`
	err := db.QueryRow(sqlStatement, sectionID).Scan(&attendanceID)
	if err != nil {
		return 0, err
	}

	return attendanceID, nil
}

// Get attendance for an individual user give the userID and attendance ID.
// It returns a slice of UserAttendance structs and any encountered error.
func (a *Attendance) GetByUser(userID int, attendanceID int) ([]UserAttendance, error) {
	db := NewDB()
	var attendance []UserAttendance

	sqlStatement := `
		SELECT 
			* 
		FROM 
			user_attendance 
		WHERE 
			user_id = $1
		AND
			attendance_id = $2
		`
	rows, err := db.Query(sqlStatement, userID, attendanceID)
	if err != nil {
		return attendance, err
	}
	defer rows.Close()

	for rows.Next() {
		var a UserAttendance
		err = rows.Scan(&a.ID, &a.UserID, &a.Status)
		if err != nil {
			return attendance, err
		}
		attendance = append(attendance, a)
	}

	return attendance, nil
}

// ** UPDATE **
// Update the status of a user_attendance record given the userID and attendanceID.
// It returns any encountered error.
func (a *Attendance) UpdateUserAttendance(userID int, attendanceID int, status string) error {
	db := NewDB()

	sqlStatement := `

		UPDATE
			user_attendance
		SET
			status = $1
		WHERE
			user_id = $2
		AND
			attendance_id = $3
	`
	_, err := db.Exec(sqlStatement, status, userID, attendanceID)
	if err != nil {
		return err
	}

	return nil
}

// ** DELETE **
// Delete all attendance records for a given class session ID.
func (a *Attendance) DeleteAll(classSessionID int) error {
	db := NewDB()

	// Get all attendanceIDs for the given classSessionID
	var attendanceIDs []int
	sqlSelect := `
        SELECT
            id
        FROM
            attendance
        WHERE
            class_session_id = $1
    `
	rows, err := db.Query(sqlSelect, classSessionID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var attendanceID int
		err = rows.Scan(&attendanceID)
		if err != nil {
			return err
		}
		attendanceIDs = append(attendanceIDs, attendanceID)
	}

	// Delete all user_attendance records for the given attendanceIDs
	for _, attendanceID := range attendanceIDs {
		err = a.DeleteUserAttendance(attendanceID)
		if err != nil {
			return err
		}
	}

	// Delete all attendance records for the given classSessionID
	sqlDelete := `
        DELETE FROM
            attendance
        WHERE
            class_session_id = $1
    `
	_, err = db.Exec(sqlDelete, classSessionID)
	if err != nil {
		return err
	}

	return nil
}

// Delete a single user_attendance record given the attendanceID.
// It returns any encountered error.
func (a *Attendance) DeleteUserAttendance(attendanceID int) error {
	db := NewDB()

	sqlStatement := `
		DELETE FROM
			user_attendance
		WHERE
			attendance_id = $1
	`
	_, err := db.Exec(sqlStatement, attendanceID)
	if err != nil {
		return err
	}

	return nil
}
