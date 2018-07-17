package main

import (
	"math/rand"
	"fmt"
	. "github.com/awebbdev/gameswithgo/evolvingpictures/apt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const winWidth, winHeight, winDepth int = 800, 600, 100

type audioState struct {
	explosionBytes []byte
	deviceID       sdl.AudioDeviceID
	audiosSpec     *sdl.AudioSpec
}

type mouseState struct {
	leftButton  bool
	rightButton bool
	x, y        int
}

func getMouseState() mouseState {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var result mouseState
	result.x = int(mouseX)
	result.y = int(mouseY)
	result.leftButton = !(leftButton == 0)
	result.rightButton = !(rightButton == 0)
	return result
}

type rgba struct {
	r, g, b byte
}

type picture struct {
	r Node
	g Node
	b Node
}

func (p *picture) String() string {
	return "R" + p.r.String() + "\nG" + p.g.String() + "\nB" + p.b.String()
}

func NewPicture() *picture {
	p := &picture{}
	p.r = GetRandomNode()
	p.g = GetRandomNode()
	p.b = GetRandomNode()

	num := rand.Intn(20)
	for i := 0; i < num; i++ {
		p.r.AddRandom(GetRandomNode())
	}
	num = rand.Intn(20)
	for i := 0; i < num; i++ {
		p.g.AddRandom(GetRandomNode())
	}

	num = rand.Intn(20)
	for i := 0; i < num; i++ {
		p.b.AddRandom(GetRandomNode())
	}

	for p.r.AddLeaf(GetRandomLeaf()){

	}
	for p.g.AddLeaf(GetRandomLeaf()){
		
	}
	for p.b.AddLeaf(GetRandomLeaf()){
		
	}
	return p
}

func (p *picture) Mutate() {
	r := rand.Intn(3)
	var nodeToMutate Node
	switch r {
	case 0:
		nodeToMutate = p.r
	case 1:
		nodeToMutate = p.g
	case 2: 
		nodeToMutate = p.b
	}
	count := nodeToMutate.NodeCount()
	r = rand.Intn(count)
	nodeToMutate, count = GetNthNode(nodeToMutate, r, 0)
	mutation := Mutate(nodeToMutate)
	if nodeToMutate == p.r {
		p.r = mutation
	}else if nodeToMutate == p.g {
		p.g = mutation
	}else if nodeToMutate == p.b {
		p.b = mutation
	}	
}

type pos struct {
	x, y float32
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c rgba, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func pixelsToTexture(renderer *sdl.Renderer, pixels []byte, w, h int) *sdl.Texture {
	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(w), int32(h))
	if err != nil {
		panic(err)
	}
	tex.Update(nil, pixels, w*4)
	return tex
}

func aptToTexture(pic *picture, w, h int, renderer *sdl.Renderer) *sdl.Texture {
	//-1.0 and 1.0
	scale := float32(255 / 2)
	offset := float32(-1.0 * scale)
	pixels := make([]byte, w*h*4)
	pixelIndex := 0
	for yi := 0; yi < h; yi++ {
		y := float32(yi)/float32(h)*2 - 1
		for xi := 0; xi < w; xi++ {
			x := float32(xi)/float32(w)*2 - 1

			r := pic.r.Eval(x, y)
			g := pic.g.Eval(x, y)
			b := pic.b.Eval(x, y)

			pixels[pixelIndex] = byte(r*scale - offset)
			pixelIndex++
			pixels[pixelIndex] = byte(g*scale - offset)
			pixelIndex++
			pixels[pixelIndex] = byte(b*scale - offset)
			pixelIndex++
			pixelIndex++

		}
	}
	return pixelsToTexture(renderer, pixels, w, h)
}

func main() {
	sdl.LogSetAllPriority(sdl.LOG_PRIORITY_VERBOSE)
	window, err := sdl.CreateWindow("Evolving Pictures", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	/*  	explosionBytes, audioSpec := sdl.LoadWAV("explode.wav")
	audioID, err := sdl.OpenAudioDevice("", false, audioSpec, nil, 0)
	if err != nil {
		panic(err)
	}
	defer sdl.FreeWAV(explosionBytes)

	audioState := audioState{explosionBytes, audioID, audioSpec} */

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	var elapsedTime float32
	var currentMouseState = getMouseState()
	var previousMouseState = currentMouseState

	rand.Seed(time.Now().UTC().UnixNano())
	pic := NewPicture()
	tex := aptToTexture(pic, 800, 600, renderer)

	for {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.TouchFingerEvent:
				touchX := e.X * float32(winWidth)
				touchY := e.Y * float32(winHeight)
				currentMouseState.x = int(touchX)
				currentMouseState.y = int(touchY)
				currentMouseState.leftButton = true
			}
		}
		currentMouseState = getMouseState()

		if previousMouseState.leftButton && !currentMouseState.leftButton {
			pic.Mutate()
			tex = aptToTexture(pic, winWidth, winHeight, renderer)
		}

		renderer.Copy(tex, nil, nil)
		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
		previousMouseState = currentMouseState
	}
}
