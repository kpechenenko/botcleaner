package bot

func convertListToSet[T comparable](arr []T, convertKey func(T) T) map[T]bool {
	m := make(map[T]bool)
	for _, v := range arr {
		m[convertKey(v)] = true
	}
	return m
}
