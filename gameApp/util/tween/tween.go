package tween

import "strategy-game/util/tween/tweentype"

type TransitionValues struct {
	X     int
	Y     int
	Angle int
}

type TweenAnimation struct {
	Type         tweentype.TweenType
	FrameCount   int
	CurrentFrame int
	DeltaValues  TransitionValues
}

func CreateTween(animationType tweentype.TweenType, time int, X int, Y int, Angle int) TweenAnimation {
	count := timeToFrames(time)
	return TweenAnimation{
		Type:         animationType,
		FrameCount:   count,
		CurrentFrame: 0,
		DeltaValues: TransitionValues{
			X / timeToFrames(time),
			Y / timeToFrames(time),
			Angle / timeToFrames(time),
		},
	}
}

func (t *TweenAnimation) Animate() TransitionValues {
	if t.CurrentFrame < t.FrameCount {
		t.CurrentFrame = t.CurrentFrame + 1
	}
	return TransitionValues{
		X:     t.DeltaValues.X * t.CurrentFrame,
		Y:     t.DeltaValues.Y * t.CurrentFrame,
		Angle: t.DeltaValues.Angle * t.CurrentFrame,
	}
}

func (t *TweenAnimation) IsEnded() bool {
	return t.CurrentFrame == t.FrameCount
}

func timeToFrames(time int) int {
	f := (float32(time) / 1000.0) * 60
	return int(f) | 1
}
