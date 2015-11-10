package main

import (
  "fmt"
  zmq "github.com/pebbe/zmq4"
  "gopkg.in/vmihailenco/msgpack.v2"
  "silba/probes"

)

func main() {
  var v probes.Probe

  subscriber, _ := zmq.NewSocket(zmq.SUB)
  defer subscriber.Close()

  subscriber.Connect("tcp://localhost:10555")
  subscriber.SetSubscribe("")
  fmt.Println("Connected")

  for {
    fmt.Println("Recving")
    //[]byte
    b, _ := subscriber.RecvBytes(0)
    err := msgpack.Unmarshal(b, &v)
    if err != nil {
      fmt.Println("Unpack error: ", err)
    }else{
      fmt.Printf("Unpack: %#v\n", v)
    }
  }
}
