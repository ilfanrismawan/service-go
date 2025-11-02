package utils

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"service/internal/shared/config"
)

// FormatCurrency formats amount as Indonesian Rupiah
// Example: 1000000 -> "Rp. 1.000.000"
func FormatCurrency(amount float64) string {
	if config.Config != nil && config.Config.Currency == "IDR" {
		return FormatIDR(amount)
	}
	// Fallback to simple format
	return fmt.Sprintf("%.2f", amount)
}

// FormatIDR formats amount as Indonesian Rupiah with thousand separators
func FormatIDR(amount float64) string {
	// Round to nearest integer
	rupiah := int64(math.Round(amount))
	
	// Convert to string and add thousand separators
	str := strconv.FormatInt(rupiah, 10)
	n := len(str)
	if n <= 3 {
		return fmt.Sprintf("Rp. %s", str)
	}
	
	// Add dots as thousand separators
	var result strings.Builder
	result.WriteString("Rp. ")
	
	// Calculate first group length (remainder when divided by 3)
	firstGroupLen := n % 3
	if firstGroupLen == 0 {
		firstGroupLen = 3
	}
	
	// Write first group
	result.WriteString(str[:firstGroupLen])
	
	// Write remaining groups with dots
	for i := firstGroupLen; i < n; i += 3 {
		result.WriteString(".")
		result.WriteString(str[i : i+3])
	}
	
	return result.String()
}

// FormatDate formats date according to configured date format
func FormatDate(t interface{}) string {
	if config.Config == nil {
		return ""
	}
	// This is a placeholder - would need to implement based on config.DateFormat
	// For now, return ISO format
	return ""
}

