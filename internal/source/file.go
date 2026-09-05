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
	// Name identifies the source for diagnostics and lookup, it should not be interpreted as a filesystem path
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
	FileName string
	Offset   uint32 // zero-based byte offset
	Line     uint32 // one-based line number
	Column   uint32 // one-based byte column
}

func (f *File) FilePos(pos Pos) (FilePos, bool) {
	if pos < f.Base || pos > f.End() {
		return FilePos{}, false
	}

	offset := f.Offset(pos)

	lo, hi := 0, len(f.lines)
	for lo+1 < hi {
		mid := lo + (hi-lo)/2

		if f.lines[mid] <= offset {
			lo = mid
		} else {
			hi = mid
		}
	}

	line := uint32(lo + 1)
	column := offset - f.lines[lo] + 1

	return FilePos{
		FileName: f.Name,
		Offset:   offset,
		Line:     line,
		Column:   column,
	}, true
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
