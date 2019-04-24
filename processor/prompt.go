package processor

import (
	"fmt"
	"gopkg.in/AlecAivazis/survey.v1"
)

//go:generate counterfeiter . promptInterface
type promptInterface interface {
	ConfirmLargeBatch(size int) bool
}

type Prompt struct {
}

func (Prompt) ConfirmLargeBatch(size int) bool {
	confirmation := false
	message := fmt.Sprintf("Do you want to process %d images?", size)
	prompt := &survey.Confirm{
		Message: message,
	}
	survey.AskOne(prompt, &confirmation, nil)
	return confirmation
}
