package helpers

import (
	"database/sql"
	"errors"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Server) EmailLoginAPI(c *gin.Context) {
	var EmailLoginInput struct {
		EmailID  string `json:"email_id"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&EmailLoginInput); err != nil {
		ReturnAPITechnicalError(c, InternalServerError, err, "could not bind to json")
		return
	}

	log.Debug().Msg("Connecting to database")
	// Open a connection to the database
	db, err := sqlx.Open("mysql", Dsn)
	if err != nil {
		ErrorLogger(err, "Failed to connect to the database")
		return
	}
	defer db.Close()

	// Verify the connection is successful
	err = db.Ping()
	if err != nil {
		ErrorLogger(err, "Error connecting to the database")
		return
	}
	log.Debug().Msg("Connected to database")

	var systemUserName string
	valid := isValidEmail(EmailLoginInput.EmailID)
	if valid {
		query := "SELECT user_name FROM system_user where email_id = ? AND password = ?"
		// Query to check if the ID exists and fetch the name
		err = db.QueryRow(query, EmailLoginInput.EmailID, EmailLoginInput.Password).Scan(&systemUserName)
		if err != nil {
			ReturnAPITechnicalError(c, InternalServerError, err, "error fetching data")
			return
		}
	} else {
		ReturnAPIFunctionalError(c, "Invalid email id format")
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		ReturnAPIFunctionalError(c, "Invalid User")
		return
	} else {
		if systemUserName != "" {
			c.AbortWithStatusJSON(Ok, gin.H{
				"message":            "Login Successful, Welcome " + systemUserName,
				"httpResponseCode":   Ok,
				"businessStatusCode": BusinessSuccess,
			})
			return
		} else {
			c.AbortWithStatusJSON(Ok, gin.H{
				"message":            "Login Not Successful",
				"httpResponseCode":   Ok,
				"businessStatusCode": BusinessSuccess,
			})
			return
		}
	}

}
