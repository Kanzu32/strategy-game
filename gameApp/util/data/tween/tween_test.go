package tween_test

import (
	"strategy-game/util/data/tween"
	"strategy-game/util/data/tween/tweentype"
	"testing"
)

func TestTweenCreate(t *testing.T) {
	animation := tween.CreateTween(tweentype.Linear, 10, 10, 20, 30)

	val := animation.GetValue()
	if val.X != 0 {
		t.Fail()
	}

	if val.Y != 0 {
		t.Fail()
	}

	if val.Angle != 0 {
		t.Fail()
	}
}

func TestTweenAnimate(t *testing.T) {
	animation := tween.CreateTween(tweentype.Linear, 10, 10, 20, 30)

	animation.Animate()

	val := animation.GetValue()
	if val.X == 0 {
		t.Fail()
	}

	if val.Y == 0 {
		t.Fail()
	}

	if val.Angle == 0 {
		t.Fail()
	}
}

func TestTweenEnd(t *testing.T) {
	animation := tween.CreateTween(tweentype.Linear, 1, 10, 20, 30)

	for i := 0; i < 100; i++ {
		animation.Animate()
	}

	val := animation.GetValue()
	if val.X < 10 {
		t.Fail()
	}

	if val.Y < 20 {
		t.Fail()
	}

	if val.Angle < 30 {
		t.Fail()
	}

	if !animation.IsEnded() {
		t.Fail()
	}
}

func TestTweenBackAndForward(t *testing.T) {
	animation := tween.CreateTween(tweentype.Back75Forward25, 500, 10, 20, 30)

	for i := 0; i < 100; i++ {
		animation.Animate()
	}

	if !animation.IsEnded() {
		t.Fail()
	}
}

func TestTweenXSin(t *testing.T) {
	animation := tween.CreateTween(tweentype.XSin, 1, 10, 20, 30)

	for i := 0; i < 100; i++ {
		animation.Animate()
	}

	if !animation.IsEnded() {
		t.Fail()
	}
}
