package advent

func CalcManhattanDistance(a, b Coord) int {
	return AbsInt(a.X-b.X) + AbsInt(a.Y-b.Y)
}
