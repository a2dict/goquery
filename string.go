package goquery

type StringFilter func(string) bool

func (f StringFilter) Do(s string) bool {
	return f(s)
}

// ContainStringFilter ...
func ContainStringFilter(strs ...string) StringFilter {
	return func(s string) bool {
		for _, t := range strs {
			if t == s {
				return true
			}
		}
		return false
	}
}

// ExcludeStringFilter ...
func ExcludeStringFilter(strs ...string) StringFilter {
	cf := ContainStringFilter(strs...)
	return func(s string) bool {
		return !cf(s)
	}
}
