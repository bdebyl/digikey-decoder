package main

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/blackjack/webcam"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
)

func ProcessAndPrint(frame []byte, width uint32, height uint32) {
	dmtx := datamatrix.NewDataMatrixReader()
	yuyv := image.NewYCbCr(image.Rect(0, 0, int(width), int(height)), image.YCbCrSubsampleRatio422)
	for i := range yuyv.Cb {
		ii := i * 4
		yuyv.Y[i*2] = frame[ii]
		yuyv.Y[i*2+1] = frame[ii+2]
		yuyv.Cb[i] = frame[ii+1]
		yuyv.Cr[i] = frame[ii+3]
	}
	img := yuyv
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	result, _ := dmtx.Decode(bmp, nil)

	if result != nil {
		partno := strings.Split(result.String(), "\x1D1P")[1]
		partno = strings.Split(partno, "\x1DK\x1D")[0]
		fmt.Println(partno)
	}
}

func main() {
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()

	format := webcam.PixelFormat(1448695129)
	//framesize := webcam.FrameSize{MinWidth: 640, MaxWidth: 640, MinHeight: 480, MaxHeight: 480}
	/*
		cmap := cam.GetControls()
		for id, c := range cmap {
			fmt.Printf("ID:%08x %-32s  Min: %4d  Max: %5d\n", id, c.Name, c.Min, c.Max)
		}
	*/
	// AutoFocus
	//cam.SetControl(webcam.ControlID(0x009a090c), 1)
	// AutoExposure
	//cam.SetControl(webcam.ControlID(0x009a0903), 1)

	// BacklightComp
	cam.SetControl(webcam.ControlID(0x0098091c), 1)
	// WhiteBal
	cam.SetControl(webcam.ControlID(0x0098091a), 6000)
	// Focus
	cam.SetControl(webcam.ControlID(0x009a090a), 250)
	// Gain
	cam.SetControl(webcam.ControlID(0x00980913), 85)
	// Brightness
	cam.SetControl(webcam.ControlID(0x00980900), 128)
	// Exposure
	cam.SetControl(webcam.ControlID(0x009a0902), 600)

	_, width, height, err := cam.SetImageFormat(format, 640, 480)
	if err != nil {
		panic(err.Error())
	}

	err = cam.StartStreaming()
	if err != nil {
		panic(err.Error())
	}
	for {
		err = cam.WaitForFrame(5)

		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			fmt.Fprint(os.Stderr, err.Error())
			continue
		default:
			panic(err.Error())
		}

		frame, err := cam.ReadFrame()
		if len(frame) != 0 {
			// Process frame
			ProcessAndPrint(frame, width, height)
		} else if err != nil {
			panic(err.Error())
		}
	}
}
