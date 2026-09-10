package handler

func (h *OrderHandler) CreatePayment(orderID int, amount float64) error {
	query := `
		INSERT INTO payments (
			order_id,
			amount,
			status
		)
		VALUES (?, ?, 'pending')
	`

	_, err := h.db.Exec(query, orderID, amount)
	if err != nil {
		return err
	}

	return nil
}
