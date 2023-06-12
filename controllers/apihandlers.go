package controllers

import (
	"coeus/models"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
)

func APIDarkThemePostHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	theme, err := new(models.Setting).ToggleDarkTheme(userID)
	if err != nil {
		fmt.Println(err)
	}

	if theme {
		session.Set("cssStyle", "style-dark")
		session.Save()
	} else {
		session.Set("cssStyle", "")
		session.Save()
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

type timezoneRequest struct {
	Timezone string `json:"timezone"`
}

func APITimezonePostHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	var req timezoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// convert the timezone to an int
	timezoneInt, err := strconv.Atoi(req.Timezone)
	if err != nil {
		fmt.Println(err)
	}

	err = new(models.Setting).UpdateTimezone(userID, timezoneInt)
	if err != nil {
		fmt.Println(err)
	}

	session.Set("timezone", req.Timezone)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func APIQuestionsPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	// Get the question from the form
	questionText := c.PostForm("questionText")

	// Get the class section id from the url
	classSessionID := c.Param("classSessionID")

	// convert classSessionIDInt to int
	classSessionIDInt, err := strconv.Atoi(classSessionID)

	// Add the question to the database
	questionID, err := new(models.Question).PostQuestion(userID, classSessionIDInt, questionText)
	if err != nil {
		fmt.Println(err)
	}

	// The frontend is going to handle the timezone formatting so we can just set it to 0 or UTC
	timezone := 0

	// GetByID the question from the database
	question, err := new(models.Question).GetByID(questionID, timezone)
	if err != nil {
		fmt.Println(err)
	}

	broadcastNewQuestion(classSessionIDInt, question)
}

func APIMessagesGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)
	classSessionID := c.Param("classSessionID")
	moderatorType := session.Get("moderatorType")

	timezone, _ := session.Get("timezone").(string)

	// Parse the class sessionID and timezone to an int
	classSessionIDInt, err := strconv.Atoi(classSessionID)
	timezoneInt, err := strconv.Atoi(timezone)

	// Get the questions from the database by the created at time
	questionsByTimeSlice, err := new(models.Question).GetAllQuestions(classSessionIDInt, "created_at", timezoneInt)
	if err != nil {
		fmt.Println(err)
	}

	// Get the questions from the database by the created at time
	questionsByVoteSlice, err := new(models.Question).GetAllQuestions(classSessionIDInt, "votes", timezoneInt)
	if err != nil {
		fmt.Println(err)
	}

	// Get the userHasVoted slice from the database
	userHasVotedSlice, err := new(models.Question).HasVotedAll(classSessionIDInt, userID)

	questionsByVoteUnanswered, err := new(models.Question).GetUnansweredQuestions(classSessionIDInt, timezoneInt, "votes")
	questionsByTimeUnanswered, err := new(models.Question).GetUnansweredQuestions(classSessionIDInt, timezoneInt, "time")

	// Append the userHasVoted bool to the questionsBy slices
	questionsByTime := new(models.Question).HasVotedAppend(questionsByTimeSlice, userHasVotedSlice)
	questionsByVote := new(models.Question).HasVotedAppend(questionsByVoteSlice, userHasVotedSlice)

	c.JSON(http.StatusOK, gin.H{
		"moderatorType":             moderatorType,
		"questionsByTime":           questionsByTime,
		"questionsByVote":           questionsByVote,
		"questionsByVoteUnanswered": questionsByVoteUnanswered,
		"questionsByTimeUnanswered": questionsByTimeUnanswered,
		"user":                      userID,
	})
}

func APIModeratorsForSectionGetHandler(c *gin.Context) {
	sectionID := c.Param("sectionID")
	sectionID = strings.Trim(sectionID, "\"")
	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		fmt.Println(err)
	}

	moderators, err := new(models.Moderator).GetAllModerators(sectionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"moderators": moderators,
	})
}

