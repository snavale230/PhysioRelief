package helpers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib" // Standard library bindings for pgx
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (s *Server) FetchAppointmentFormAudit(c *gin.Context) {
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

	var viewListData []map[string]interface{}

	query := `SELECT * FROM appointment_audit`

	rows, err := db.Queryx(query)
	if err != nil {
		ReturnAPITechnicalError(c, InternalServerError, err, "Error executing the query in appoinment audit table")
		return
	}
	defer func(rows *sqlx.Rows) {
		err2 := rows.Close()
		if err2 != nil {
			ErrorLogger(err2, "error while closing rows")
		}
	}(rows)
	for rows.Next() {
		r := make(map[string]interface{})
		err = rows.MapScan(r)
		if err != nil {
			ReturnAPITechnicalError(c, InternalServerError, err, "Error scanning data in appoinment table")
			return
		}

		convertJSONBColumns(r)
		viewListData = append(viewListData, r)
	}
	if err = rows.Err(); err != nil {
		ReturnAPITechnicalError(c, InternalServerError, err, "Error iterating over rows")
		return
	}

	c.AbortWithStatusJSON(Ok, gin.H{
		"appoinmentAuditList": viewListData,
		"httpResponseCode":    Ok,
		"businessStatusCode":  BusinessSuccess,
	})

}

// convertJSONBColumns converts JSONB byte slices to string maps
func convertJSONBColumns(m map[string]interface{}) {
	InfoLogger("The task to convertJSONBColumns has started")
	for k, v := range m {
		if b, ok := v.([]byte); ok {
			// Check if the field represents a JSONB column
			if isJSON(b) {
				var jsonMap map[string]interface{}
				if err := json.Unmarshal(b, &jsonMap); err == nil {
					m[k] = jsonMap
				} else {
					// Handle unmarshaling error
					m[k] = string(b) // Fallback to string representation if JSON parsing fails
				}
			} else {
				m[k] = string(b) // Convert byte slice to string
			}
		}
	}
	InfoLogger("The task to convertJSONBColumns has Completed")
}

// isJSON checks if a byte slice is valid JSON
func isJSON(b []byte) bool {
	// InfoLogger("The task to isJSON is Executed")
	var js json.RawMessage
	return json.Unmarshal(b, &js) == nil
}
