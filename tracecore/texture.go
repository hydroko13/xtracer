package tracecore

import (

	"image"
	"image/color"
	_ "image/png"
	"math"
	"os"
)

type Texture struct {
	_data *[]uint8
	width int
	height int
}



func LoadTexture(filePath string) (*Texture, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return &Texture{}, err
	}

	defer f.Close()

	
	
	img, _, err := image.Decode(f)

	if err != nil {
		return &Texture{}, err
	}

	w, h := img.Bounds().Dx(), img.Bounds().Dy()

	pixData := make([]uint8, w*h*4)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			i := (x*h+y)*4
			clr := img.At(x, y)
			r, g, b, a := clr.RGBA()
			pixData[i] = uint8(r)
			pixData[i+1] = uint8(g)
			pixData[i+2] = uint8(b)
			pixData[i+3] = uint8(a)
		}
	}

	return &Texture{
		&pixData,
		w,
		h,
	}, nil

}

func (tex *Texture) Width() int {
	return tex.width
}

func (tex *Texture) Height() int {
	return tex.height
}

func (tex *Texture) Sample(x float64, y float64) color.RGBA {

	ix := max(min(int(math.Floor(x*float64(tex.width))), tex.width-1), 0)
	iy := max(min(int(math.Floor(y*float64(tex.height))), tex.height-1), 0)

	
	i := (ix*tex.height+iy)*4
	

	return color.RGBA{(*tex._data)[i], (*tex._data)[i+1], (*tex._data)[i+2], (*tex._data)[i+3]}
}