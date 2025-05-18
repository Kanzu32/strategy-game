package sprite

import (
	"image"
	"strategy-game/game/singletons"

	"github.com/hajimehoshi/ebiten/v2"
)

type Frame struct {
	N    int `json:"tileid"`
	Time int `json:"duration"` // milliseconds
}

type Sprite struct {
	image               *ebiten.Image
	Animations          map[string][]Frame
	width               int
	height              int
	framesX             int
	framesY             int
	currentAnimation    string
	currentFrame        int
	animationStartFrame int
}

func (s *Sprite) Width() int {
	return s.width
}

func (s *Sprite) Height() int {
	return s.height
}

func timeToFrames(time int) int {
	f := (float32(time) / 1000.0) * 60
	return int(f) | 1
}

func NewSprite(img *ebiten.Image, w int, h int) Sprite {
	s := Sprite{}
	s.width = w
	s.height = h
	s.image = img
	s.framesX = img.Bounds().Dx() / w
	s.framesY = img.Bounds().Dy() / h
	s.Animations = make(map[string][]Frame)
	s.currentFrame = 0
	return s
}

func (s *Sprite) AnimationProgress() float64 {
	time := 0
	for _, frame := range s.Animations[s.currentAnimation] {
		time += frame.Time
	}
	return float64(timeToFrames(time)) / float64(s.currentFrame+1)
}

func (s *Sprite) Animate() *ebiten.Image {
	frames, ok := s.Animations[s.currentAnimation]
	if !ok || s.currentFrame < 0 || s.currentFrame > len(frames) {
		return nil
	}
	var f Frame
	if (singletons.FrameCount-s.animationStartFrame)%timeToFrames(frames[s.currentFrame].Time) != 0 {
		f = frames[s.currentFrame]
	} else {
		s.currentFrame = (s.currentFrame + 1) % len(frames)
		f = frames[s.currentFrame]
	}
	sub := s.image.SubImage(image.Rect(
		(f.N%s.framesX)*s.width,
		(f.N/s.framesX)*s.height,
		(f.N%s.framesX+1)*s.width,
		(f.N/s.framesX+1)*s.height,
	)).(*ebiten.Image)
	return sub
}

func (s *Sprite) SetAnimation(animationName string) {
	_, ok := s.Animations[animationName]
	if !ok {
		panic("UNKNOWN ANIMATION")
	}
	s.currentAnimation = animationName
	s.currentFrame = 0
	s.animationStartFrame = singletons.FrameCount - 1
	// log.Println("ANIMATION SET ", animationName)
}

func (s *Sprite) AddAnimation(animationName string, frames []Frame) {
	s.Animations[animationName] = frames
}
