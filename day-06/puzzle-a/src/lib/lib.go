package lib

type Coord struct {
	Row int
	Col int
}

type Cell int

const (
	Empty Cell = iota
	Visited
	Blocked
)

func (c Coord) MoveOne(dir Coord) Coord {
	return Coord{
		Row: c.Row + dir.Row,
		Col: c.Col + dir.Col,
	}
}

func (c Coord) IsValid(dimension Coord) bool {
	if c.Row < 0 || c.Row >= dimension.Row {
		return false
	}
	if c.Col < 0 || c.Col >= dimension.Col {
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
