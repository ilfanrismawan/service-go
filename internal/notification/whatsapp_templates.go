package notification

import (
	"fmt"
	"service/internal/utils"
	"time"
)

// WhatsAppTemplateType represents the type of WhatsApp template
type WhatsAppTemplateType string

const (
	TemplateOrderConfirmation WhatsAppTemplateType = "order_confirmation"
	TemplateStatusUpdate      WhatsAppTemplateType = "status_update"
	TemplatePaymentReminder   WhatsAppTemplateType = "payment_reminder"
	TemplatePickupNotification WhatsAppTemplateType = "pickup_notification"
	TemplateDeliveryNotification WhatsAppTemplateType = "delivery_notification"
	TemplatePromoMessage      WhatsAppTemplateType = "promo_message"
)

// GetWhatsAppTemplate returns formatted WhatsApp message based on template type
func GetWhatsAppTemplate(templateType WhatsAppTemplateType, data map[string]interface{}) string {
	switch templateType {
	case TemplateOrderConfirmation:
		return getOrderConfirmationTemplate(data)
	case TemplateStatusUpdate:
		return getStatusUpdateTemplate(data)
	case TemplatePaymentReminder:
		return getPaymentReminderTemplate(data)
	case TemplatePickupNotification:
		return getPickupNotificationTemplate(data)
	case TemplateDeliveryNotification:
		return getDeliveryNotificationTemplate(data)
	case TemplatePromoMessage:
		return getPromoMessageTemplate(data)
	default:
		if msg, ok := data["message"].(string); ok {
			return msg
		}
		return ""
	}
}

// getOrderConfirmationTemplate returns order confirmation template
func getOrderConfirmationTemplate(data map[string]interface{}) string {
	orderNo := getString(data, "order_number", "N/A")
	customerName := getString(data, "customer_name", "Pelanggan")
	branchName := getString(data, "branch_name", "Cabang")
	estimatedCost := getFloat(data, "estimated_cost", 0)
	estimatedDays := getInt(data, "estimated_days", 0)

	return fmt.Sprintf(`Halo %s! 👋

Terima kasih sudah mempercayakan iPhone Anda pada kami.

📋 *Detail Order*
Order Number: %s
Cabang: %s
Estimasi Biaya: %s
Estimasi Waktu: %d hari

Kami akan mengupdate status order Anda secara berkala. 
Terima kasih! 🙏`, customerName, orderNo, branchName, utils.FormatIDR(estimatedCost), estimatedDays)
}

// getStatusUpdateTemplate returns order status update template
func getStatusUpdateTemplate(data map[string]interface{}) string {
	orderNo := getString(data, "order_number", "N/A")
	status := getString(data, "status", "")
	statusText := getStatusText(status)
	notes := getString(data, "notes", "")

	msg := fmt.Sprintf(`📱 *Update Order #%s*

Status: %s`, orderNo, statusText)

	if notes != "" {
		msg += fmt.Sprintf("\n\nCatatan: %s", notes)
	}

	msg += "\n\nTerima kasih!"
	return msg
}

// getPaymentReminderTemplate returns payment reminder template
func getPaymentReminderTemplate(data map[string]interface{}) string {
	orderNo := getString(data, "order_number", "N/A")
	amount := getFloat(data, "amount", 0)
	dueDate := getTime(data, "due_date", time.Now().Add(24*time.Hour))

	return fmt.Sprintf(`💳 *Pengingat Pembayaran*

Order Number: %s
Jumlah: %s
Batas Waktu: %s

Mohon segera lakukan pembayaran agar order Anda dapat diproses.
Terima kasih! 🙏`, orderNo, utils.FormatIDR(amount), dueDate.Format("02/01/2006 15:04"))
}

// getPickupNotificationTemplate returns pickup notification template
func getPickupNotificationTemplate(data map[string]interface{}) string {
	orderNo := getString(data, "order_number", "N/A")
	courierName := getString(data, "courier_name", "Kurir kami")
	eta := getString(data, "eta", "segera")

	return fmt.Sprintf(`🚚 *Notifikasi Pickup*

Order #%s akan dijemput oleh %s.
Estimasi waktu kedatangan: %s

Pastikan iPhone Anda sudah siap untuk dijemput.
Terima kasih! 🙏`, orderNo, courierName, eta)
}

// getDeliveryNotificationTemplate returns delivery notification template
func getDeliveryNotificationTemplate(data map[string]interface{}) string {
	orderNo := getString(data, "order_number", "N/A")
	courierName := getString(data, "courier_name", "Kurir kami")
	eta := getString(data, "eta", "segera")

	return fmt.Sprintf(`📦 *Notifikasi Pengiriman*

Order #%s sedang dalam perjalanan ke alamat Anda.
Kurir: %s
Estimasi waktu: %s

Pastikan ada yang menerima di alamat tujuan.
Terima kasih! 🙏`, orderNo, courierName, eta)
}

// getPromoMessageTemplate returns promotional message template
func getPromoMessageTemplate(data map[string]interface{}) string {
	title := getString(data, "title", "Penawaran Spesial")
	message := getString(data, "message", "")
	validUntil := getTime(data, "valid_until", time.Now().Add(7*24*time.Hour))

	return fmt.Sprintf(`🎉 *%s*

%s

Penawaran berlaku hingga: %s

Terima kasih! 🙏`, title, message, validUntil.Format("02/01/2006"))
}

// Helper functions for template data extraction
func getString(data map[string]interface{}, key, defaultValue string) string {
	if val, ok := data[key].(string); ok && val != "" {
		return val
	}
	return defaultValue
}

func getFloat(data map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := data[key].(float64); ok {
		return val
	}
	return defaultValue
}

func getInt(data map[string]interface{}, key string, defaultValue int) int {
	if val, ok := data[key].(int); ok {
		return val
	}
	return defaultValue
}

func getTime(data map[string]interface{}, key string, defaultValue time.Time) time.Time {
	if val, ok := data[key].(time.Time); ok {
		return val
	}
	return defaultValue
}

// getStatusText returns Indonesian text for order status
func getStatusText(status string) string {
	statusMap := map[string]string{
		"pending_pickup": "Menunggu Penjemputan",
		"on_pickup":      "Sedang Dijemput",
		"in_service":     "Sedang Dikerjakan",
		"ready":          "Siap Diambil",
		"delivered":      "Sudah Dikirim",
		"completed":      "Selesai",
		"cancelled":      "Dibatalkan",
	}
	if text, ok := statusMap[status]; ok {
		return text
	}
	return status
}

