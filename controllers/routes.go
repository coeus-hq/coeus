package controllers

import (
	"github.com/gin-gonic/gin"
)

func PublicRoutes(g *gin.RouterGroup) {
	g.Use(CheckDatabaseType)

	g.GET("/sign-in", SignInGetHandler)
	g.POST("/sign-in", SignInPostHandler)
	g.GET("/create-account", CreateAccountGetHandler)
	g.POST("/create-account", CreateAccountPostHandler)
	g.GET("/forgot-password", ForgotPasswordGetHandler)

	// If the org is not set up (onboardingComplete), redirect to onboarding page
	g.GET("/onboarding", OnboardingGetHandler)
}

func PrivateRoutes(g *gin.RouterGroup) {

	g.GET("/logout", LogoutGetHandler)

	// STUDENT AND STAFF ROUTES
	studentAndStaffRoutes := g.Group("/")
	studentAndStaffRoutes.Use(AdminForbidden())
	{
		studentAndStaffRoutes.GET("/", MyCoursesGetHandler)
		studentAndStaffRoutes.DELETE("/:sectionID", MyCoursesDeleteHandler)
		studentAndStaffRoutes.GET("/session", SessionGetHandler)
		studentAndStaffRoutes.GET("/settings", SettingsGetHandler)
		studentAndStaffRoutes.POST("/settings", SettingsPostHandler)
		studentAndStaffRoutes.GET("/course-section/:courseID", SectionGetHandler)
		studentAndStaffRoutes.POST("/course-section", SectionPostHandler)
		studentAndStaffRoutes.GET("/course-search", CourseSearchGetHandler)
		studentAndStaffRoutes.POST("/course-search", CourseSearchPostHandler)
		studentAndStaffRoutes.GET("/class-session/:sectionID/:classSessionID", ClassSessionGetHandler)
	}

	// ADMIN ROUTES
	adminRoutes := g.Group("/admin")
	adminRoutes.Use(AdminRequired)
	{
		adminRoutes.GET("/settings", AdminSettingsGetHandler)
		adminRoutes.GET("", AdminUsersGetHandler)
	}

	// INSTRUCTOR ROUTES
	instructorRoutes := g.Group("/instructor")
	instructorRoutes.Use(InstructorRequired)
	{
		instructorRoutes.GET("", InstructorCoursesGetHandler)
		instructorRoutes.GET("/attendance", InstructorAttendanceGetHandler)
		instructorRoutes.GET("/attendance/:classSessionID", InstructorAttendanceGetHandler)
	}

}
