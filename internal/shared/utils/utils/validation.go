package utils

import (
	"encoding/base64"
	"fmt"
	"math"
	"regexp"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// InitValidator initializes the validator
func InitValidator() {
	validate = validator.New()

	// Register custom validators
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("phone", validatePhone)
}

// ValidateStruct validates a struct using go-playground/validator
func ValidateStruct(s interface{}) error {
	if validate == nil {
		InitValidator()
	}

	return validate.Struct(s)
}

// validatePassword validates password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Password must be at least 6 characters
	if len(password) < 6 {
		return false
	}

	// Password must contain at least one letter and one number
	hasLetter := false
	hasNumber := false

	for _, char := range password {
		if char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' {
			hasLetter = true
		}
		if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
}

// validatePhone validates Indonesian phone number
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()

	// Remove spaces and special characters
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Check if phone starts with +62 or 08
	if strings.HasPrefix(phone, "+62") {
		phone = "0" + phone[3:]
	}

	// Check if phone starts with 08 and has 10-13 digits
	if strings.HasPrefix(phone, "08") && len(phone) >= 10 && len(phone) <= 13 {
		// Check if all remaining characters are digits
		for _, char := range phone[2:] {
			if char < '0' || char > '9' {
				return false
			}
		}
		return true
	}

	return false
}

// GetValidationErrors returns formatted validation errors
func GetValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			field := strings.ToLower(e.Field())
			errors[field] = getValidationMessage(e)
		}
	}

	return errors
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", e.Field(), e.Param())
	case "password":
		return fmt.Sprintf("%s must be at least 6 characters and contain letters and numbers", e.Field())
	case "phone":
		return fmt.Sprintf("%s must be a valid Indonesian phone number", e.Field())
	default:
		return fmt.Sprintf("%s is invalid", e.Field())
	}
}

// ValidateRequiredFields validates that required fields are not empty
func ValidateRequiredFields(data interface{}, requiredFields []string) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a struct")
	}

	for _, fieldName := range requiredFields {
		field := v.FieldByName(fieldName)
		if !field.IsValid() {
			continue
		}

		if field.Kind() == reflect.String && field.String() == "" {
			return fmt.Errorf("%s is required", fieldName)
		}

		if field.Kind() == reflect.Ptr && field.IsNil() {
			return fmt.Errorf("%s is required", fieldName)
		}
	}

	return nil
}

// SanitizeString removes leading/trailing whitespace and converts to lowercase
func SanitizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// stripXSS removes common XSS vectors like <script> tags and javascript: URIs
var scriptTagRegex = regexp.MustCompile(`(?is)<\s*script[^>]*>.*?<\s*/\s*script\s*>`)
var jsUriRegex = regexp.MustCompile(`(?is)javascript:\s*`)

// SanitizeXSSString cleans potentially dangerous content
func SanitizeXSSString(s string) string {
	clean := scriptTagRegex.ReplaceAllString(s, "")
	clean = jsUriRegex.ReplaceAllString(clean, "")
	return clean
}

// SanitizeStructStrings walks through struct string fields and applies XSS sanitization
func SanitizeStructStrings(v interface{}) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}
	val = val.Elem()
	if val.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(SanitizeXSSString(field.String()))
		case reflect.Struct:
			f := val.Field(i).Addr().Interface()
			SanitizeStructStrings(f)
		case reflect.Ptr:
			if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
				SanitizeStructStrings(field.Interface())
			}
		}
	}
}

// FormatPhone formats Indonesian phone number
func FormatPhone(phone string) string {
	// Remove spaces and special characters
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Convert +62 to 08
	if strings.HasPrefix(phone, "+62") {
		phone = "0" + phone[3:]
	}

	return phone
}

// GenerateOrderNumber generates a unique order number
func GenerateOrderNumber() string {
	// Format: ORD-YYYYMMDD-HHMMSS-XXXX
	// This is a simple implementation - in production, you might want to use a more sophisticated approach
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%06d", now.Format("20060102"), now.Unix()%1000000)
}

// GenerateInvoiceNumber generates a unique invoice number
func GenerateInvoiceNumber() string {
	// Format: INV-YYYYMMDD-HHMMSS-XXXX
	now := time.Now()
	return fmt.Sprintf("INV-%s-%06d", now.Format("20060102"), now.Unix()%1000000)
}

// CalculateDistance calculates distance between two coordinates using Haversine formula
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Haversine formula
	dlat := lat2Rad - lat1Rad
	dlon := lon2Rad - lon1Rad

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// ParseInt parses string to int
func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ParseFloat parses string to float64
func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// ParseBool parses string to bool
func ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// EncodeBase64 encodes string to base64
func EncodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// DecodeBase64 decodes base64 string
func DecodeBase64(s string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
