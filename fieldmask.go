package fieldmask

import "strconv"

type Slice[T ~string] string

func (f Slice[T]) Index(i int) T {
	is := "[" + strconv.Itoa(i) + "]"
	if f == "" {
		return T(is)
	}
	return T(string(f) + "." + is)
}
