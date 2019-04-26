package processor

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
)

//go:generate counterfeiter . NotifierInterface
type NotifierInterface interface {
	Success(path string, imageNumber int, totalImages int)
	Error(err error, path string, imageNumber int, totalImages int)
}

type Notifier struct {
	Logger *logrus.Logger
}

func NewNotifier() Notifier {
	l := logrus.New()
	logrus.SetOutput(colorable.NewColorableStdout())

	return Notifier{
		Logger: l,
	}
}

func (n Notifier) Success(path string, imageNumber int, totalImages int) {
	n.Logger.WithFields(logrus.Fields{
		"image": fmt.Sprintf("%d/%d", imageNumber, totalImages),
		"input": path,
	}).Info("Processed image")
}

func (n Notifier) Error(err error, path string, imageNumber int, totalImages int) {
	n.Logger.WithFields(logrus.Fields{
		"image": fmt.Sprintf("%d/%d", imageNumber, totalImages),
		"input": path,
	}).Error(err)
}
