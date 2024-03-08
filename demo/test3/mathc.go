package mathc

type Sumifier interface {
	Add(a, b int32) int32
}

type Sumer struct {
	id int32
}

func (math Sumer) Add(a, b int32) int32 {
	return a + b
}

type SumerPointer struct {
	id int32
}

func (math *SumerPointer) Add(a, b int32) int32 {
	return a + b
}
