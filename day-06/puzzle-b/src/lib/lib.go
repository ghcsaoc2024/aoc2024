package lib

type Coord struct {
	Row int
	Col int
}

type Cell int

const (
	Empty Cell = iota
	Blocked
)

func NextCoords(current, dir Coord) Coord {
	return Coord{
		Row: current.Row + dir.Row,
		Col: current.Col + dir.Col,
	}
}

func IsValidCoord(coords, dimensions Coord) bool {
	if coords.Row < 0 || coords.Row >= dimensions.Row {
		return false
	}
	if coords.Col < 0 || coords.Col >= dimensions.Col {
		return false
	}
	return true
}

func TurnRight(dir Coord) Coord {
	return Coord{
		Row: dir.Col,
		Col: -dir.Row,
	}
}
