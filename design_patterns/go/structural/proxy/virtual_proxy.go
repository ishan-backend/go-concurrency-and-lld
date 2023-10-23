package main

import "fmt"

type Image interface {
	Drawing()
}

type Bitmap struct {
	filename string
}

func (b *Bitmap) Drawing() {
	fmt.Println("drawing image", b.filename)
}

func NewBitmap(filename string) *Bitmap {
	fmt.Println("loading image from filename", filename)
	return &Bitmap{filename: filename}
}

func DrawImage(image Image) { // proxy
	fmt.Println("about to draw image")
	image.Drawing()
	fmt.Println("done drawing image")
}

/*
	Lazy Bitmap - Virtual Proxy
*/
type LazyBitmap struct {
	fileName string
	bitmap   *Bitmap // this bitmap will be lazily constructed
}

func NewLazyBitmap(filename string) *LazyBitmap {
	return &LazyBitmap{fileName: filename} // we don't pass bitmap as it will be constructed lazily i.e. when required
}

func (l *LazyBitmap) Drawing() {
	if l.bitmap == nil {
		l.bitmap = NewBitmap(l.fileName)
	}
	l.bitmap.Drawing()
}

func main() {
	// step 1: load bitmap
	bm := NewBitmap("okay.txt")
	// step 2: draw bitmap
	DrawImage(bm) // proxy

	// ideally, there should be no need to load bm unless you have to draw image
	// LazyBitMap wraps an ordinary BitMap & impls image interface and Draw() method,
	lbm := NewLazyBitmap("okay2.txt")
	// pass lbm into proxy
	DrawImage(lbm)
}
