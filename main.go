package main

import (
	feedback "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend"
	history "github.com/ercanaziz/notiFY/Sema-Durgut"
)

func main() {

	go feedback.Start()
	history.Start()
	select {}
}
