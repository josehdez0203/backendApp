package logger

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func L_Info(m string) {
	infoColor := color.New(color.FgHiBlue).SprintFunc()
	fechaColor := color.New(color.FgYellow).SprintFunc()
	fecha := time.Now()
	d := fecha.Format(time.DateOnly)
	t := fecha.Format(time.TimeOnly)
	fmt.Printf("%s %s INFO: %s\n", fechaColor(d), fechaColor(t), infoColor(m))
}

func L_Error(m string) {
	errorColor := color.New(color.FgRed).SprintFunc()
	fechaColor := color.New(color.FgHiRed).SprintFunc()
	fecha := time.Now()
	d := fecha.Format(time.DateOnly)
	t := fecha.Format(time.TimeOnly)
	res := errorColor(" ERROR ")
	fmt.Printf("%s %s %s %s\n", fechaColor(d), fechaColor(t), res, errorColor(m))
}
