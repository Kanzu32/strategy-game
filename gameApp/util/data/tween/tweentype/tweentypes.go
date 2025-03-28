package tweentype

type TweenType uint8

//go:generate stringer -type=TweenType
const (
	Linear TweenType = iota + 1
)
