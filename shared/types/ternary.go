package types

type Ternary int

const (
	Ternary_Positive Ternary = iota
	Ternary_Neutral
	Ternary_Negative
)

func (t Ternary) String() string {
	switch t {
	case Ternary_Negative:
		return "Ternary_Negative"
	case Ternary_Positive:
		return "Ternary_Positive"
	case Ternary_Neutral:
		return "Ternary_Neutral"
	default:
		return "invalid"
	}
}

func (t Ternary) IsPositive() bool {
	return t == Ternary_Positive
}

func (t Ternary) IsNeutral() bool {
	return t == Ternary_Neutral
}

func (t Ternary) IsNegative() bool {
	return t == Ternary_Negative
}
