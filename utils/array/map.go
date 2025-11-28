package array

func ToMapped[T, U any](seq []T, f func(T) U) []U {
	var result []U
	for a := range seq {
		result = append(result, f(seq[a]))
	}
	return result
}
