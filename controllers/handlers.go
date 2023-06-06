package controllers

import (
	"coeus/models"
	"fmt"
	"time"

	"log"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"

	helpers "coeus/helpers"
)

// Catch all / 404 page
func NotFoundHandler(c *gin.Context) {
	RenderTemplate(c, http.StatusNotFound, "404.html", gin.H{})
}

func OnboardingGetHandler(c *gin.Context) {

	// Delete the session cookie when the user is on the onboarding page
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	RenderTemplate(c, http.StatusOK, "onboarding.html", gin.H{
		"onboarding": true,
	})
}

func SignInGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	session.Save()

	if userID != nil {
		c.Redirect(http.StatusSeeOther, "/")
	} else {

		// CheckAPIKey function to check if the API key exists
		canResetPassword, _ := new(models.Organization).CheckAPIKey()

		RenderTemplate(c, http.StatusOK, "sign-in.html", gin.H{
			"canResetPassword": canResetPassword,
		})
	}
}

func SignInPostHandler(c *gin.Context) {
	// Clear any existing session
	session := sessions.Default(c)
	session.Clear()
	err := session.Save()
	if err != nil {
		log.Printf("Failed to clear session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}

	// Delete the session cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	userID := session.Get("userID")

	// Check if the database is demo
	isDemo, err := new(models.Organization).GetStatus()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Set the isDemo to the session
	session.Set("isDemo", isDemo)

	if userID != nil {

		RenderTemplate(c, http.StatusOK, "sign-in.html", gin.H{
			"content": "Please logout first",
		})
		return
	}

	// Get the admin id
	adminID, err := new(models.Organization).GetAdminID()
	if err != nil {
		fmt.Println(err)
	}

	username := c.PostForm("username")
	password := c.PostForm("password")

	if helpers.EmptyUserPass(username, password) {
		RenderTemplate(c, http.StatusOK, "sign-in.html", gin.H{
			"content": "Parameters can't be empty",
		})
	}

	u := new(models.User)
	id := u.Authenticate(username, password)

	if id > 0 {
		session.Set("userID", id)

		// Get the organization id
		organizationID, err := new(models.Organization).GetOrganizationID()
		if err != nil {
			fmt.Println(err)
		}

		// Get the organization from the database
		organization, err := new(models.Organization).Get(organizationID)
		if err != nil {
			fmt.Println(err)
		}

		// set the organization logo and name to the session
		session.Set("organizationLogo", organization.LogoPath)
		session.Set("organizationName", organization.Name)

		var s models.Setting
		setting := new(models.Setting)
		s, _ = setting.Get(id)

		darkMode := s.DarkTheme
		if darkMode {
			session.Set("cssStyle", "style-dark")
		}

		// Set the timezone to the session
		session.Set("timezone", s.TimezoneOffset)

		if err := session.Save(); err != nil {
			log.Println("Failed to save session:", err)
		}

		// The following code is used to redirect the user to the correct page
		// after sign-in depending on their role

		// Check if the user is an instructor
		IsInstructor, err := new(models.Moderator).IsInstructor(id)
		if err != nil {
			fmt.Println(err)
		}

		// If the user is an instructor, redirect to the instructor dashboard.
		if IsInstructor {
			session.Set("isInstructor", true)
			session.Save()

			c.Redirect(http.StatusSeeOther, "/")
		} else if id == adminID {
			// If the user is a superadmin, redirect to the admin dashboard.
			session.Set("isAdmin", true)
			session.Save()

			c.Redirect(http.StatusSeeOther, "/admin")
		} else {
			// If user is not an admin, and redirect to the my-courses page
			c.Redirect(http.StatusSeeOther, "/")
		}

	} else {

		isDemo, err := new(models.Organization).GetStatus()
		if err != nil {
			fmt.Println(err)
			return
		}

		var content string

		if isDemo == "true" {
			content = "Incorrect username or password. For soft launch: Database may not be seeded yet (scroll down)."
		} else {
			content = "Incorrect username or password."
		}

		RenderTemplate(c, http.StatusOK, "sign-in.html", gin.H{
			"content": content,
		})
	}
}

func LogoutGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		log.Println("Invalid session token")
		return
	}

	// Clear the session
	session.Clear()
	err := session.Save()
	if err != nil {
		log.Printf("Failed to clear session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}

	// Delete the session cookie
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "session",
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
		MaxAge:  -1,
	})

	c.Redirect(http.StatusSeeOther, "/")
}

