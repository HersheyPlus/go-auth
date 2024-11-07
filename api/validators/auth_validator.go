package validators

import (
	"github.com/HersheyPlus/go-auth/dto"
	"strings"
)

func ValidateRegisterFields(req *dto.UserRegisterRequest) string {
	var missingFields []string

	if isEmptyField(req.Username) {
		missingFields = append(missingFields, "username")
	}
	if isEmptyField(req.Phone) {
		missingFields = append(missingFields, "phone")
	}
	if isEmptyField(req.Email) {
		missingFields = append(missingFields, "email")
	}
	if isEmptyField(req.Password) {
		missingFields = append(missingFields, "password")
	}

	if len(req.Password) < 8 {
		missingFields = append(missingFields, "password must be at least 8 characters")
	}
	

	if len(missingFields) > 0 {
		return "Missing required fields: " + strings.Join(missingFields, ", ")
	}

	return ""
}


func isEmptyField(field string) bool {
	return field == ""
}