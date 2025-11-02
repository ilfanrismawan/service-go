package utils

import (
	"service/internal/config"
)

// GetErrorMessage returns error message in configured language
func GetErrorMessage(key string) string {
	if config.Config == nil {
		return key
	}
	
	lang := config.Config.DefaultLanguage
	if lang == "id-ID" {
		return GetIndonesianError(key)
	}
	// Default to English
	return GetEnglishError(key)
}

// GetIndonesianError returns Indonesian error messages
func GetIndonesianError(key string) string {
	messages := map[string]string{
		"validation_error":          "Validasi gagal",
		"unauthorized":              "Tidak memiliki izin",
		"forbidden":                 "Akses ditolak",
		"not_found":                 "Data tidak ditemukan",
		"invalid_input":             "Input tidak valid",
		"internal_error":            "Terjadi kesalahan internal",
		"user_not_found":            "Pengguna tidak ditemukan",
		"invalid_password":          "Password tidak valid",
		"order_not_found":           "Order tidak ditemukan",
		"branch_not_found":          "Cabang tidak ditemukan",
		"payment_not_found":         "Pembayaran tidak ditemukan",
		"invalid_token":             "Token tidak valid",
		"token_expired":             "Token telah kadaluarsa",
		"email_exists":              "Email sudah terdaftar",
		"phone_exists":              "Nomor telepon sudah terdaftar",
		"registration_failed":       "Registrasi gagal",
		"login_failed":              "Login gagal",
		"payment_failed":            "Pembayaran gagal",
		"order_creation_failed":     "Pembuatan order gagal",
		"invalid_id":                "Format ID tidak valid",
		"invalid_order_id":          "Format Order ID tidak valid",
		"invalid_payment_id":        "Format Payment ID tidak valid",
	}
	
	if msg, ok := messages[key]; ok {
		return msg
	}
	return key
}

// GetEnglishError returns English error messages (default)
func GetEnglishError(key string) string {
	// Return key as-is for English (can be expanded later)
	return key
}