func CreateAccountPostHandler(c *gin.Context) {
	session := sessions.Default(c)

	// Get form input
	firstName := c.PostForm("firstName")
	lastName := c.PostForm("lastName")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Add user to the database
	userID, err := new(models.User).Add(email, password, lastName, firstName)
	if err != nil {
		log.Println("Failed to add user:", err)

		// Render an error message or redirect to an error page
		c.Redirect(http.StatusSeeOther, "/create-account")
		return
	}

	// parse the userID from an int64 to an int
	userIDInt, err := strconv.Atoi(fmt.Sprintf("%d", userID))
	if err != nil {
		log.Println("Failed to parse userID:", err)
	}

	// Get the organization id
	organizationID, err := new(models.Organization).GetOrganizationID()
	if err != nil {
		fmt.Println(err)
	}

	// Add the user to the default organization
	err = new(models.User).AddUserToOrganization(userIDInt, organizationID)
	if err != nil {
		log.Println("Failed to add user to organization:", err)
	}

	// Get the organization from the database
	organization, err := new(models.Organization).Get(organizationID)
	if err != nil {
		fmt.Println(err)
	}

	// Add a default setting for the user
	_, err = new(models.Setting).Add(userIDInt)

	// Get the users settings
	settings, err := new(models.Setting).Get(userIDInt)
	session.Set("timezone", settings.TimezoneOffset)

	// set the organization logo to the session
	session.Set("organizationLogo", organization.LogoPath)

	// Set the user ID in the session
	session.Set("userID", userIDInt)
	session.Save()

	// send a json response
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func CreateAccountGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")
	session.Save()

	if userID != nil {
		c.Redirect(http.StatusSeeOther, "/")
	} else {
		RenderTemplate(c, http.StatusOK, "create-account.html", gin.H{
			"content": "",
		})
	}
}

func ForgotPasswordGetHandler(c *gin.Context) {
	RenderTemplate(c, http.StatusOK, "forgot-password.html", gin.H{})
}

func SettingsGetHandler(c *gin.Context) {
	var s models.Setting

	session := sessions.Default(c)
	userID := session.Get("userID")

	// Add a check to ensure the type assertion won't panic
	userIDint, ok := userID.(int)
	if !ok {
		// Handle the case where userID is not an int
		fmt.Println("userID is not an int")
		return
	}

	setting := new(models.Setting)
	s, _ = setting.Get(userID.(int))

	var checked string
	if s.DarkTheme {
		checked = "checked"
	} else {
		checked = ""
	}

	// Convert int to int64
	userIDint64 := int64(userIDint)

	user, err := new(models.User).Get(userIDint64)
	if err != nil {
		fmt.Println(err)
	}

	RenderTemplate(c, http.StatusOK, "settings.html", gin.H{
		"content":   "Settings page",
		"checked":   checked,
		"timezone":  s.TimezoneOffset,
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
	})
}

func SettingsPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)
	setting := new(models.Setting)
	var cssStyle string

	// Get form input
	darkModeToggle := c.PostForm("dark-mode-toggle")

	// Get the user's timezone
	timezone := c.PostForm("timezone")

	// Set timezone to session
	session.Set("timezone", timezone)

	// Changing the value of darkModeToggle to a string that will be used in sql
	if darkModeToggle == "on" {
		darkModeToggle = "true"
		cssStyle = "style-dark"
	} else {
		darkModeToggle = "false"
	}

	// Convert timezone to int
	timezoneInt, err := strconv.Atoi(timezone)

	// Update the user's settings
	err = setting.Update(userID, "theme", darkModeToggle, timezoneInt)
	session.Set("cssStyle", cssStyle)
	session.Save()
	if err != nil {
		log.Println("Failed to save session:", err)
		return
	}

	// Redirect to the settings page
	c.Redirect(http.StatusSeeOther, "/settings")
}

func SessionGetHandler(c *gin.Context) {
	RenderTemplate(c, http.StatusOK, "session.html", gin.H{
		"content": "Top questions:",
	})
}

func MyCoursesGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	initials, err := new(models.User).GetUserInitials(userID)
	if err != nil {
		log.Println("Failed to get initials:", err)
		return
	}

	// Set the initials in the session
	session.Set("userInitials", initials)
	session.Save()

	userInitials := session.Get("userInitials")

	// Get the courses from the database
	courses, err := new(models.Course).GetByUserId(userID)
	if err != nil {
		log.Println("Failed to get courses:", err)
		return
	}

	RenderTemplate(c, http.StatusOK, "my-courses.html", gin.H{
		"courses":      courses,
		"user":         userID,
		"userInitials": userInitials,
	})
}

