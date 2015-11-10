package pub

import (
  zmq "github.com/pebbe/zmq4"
  "fmt"
)


type ZeroMQ struct {
  ip_address string
  port int
}

var ZeroMQPub ZeroMQ

func (z ZeroMQ) PublishData(quit <-chan bool) (chan<- []byte) {

  fmt.Println("PublishData ZeroMQ")

  publish := make(chan []byte)


  go func() {

    publisher, _ := zmq.NewSocket(zmq.PUB)
    defer publisher.Close()

    tcp_conn := fmt.Sprintf("tcp://%s:%d", z.ip_address, z.port)
    publisher.Bind(tcp_conn)

    for {
      select {
      case d := <-publish:
        //publish data.
        publisher.SendBytes(d, 0)

      case <-quit:
          return
        }
    }
  }()

  return publish
}

func (zmq ZeroMQ) Configure(ip string, port int) {

  zmq.ip_address  = ip
  zmq.port        = port
}

func (zmq ZeroMQ) String() string {
  return fmt.Sprintf("[ZMQ ip:%s port: %d]\n", zmq.ip_address, zmq.port)
}

func init() {
  //default
  ZeroMQPub = ZeroMQ{"0.0.0.0", 10555}
}
