package utils

// CalculatePPN calculates PPN 11% for Indonesian tax
func CalculatePPN(amount float64) float64 {
	return amount * 0.11
}

// CalculateSubtotal calculates subtotal from total amount (reverse PPN calculation)
func CalculateSubtotal(totalAmount float64) float64 {
	// If total includes PPN: subtotal = total / 1.11
	return totalAmount / 1.11
}

// CalculateAmountWithTax calculates total amount including PPN 11%
func CalculateAmountWithTax(subtotal float64) float64 {
	return subtotal + CalculatePPN(subtotal)
}
