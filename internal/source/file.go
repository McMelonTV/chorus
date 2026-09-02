package source

import "math"

type FileSet struct {
	files      []*File
	fileByName map[string]*File
	nextBase   Pos
}

func NewFileSet() *FileSet {
	return &FileSet{
		fileByName: make(map[string]*File),
		nextBase:   1,
	}
}

func (fs *FileSet) AddFile(name string, data []byte) (*File, bool) {
	if fs.nextBase == 0 {
		fs.nextBase = 1
	}

	if fs.fileByName == nil {
		fs.fileByName = make(map[string]*File)
	}

	if _, ok := fs.fileByName[name]; ok {
		return nil, false
	}

	if uint64(fs.nextBase)+uint64(len(data))+1 > math.MaxUint32 {
		return nil, false
	}

	file := &File{
		Name: name,
		Data: data,
		Base: fs.nextBase,

		lines: buildLines(data),
	}

	fs.files = append(fs.files, file)
	fs.fileByName[name] = file

	fs.nextBase += Pos(len(data)) + 1

	return file, true
}

func buildLines(data []byte) []uint32 {
	lines := []uint32{0} // first (0-th in this array) line begins at index 0

	for i, b := range data {
		if b == '\n' {
			lines = append(lines, uint32(i+1))
		}
	}

	return lines
}

type File struct {
	Name string
	Data []byte
	Base Pos

	lines []uint32
}

func (f *File) Pos(offset uint32) Pos {
	return f.Base + Pos(offset)
}

func (f *File) Offset(pos Pos) uint32 {
	return uint32(pos - f.Base)
}

func (f *File) End() Pos {
	return f.Base + Pos(len(f.Data))
}

type FilePos struct {
	FileName             string
	Offset, Line, Column uint32
}

type Pos uint32

const NoPos Pos = 0

type Span struct {
	Start  Pos
	Length uint32
}

func NewSpan(start, end Pos) (Span, bool) {
	if end < start || start == 0 || end == 0 {
		return Span{}, false
	}

	return Span{Start: start, Length: uint32(end - start)}, true
}

func (s Span) End() Pos {
	return s.Start + Pos(s.Length)
}
