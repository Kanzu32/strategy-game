package tweentype

type TweenType uint8

//go:generate stringer -type=TweenType
const (
	StrightLinear TweenType = iota + 1
)
