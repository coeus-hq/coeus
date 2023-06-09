package models

var ddl_blank = []string{
	// ******************************************************************************
	// Create tables
	// ******************************************************************************
	`DROP TABLE IF EXISTS user`,
	`DROP TABLE IF EXISTS site`,
	`DROP TABLE IF EXISTS setting`,
	`DROP TABLE IF EXISTS course`,
	`DROP TABLE IF EXISTS section`,
	`DROP TABLE IF EXISTS enrollment`,
	`DROP TABLE IF EXISTS class_session`,
	`DROP TABLE IF EXISTS question`,
	`DROP TABLE IF EXISTS participants`,
	`DROP TABLE IF EXISTS schedule`,
	`DROP TABLE IF EXISTS vote`,
	`DROP TABLE IF EXISTS moderator`,
	`DROP TABLE IF EXISTS user_organization`,
	`DROP TABLE IF EXISTS organization`,
	`DROP TABLE IF EXISTS verify_user`,
	`DROP TABLE IF EXISTS attendance`,
	`DROP TABLE IF EXISTS user_attendance`,
	`DROP TABLE IF EXISTS is_admin`,
	`DROP TABLE IF EXISTS schema_version`,

	`CREATE TABLE schema_version (
        version INTEGER PRIMARY KEY,
        applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`,

	`CREATE TABLE is_admin (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL
    );`,

	`CREATE TABLE user(
       id INTEGER PRIMARY KEY AUTOINCREMENT,
       email VARCHAR(128) NOT NULL,
       hash TEXT NOT NULL,
       last_name VARCHAR(40) NOT NULL,
       first_name VARCHAR(40) NOT NULL,
       created_at TEXT NOT NULL,
       updated_at TEXT NOT NULL
    )`,

	`CREATE TABLE user_organization (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        organization_id INTEGER NOT NULL
    );`,

	`CREATE TABLE setting(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER NOT NULL,
      theme BOOLEAN NOT NULL,
      timezone_offset INTEGER NOT NULL
    )`,

	`CREATE TABLE organization(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name VARCHAR(256) NOT NULL,
      organization_timezone VARCHAR(5),
      logo_path TEXT,
      api_key TEXT,
      email TEXT,
      onboarding_complete BOOLEAN NOT NULL,
      created_at TEXT NOT NULL,
      updated_at TEXT NOT NULL,
	  is_demo BOOLEAN NOT NULL
    )`,

	`CREATE TABLE course(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      number TEXT NOT NULL,
      title TEXT NOT NULL,
      start_date TEXT NOT NULL,
      end_date TEXT NOT NULL,
      semester TEXT NOT NULL,
      year INTEGER NOT NULL,
      created_at TEXT NOT NULL,
      updated_at TEXT NOT NULL
    )`,

	`CREATE TABLE section(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      course_id INTEGER NOT NULL REFERENCES course(id),
      name VARCHAR(256) NOT NULL,
      created_at TEXT NOT NULL,
      updated_at TEXT NOT NULL
    )`,

	`CREATE TABLE enrollment(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      section_id INTEGER NOT NULL REFERENCES section(id),
      user_id INTEGER NOT NULL REFERENCES user(id)
    )`,

	`CREATE TABLE schedule(
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      section_id INTEGER NOT NULL REFERENCES section(id),
      day STRING NOT NULL,
      timeslot STRING NOT NULL
    )`,

	`CREATE TABLE class_session(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        section_id INTEGER NOT NULL,
        schedule_id INTEGER NOT NULL,
        in_progress BOOLEAN NOT NULL,
        data TEXT NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
    )`,

	`CREATE TABLE attendance(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        section_id INTEGER NOT NULL,
        instructor_id INTEGER NOT NULL,
        class_session_id INTEGER NOT NULL,
        attended_at TEXT NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
    )`,

	`CREATE TABLE user_attendance(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        attendance_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        status TEXT CHECK(status IN ('present', 'absent', 'excused', 'late'))
    )`,

	`CREATE TABLE question(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        session_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        text VARCHAR(140) NOT NULL,
        votes INTEGER NOT NULL,
        answered BOOLEAN NOT NULL,
        created_at TEXT NOT NULL,
        updated_at TEXT NOT NULL
    )`,

	`CREATE TABLE participants(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        session_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        joined_at TEXT NOT NULL
    )`,

	`CREATE TABLE vote(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        question_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL
    )`,

	`CREATE TABLE moderator(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        section_id INTEGER,
		type TEXT CHECK(type IN ('student', 'moderator', 'teacher assistant', 'instructor'))    )`,

	`CREATE TABLE verify_user(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        email TEXT NOT NULL,
        token TEXT NOT NULL,
        expiration TEXT NOT NULL,
        created_at TEXT NOT NULL,
        status TEXT CHECK(status IN ('pending', 'verified', 'expired'))
    )`,
}
