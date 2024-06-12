package helpers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
)

func (s *Server) BookSessionAPI(c *gin.Context) {
	var BookSessionInput struct {
		EmailID  string `json:"email_id"`
		Name     string `json:"name"`
		MobileNo string `json:"mobile"`
		City     string `json:"city"`
		Service  string `json:"service"`
	}

	if err := c.ShouldBindJSON(&BookSessionInput); err != nil {
		ReturnAPITechnicalError(c, InternalServerError, err, "could not bind to json")
		return
	}

	valid := isValidEmail(BookSessionInput.EmailID)
	if valid {
		err := sessionEmailSending(BookSessionInput.EmailID, BookSessionInput.MobileNo, BookSessionInput.Name, BookSessionInput.City, BookSessionInput.Service)
		if err != nil {
			ReturnAPITechnicalError(c, InternalServerError, err, "email sending error")
			return
		}
		err = sessionAudit(BookSessionInput.Name, BookSessionInput.EmailID, BookSessionInput.City, BookSessionInput.Service)
		if err != nil {
			ReturnAPITechnicalError(c, InternalServerError, err, "error inserting data into session audit table data")
			return
		}
		c.AbortWithStatusJSON(Ok, gin.H{
			"message":            "Session Booked Successfully!",
			"httpResponseCode":   Ok,
			"businessStatusCode": BusinessSuccess,
		})
		return
	} else {
		ReturnAPIFunctionalError(c, "Invalid email id format")
		return
	}

}
