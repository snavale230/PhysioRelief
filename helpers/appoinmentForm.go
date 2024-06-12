package helpers

import (
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
)

func (s *Server) AppointmentFormAPI(c *gin.Context) {
	var AppoinmentFormInput struct {
		EmailID        string `json:"email_id"`
		Name           string `json:"name"`
		MobileNo       string `json:"mobile"`
		AppoinmentData string `json:"appoinment_date"`
		AppoinmentTime string `json:"appoinment_time"`
		Address        string `json:"address"`
		Description    string `json:"description"`
		Service        string `json:"service"`
	}

	if err := c.ShouldBindJSON(&AppoinmentFormInput); err != nil {
		ReturnAPITechnicalError(c, InternalServerError, err, "could not bind to json")
		return
	}

	// Parse the date and time strings into time.Time objects
	date, err := time.Parse("2006-01-02", AppoinmentFormInput.AppoinmentData)
	if err != nil {
		ReturnAPITechnicalError(c, InternalServerError, err, "error date parsing")
		return
	}

	timeParsed, err := time.Parse("15:04", AppoinmentFormInput.AppoinmentTime)
	if err != nil {
		ReturnAPITechnicalError(c, InternalServerError, err, "error time parsing")
		return
	}

	// Concatenate the date and time objects
	appoinmentDateAndTime := time.Date(date.Year(), date.Month(), date.Day(), timeParsed.Hour(), timeParsed.Minute(), timeParsed.Second(), 0, time.UTC)

	valid := isValidEmail(AppoinmentFormInput.EmailID)
	if valid {
		err := appoinmentEmailSending(AppoinmentFormInput.EmailID, AppoinmentFormInput.MobileNo, AppoinmentFormInput.Name, appoinmentDateAndTime, AppoinmentFormInput.Address, AppoinmentFormInput.Description, AppoinmentFormInput.Service)
		if err != nil {
			ReturnAPITechnicalError(c, InternalServerError, err, "email sending error")
			return
		}

		err = appoinmentAudit(AppoinmentFormInput.Name, AppoinmentFormInput.MobileNo, AppoinmentFormInput.EmailID, AppoinmentFormInput.Service, AppoinmentFormInput.Address, appoinmentDateAndTime, AppoinmentFormInput.Description)
		if err != nil {
			ReturnAPITechnicalError(c, InternalServerError, err, "error inserting data into appoinment audit table data")
			return
		}
		c.AbortWithStatusJSON(Ok, gin.H{
			"message":            "Appoinment Booked Successfully!",
			"httpResponseCode":   Ok,
			"businessStatusCode": BusinessSuccess,
		})
		return
	} else {
		ReturnAPIFunctionalError(c, "Invalid email id format")
		return
	}

}
