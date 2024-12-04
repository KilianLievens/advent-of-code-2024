package advent

func Transpose2D[T any](iMatrix *[][]T) [][]T {
	matrix := *iMatrix
	newMatrix := make([][]T, len(matrix[0]))

	for y := range newMatrix {
		newMatrix[y] = make([]T, len(matrix))
	}

	for y, row := range matrix {
		for x, s := range row {
			newMatrix[x][y] = s
		}
	}

	return newMatrix
}
