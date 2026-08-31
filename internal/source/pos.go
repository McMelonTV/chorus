package source

type Pos uint32

type Span struct {
	Start  Pos
	Length uint32
}
