package std

var sboolTrueMap = map[string]bool{
	"true": true,
	"TRUE": true,
	"yes":  true,
	"YES":  true,
	"1":    true,
}

type Sbool struct {
	Str   string
	Valid bool
}

func NewSbool(s string) *Sbool {
	return &Sbool{
		Str:   s,
		Valid: true,
	}
}

func (b *Sbool) Value() bool {
	if b == nil {
		return false
	} else if !b.Valid {
		return false
	}
	return sboolTrueMap[b.Str]
}

// Stringer interface
func (b *Sbool) String() string {
	if b.Value() {
		return "true"
	} else {
		return "false"
	}
}
