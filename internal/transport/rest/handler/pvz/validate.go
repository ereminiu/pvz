package pvz

func orderByIsValid(field string) bool {
	allowed := map[string]bool{
		"":           true,
		"order_id":   true,
		"user_id":    true,
		"expire_at":  true,
		"price":      true,
		"status":     true,
		"created_at": true,
		"updated_at": true,
	}

	return allowed[field]
}