func APIAddModeratorPostHandler(c *gin.Context) {
	email := c.Param("email")
	SectionID := c.Param("sectionID")

	// Get user by email to check if they exist
	_, err := new(models.User).GetUserId(email)
	if err != nil {
		fmt.Println(err)

		// Send a 400 status code to the client if the user doesn't exist
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "user doesn't exist",
		})

		return
	}

	// Parse the section id to an int
	SectionIDInt, err := strconv.Atoi(SectionID)
	if err != nil {
		fmt.Println(err)
	}

	moderatorUserID, err := new(models.User).GetUserId(email)
	if err != nil {
		fmt.Println(err)
	}

	// Add the new moderator to the database
	_, err = new(models.Moderator).Update(moderatorUserID, SectionIDInt, "teacher assistant")
	if err != nil {
		fmt.Println(err)
	}

	//send a 200 status code to the client if the adding was successful
	c.JSON(http.StatusOK, gin.H{
		"status": "updated moderator to teacher assistant",
	})
}

func APIRemoveModeratorDeleteHandler(c *gin.Context) {
	userID := c.Param("userID")
	sectionID := c.Param("sectionID")

	// Parse the section and user id to an int
	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		fmt.Println(err)
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println(err)
	}

	// Remove the moderator from the database
	err = new(models.Moderator).Delete(userIDInt, sectionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	//send a 200 status code to the client if the adding was successful
	c.JSON(http.StatusOK, gin.H{
		"status": "removed moderator",
	})

}

func APIStartSessionPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	instructorIDInt := session.Get("userID").(int)
	sectionID := c.Param("sectionID")

	// convert sectionID to int
	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		fmt.Println(err)
	}

	// Start the class session in the database
	classSessionID, err := new(models.ClassSession).Start(sectionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	// Create an attendance table for the class session
	attendanceID, err := new(models.Attendance).Add(classSessionID, sectionIDInt, instructorIDInt)
	if err != nil {
		fmt.Println(err)
	}

	constructStartSession(sectionIDInt, attendanceID)
}

