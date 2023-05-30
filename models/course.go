package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type Course struct {
	ID        int
	Number    string
	Title     string
	StartDate string
	EndDate   string
	Semester  string
	Year      int
	CreatedAt string
	UpdatedAt string
}

// ** CREATE **
// AddCourseAndSections takes a course number, title, start date, end date, semester, year, and number of sections.
// It returns a Course id and any encountered errors.
func (c Course) AddCourseAndSections(courseNumber, courseTitle, semester, courseStartDate, courseEndDate, year string, numberOfSections int) (int, error) {
	db := NewDB()
	sqlStatement := `
        INSERT INTO
            course
        VALUES(
			NULL,
            $1,
            $2, 
            $3, 
            $4, 
            $5, 
            $6,
            datetime('now'),
            datetime('now')
        )
        RETURNING
            id;
    `
	var courseID int
	err := db.QueryRow(sqlStatement, courseNumber, courseTitle, courseStartDate, courseEndDate, semester, year).Scan(&courseID)
	if err != nil {
		return courseID, err
	}

	for i := 1; i <= numberOfSections; i++ {
		sqlStatement = `
			INSERT INTO
				section
			VALUES(
				NULL,
				$1,
				$2,
				datetime('now'),
				datetime('now'))
			RETURNING
				id;
			`
		var sectionID int
		err = db.QueryRow(sqlStatement, courseID, i).Scan(&sectionID)
		if err != nil {
			return courseID, errors.New("error adding section to course")
		}

	}

	return courseID, nil
}

