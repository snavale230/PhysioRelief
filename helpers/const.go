package helpers

import "net/http"

const (
	InternalServerError = http.StatusInternalServerError
	Ok                  = http.StatusOK
	BadRequest          = http.StatusBadRequest
	Unauthorized        = http.StatusUnauthorized
	UserLead            = "user_lead"
	Dsn                 = "u190031182_physio_relief:PhysioRelief@123@tcp(thephysiorelief.com:3306)/u190031182_physio_relief"
	BusinessSuccess     = 2
	BusinessFailure     = 1
)
