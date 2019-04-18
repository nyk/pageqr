package pageqr

// StringSet is a map type that represents a lookup set of strings
type StringSet map[string]bool

// NewStringSet constructs a new StringSet instance
func NewStringSet(slice []string) StringSet {

	m := make(StringSet)
	for _, val := range slice {
		m[val] = true
	}

	return m

}
