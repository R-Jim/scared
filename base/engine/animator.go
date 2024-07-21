package engine

type RenderLayer string

type RenderPosition struct {
	X int
	Y int
}

type Frame struct {
	Image          []byte
	RenderPosition *RenderPosition // Pointer the give the ability to change frame position without creating a new animation
	RenderLayer    RenderLayer
}

type Animation struct {
	index  int
	isLoop bool
	frames []Frame
}

func NewAnimation(isLoop bool, frames ...Frame) Animation {
	return Animation{
		isLoop: isLoop,
		frames: frames,
	}
}

func (a *Animation) Frame(fps int) Frame {
	frame := a.frames[a.index]

	a.index++
	if a.isLoop && len(a.frames) <= a.index {
		a.index = 0
	}

	return frame
}

func (a *Animation) IsCompleted() bool {
	return !a.isLoop && len(a.frames) <= a.index
}

type Animator interface {
	GetHook() Hook
	Frame() []Frame
}
