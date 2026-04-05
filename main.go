package main

import (
	feedback "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend"
	history "github.com/ercanaziz/notiFY/Sema-Durgut"
)

func main() {

	feedback.Start()
	history.Start()
	select {}
}
