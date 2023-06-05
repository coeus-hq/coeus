package models

import (
	_ "coeus/globals"
	"fmt"
	"testing"
)

func TestUserAuthenticate(t *testing.T) {
	var match int

	user := new(User)
	match = user.Authenticate("student@coeus.education", "coeus")

	if match == 0 {
		t.Fatalf("Authentication FAILED")
	}
}

func TestUserGet(t *testing.T) {
	user := new(User)

	// Get Details for a user given an ID
	_, error := user.Get(1)
	if error != nil {
		t.Fatal(error)
	}
}

func TestGetAllUser(t *testing.T) {
	user := new(User)

	// Get all users
	_, error := user.GetAll()
	if error != nil {
		t.Fatal(error)
	}
}

func TestUserAdd(t *testing.T) {
	u := new(User)

	// Add a new user
	id, err := u.Add("email@addexample.com", "password123", "Doe", "John")
	if err != nil {
		t.Fatal("Failed to add user")
	}
	u.Delete(id)

}

func TestUserUpdate(t *testing.T) {

	u := new(User)

	//Add a new user
	id, err := u.Add("email@updateexample.com", "password123", "Doe", "John")
	if err != nil {
		t.Fatal("Failed to add user")
	}

	//	Update the user
	err = u.Update(id, "email@updateexample.com", "Doe", "Jane", "new")
	if err != nil {
		t.Fatal("Failed to delete user")
	}

	// Delete the user
	err = u.Delete(id)
	if err != nil {
		t.Fatal("Failed to delete user")
	}

}

func TestUserDelete(t *testing.T) {
	u := new(User)

	// Add a new user
	id, err := u.Add("email@deleteexample.com", "password123", "Doe", "John")
	if err != nil {
		t.Fatal("Failed to add user")
	}

	// Delete the user
	err = u.Delete(id)
	if err != nil {
		t.Fatal("Failed to delete user")
	}

	// Check if user was deleted
	_, err = u.Get(id)
	if err == nil {
		t.Fatal("User was not deleted")
	}
}

func TestUserCount(t *testing.T) {
	user := new(User)

	// Get count of all users
	_, error := user.Count()
	if error != nil {
		t.Fatal(error)
	}
}

func TestSettingAdd(t *testing.T) {
	s := new(Setting)

	// Add a new setting
	_, err := s.Add(999)
	if err != nil {
		t.Fatalf("Failed to add setting: %v", err)
	}

	// Delete the setting
	err = s.Delete(999)
	if err != nil {
		t.Fatal("Failed to delete setting")
	}
}

func TestSettingDelete(t *testing.T) {
	s := new(Setting)

	// Add a new setting
	_, err := s.Add(999)
	if err != nil {
		t.Fatalf("Failed to add setting: %v", err)
	}

	// Delete the setting
	err = s.Delete(999)
	if err != nil {
		t.Fatal("Failed to delete setting")
	}
}

func TestSettingGet(t *testing.T) {
	setting := new(Setting)

	// Get Details for a user given an ID
	_, error := setting.Get(1)
	if error != nil {
		t.Fatal(error)
	}
}

func TestSettingUpdate(t *testing.T) {

	s := new(Setting)

	// Update settings for user id 1
	error := s.Update(2, "theme", "false", -360)

	if error != nil {
		t.Fatal(error)
	}
}

func TestCourseCount(t *testing.T) {
	course := new(Course)

	// Get count of all courses
	_, error := course.Count()
	if error != nil {
		t.Fatal(error)
	}
}

func TestCourseGet(t *testing.T) {
	course := new(Course)

	// Get a course with a known course number
	result, err := course.Get("CSI 1113")
	if err != nil {
		t.Fatal(err)
	}

	if result.Number != "CSI 1113" {
		t.Fatalf("Expected course number CS101, but got %s", result.Number)
	}
}

