package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)
const winWidth, winHeight int = 800, 600

type color struct {
	r, g, b byte
}

type pos struct{
	x, y float32
}

type ball struct{
	pos
	radius	int
	xv		float32
	yv		float32
	color 	color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius ; y < ball.radius ; y++ {
		for x := -ball.radius ; x < ball.radius ; x ++ {
			if x * x + y * y < ball.radius * ball.radius {
				setPixel(int(ball.x) + x, int(ball.y)+y, color{255, 255, 255}, pixels)
			} 
		}
	}
}

func (ball *ball) update() {
	ball.x += ball.xv
	ball.y += ball.yv

	if ball.y < 0 {
		ball.yv = -ball.yv
	}else if int(ball.y) > winHeight {
		ball.yv = -ball.yv
	}
}

type paddle struct{
	pos
	w		int
	h		int
	color	color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w / 2
	startY := int(paddle.y) - paddle.h / 2

	for y := 0 ; y < paddle.h ; y++ {
		for x := 0 ; x < paddle.w ; x ++ {
			setPixel(startX + x, startY + y, color{255, 255, 255}, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8) {

}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y * winWidth + x) * 4

	if index < len(pixels) - 4 && index >= 0{
		pixels[index] = c.r
		pixels[index + 1] = c.g
		pixels[index + 2] = c.b
	}
}

func main() {
	window, err := sdl.CreateWindow("Testing SDL Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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
	
		tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
		if err != nil{
			fmt.Println(err)
			return
		}
		defer tex.Destroy()
	
		pixels := make([]byte, winWidth*winHeight*4)
	


		player1 := paddle{pos{100, 100}, 20, 100, color{255,255,255}}
		ball := ball{pos{300,300}, 20, 0, 0, color{255, 255, 255}}

		//keyState := sdl.GetKeyboardState()

		for {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch event.(type) {
				case *sdl.QuitEvent:
					return
				}
			}
			player1.draw(pixels)
			ball.draw(pixels)
			tex.Update(nil, pixels, winWidth * 4)
			renderer.Copy(tex, nil, nil)
			renderer.Present()


			sdl.Delay(16)
		}
}