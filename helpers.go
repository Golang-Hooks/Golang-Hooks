package hooks

// insert inserts an element at a specific index.
func insert[T any](a []T, index int, value T) []T {
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}