func APIEndSessionPostHandler(c *gin.Context) {
	classSessionID := c.Param("classSessionID")

	// convert classSessionIDInt to int
	classSessionIDInt, err := strconv.Atoi(classSessionID)

	// Get the section id from the database by the class session id
	sectionID, err := new(models.ClassSession).GetSectionID(classSessionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	// Get the attendance data for the instructor by their user id
	attendanceID, err := new(models.Attendance).GetAttendanceIDBySectionID(sectionID)
	if err != nil {
		fmt.Println(err)
	}

	// Gets the list of users who are enrolled in the section
	enrolledUsers, err := new(models.Enrollment).GetEnrolledUsersBySectionID(sectionID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the list of users in the participants table
	participantUsers, err := new(models.Participant).GetParticipants(classSessionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	// Find users who are enrolled but not participating
	nonParticipantUserIDs := findAbsentUsers(enrolledUsers, participantUsers)

	// Create new user_attendance entries for those users
	for _, userID := range nonParticipantUserIDs {
		_, err = new(models.Attendance).AddUserAttendance(attendanceID, userID, "absent")
		if err != nil {
			fmt.Println(err)
		}
	}

	// End the class session in the database
	err = new(models.ClassSession).End(classSessionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	constructEndSession(classSessionIDInt, sectionID)
}

// findAbsentUsers returns a slice of user ids who are enrolled but not participating it is used in APIEndSessionPostHandler.
func findAbsentUsers(enrolledUsers []models.Enrollment, participantUsers []models.Participant) []int {
	nonParticipantUserIDs := make([]int, 0)

	for _, enrolledUser := range enrolledUsers {
		found := false
		for _, participantUser := range participantUsers {
			if enrolledUser.UserId == participantUser.UserID {
				found = true
				break
			}
		}
		if !found {
			nonParticipantUserIDs = append(nonParticipantUserIDs, enrolledUser.UserId)
		}
	}

	return nonParticipantUserIDs
}

func APIGetUserGetHandler(c *gin.Context) {

	// Get the users from the database
	users, err := new(models.User).GetAll()
	if err != nil {
		fmt.Println(err)
	}

	isDemo, err := new(models.Organization).GetStatus()
	if err != nil {
		fmt.Println(err)
		return
	}

	// If the database is for softlaunch, then remove the admin user from the slice of users structs
	// this is done so that the admin user cannot be seen in the users table in the softlaunch database
	// preventing someone from changing the admin credentials
	if isDemo == "true" {
		filteredUsers := make([]struct {
			ID                   int64
			Email                string
			LastName             string
			FirstName            string
			CreatedAt            string
			UpdatedAt            string
			HighestModeratorType string
		}, 0)
		for _, user := range users {
			if user.Email != "admin@coeus.education" {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func APIAddUserPostHandler(c *gin.Context) {

	// Struct to hold the JSON data
	type UserData struct {
		Email         string `json:"email"`
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		Password      string `json:"password"`
		ModeratorType string `json:"moderatorType"`
	}

	// Parse the JSON data from the request body
	var userData UserData
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the added user info from the struct
	email := userData.Email
	firstName := userData.FirstName
	lastName := userData.LastName
	password := userData.Password
	moderatorType := userData.ModeratorType

	// Add the user to the database
	newUserID, err := new(models.User).Add(email, password, lastName, firstName)
	if err != nil {
		fmt.Println(err)
	}

	// Parse the newUserID from an int64 to an int
	newUserIDInt := int(newUserID)

	// Add the user mod status to the database with the section id as NULL
	_, err = new(models.Moderator).AdminAdd(newUserIDInt, moderatorType)
	if err != nil {
		fmt.Println(err)
	}

	// Get the organization id
	organizationID, err := new(models.Organization).GetOrganizationID()
	if err != nil {
		fmt.Println(err)
	}

	// Add the user to the default organization
	err = new(models.User).AddUserToOrganization(newUserIDInt, organizationID)
	if err != nil {
		log.Println("Failed to add user to organization:", err)
	}

	// Add a setting for the user
	_, err = new(models.Setting).Add(newUserIDInt)
	if err != nil {
		fmt.Println(err)
	}
}

func APIUpdateUserPutHandler(c *gin.Context) {

	// Struct to hold the JSON data
	type UserData struct {
		UserId        string `json:"userId"`
		Email         string `json:"email"`
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		Password      string `json:"password"`
		ModeratorType string `json:"moderatorType"`
	}

	// Parse the JSON data from the request body
	var userData UserData
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the edited user info from the struct
	userID := userData.UserId
	email := userData.Email
	firstName := userData.FirstName
	lastName := userData.LastName
	password := userData.Password
	moderatorType := userData.ModeratorType

	// Parse the userID to an int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println(err)
	}

	// If the moderatorType is not empty, then update the user mod status (this is mainly for admins)
	if moderatorType != "" {
		// Update the user mod status to the database with the section id as NULL
		_, err = new(models.Moderator).AdminUpdate(userIDInt, moderatorType)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Parse the userID to an int64
	userIDInt64, err := strconv.ParseInt(userID, 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	// Update the user in the database
	err = new(models.User).Update(userIDInt64, email, lastName, firstName, password)
	if err != nil {
		fmt.Println(err)
	}
}

func APIDeleteUserDeleteHandler(c *gin.Context) {
	// Get the user from the form
	userID := c.Param("ID")

	// Parse the userID to an int64
	userIDInt64, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		fmt.Println(err)
	}

	// Delete the user to the database
	err = new(models.User).Delete(userIDInt64)
	if err != nil {
		fmt.Println(err)
	}

	// Parse the userID to an int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Println(err)
	}

	// Delete the setting for the user
	err = new(models.Setting).Delete(userIDInt)
	if err != nil {
		fmt.Println(err)
	}
}

func APIAddAdminPostHandler(c *gin.Context) {

	type Admin struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	admin := &Admin{}

	// Bind the JSON to the struct
	if err := c.ShouldBindJSON(admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Add the admin to the database
	adminID, err := new(models.User).Add(admin.Email, admin.Password, admin.LastName, admin.FirstName)
	if err != nil {
		// Log the error and respond with an error status and message
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to add admin."})
		return
	}

	// convert adminID to int
	adminIDInt := int(adminID)

	// Adds the admin to a flag table to indicate that they are an admin during sign-in
	err = new(models.Organization).SetAdmin(adminIDInt)
	if err != nil {
		fmt.Println(err)
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"success": true, "adminName": admin.FirstName})
}

func APIOnboardingPostHandler(c *gin.Context) {
	session := sessions.Default(c)

	type Organization struct {
		ID                   int    `json:"id"`
		Name                 string `json:"name"`
		OrganizationTimezone string `json:"timezone"`
		LogoPath             string `json:"logoPath"`
		APIKey               string `json:"apiKey"`
		Email                string `json:"email"`
		CreatedAt            string `json:"created_at"`
		UpdatedAt            string `json:"updated_at"`
	}

	// Create a new Organization
	var org Organization

	// Parse multipart form data
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the values from the form
	org.Name = c.PostForm("name")
	org.OrganizationTimezone = c.PostForm("timezone")
	org.APIKey = c.PostForm("apiKey")
	org.Email = c.PostForm("email")

	// If the API key exists, write it to the .env file
	if org.APIKey != "" {

		// Set usingSendGrid to true in the session
		session.Set("usingSendGrid", true)

		// Write the API key and email to the .env file
		err = writeAPIKeyToEnv(org.APIKey, org.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		session.Set("usingSendGrid", false)
	}

	// Process the file
	file, err := c.FormFile("logoPath")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the file extension
	fileExtension := filepath.Ext(file.Filename)

	// Create the file name
	fileName := "logo" + fileExtension

	// Remove existing logo files if they exist
	existingFiles, err := filepath.Glob("./views/static/logo/logo.*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for _, f := range existingFiles {
		err := os.Remove(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Save the file to the server
	err = c.SaveUploadedFile(file, "./views/static/logo/"+fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// The file path is set as the logoPath of the organization
	org.LogoPath = "/static/logo/" + fileName

	// Add the organization to the database
	organizationID, err := new(models.Organization).Add(org.Name, org.OrganizationTimezone, org.LogoPath, org.APIKey, org.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success response
	c.JSON(http.StatusOK, gin.H{"success": true, "organizationID": organizationID})
}

func APIOrganizationPostHandler(c *gin.Context) {
	session := sessions.Default(c)

	file, err := c.FormFile("upload-logo-input")
	if err == nil {

		// Get the file extension
		fileExtension := filepath.Ext(file.Filename)

		// Create the file name
		fileName := "logo" + fileExtension

		// Remove existing logo files if they exist
		existingFiles, err := filepath.Glob("./views/static/logo/logo.*")
		if err != nil {
			fmt.Println(err)
		}
		for _, f := range existingFiles {
			err := os.Remove(f)
			if err != nil {
				fmt.Println(err)
			}
		}

		// Save the file to the server
		err = c.SaveUploadedFile(file, "./views/static/logo/"+fileName)
		if err != nil {
			fmt.Println(err)
		}

		// Get the organization id
		organizationID, err := new(models.Organization).GetOrganizationID()
		if err != nil {
			fmt.Println(err)
		}

		// Append the logo path of /static/logo/ to the file name
		fileName = "/static/logo/" + fileName

		// Update the logo in the database
		err = new(models.Organization).StoreLogo(organizationID, fileName)
		if err != nil {
			fmt.Println(err)
		}

		// Get the organization from the database
		organization, err := new(models.Organization).Get(organizationID)
		if err != nil {
			fmt.Println(err)
		}

		// set the organization logo to the session
		session.Set("organizationLogo", organization.LogoPath)
	}

	if c.PostForm("org-settings-sendgrid-api-key") != "" {

		apiKey := c.PostForm("org-settings-sendgrid-api-key")
		email := c.PostForm("org-settings-sendgrid-email")

		// Write the API key and email to the .env file
		err = writeAPIKeyToEnv(apiKey, email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

	}

	session.Save()
}

// HELPER FUNCTION
// writeAPIKeyToEnv writes the API key and organization email to the .env file
func writeAPIKeyToEnv(apiKey string, email string) error {
	// Read the existing .env file
	data, err := ioutil.ReadFile("globals/.env")
	if err != nil {
		return err
	}

	// Convert to string
	envData := string(data)

	// Check if the SENDGRID_API_KEY already exists
	if strings.Contains(envData, "SENDGRID_API_KEY") {
		// If it does, replace the existing key with the new one
		re := regexp.MustCompile(`(?m)SENDGRID_API_KEY='.*'`)
		envData = re.ReplaceAllString(envData, "SENDGRID_API_KEY='"+apiKey+"'")
	} else {
		// If it doesn't, add the new API key
		envData += "\nSENDGRID_API_KEY='" + apiKey + "'"
	}

	// Check if the SENDGRID_ORGANIZATION_EMAIL already exists
	if strings.Contains(envData, "SENDGRID_ORGANIZATION_EMAIL") {
		// If it does, replace the existing email with the new one
		re := regexp.MustCompile(`(?m)SENDGRID_ORGANIZATION_EMAIL='.*'`)
		envData = re.ReplaceAllString(envData, "SENDGRID_ORGANIZATION_EMAIL='"+email+"'")
	} else {
		// If it doesn't, add the new email
		envData += "\nSENDGRID_ORGANIZATION_EMAIL='" + email + "'"
	}

	// Write the new data back to the .env file
	err = ioutil.WriteFile("globals/.env", []byte(envData), 0644)
	if err != nil {
		return err
	}

	return nil
}

func APICoursesGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	// Get the course and section data
	courseSections, err := new(models.Course).GetCourseSectionsByIntructor(userID)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"courseSections": courseSections,
	})
}

type Course struct {
	CourseNumber    string   `json:"courseNumber"`
	CourseTitle     string   `json:"courseTitle"`
	CourseSemester  string   `json:"courseSemester"`
	CourseStartDate string   `json:"courseStartDate"`
	CourseEndDate   string   `json:"courseEndDate"`
	CourseYear      string   `json:"courseYear"`
	CourseSections  string   `json:"courseSections"`
	Schedules       []string `json:"schedules"`
}

func APICoursesPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	var course Course

	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the course sections to an int
	courseSectionsInt, err := strconv.Atoi(course.CourseSections)
	if err != nil {
		fmt.Println(err)
	}

	// Add the course to the database
	courseID, err := new(models.Course).AddCourseAndSections(course.CourseNumber, course.CourseTitle, course.CourseSemester, course.CourseStartDate, course.CourseEndDate, course.CourseYear, courseSectionsInt)
	if err != nil {
		fmt.Println(err)
	}

	sectionIDs, err := new(models.Course).GetSectionIds(courseID)
	if err != nil {
		fmt.Println(err)
	}

	// For each section id, add the instructor to an enrollment table
	// and create a schedule for the section
	for i, sectionID := range sectionIDs {
		err = new(models.Section).AddEnrollment(sectionID, userID)
		if err != nil {
			fmt.Println(err)
		}

		// Add the schedule to the database
		scheduleID, err := new(models.Section).CreateSchedule(sectionID, course.Schedules[i])
		if err != nil {
			fmt.Println(err)
		}

		// Add class sessions to the database
		err = new(models.Section).AddClassSession(sectionID, scheduleID)
		if err != nil {
			fmt.Println(err)
		}

		// Add a moderator status for that instructor given the section id
		_, err = new(models.Moderator).Add(userID, sectionID, "instructor")
		if err != nil {
			fmt.Println(err)
		}
	}
}

func APICourseSectionDeleteHandler(c *gin.Context) {

	// Get the course id and section number from the form
	courseID := c.Param("courseID")
	sectionNumber := c.Param("sectionNumber")

	// Parse the course id and section number to an int
	courseIDInt, err := strconv.Atoi(courseID)
	sectionNumberInt, err := strconv.Atoi(sectionNumber)

	// Delete the section from the database
	err = new(models.Section).DeleteByCourseIdAndSection(courseIDInt, sectionNumberInt)
	if err != nil {
		fmt.Println(err)
	}
}

func APICourseSectionPutHandler(c *gin.Context) {
	// Define a struct to hold the JSON data
	type CourseData struct {
		CourseID        int    `json:"courseID"`
		SectionID       int    `json:"sectionID"`
		CourseNumber    string `json:"courseNumber"`
		CourseTitle     string `json:"courseTitle"`
		Semester        string `json:"semester"`
		Year            string `json:"year"`
		SectionName     string `json:"sectionName"`
		CourseStartDate string `json:"courseStartDate"`
		CourseEndDate   string `json:"courseEndDate"`
		ScheduleDays    string `json:"scheduleDays"`
		ScheduleTime    string `json:"scheduleTime"`
	}

	// Parse the JSON request body
	var courseData CourseData
	if err := c.ShouldBindJSON(&courseData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the course and section
	err := new(models.Course).UpdateCourseAndSection(courseData.CourseID, courseData.SectionID, courseData.CourseNumber, courseData.CourseTitle, courseData.Semester, courseData.Year, courseData.SectionName, courseData.CourseStartDate, courseData.CourseEndDate, courseData.ScheduleDays, courseData.ScheduleTime)
	if err != nil {
		fmt.Println(err)
	}
}

type EmailRequest struct {
	Email string `json:"email"`
}

func APIPasswordResetSendEmailPostHandler(c *gin.Context) {
	var emailRequest EmailRequest

	if err := c.ShouldBindJSON(&emailRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := emailRequest.Email

	userID, err := new(models.User).GetUserId(email)

	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Email does not exist",
		})
		return
	}

	// remove any existing tokens
	err = new(models.VerifyUser).DeleteToken(userID)
	if err != nil {
		log.Println("Failed to delete token:", err)
	}

	_, err = new(models.VerifyUser).CreateToken(email)
	if err != nil {
		log.Println("Failed to send email:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"userID": userID,
	})
}

type PinRequest struct {
	Pin string `json:"pin"`
}

func APIPasswordResetVerifyPinPostHandler(c *gin.Context) {
	var pinRequest PinRequest

	if err := c.ShouldBindJSON(&pinRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pin := pinRequest.Pin

	_, userID, err := new(models.VerifyUser).MatchToken(pin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Pin does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"userID": userID,
	})
}

type PasswordRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func APIPasswordResetPostHandler(c *gin.Context) {

	var passwordRequest PasswordRequest

	if err := c.ShouldBindJSON(&passwordRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := passwordRequest.Email
	password := passwordRequest.Password

	// Get the user id from the email
	userID, err := new(models.User).GetUserId(email)
	if err != nil {
		fmt.Println(err)
	}

	// Update the user's password
	err = new(models.User).UpdatePassword(userID, password)
	if err != nil {
		fmt.Println(err)
	}

	// Delete the token
	err = new(models.VerifyUser).DeleteToken(userID)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})

}

func APIAttendanceSectionsGetHandler(c *gin.Context) {
	// Get the user id from the session
	instructorID := sessions.Default(c).Get("userID").(int)

	// Get the course id from the url
	courseID := c.Param("courseID")

	// Parse the course id to an int
	courseIDInt, err := strconv.Atoi(courseID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the attendance data for the instructor by their user id
	sections, err := new(models.Attendance).GetSectionsByInstructorAndCourse(instructorID, courseIDInt)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"sections": sections,
	})
}

func APIGetAttendanceIDBySectionIDHandler(c *gin.Context) {

	// Get the section id from the url
	sectionID := c.Param("sectionID")

	// Parse the section id to an int
	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the attendance data for the instructor by their user id
	attendanceID, err := new(models.Attendance).GetAttendanceIDBySectionID(sectionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"attendanceID": attendanceID,
	})
}

func APIAttendanceStudentsGetHandler(c *gin.Context) {

	// Get the attendance id from the url
	attendanceID := c.Param("attendanceID")

	// Parse the attendance id to an int
	attendanceIDInt, err := strconv.Atoi(attendanceID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the attendance data for the instructor by their user id
	students, err := new(models.AttendanceRecord).GetStudentsByAttendance(attendanceIDInt)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"students": students,
	})
}

func APIMarkPresentPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	studentID := session.Get("userID").(int)
	attendanceID := c.Param("attendanceID")

	// Parse the attendance id to an int
	attendanceIDInt, err := strconv.Atoi(attendanceID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the attendance data for the instructor by their user id
	_, err = new(models.Attendance).AddUserAttendance(attendanceIDInt, studentID, "present")
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
