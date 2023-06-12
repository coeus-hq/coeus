package controllers

import (
	"github.com/gin-gonic/gin"
)

func APIRoutes(g *gin.RouterGroup) {
	g.POST("/api/settings/dark-theme", APIDarkThemePostHandler)
	g.PUT("/api/settings/timezone", APITimezonePostHandler)

	g.GET("/api/questions/:classSessionID", APIMessagesGetHandler)
	g.POST("/api/questions/:classSessionID", APIQuestionsPostHandler)
	g.POST("/api/vote-up/:questionID", VoteUpPostHandler)
	g.POST("/api/mark-question/:questionID", MarkQuestionPostHandler)
	g.PUT("/api/add-moderator/:email/:sectionID", APIAddModeratorPostHandler)
	g.DELETE("/api/remove-moderator/:userID/:sectionID", APIRemoveModeratorDeleteHandler)
	g.GET("/api/moderators/:sectionID", APIModeratorsForSectionGetHandler)

	g.POST("/api/start-session/:sectionID", APIStartSessionPostHandler)
	g.POST("/api/end-session/:classSessionID", APIEndSessionPostHandler)

	g.GET("/api/user", APIGetUserGetHandler)
	g.POST("/api/user", APIAddUserPostHandler)
	g.PUT("/api/user", APIUpdateUserPutHandler)
	g.DELETE("/api/user/:ID", APIDeleteUserDeleteHandler)

	g.POST("/api/admin", APIAddAdminPostHandler)
	g.POST("/api/onboarding", APIOnboardingPostHandler)

	g.PUT("/api/organization", APIOrganizationPostHandler)

	g.GET("/api/course", APICoursesGetHandler)
	g.POST("/api/course", APICoursesPostHandler)

	g.DELETE("/api/course/section/:courseID/:sectionNumber", APICourseSectionDeleteHandler)
	g.PUT("/api/course/section", APICourseSectionPutHandler)

	g.POST("/api/password-reset/send-email", APIPasswordResetSendEmailPostHandler)
	g.POST("/api/password-reset/verify-pin", APIPasswordResetVerifyPinPostHandler)
	g.POST("/api/password-reset", APIPasswordResetPostHandler)

	g.GET("/api/attendance/:courseID", APIAttendanceSectionsGetHandler)
	g.GET("/api/attendance/id/:sectionID", APIGetAttendanceIDBySectionIDHandler)
	g.GET("/api/attendance/students/:attendanceID", APIAttendanceStudentsGetHandler)
	g.POST("/api/attendance/mark-present/:attendanceID", APIMarkPresentPostHandler)

}
