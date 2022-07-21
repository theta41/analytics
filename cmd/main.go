package main

import (
	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/analytics/internal/app"
)

func main() {
	a := app.NewApp()
	err := a.Run()
	if err != nil {
		logrus.Panic(err)
	}
}
