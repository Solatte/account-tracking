package utils

type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return "string representation of flag"
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func BoolPointer(b bool) *bool {
	return &b
}

func Uint64Ptr(val uint64) *uint64 {
	return &val
}
