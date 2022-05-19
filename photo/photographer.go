package photo

import (
	_ "embed"
	"fmt"
	"github.com/areknoster/attendgo/domain"
	pigo "github.com/esimov/pigo/core"
	"image"
	"math"
	"time"
)

type Capturer interface {
	Capture() domain.Photo
}

//go:embed classifiers/facefinder
var faceFinderClassifier []byte

type FacePhotographer struct {
	pigo     *pigo.Pigo
	capturer Capturer
	pub      domain.Publisher
}

func NewFacePhotographer(capturer Capturer, pub domain.Publisher) (*FacePhotographer, error) {
	pg := pigo.NewPigo()
	pg, err := pg.Unpack(faceFinderClassifier)
	if err != nil {
		return nil, fmt.Errorf("unpack pigo classifier: %w", err)
	}
	return &FacePhotographer{
		pigo:     pg,
		capturer: capturer,
		pub:      pub,
	}, nil

}

// Globals are awful but no reason to bother in this project
var (
	PhotoInterval               = time.Second
	MaxAttempts                 = 20
	DetectionCertaintyThreshold = float32(0.8)
)

func (fp *FacePhotographer) StartCapturing() {
	for i := 0; i < MaxAttempts; i++ {
		photo := fp.capturer.Capture()
		switch fp.facesNumber(photo.Img) {
		case 0:
			fp.pub.Publish(domain.ErrNoFace)
		case 1:

			fp.pub.Publish(domain.EventFacePhotoTaken{Photo: photo})
			return
		default:
			fp.pub.Publish(domain.ErrTooManyFaces)
		}
		time.Sleep(PhotoInterval)
	}
	fp.pub.Publish(domain.ErrFacePhotoNotTaken)
}

func (fp *FacePhotographer) facesNumber(img image.Image) int {
	pigoImg := pigo.RgbToGrayscale(img)
	cols, rows := img.Bounds().Max.X, img.Bounds().Max.Y

	detections := fp.pigo.RunCascade(pigo.CascadeParams{
		MinSize:     cols / 5,
		MaxSize:     int(math.Min(float64(cols), float64(rows)) * 0.9),
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: pigoImg,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}, 0.0)
	detections = fp.pigo.ClusterDetections(detections, 0.2)
	faceNumber := 0
	for _, detection := range detections {
		if detection.Q > DetectionCertaintyThreshold {
			faceNumber++
		}
	}

	return faceNumber
}
