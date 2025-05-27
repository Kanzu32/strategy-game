package sprite_test

import (
	"strategy-game/game/singletons"
	"strategy-game/util/data/sprite"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
)

func TestSpriteCreate(t *testing.T) {
	spr := sprite.NewSprite(ebiten.NewImage(1, 1), 1, 1)

	if spr.Height() != 1 || spr.Width() != 1 {
		t.Fail()
	}

	println(spr.AnimationProgress())
	if spr.AnimationProgress() != 0.0 {
		t.Fail()
	}
}

func TestSpriteSetAnimation(t *testing.T) {
	spr := sprite.NewSprite(ebiten.NewImage(1, 1), 1, 1)

	if spr.Height() != 1 || spr.Width() != 1 {
		t.Fail()
	}

	spr.AddAnimation("test", []sprite.Frame{sprite.Frame{N: 1, Time: 1000}})

	spr.SetAnimation("test")

	if spr.AnimationProgress() != 0.0 {
		t.Fail()
	}
}

func TestSpriteAnimate(t *testing.T) {
	spr := sprite.NewSprite(ebiten.NewImage(1, 1), 1, 1)

	if spr.Height() != 1 || spr.Width() != 1 {
		t.Fail()
	}

	spr.AddAnimation("test", []sprite.Frame{
		sprite.Frame{N: 0, Time: 1000},
		sprite.Frame{N: 0, Time: 1000},
	})

	spr.SetAnimation("test")

	for i := 0; i < 100; i++ {
		singletons.FrameCount++
		res := spr.Animate()
		if res == nil {
			t.Fail()
		}
	}

	if spr.AnimationProgress() == 0.0 {
		t.Fail()
	}

	singletons.FrameCount = 0
}