// ** READ **
// For admin course table to get all
// Count returns the number of courses in the database.
// It returns the number of courses and any encountered errors.
func (c Course) Count() (int, error) {
	var count int
	db := NewDB()
	sqlStatement := `
	SELECT
	 COUNT(*)
	FROM
	 course
	;`
	row := db.QueryRow(sqlStatement)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountByInstructor takes an instructor id and returns the number of courses.
// It returns the number of courses and any encountered errors.
func (c Course) CountByInstructor(instructorID int) (int, error) {
	var count int
	db := NewDB()
	sqlStatement := `
	SELECT
		COUNT(DISTINCT course.id)
	FROM
		course
	JOIN
		section
	ON
		course.id = section.course_id
	JOIN
		moderator
	ON
		section.id = moderator.section_id
	AND
		moderator.user_id = ?
	AND
		moderator.type = 'instructor';`
	row := db.QueryRow(sqlStatement, instructorID)
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// For admin course table to get all courses
// GetCourseSections returns a slice of maps containing all course information combind with section name, section id.
// It returns the slice of maps and any encountered errors.
func (c Course) GetCourseSections() ([]map[string]string, error) {
	db := NewDB()
	var sqlStatement string
	var data []map[string]string

	sqlStatement = fmt.Sprintf(`
	SELECT
		course.id AS course_id,
		course.number,
		course.title,
		course.semester,
		course.start_date,
		course.end_date,
		course.year,
		section.name,
		section.id AS section_id,
		COALESCE(GROUP_CONCAT(schedule.day || ' | ' || schedule.timeslot), '') AS schedule,
		COALESCE(num_students, 0) AS num_students
	FROM
 		course
    JOIN
    	section
    ON
		course.id = section.course_id
    LEFT JOIN
    	schedule
    ON
		section.id = schedule.section_id
    LEFT JOIN (
        SELECT
            section_id,
            COUNT(id) AS num_students
        FROM
            enrollment
        GROUP BY
            section_id
    ) AS
		enrollment_count
    ON
		section.id = enrollment_count.section_id
	GROUP BY
    	section.id
	;`)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var courseID int
		var number string
		var title string
		var semester string
		var startDate string
		var endDate string
		var year int
		var name string
		var sectionID int
		var schedule string
		var numStudents int
		err = rows.Scan(&courseID, &number, &title, &semester, &startDate, &endDate, &year, &name, &sectionID, &schedule, &numStudents)
		if err != nil {
			return nil, err
		}
		row := map[string]string{
			"courseID":    strconv.Itoa(courseID),
			"number":      number,
			"title":       title,
			"semester":    semester,
			"startDate":   startDate,
			"endDate":     endDate,
			"year":        strconv.Itoa(year),
			"name":        name,
			"sectionID":   strconv.Itoa(sectionID),
			"schedule":    schedule,
			"numStudents": strconv.Itoa(numStudents),
		}
		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetCourseSectionsByIntructor returns a slice of maps containing all course information combind with section name, section id.
// It returns the slice of maps and any encountered errors.
func (c Course) GetCourseSectionsByIntructor(userID int) ([]map[string]string, error) {
	db := NewDB()
	var data []map[string]string

	sqlStatement := `
	SELECT
		course.id AS course_id,
		course.number,
		course.title,
		course.semester,
		course.start_date,
		course.end_date,
		course.year,
		section.name,
		section.id AS section_id,
		COALESCE(GROUP_CONCAT(schedule.day || ' | ' || schedule.timeslot), '') AS schedule,
		COALESCE(num_students, 0) AS num_students
	FROM
		course
	JOIN
		section
	ON
		course.id = section.course_id
	JOIN
		moderator
	ON
		section.id = moderator.section_id
	AND
		moderator.user_id = $1
	AND
		moderator.type = 'instructor'
	LEFT JOIN
		schedule
	ON
		section.id = schedule.section_id
	LEFT JOIN (
	SELECT
		section_id,
		COUNT(id) AS num_students
	FROM
		enrollment
	GROUP BY
		section_id
	) AS
		enrollment_count
	ON
		section.id = enrollment_count.section_id
	GROUP BY
		section.id
	;`
	rows, err := db.Query(sqlStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var courseID int
		var number string
		var title string
		var semester string
		var startDate string
		var endDate string
		var year int
		var name string
		var sectionID int
		var schedule string
		var numStudents int
		err = rows.Scan(&courseID, &number, &title, &semester, &startDate, &endDate, &year, &name, &sectionID, &schedule, &numStudents)
		if err != nil {
			return nil, err
		}
		row := map[string]string{
			"courseID":    strconv.Itoa(courseID),
			"number":      number,
			"title":       title,
			"semester":    semester,
			"startDate":   startDate,
			"endDate":     endDate,
			"year":        strconv.Itoa(year),
			"name":        name,
			"sectionID":   strconv.Itoa(sectionID),
			"schedule":    schedule,
			"numStudents": strconv.Itoa(numStudents),
		}
		data = append(data, row)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Get takes a course number.
// It returns a course struct.
func (c Course) Get(CourseNumber string) (Course, error) {
	var ID int
	var Number string
	var Title string
	var StartDate string
	var EndDate string
	var Semester string
	var Year int
	var CreatedAt string
	var UpdatedAt string
	db := NewDB()
	sqlStatement := `
	SELECT
	 id,
	 number,
	 title,
	 start_date, 
	 end_date, 
	 semester,  
	 year, 
	 created_at, 
	 updated_at 
	FROM
	 course
	WHERE
	 number
	  = $1;
	  `
	row := db.QueryRow(sqlStatement, CourseNumber)

	switch err := row.Scan(&ID, &Number, &Title, &StartDate, &EndDate, &Semester, &Year, &CreatedAt, &UpdatedAt); err {
	case sql.ErrNoRows:
		return Course{}, err
	case nil:
		course := Course{ID, Number, Title, StartDate, EndDate, Semester, Year, CreatedAt, UpdatedAt}
		return course, err
	default:
		return Course{}, err
	}
}

// GetByUserId takes a user id.
// It returns a slice of maps containing the course number, title, semester, year,
// section name, section id, days, and timeslot.
func (c Course) GetByUserId(userID int) ([]map[string]string, error) {
	db := NewDB()
	var sqlStatement string
	var data []map[string]string

	sqlStatement = fmt.Sprintf(`
	SELECT
    course.number,
    course.title,
    course.semester,
    course.year,
    section.name,
    section.id,
    class_session.id AS class_session_id,
    moderator.type AS moderator_type,
    CAST(class_session.in_progress AS INTEGER) AS in_progress,
    schedule.day,
    MIN(schedule.timeslot) AS timeslot
FROM
    user
    JOIN
    enrollment
    ON user.id = enrollment.user_id
    JOIN
    section
    ON enrollment.section_id = section.id
    JOIN
    course
    ON section.course_id = course.id
    JOIN
    schedule
    ON section.id = schedule.section_id
    JOIN
    class_session
    ON section.id = class_session.section_id
    LEFT JOIN
    moderator
    ON user.id = moderator.user_id AND section.id = moderator.section_id
WHERE
    user.id = %d
GROUP BY
    course.number,
    course.title,
    course.semester,
    course.year,
    section.name,
    moderator.type
ORDER BY
    section.name ASC;

	`, userID)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var number, title, semester, year, name, sectionID, classSessionID, days, timeslot string
		var moderatorType sql.NullString
		var inProgress bool

		if err := rows.Scan(&number, &title, &semester, &year, &name, &sectionID, &classSessionID, &moderatorType, &inProgress, &days, &timeslot); err != nil {
			return nil, err
		}
		item := make(map[string]string)
		item["number"] = number
		item["title"] = title
		item["semester"] = semester
		item["year"] = year
		item["name"] = name
		item["sectionID"] = sectionID
		item["classSessionID"] = classSessionID
		if moderatorType.Valid { // Check if the value is not NULL
			item["moderatorType"] = moderatorType.String
		} else {
			item["moderatorType"] = ""
		}
		item["inProgress"] = strconv.FormatBool(inProgress)
		item["days"] = days
		item["timeslot"] = timeslot
		data = append(data, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

// Search takes a course identifier string.
// It returns a slice of Course structs.
func (c Course) Search(courseIdentifier string) ([]Course, error) {
	db := NewDB()
	sqlStatement := fmt.Sprintf(`
	SELECT
		id,
		number,
		title,
		start_date,
		end_date,
		semester,
		year,
		created_at,
		updated_at
	FROM
		course
	WHERE
		number
		LIKE
		'%%%s%%'
		OR
		title
		LIKE
		'%%%s%%'
		;
		`, courseIdentifier, courseIdentifier)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []Course
	for rows.Next() {
		var c Course
		err := rows.Scan(&c.ID, &c.Number, &c.Title, &c.StartDate, &c.EndDate, &c.Semester, &c.Year, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, errors.New("unsucessful scan of course row")
		}
		courses = append(courses, c)
	}

	return courses, err
}

// GetBySectionId takes a section id.
// It returns a Course struct and any encounted errors.
func (c Course) GetBySectionId(sectionID int) (Course, error) {
	var ID int
	var Number string
	var Title string
	var StartDate string
	var EndDate string
	var Semester string
	var Year int
	var CreatedAt string
	var UpdatedAt string
	db := NewDB()
	sqlStatement := `
	SELECT
	 course.id,
	 course.number,
	 course.title,
	 course.start_date,
	 course.end_date,
	 course.semester,
	 course.year,
	 course.created_at,
	 course.updated_at
	FROM
	 course
		JOIN
		section
		ON course.id = section.course_id
	WHERE
		section.id = $1;
	`
	row := db.QueryRow(sqlStatement, sectionID)

	switch err := row.Scan(&ID, &Number, &Title, &StartDate, &EndDate, &Semester, &Year, &CreatedAt, &UpdatedAt); err {
	case sql.ErrNoRows:
		return Course{}, err
	case nil:
		course := Course{ID, Number, Title, StartDate, EndDate, Semester, Year, CreatedAt, UpdatedAt}
		return course, err
	default:
		return Course{}, err
	}
}

// GetSectionIds returns an array of section ids for a given course id.
// It returns an array of section ids and any encountered errors.
func (c Course) GetSectionIds(courseID int) ([]int, error) {
	db := NewDB()
	sqlStatement := `
	SELECT
		id
	FROM
		section
	WHERE
		course_id = $1;
	`
	rows, err := db.Query(sqlStatement, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sectionIDs []int
	for rows.Next() {
		var sectionID int
		err := rows.Scan(&sectionID)
		if err != nil {
			return nil, errors.New("unsucessful scan of section id row")
		}
		sectionIDs = append(sectionIDs, sectionID)
	}

	return sectionIDs, err
}

// GetCourseByID takes a course id.
// It returns a Course struct and any encounted errors.
func (c Course) GetCourseByID(courseID int) (Course, error) {
	var ID int
	var Number string
	var Title string
	var StartDate string
	var EndDate string
	var Semester string
	var Year int
	var CreatedAt string
	var UpdatedAt string
	db := NewDB()
	sqlStatement := `
	SELECT
	 id,
	 number,
	 title,
	 start_date,
	 end_date,
	 semester,
	 year,
	 created_at,
	 updated_at
	FROM
	 course
	WHERE
	 id = $1;
	`
	row := db.QueryRow(sqlStatement, courseID)

	switch err := row.Scan(&ID, &Number, &Title, &StartDate, &EndDate, &Semester, &Year, &CreatedAt, &UpdatedAt); err {
	case sql.ErrNoRows:
		return Course{}, err
	case nil:
		course := Course{ID, Number, Title, StartDate, EndDate, Semester, Year, CreatedAt, UpdatedAt}
		return course, err
	default:
		return Course{}, err
	}
}

// ** UPDATE **
// UpdateCourseAndSection takes a course id, section id, course number, course title, semester, year, and section name and updates the course and section data.
// It returns any encountered errors.
func (c Course) UpdateCourseAndSection(courseID, sectionID int, courseNumber, courseTitle, semester, year, sectionName, courseStartDate, courseEndDate, scheduleDays, scheduleTime string) error {

	db := NewDB()
	sqlStatement := `
    UPDATE
        course
    SET
        number = $1,
        title = $2,
        semester = $3,
        year = $4,

        start_date = $5,
        end_date = $6,
        updated_at = datetime('now')
    WHERE
        id = $7;
    `
	_, err := db.Exec(sqlStatement, courseNumber, courseTitle, semester, year, courseStartDate, courseEndDate, courseID)
	if err != nil {
		return err
	}

	sqlStatement = `
    UPDATE
        section
    SET
        name = $1,
        updated_at = datetime('now')
    WHERE
        id = $2
    AND
        course_id = $3;
    `
	_, err = db.Exec(sqlStatement, sectionName, sectionID, courseID)
	if err != nil {
		return err
	}

	// Update the schedule data

	sqlStatement = `
	UPDATE
		schedule
	SET
		day = $1,
		timeslot = $2
	WHERE
		section_id = $3
	`
	_, err = db.Exec(sqlStatement, scheduleDays, scheduleTime, sectionID)
	if err != nil {
		return err
	}

	return nil
}
