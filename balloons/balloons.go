package main

import (
	"image/png"
	"os"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/awebbdev/gameswithgo/noise"
)

const winWidth, winHeight int = 800, 600

type texture struct {
	pos
	pixels []byte
	w, h, pitch int
}

type rgba struct {
	r, g, b byte
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
	index := (y * winWidth + x) * 4

	if index < len(pixels) - 4 && index >= 0{
		pixels[index] = c.r
		pixels[index + 1] = c.g
		pixels[index + 2] = c.b
	}
}

func (tex *texture) draw(pixels []byte) {
	for y := 0; y < tex.h; y++ {
		for x := 0; x < tex.w; x++ {
			screenY := y + int(tex.y)
			screenX := x + int(tex.x)
			if screenX >= 0 && screenX < winWidth && screenY >= 0 && screenY < winHeight {
				texIndex := y*tex.pitch + x*4 
				screenIndex := screenY*winWidth*4 + screenX*4

				pixels[screenIndex] = tex.pixels[texIndex]
				pixels[screenIndex+1] = tex.pixels[texIndex+1]
				pixels[screenIndex+2] = tex.pixels[texIndex+2]
				pixels[screenIndex+3] = tex.pixels[texIndex+3]
			}
		}
	}
}

func (tex *texture) drawAlpha(pixels []byte) {
	for y := 0; y < tex.h; y++ {
		for x := 0; x < tex.w; x++ {
			screenY := y + int(tex.y)
			screenX := x + int(tex.x)
			if screenX >= 0 && screenX < winWidth && screenY >= 0 && screenY < winHeight {
				texIndex := y*tex.pitch + x*4 
				screenIndex := screenY*winWidth*4 + screenX*4

				srcR := int(tex.pixels[texIndex])
				srcG := int(tex.pixels[texIndex+1])
				srcB := int(tex.pixels[texIndex+2])
				srcA := int(tex.pixels[texIndex+3])

				dstR := int(pixels[screenIndex])
				dstG := int(pixels[screenIndex+1])
				dstB := int(pixels[screenIndex+2])

				rstR := (srcR*255 + dstR*(255-srcA)) /255
				rstG := (srcG*255 + dstG*(255-srcA)) /255
				rstB := (srcB*255 + dstB*(255-srcA)) /255

				pixels[screenIndex] = byte(rstR)
				pixels[screenIndex+1] = byte(rstG)
				pixels[screenIndex+2] = byte(rstB)
				//pixels[screenIndex+3] = tex.pixels[texIndex+3]
			}
		}
	}
}

func loadBalloons() []texture {
	balloonStrs := []string{"balloon_red.png", "balloon_blue.png", "balloon_green.png"}
	balloonTextures := make([]texture, len(balloonStrs))
	for i, bstr := range balloonStrs{
		infile, err := os.Open(bstr)
		if err != nil{
			panic(err)
		}
		defer infile.Close()

		img, err := png.Decode(infile)
		if err != nil{
			panic(err)
		}

		w := img.Bounds().Max.X
		h := img.Bounds().Max.Y

		balloonPixels := make([]byte,w*h*4)
		bIndex := 0
		for y := 0; y < h;y++{
			for x := 0; x < w; x++{
				r,g,b,a  := img.At(x,y).RGBA()
				balloonPixels[bIndex] = byte(r/256)
				bIndex++
				balloonPixels[bIndex] = byte(g/256)
				bIndex++
				balloonPixels[bIndex] = byte(b/256)
				bIndex++
				balloonPixels[bIndex] = byte(a/256)
				bIndex++
			}
		}
		balloonTextures[i] = texture{pos{float32(i*60),float32(i*60)}, balloonPixels, w, h, w * 4}
	}
	return balloonTextures
}

func lerp(b1, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct*(float32(b2)-float32( b1)))
}

func colorLerp(c1,c2 rgba, pct float32) rgba {
	return rgba{lerp(c1.r, c2.r, pct), lerp(c1.g, c2.g, pct), lerp(c1.b, c2.b, pct)}
}

func getGradient(c1, c2 rgba) []rgba {
	result := make( []rgba, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}
	return result
}

func getDualGradient(c1, c2, c3, c4 rgba) []rgba {
	result := make( []rgba, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		if pct < 0.5 {			
			result[i] = colorLerp(c1, c2, pct*float32(2))
		} else {
			result[i] = colorLerp(c3, c4, pct*float32(1.5)- float32(0.5))
		}
	}
	return result
}

func clamp(min, max, v int) int {
	if v < min {
		v = min
	}else if v > max {
		v = max
	}
	return v
}

func rescaleAndDraw(noise []float32,min, max float32, gradient []rgba, w, h int) []byte{
	result := make([]byte, w*h*4)
	scale := 255.0 / (max - min)
	offset := min * scale

	for i := range noise {
		noise[i] = noise[i]*scale - offset
		c := gradient[clamp(0,255, int(noise[i]))]
		p := i*4
		result[p] = c.r
		result[p+1] = c.g
		result[p+2] = c.b
	}
	return result
}

func main(){

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

	cloudNoise, min, max := noise.MakeNoise(noise.FBM, 0.009, 0.5, 3, 3, winWidth, winHeight)
	cloudGradient := getGradient(rgba{0,0,255}, rgba{255,255,255})
	cloudPixels := rescaleAndDraw(cloudNoise, min, max, cloudGradient, winWidth, winHeight)
	cloudTexture := texture{pos{0,0}, cloudPixels, winWidth, winHeight, winWidth * 4}

	pixels := make([]byte, winWidth*winHeight*4)
	balloonTextures := loadBalloons()
	dir := 1
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent(){
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		cloudTexture.draw(pixels)
		for _, tex := range balloonTextures {
			tex.drawAlpha(pixels)
		}

		balloonTextures[1].x += float32(1*dir)
		if balloonTextures[1].x > 400 || balloonTextures[1].x < 0 {
			dir = dir * -1
		}
		tex.Update(nil, pixels, winWidth * 4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(16)
	}
}