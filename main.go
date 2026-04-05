package main

import (
    feedback "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend"
    "sync"
)

func main() {
       
    go feedback.Start()
    select {}   
}      
       