package controllers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RenderTemplate(c *gin.Context, code int, templateName string, data gin.H) {
	session := sessions.Default(c)

	data["isInstructor"] = session.Get("isInstructor")
	data["managementTemplates"] = session.Get("managementTemplates")
	data["isAdmin"] = session.Get("isAdmin")
	data["usingSendGrid"] = session.Get("usingSendGrid")
	data["userInitials"] = session.Get("userInitials")
	data["sectionID"] = session.Get("sectionID")
	data["moderatorType"] = session.Get("moderatorType")
	data["classSessionID"] = session.Get("classSessionID")
	data["user"] = session.Get("userID")
	data["cssStyle"] = session.Get("cssStyle")
	data["timezone"] = session.Get("timezone")
	data["organizationLogo"] = session.Get("organizationLogo")
	data["organizationName"] = session.Get("organizationName")
	data["isDemo"] = session.Get("isDemo")
	c.HTML(http.StatusOK, templateName, data)
}
