package slice

type Predicate[T any] func(item T) bool

func Chunk[T any](slice []T, chunkSize int) [][]T {
	var chunks [][]T

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Unique[T comparable](s []T) []T {
	inResult := make(map[T]bool)
	var result []T
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func Map[T1, T2 any](input []T1, f func(T1) T2) (output []T2) {
	output = make([]T2, 0, len(input))
	for _, v := range input {
		output = append(output, f(v))
	}
	return output
}

func Filter[T any](slice []T, p Predicate[T]) []T {
	var n []T
	for _, e := range slice {
		if p(e) {
			n = append(n, e)
		}
	}
	return n
}

func Find[T any](collection []T, p Predicate[T]) (T, bool) {
	for _, item := range collection {
		if p(item) {
			return item, true
		}
	}
	var result T
	return result, false
}
