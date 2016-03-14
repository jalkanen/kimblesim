package main

type Player interface {
	Move(board *Board, roll int) bool
	Color() Color
	Idx() int
}

type GenericPlayer struct {
	Col Color
}

func (p *GenericPlayer) Idx() int {
	return p.Color().Idx()
}

func (p *GenericPlayer) Color() Color {
	return p.Col
}

// Rolled a six, so start the game
func (p *GenericPlayer) start(b *Board) bool {
	if b.base[p.Idx()] <= 0 {
		return false
	} // No more pieces in base
	sp := startPos(p.Col)

	switch b.colors[sp] {
	case Unoccupied:
		b.colors[sp] = p.Col
		b.base[p.Idx()]--
	case p.Col:
		// Occupied with own color
		return false
	default:
		other := b.colors[sp]
		b.base[other.Idx()]++
		b.colors[sp] = p.Col
		b.base[p.Idx()]--
	}

	return true
}

type RandomPlayer struct {
	GenericPlayer
}

func NewRandom(col Color) *RandomPlayer {
	return &RandomPlayer{
		GenericPlayer: GenericPlayer{Col: col},
	}
}

func (p *RandomPlayer) Move(board *Board, roll int) bool {

	if roll == 6 {
		if p.start(board) {
			return true
		}
	}

	// OK, so pick a random move then

	for i, v := range board.colors {
		if v == p.Color() {
			if board.Move(i, roll) {
				return true
			}
		}
	}

	return false
}
