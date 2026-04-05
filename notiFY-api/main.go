package main

import (
    "github.com/ercanaziz/notiFY/Ercan-Aziz/Backend"
    "github.com/ercanaziz/notiFY/Sema-Durgut"
    "github.com/ercanaziz/notiFY/Dogukan-Dursoy"
    "github.com/ercanaziz/notiFY/Betul-Erkoc"
    "github.com/ercanaziz/notiFY/Nisanur-Sutcu"
)

func main() {
    go product.Start()     // fiyat çekme
    go feedback.Start()         // HTTP API
    go worker.Start()      // Kafka/RabbitMQ worker
    go scheduler.Start()   // zamanlayıcı
    go notifier.Start()       // bildirim (son çalışan, block eder)
}