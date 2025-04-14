package validate

func CheckPattern(pattern map[string]string) bool {
	allowed := map[string]bool{
		"status":  true,
		"packing": true,
	}
	used := make(map[string]bool)

	for key := range pattern {
		if !allowed[key] || used[key] {
			return false
		}

		used[key] = true
	}

	return true
}