func MyCoursesDeleteHandler(c *gin.Context) {
	sectionID := c.Param("sectionID")

	// convert sectionID to int
	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		log.Println("Failed to convert sectionID to int:", err)
		return
	}

	// Get the courses from the database
	err = new(models.Section).DeleteByID(sectionIDInt)
	if err != nil {
		log.Println("Failed to delete section:", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func CourseSearchGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")

	RenderTemplate(c, http.StatusOK, "course-search.html", gin.H{
		"content": "Course search",
		"user":    userID,
	})
}

func CourseSearchPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	// Get the course identifier from the form
	courseIdentifier := c.PostForm("courseIdentifier")

	// Get the courses from the database
	courses, _ := new(models.Course).Search(courseIdentifier)

	// Pass the course data and user ID to the template
	RenderTemplate(c, http.StatusOK, "course-results.html", gin.H{
		"courses":          courses,
		"courseIdentifier": courseIdentifier,
		"user":             userID,
	})
}

func SectionGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID")

	courseID := c.Param("courseID")
	courseIDInt, err := strconv.Atoi(courseID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the section IDs that the user is enrolled in
	enrolledSectionIDs, err := new(models.Section).GetEnrolledSections(userID.(int))
	if err != nil {
		fmt.Println(err)
	}

	// Get a list of sections for the course
	sections, err := new(models.Section).GetByCourse(courseIDInt)
	if err != nil {
		fmt.Println(err)
	}

	// Mark the sections that the user is enrolled in
	for i, section := range sections {
		for _, sectionID := range enrolledSectionIDs {
			if section.ID == sectionID {
				sections[i].Enrolled = true
			}
		}
	}

	// Get course by ID
	course, err := new(models.Course).GetCourseByID(courseIDInt)
	if err != nil {
		fmt.Println(err)
	}

	RenderTemplate(c, http.StatusOK, "course-sections.html", gin.H{
		"sections": sections,
		"course":   course,
		"user":     userID,
	})
}

func SectionPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	// Get the section id from the form we are getting the section id from the input value
	sectionId := c.PostForm("sectionSelect")

	// Parse the section id to an int
	sectionIdInt, err := strconv.Atoi(sectionId)
	if err != nil {
		fmt.Println(err)
	}

	// Add a moderator status for that student given the section id
	_, err = new(models.Moderator).Add(userID, sectionIdInt, "student")
	if err != nil {
		fmt.Println(err)
	}

	// Delete all enrollments for the user associated with the course
	err = new(models.Section).DeleteBySectionId(sectionIdInt, userID)
	if err != nil {
		fmt.Println(err)
	}

	// Post the new enrollment to the database
	err = new(models.Section).AddEnrollment(sectionIdInt, userID)
	if err != nil {
		fmt.Println(err)
	}

	// Redirect to the my courses page
	c.Redirect(http.StatusSeeOther, "/")
}

func ClassSessionGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	// Reset moderatorType in session as a user can have multiple moderator types
	session.Set("moderatorType", "")

	// Get the section id and session id from the url
	sectionID := c.Param("sectionID")
	classSessionID := c.Param("classSessionID")

	// Set the section id to the session
	session.Set("sectionID", sectionID)
	session.Set("classSessionID", classSessionID)
	session.Save()

	// Parse the section id to an int
	sectionIDInt, err := strconv.Atoi(sectionID)
	if err != nil {
		fmt.Println(err)
	}

	// Parse the class session id to an int
	classSessionIDInt, err := strconv.Atoi(classSessionID)
	if err != nil {
		fmt.Println(err)
	}

	// Add participant to the class session to the participants table though join
	_, err = new(models.ClassSession).Join(classSessionIDInt, userID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the participant count for the class session
	participantCount, err := new(models.ClassSession).GetParticipantCount(classSessionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	// Get the course info from the database by the class session id
	courseInfo, err := new(models.Course).GetBySectionId(sectionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	// Get the users moderator status for the section
	moderatorStatus, err := new(models.Moderator).GetStatus(userID, sectionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	// Get the timezone from the settings table
	settings, err := new(models.Setting).Get(userID)
	if err != nil {
		fmt.Println(err)
	}

	// Get schedule info for the class session
	scheduleInfo, err := new(models.Schedual).GetSchedualBySectionID(sectionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	session.Set("moderatorType", moderatorStatus.Type)
	session.Save()

	RenderTemplate(c, http.StatusOK, "class-session.html", gin.H{
		"user":              userID,
		"moderatorStatus":   moderatorStatus,
		"courseInfo":        courseInfo,
		"scheduleInfo":      scheduleInfo,
		"classSessionIDInt": classSessionIDInt,
		"timezone":          settings.TimezoneOffset,
		"participantCount":  participantCount,
	})

	constructParticipantJoined(classSessionIDInt, participantCount)
}
func VoteUpPostHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	// Get the class section id from the session
	classSessionID := session.Get("classSessionID").(string)

	//parse the class section id to an int
	classSessionIDInt, err := strconv.Atoi(classSessionID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the question id from the form
	questionID := c.Param("questionID")

	// Parse the question id to an int
	questionIDInt, err := strconv.Atoi(questionID)
	if err != nil {
		fmt.Println(err)
	}
	hasVoted, err := new(models.Question).HasVoted(questionIDInt, userID)
	if err != nil {
		fmt.Println(err)
	}
	if hasVoted {
		c.JSON(http.StatusOK, gin.H{
			"status": "already voted",
		})
		return
	}

	// Add the question to the database
	err = new(models.Question).VoteQuestion(questionIDInt, userID)
	if err != nil {
		fmt.Println(err)
	}
	//send a 200 status code to the client if the vote was successful
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	// Get the updated vote count for the question
	updatedVoteCount, err := new(models.Question).GetVoteCount(questionIDInt)
	if err != nil {
		log.Println("Error getting vote count:", err)
		return
	}

	// Broadcast the "vote-up" action to all WebSocket connections
	constructVoteUp(classSessionIDInt, userID, questionIDInt, updatedVoteCount)
}

func MarkQuestionPostHandler(c *gin.Context) {
	session := sessions.Default(c)

	// Get the question id from the form
	questionID := c.Param("questionID")

	// Get the class section id from the session
	classSessionID := session.Get("classSessionID").(string)

	//parse the class section id to an int
	classSessionIDInt, err := strconv.Atoi(classSessionID)
	if err != nil {
		fmt.Println(err)
	}

	// Parse the question id to an int
	questionIDInt, err := strconv.Atoi(questionID)
	if err != nil {
		fmt.Println(err)
	}

	// Add the question to the database
	err = new(models.Question).MarkQuestion(questionIDInt)
	if err != nil {
		fmt.Println(err)
	}

	//send a 200 status code to the client if the vote was successful
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

	constructMarkQuestion(classSessionIDInt, questionIDInt)
}

func AdminSettingsGetHandler(c *gin.Context) {
	RenderTemplate(c, http.StatusOK, "organization-settings.html", gin.H{})
}

func AdminUsersGetHandler(c *gin.Context) {

	session := sessions.Default(c)

	// Get the admin ID
	adminID, err := new(models.Organization).GetAdminID()
	if err != nil {
		fmt.Println(err)
	}

	// Get the users from the database
	users, err := new(models.User).GetAll()
	if err != nil {
		fmt.Println(err)
	}

	//Get user count from the database
	userCount, err := new(models.User).Count()
	if err != nil {
		fmt.Println(err)
	}

	// Get the organization id
	organizationID, err := new(models.Organization).GetOrganizationID()
	if err != nil {
		fmt.Println(err)
	}

	// Get the organization from the database
	organization, err := new(models.Organization).Get(organizationID)
	if err != nil {
		fmt.Println(err)
	}

	// set the admin initials to the session
	adminInitials, err := new(models.User).GetUserInitials(adminID)
	if err != nil {
		fmt.Println(err)
	}

	session.Set("userInitials", adminInitials)
	session.Save()

	RenderTemplate(c, http.StatusOK, "user-management.html", gin.H{
		"users":        users,
		"count":        userCount,
		"organization": organization,
	})

}

func InstructorCoursesGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	//Get the course and section count for the instructor
	courseCount, err := new(models.Course).CountByInstructor(userID)
	if err != nil {
		fmt.Println(err)
	}
	sectionCount, err := new(models.Section).CountByInstructor(userID)
	if err != nil {
		fmt.Println(err)
	}

	// Get the course and section data
	courseSections, err := new(models.Course).GetCourseSectionsByIntructor(userID)
	if err != nil {
		fmt.Println(err)
	}

	RenderTemplate(c, http.StatusOK, "course-management.html", gin.H{
		"courseCount":    courseCount,
		"sectionCount":   sectionCount,
		"courseSections": courseSections,
	})
}

func InstructorAttendanceGetHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("userID").(int)

	// Get the course data for the instructor by their user id
	courseInfo, err := new(models.CourseSection).GetCoursesByInstructor(userID)
	if err != nil {
		fmt.Println(err)
	}

	RenderTemplate(c, http.StatusOK, "attendance.html", gin.H{
		"courseInfo": courseInfo,
	})
}
