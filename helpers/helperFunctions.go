package helpers

import (
	"database/sql"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func ReturnAPITechnicalError(c *gin.Context, statusCode int, err error, msg string) {
	msg = strings.ToLower(msg)
	ErrorLogger(err, msg)
	switch statusCode {
	case InternalServerError:
		msg = "internal server error"
	case BadRequest:
		msg = "bad request"
	case Unauthorized:
		msg = "unauthorized user"
	default:
		// Do nothing
		msg = ""
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"errorCode": statusCode, "message": msg})
}

func ErrorLogger(err error, msg string) {
	msg = strings.ToLower(msg)
	log.Error().Err(err).Msg(msg)
}

func isValidEmail(email string) bool {
	// Define the email regex pattern
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regex
	regex := regexp.MustCompile(pattern)

	// Use the regex to match the email
	return regex.MatchString(email)
}

func InfoLogger(msg string) {
	msg = strings.ToLower(msg)
	log.Info().Msg(msg)
}

type TechnicallySuccessFullResponseButNotFunctionally struct {
	HttpResponseCode      int     `json:"httpResponseCode,omitempty"`
	BusinessStatusCode    int     `json:"businessStatusCode,omitempty"`
	BusinessStatusSubCode float64 `json:"businessStatusSubCode,omitempty"`
	Msg                   string  `json:"message,omitempty"`
}

func NewTechnicallySuccessFullResponseButNotFunctionally(message string, businessStatusSubCode ...float64) *TechnicallySuccessFullResponseButNotFunctionally {
	var businessStatusSubCodeFinal float64
	businessStatusSubCodeFinal = 0.0
	if len(businessStatusSubCode) > 0 {
		businessStatusSubCodeFinal = businessStatusSubCode[0]
	}
	return &TechnicallySuccessFullResponseButNotFunctionally{
		HttpResponseCode:      Ok,
		BusinessStatusCode:    BusinessFailure,
		BusinessStatusSubCode: businessStatusSubCodeFinal,
		Msg:                   message,
	}
}

func ReturnAPIFunctionalError(c *gin.Context, msg string, businessStatusSubCode ...float64) {
	InfoLogger(msg)

	var technicallySuccessFullFunctionallyNot *TechnicallySuccessFullResponseButNotFunctionally
	if len(businessStatusSubCode) > 0 {
		technicallySuccessFullFunctionallyNot = NewTechnicallySuccessFullResponseButNotFunctionally(msg, businessStatusSubCode...)
	} else {
		technicallySuccessFullFunctionallyNot = NewTechnicallySuccessFullResponseButNotFunctionally(msg)
	}
	c.AbortWithStatusJSON(Ok, technicallySuccessFullFunctionallyNot)
}

func sessionAudit(userName string, userEmail string, city string, services string) error {
	log.Debug().Msg("Connecting to database")
	// Open a connection to the database
	db, err := sqlx.Open("mysql", Dsn)
	if err != nil {
		ErrorLogger(err, "Failed to connect to the database")
		return err
	}
	defer db.Close()

	// Verify the connection is successful
	err = db.Ping()
	if err != nil {
		ErrorLogger(err, "Error connecting to the database")
		return err
	}
	log.Debug().Msg("Connected to database")

	_, err = db.Exec(`INSERT INTO session_audit (session_audit_id, user_name, user_email, city, service) VALUES (?, ?, ?, ?, ?)`, uuid.New().String(), userName, userEmail, city, services)
	if err != nil {
		ErrorLogger(err, "Error inserting data into the session audit table")
		return err
	}
	return nil
}

func appoinmentAudit(userName string, userMobile string, userEmail string, service string, userAddress string, appoinmentDate time.Time, description string) error {
	log.Debug().Msg("Connecting to database")
	// Open a connection to the database
	db, err := sql.Open("mysql", Dsn)
	if err != nil {
		ErrorLogger(err, "Failed to connect to the database")
		return err
	}
	defer db.Close()

	// Verify the connection is successful
	err = db.Ping()
	if err != nil {
		ErrorLogger(err, "Error connecting to the database")
		return err
	}
	log.Debug().Msg("Connected to database")
	_, err = db.Exec(`INSERT INTO appointment_audit (appointment_audit_id, user_name, user_mobile, user_email, service, user_address, appointment_date_and_time, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, uuid.New().String(), userName, userMobile, userEmail, service, userAddress, appoinmentDate, description)
	if err != nil {
		ErrorLogger(err, "Error inserting data into the appoinment audit table")
		return err
	}
	return nil

}
