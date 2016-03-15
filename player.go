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

func (p *GenericPlayer) diffToHomeStretch(pos int) int {
	diff := startPos(p.Color())-pos

	if diff <= 0 {
		diff = BoardSize+diff
	}

	return diff
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

// Moves the first available piece
type FirstMover struct {
	GenericPlayer
}

func NewFirstMover(col Color) *FirstMover {
	return &FirstMover{
		GenericPlayer: GenericPlayer{Col: col},
	}
}

func (p *FirstMover) Move(board *Board, roll int) bool {

	for i := BoardSize-1; i >= 0; i-- {
		pos := (startPos(p.Col) - i)

		if pos < 0 {
			pos = BoardSize+pos
		}
		v := board.colors[pos]

		if v == p.Color() {
			if board.Move(pos, roll) {
				return true
			}
		}
	}

	if roll == 6 {
		return p.start(board)
	}

	return false
}

// Moves the last available piece
type LastMover struct {
	GenericPlayer
}

func NewLastMover(col Color) *LastMover {
	return &LastMover{
		GenericPlayer: GenericPlayer{Col: col},
	}
}

func (p *LastMover) Move(board *Board, roll int) bool {

	for i := 0; i < BoardSize; i++ {
		pos := (startPos(p.Col) + i)

		if pos >= BoardSize {
			pos = pos - BoardSize
		}
		//fmt.Printf("pos=%d\n",pos)
		v := board.colors[pos]

		if v == p.Color() {
			if board.Move(pos, roll) {
				return true
			}
		}
	}

	if roll == 6 {
		return p.start(board)
	}

	return false
}

// Attempts to always prefer eating move, then moves forward piece
type Eater struct {
	GenericPlayer
	FirstMover
}

func NewEater(col Color) *Eater {
	return &Eater{
		GenericPlayer: GenericPlayer{Col: col},
		FirstMover: *NewFirstMover(col),
	}
}

func (p *Eater) Move(board *Board, roll int) bool {

	if roll == 6 && board.colors[startPos(p.Col)] != Unoccupied && board.colors[startPos(p.Col)] != p.Col {
		if p.start(board) {
			return true
		}
	}

	for i := 0; i < BoardSize; i++ {
		pos := (startPos(p.Col) + i)

		if pos >= BoardSize {
			pos = pos - BoardSize
		}

		v := board.colors[pos]

		if v == p.Color() {
			newpos := (pos + roll)%BoardSize

			if board.colors[newpos] != Unoccupied && board.colors[newpos] != p.Color() {
				return board.Move(pos,roll)
			}
		}
	}

	return p.FirstMover.Move(board,roll)
}
