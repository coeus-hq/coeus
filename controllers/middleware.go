package controllers

import (
	"coeus/models"
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"net/http"
)

func CheckDatabaseType(c *gin.Context) {
	session := sessions.Default(c)

	isDemo, err := new(models.Organization).GetStatus()
	if err != nil {
		fmt.Println(err)
		return
	}

	session.Set("isDemo", isDemo)
	session.Save()

	c.Next()
}

// If the organization is not set, redirect to the organization page
func OrganizationRequired(c *gin.Context) {
	// check if organization exists OrganizationExists
	organizationExists := new(models.Organization).OrganizationExists()

	if !organizationExists {
		// Check if the current route is not the onboarding page
		if c.Request.URL.Path != "/onboarding" {
			c.Abort()
			c.Redirect(http.StatusSeeOther, "/onboarding")
			return
		}
	}
	c.Next()
}

func PreventOnboardingIfOrganizationExists(c *gin.Context) {
	// Check if organization exists
	organizationExists := new(models.Organization).OrganizationExists()

	// If organization exists and the user tries to access onboarding page, redirect them
	if organizationExists && c.Request.URL.Path == "/onboarding" {
		c.Redirect(http.StatusSeeOther, "/")
		c.Abort()
		return
	}

	c.Next()
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("userID")

	if user == nil {
		c.Abort()
		c.Redirect(http.StatusMovedPermanently, "/sign-in")
		return
	}
	c.Next()
}

// Prevents non Admins from accessing the route
func AdminRequired(c *gin.Context) {
	session := sessions.Default(c)
	isAdmin := session.Get("isAdmin")

	if isAdmin == nil || !isAdmin.(bool) {
		c.Abort()
		c.Redirect(http.StatusMovedPermanently, "/sign-in")
		return
	}
	c.Next()
}

// Prevents Admins from accessing the route
func AdminForbidden() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		isAdmin := session.Get("isAdmin")

		session.Set("managementTemplates", false)
		session.Save()

		if isAdmin != nil && isAdmin.(bool) {
			c.HTML(http.StatusNotFound, "404.html", gin.H{})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Prevents non Instructors from accessing the route
func InstructorRequired(c *gin.Context) {
	session := sessions.Default(c)
	isInstructor := session.Get("isInstructor")

	if isInstructor == nil || !isInstructor.(bool) {
		c.Abort()
		c.Redirect(http.StatusMovedPermanently, "/sign-in")
		return
	} else {
		session.Set("managementTemplates", true)
		session.Save()
	}

	c.Next()
}
