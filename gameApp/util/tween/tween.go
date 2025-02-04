package tween

import (
	"strategy-game/util/tween/tweentype"
)

type TransitionValues struct {
	X     float64
	Y     float64
	Angle float64
}

type TweenAnimation struct {
	Type             tweentype.TweenType
	FrameTime        float64
	CurrentFrameTime float64
	DeltaValues      TransitionValues
}

func CreateTween(animationType tweentype.TweenType, time int, X int, Y int, Angle int) TweenAnimation {
	return TweenAnimation{
		Type:             animationType,
		FrameTime:        timeToFrames(time),
		CurrentFrameTime: 0.0,
		DeltaValues: TransitionValues{
			float64(X) / timeToFrames(time),
			float64(Y) / timeToFrames(time),
			float64(Angle) / timeToFrames(time),
		},
	}
}

func (t *TweenAnimation) Animate() TransitionValues {
	if t.CurrentFrameTime < t.FrameTime {
		t.CurrentFrameTime = t.CurrentFrameTime + 1
	}
	return TransitionValues{
		X:     t.DeltaValues.X * t.CurrentFrameTime,
		Y:     t.DeltaValues.Y * t.CurrentFrameTime,
		Angle: t.DeltaValues.Angle * t.CurrentFrameTime,
	}
}

func (t *TweenAnimation) GetValue() TransitionValues {
	return TransitionValues{
		X:     t.DeltaValues.X * t.CurrentFrameTime,
		Y:     t.DeltaValues.Y * t.CurrentFrameTime,
		Angle: t.DeltaValues.Angle * t.CurrentFrameTime,
	}
}

func (t *TweenAnimation) IsEnded() bool {
	return t.CurrentFrameTime >= t.FrameTime
}

func timeToFrames(time int) float64 {
	f := (float64(time) / 1000.0) * 60
	return f
}
