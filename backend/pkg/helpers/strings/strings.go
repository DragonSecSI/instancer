package strings

type Strings struct {
	FuzzyTime func(int) string
}

func NewStringsHelper() Strings {
	return Strings{
		FuzzyTime: fuzzyTime,
	}
}