func TestCourseGetByUserId(t *testing.T) {
	course := new(Course)

	// Get courses by user ID
	courses, err := course.GetByUserId(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(courses) == 0 {
		t.Fatalf("Expected courses, but got %d", len(courses))
	}
}

func TestCourseSearch(t *testing.T) {
	course := new(Course)

	// Search for courses using a known course identifier
	results, err := course.Search("CSI 1113")
	if err != nil {
		t.Fatal(err)
	}

	if len(results) == 0 {
		t.Fatalf("Expected courses, but got %d", len(results))
	}
}

func TestCourseGetBySectionId(t *testing.T) {
	course := new(Course)

	// Get a course with a known section ID
	result, err := course.GetBySectionId(3)
	if err != nil {
		t.Fatal(err)
	}

	if result.ID != 2 {
		t.Fatalf("Expected course with ID 1, but got %d", result.ID)
	}
}

func TestGetCourseSections(t *testing.T) {
	course := new(Course)

	// Get sections for a course with the course data
	_, err := course.GetCourseSections()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateSchedule(t *testing.T) {
	s := new(Section)

	// Create a schedule for a user
	_, err := s.CreateSchedule(23, "TEST | TEST")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddCourseAndSections(t *testing.T) {
	c := new(Course)

	// Add a new course
	_, err := c.AddCourseAndSections("CSI 9999", "Test Course", "Fall", "2023-05-14", "2023-05-14", "2023", 2)

	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateCourseAndSection(t *testing.T) {
	c := new(Course)
	s := new(Section)

	// Add a new course
	ID, err := c.AddCourseAndSections("CSI 9999", "Test Course", "Fall", "2023-05-14", "2023-05-14", "2023", 2)

	// THIS FUNC IS ONLY FOR TESTING PURPOSES
	// Get the section ID of the new course based on the course ID and section number
	sectionID, err := s.GetSectionIDByCourseIDAndSectionNumber(ID, 1)

	// Update a course and its section
	err = c.UpdateCourseAndSection(ID, sectionID, "CSI 1113", "Test Course", "Fall", "2023", "5", "2022 Feb", "2022 May", "M, F", "12:00:00 PM")

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSectionIds(t *testing.T) {
	c := new(Course)

	// Get section IDs for a course
	_, err := c.GetSectionIds(1)

	if err != nil {
		t.Fatal(err)
	}
}

func TestAddClassSession(t *testing.T) {
	s := new(Section)

	// Add a class session
	err := s.AddClassSession(1, 3)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteByCourseIdAndSection(t *testing.T) {
	c := new(Section)

	// Delete a course and its sections
	err := c.DeleteByCourseIdAndSection(1, 1)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSectionCount(t *testing.T) {
	section := new(Section)

	// Get count of all sections
	_, error := section.Count()
	if error != nil {
		t.Fatal(error)
	}
}

func TestSectionGetByCourse(t *testing.T) {
	db := NewDB()
	// Insert a test section with a known CourseId
	_, err := db.Exec("INSERT INTO section (course_id, name, created_at, updated_at) VALUES ($1, 'Test Section', datetime('now'), datetime('now'))", 1)

	if err != nil {
		t.Fatal(err)
	}

	s := new(Section)
	sections, err := s.GetByCourse(1)
	if err != nil {
		t.Fatal(err)
	}

	if len(sections) == 0 {
		t.Fatalf("Expected sections, but got %d", len(sections))
	}

	// Clean up the test data
	_, err = db.Exec("DELETE FROM section WHERE course_id = $1", 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddEnrollment(t *testing.T) {

	db := NewDB()
	s := new(Section)

	// Add enrollment for the section CS50 1 and user Asim
	err := s.AddEnrollment(1, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the test data
	_, err = db.Exec("DELETE FROM enrollment WHERE section_id=$1 AND user_id=$2", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetEnrolledUsersBySectionID(t *testing.T) {

	s := new(Enrollment)

	// Get enrolled users for a section
	_, err := s.GetEnrolledUsersBySectionID(3)

	if err != nil {
		t.Fatal(err)
	}
}

func TestPostQuestion(t *testing.T) {

	s := new(Question)

	_, err := s.PostQuestion(1, 2, "This is a test question")

	if err != nil {
		t.Errorf("postQuestion failed with error: %v", err)
	}
}

func TestQuestionGetByID(t *testing.T) {

	s := new(Question)

	// Insert a test question
	s.PostQuestion(1, 2, "test")

	_, err := s.GetByID(1, 100)
	if err != nil {
		t.Errorf("getByID failed with error: %v", err)
	}

	// Clean up the test data
	db := NewDB()
	_, err = db.Exec("DELETE FROM question WHERE text = $1", "test")

}

func TestGetAllQuestions(t *testing.T) {

	s := new(Question)

	questions, err := s.GetAllQuestions(10, "votes", 300)

	if err != nil {
		t.Errorf("getAllQuestions failed with error: %v", err)
	}

	if len(questions) == 0 {
		t.Fatalf("Expected questions, but got %d", len(questions))
	}
}

func TestGetUnansweredQuestions(t *testing.T) {

	_, err := new(Question).GetUnansweredQuestions(3, 300, "votes")
	if err != nil {
		t.Errorf("getUnansweredQuestions failed with error: %v", err)
	}

}

func TestJoin(t *testing.T) {

	_, err := new(ClassSession).Join(1, 1)

	if err != nil {
		t.Errorf("join failed with error: %v", err)
	}
}

func TestGetParticipantCount(t *testing.T) {

	_, err := new(ClassSession).GetParticipantCount(1)

	if err != nil {
		t.Errorf("getParticipants failed with error: %v", err)
	}
}

func TestGetParticipants(t *testing.T) {

	_, err := new(Participant).GetParticipants(1)

	if err != nil {
		t.Errorf("getParticipants failed with error: %v", err)
	}
}

func TestGetSectionID(t *testing.T) {

	_, err := new(ClassSession).GetSectionID(1)

	if err != nil {
		t.Errorf("getSectionID failed with error: %v", err)
	}
}

func TestVoteQuestion(t *testing.T) {

	err := new(Question).VoteQuestion(1, 1)

	if err != nil {
		t.Errorf("voteQuestion failed with error: %v", err)
	}
}

func TestGetVoteCount(t *testing.T) {

	_, err := new(Question).GetVoteCount(1)

	if err != nil {
		t.Errorf("getVoteCount failed with error: %v", err)
	}
}

func TestGetBySectionId(t *testing.T) {

	_, err := new(Course).GetBySectionId(3)

	if err != nil {
		t.Errorf("getBySectionId failed with error: %v", err)
	}
}

func TestGetUserInitials(t *testing.T) {

	_, err := new(User).GetUserInitials(1)

	if err != nil {
		t.Errorf("getUserInitials failed with error: %v", err)
	}
}

func TestQuestionVoteQuestion(t *testing.T) {
	q := new(Question)

	// Add a vote to a question with a known question ID and user ID
	err := q.VoteQuestion(1, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestQuestionHasVoted(t *testing.T) {
	q := new(Question)
	// Check if a user has voted for a question with a known question ID and user ID
	result, err := q.HasVoted(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Fatalf("Expected user to have voted, but got %v", result)

	}

}
func TestAddModerator(t *testing.T) {

	err := new(Moderator).Delete(999, 3)
	if err != nil {
		t.Errorf("deleteModerator failed with error: %v", err)
	}

	_, err = new(Moderator).Add(999, 3, "student")
	if err != nil {
		t.Errorf("addModerator failed with error: %v", err)
	}

	err = new(Moderator).Delete(999, 3)
	if err != nil {
		t.Errorf("deleteModerator failed with error: %v", err)
	}
}

func TestUpdateModerator(t *testing.T) {

	_, err := new(Moderator).Add(999, 3, "student")
	if err != nil {
		t.Errorf("addModerator failed with error: %v", err)
	}

	_, err = new(Moderator).Update(999, 3, "teacher assistant")
	if err != nil {
		t.Errorf("updateModerator failed with error: %v", err)
	}

	err = new(Moderator).Delete(999, 3)
	if err != nil {
		t.Errorf("deleteModerator failed with error: %v", err)
	}

}

func TestGetModerator(t *testing.T) {

	_, err := new(Moderator).Add(999, 3, "student")
	if err != nil {
		t.Errorf("addModerator failed with error: %v", err)
	}

	_, err = new(Moderator).Get(999)
	if err != nil {
		t.Errorf("getModerator failed with error: %v", err)
	}

}

func TestGetStatusModerator(t *testing.T) {

	err := new(Moderator).Delete(999, 3)
	if err != nil {
		t.Errorf("deleteModerator failed with error: %v", err)
	}

	_, err = new(Moderator).Add(999, 3, "student")
	if err != nil {
		t.Errorf("addModerator failed with error: %v", err)
	}

	_, err = new(Moderator).GetStatus(999, 3)
	if err != nil {
		t.Errorf("getStatusModerator failed with error: %v", err)
	}

	err = new(Moderator).Delete(999, 3)
	if err != nil {
		t.Errorf("deleteModerator failed with error: %v", err)
	}
}

func TestGetInProgress(t *testing.T) {

	_, err := new(ClassSession).GetInProgress(1)
	if err != nil {
		t.Errorf("getInProgress failed with error: %v", err)
	}
}
func TestStart(t *testing.T) {

	_, err := new(ClassSession).Start(1)
	if err != nil {
		t.Errorf("start failed with error: %v", err)
	}
}

func TestEnd(t *testing.T) {

	_, err := new(ClassSession).Start(1)
	if err != nil {
		t.Errorf("start failed with error: %v", err)
	}

	err = new(ClassSession).End(1)
	if err != nil {
		t.Errorf("end failed with error: %v", err)
	}
}

func TestDeleteModerator(t *testing.T) {

	// First delete any existing moderator
	err := new(Moderator).Delete(999, 3)
	if err != nil {
		t.Errorf("deleteModerator failed with error: %v", err)
	}

	// Add a moderator with a known user ID and section ID
	_, err = new(Moderator).Add(999, 3, "student")
	if err != nil {
		t.Errorf("addModerator failed with error: %v", err)
	}

	// Than actually delete the added moderator
	err = new(Moderator).Delete(999, 3)
	if err != nil {
		t.Errorf("deleteModerator failed with error: %v", err)
	}
}

func TestIsInstructor(t *testing.T) {

	// First delete any existing moderator
	err := new(Moderator).Delete(999, 3)
	if err != nil {
		t.Errorf("deleteModerator failed with error: %v", err)
	}

	// Add a moderator with a known user ID and section ID
	_, err = new(Moderator).Add(999, 3, "instructor")
	if err != nil {
		t.Errorf("addModerator failed with error: %v", err)
	}

	// Check if a user with a known user ID is an instructor
	result, err := new(Moderator).IsInstructor(999)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Fatalf("Expected user to be an instructor, but got %v", result)
	}

	// Than actually delete the added moderator
	err = new(Moderator).Delete(999, 3)
	if err != nil {
		t.Errorf("deleteModerator failed with error: %v", err)
	}
}

func TestCountByInstructor(t *testing.T) {
	// Count the number of courses and sections taught by an instructor with a known user ID
	_, err := new(Course).CountByInstructor(3)
	if err != nil {
		fmt.Println(err)
	}
	_, err = new(Section).CountByInstructor(3)
	if err != nil {
		fmt.Println(err)
	}
}

func TestGetCourseSectionsByIntructor(t *testing.T) {
	// Get the sections taught by an instructor with a known user ID
	_, err := new(Course).GetCourseSectionsByIntructor(3)
	if err != nil {
		fmt.Println(err)
	}
}

func TestHasVoted(t *testing.T) {

	// Check if a user has voted for a question with a known question ID and user ID
	result, err := new(Question).HasVoted(1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !result {
		t.Fatalf("Expected user to have voted, but got %v", result)
	}
}

func TestHasVotedAll(t *testing.T) {

	// Get the userHasVoted slice from the database
	result, err := new(Question).HasVotedAll(3, 1)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatalf("Expected slice of structs with userHasVoted and a question id, but got %v", result)
	}
}

func TestHasVotedAppend(t *testing.T) {

	// Get the questions from the database by the created at time
	questionsByTimeSlice, err := new(Question).GetAllQuestions(3, "created_at", 300)
	if err != nil {
		fmt.Println(err)
	}

	// Get the userHasVoted slice from the database
	userHasVotedSlice, err := new(Question).HasVotedAll(3, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Append the userHasVoted slice to the questionsByTime slice
	result := new(Question).HasVotedAppend(questionsByTimeSlice, userHasVotedSlice)
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatalf("Expected questions and voted slice to be combined, but got %v", result)
	}
}

func TestGetStatus(t *testing.T) {

	_, err := new(Organization).GetStatus()
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetOrganizationID(t *testing.T) {

	// Get the organization ID from the database
	_, err := new(User).GetOrganizationID()
	if err != nil {
		t.Fatal(err)
	}

}

func TestAddUserToOrganization(t *testing.T) {

	// Add a user to an organization with a known user ID and organization ID
	err := new(User).AddUserToOrganization(1, 1)
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetOrganization(t *testing.T) {

	// Add an organization to the database
	ID, err := new(Organization).Add("Test Organization", "0", "logo.svg", "1234", "user@email")
	if err != nil {
		t.Fatal(err)
	}

	// Get the organization from the database
	organization, err := new(Organization).Get(ID)
	if err != nil {
		t.Fatal(err)
		t.Fatalf("Expected organization, but got %v", organization)
	}
}

func TestAddOrganization(t *testing.T) {

	// Add an organization to the database
	_, err := new(Organization).Add("Test Organization", "0", "logo.svg", "1234", "user@email")
	if err != nil {
		t.Fatal(err)
	}

}

func TestSetAdmin(t *testing.T) {

	// Set a user as an admin with a known user ID
	err := new(Organization).SetAdmin(999)
	if err != nil {
		t.Fatal(err)
	}

	// Remove a user as an admin with a known user ID
	err = new(Organization).DeleteAdmin(999)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAdminID(t *testing.T) {

	// Set a user as an admin with a known user ID
	err := new(Organization).SetAdmin(999)
	if err != nil {
		t.Fatal(err)
	}

	// Get the admin ID from the database
	_, err = new(Organization).GetAdminID()
	if err != nil {
		t.Fatal(err)
	}

	// Remove a user as an admin with a known user ID
	err = new(Organization).DeleteAdmin(999)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdminDelete(t *testing.T) {

	// Delete the user with a known user ID
	err := new(Moderator).AdminDelete(999)
	if err != nil {
		t.Fatal(err)
	}

	// Add a user with a known user ID
	_, err = new(Moderator).AdminAdd(999, "teacher assistant")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the user with a known user ID
	err = new(Moderator).AdminDelete(999)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdminAdd(t *testing.T) {

	err := new(Moderator).AdminDelete(999)
	if err != nil {
		t.Fatal(err)
	}

	_, err = new(Moderator).AdminAdd(999, "teacher assistant")
	if err != nil {
		t.Fatal(err)
	}

	err = new(Moderator).AdminDelete(999)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdminUpdate(t *testing.T) {

	err := new(Moderator).AdminDelete(999)
	if err != nil {
		t.Fatal(err)
	}

	// Add a user with a known user ID
	_, err = new(Moderator).AdminAdd(999, "teacher assistant")
	if err != nil {
		t.Fatal(err)
	}

	// Update the user with a known user ID
	_, err = new(Moderator).AdminUpdate(999, "instructor")
	if err != nil {
		t.Fatal(err)
	}

	err = new(Moderator).AdminDelete(999)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateToken(t *testing.T) {

	_, err := new(VerifyUser).CreateToken("whalencollin@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

}

func TestMatchToken(t *testing.T) {

	ID, err := new(VerifyUser).CreateToken("whalencollin@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	token, err := new(VerifyUser).GetToken(ID)

	bool, _, err := new(VerifyUser).MatchToken(token)
	if err != nil {
		t.Fatal(err)
	}
	if bool == false {
		t.Fatal("Expected true, but got false")
	}
}

func TestDeleteToken(t *testing.T) {

	// Create a token
	ID, err := new(VerifyUser).CreateToken("whalencollin@gmail.com")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the token
	err = new(VerifyUser).DeleteToken(ID)
	if err != nil {
		t.Fatal(err)
	}

}

func TestCheckAPIKey(t *testing.T) {

	// CheckAPIKey function to check if the API key exists true or false is considered a pass
	_, err := new(Organization).CheckAPIKey()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdatePassword(t *testing.T) {

	// Update the user's password
	err := new(User).UpdatePassword(999, "Ch0ngeme!")
	if err != nil {
		fmt.Println(err)
	}

}

func TestAttendanceAdd(t *testing.T) {

	// Add a user to the attendance table
	_, err := new(Attendance).Add(3, 1, 3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAttendanceAddUserAttendance(t *testing.T) {

	// Add a user to the attendance table
	ID, err := new(Attendance).Add(3, 1, 3)
	if err != nil {
		t.Fatal(err)
	}

	//convert ID from int64 to int
	IDint := int(ID)

	// Add a user to the user attendance table
	_, err = new(Attendance).AddUserAttendance(999, IDint, "present")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAttendanceGetByClassSession(t *testing.T) {

	// Get the attendance for a class session
	_, err := new(Attendance).GetByClassSession(3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAttendanceIDBySectionID(t *testing.T) {

	_, err := new(Attendance).Add(3, 3, 3)
	if err != nil {
		t.Fatal(err)
	}

	// Get the attendance ID for a section
	_, err = new(Attendance).GetAttendanceIDBySectionID(3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAttendanceGetByUser(t *testing.T) {

	// Add a user to the attendance table
	ID, err := new(Attendance).Add(3, 1, 3)
	if err != nil {
		t.Fatal(err)
	}

	//convert ID from int64 to int
	IDint := int(ID)

	// Add a user to the user attendance table
	_, err = new(Attendance).AddUserAttendance(999, IDint, "present")
	if err != nil {
		t.Fatal(err)
	}

	// Get the attendance for a user
	_, err = new(Attendance).GetByUser(999, IDint)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdateUserAttendance(t *testing.T) {

	// Add a user to the attendance table
	ID, err := new(Attendance).Add(3, 1, 3)
	if err != nil {
		t.Fatal(err)
	}

	//convert ID from int64 to int
	IDint := int(ID)

	// Add a user to the user attendance table
	_, err = new(Attendance).AddUserAttendance(999, IDint, "absent")
	if err != nil {
		t.Fatal(err)
	}

	// Update a user's attendance
	err = new(Attendance).UpdateUserAttendance(999, IDint, "present")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAttendanceDeleteAll(t *testing.T) {

	// Add a user to the attendance table
	ID, err := new(Attendance).Add(3, 1, 3)
	if err != nil {
		t.Fatal(err)
	}

	//convert ID from int64 to int
	IDint := int(ID)

	// Add a user to the user attendance table
	_, err = new(Attendance).AddUserAttendance(999, IDint, "absent")
	if err != nil {
		t.Fatal(err)
	}

	// Delete all attendance
	err = new(Attendance).DeleteAll(3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAttendanceDeleteUserAttendance(t *testing.T) {

	// Add a user to the attendance table
	ID, err := new(Attendance).Add(3, 1, 3)
	if err != nil {
		t.Fatal(err)
	}

	//convert ID from int64 to int
	IDint := int(ID)

	// Add a user to the user attendance table
	_, err = new(Attendance).AddUserAttendance(999, IDint, "absent")
	if err != nil {
		t.Fatal(err)
	}

	// Delete a user's attendance
	err = new(Attendance).DeleteUserAttendance(999)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetByInstructor(t *testing.T) {

	// Get the attendance for an instructor
	_, err := new(Attendance).GetByInstructor(3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetCoursesByInstructor(t *testing.T) {

	// Get the courses for an instructor
	_, err := new(CourseSection).GetCoursesByInstructor(3)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetSectionsByInstructorAndCourse(t *testing.T) {

	// Get the sections for an instructor and course
	_, err := new(CourseSection).GetSectionsByInstructorAndCourse(3, 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetStudentsByAttendance(t *testing.T) {

	// Get the students by attendance
	_, err := new(AttendanceRecord).GetStudentsByAttendance(3)
	if err != nil {
		t.Fatal(err)
	}
}
