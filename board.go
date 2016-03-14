package main

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	Unoccupied = 0
	Players    = 4
	BoardSize  = Players * 7
	Red        = 1
	Purple     = 2
	White      = 3
	Blue       = 4
)

var Names = []string{".", "Red", "Purple", "White", "Blue"}

type Color int8

func (p Color) Abbr() string {
	return string(Names[p][0])
}

func (p Color) Name() string {
	return Names[p]
}

func (p Color) Idx() int {
	return int(p) - 1
}

type Board struct {
	colors [BoardSize]Color
	home   [Players][4]Color
	base   [Players]int
}

var (
	board Board
)

func NewBoard() *Board {
	b := Board{}
	for i := range b.base {
		b.base[i] = 4
	}
	return &b
}

func startPos(player Color) int {
	return int((player.Idx()) * 7)
}

func (b *Board) IsOccupied(pos int) bool {
	return b.colors[pos] != Unoccupied
}

// Returns false, if move would be illegal
func (b *Board) Move(pos int, steps int) bool {
	ply := b.colors[pos]

	finalpos := (pos + steps) % BoardSize

	if finalpos > startPos(ply) {
		diff := finalpos - startPos(ply)

		if diff < steps && diff <= 4 && b.home[ply.Idx()][diff-1] == Unoccupied { // This was complete round and it fits
			b.home[ply.Idx()][diff-1] = ply
			b.colors[pos] = Unoccupied
			if *verbose {
				fmt.Printf("Moved ply %d to home\n", ply)
			}
			return true
		}
	}

	switch b.colors[finalpos] {
	case Unoccupied:
		b.colors[pos] = Unoccupied
		b.colors[finalpos] = ply
	case ply:
		return false
	default:
		other := b.colors[finalpos]
		b.base[other.Idx()]++
		b.colors[pos] = Unoccupied
		b.colors[finalpos] = ply
	}

	if *verbose {
		fmt.Printf("Moved ply %d to %d\n", ply, finalpos)
	}

	return true // Move was legal
}

func (b *Board) String() string {
	var buf bytes.Buffer

	buf.WriteString(" Board ")
	for i, p := range b.colors {
		if i%7 == 0 && p == Unoccupied {
			buf.WriteRune(',')
		} else {
			buf.WriteString(p.Abbr())
		}
	}

	buf.WriteString("\n")
	for i := 1; i <= Players; i++ {
		buf.WriteString(fmt.Sprintf("%6s ", Names[i]))
		for p := range b.home[i] {
			buf.WriteString(b.home[i-1][p].Abbr())
		}

		buf.WriteString("   base ")

		buf.WriteString(strings.Repeat(Color(i).Abbr(), b.base[i-1]))

		buf.WriteRune('\n')
	}

	return buf.String()
}
