package advent

func RotateRight2D[T any](matrix [][]T) [][]T {
	Assert(len(matrix) > 1 && len(matrix[0]) > 1, "Input matrix too small")

	newMatrix := [][]T{}
	for x := 0; x < len(matrix); x++ {
		line := []T{}
		for y := len(matrix) - 1; y >= 0; y-- {
			line = append(line, matrix[y][x])
		}
		newMatrix = append(newMatrix, line)
	}

	Assert(len(newMatrix) == len(matrix[0]) && len(newMatrix[0]) == len(matrix), "Rotated matrix dimensions incorrect")

	return newMatrix
}
