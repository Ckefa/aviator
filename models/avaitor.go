package models

type Aviator struct {
	Round int
}

func NewAvaitor() *Aviator {
	return &Aviator{
		Round: 0,
	}
}

type Game struct {
	ID      int
	Players int
	Burst   float32
}

func (a *Aviator) NewGame() *Game {
	return &Game{
		ID: a.Round + 1,
	}
}
