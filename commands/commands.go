package commands

import (
  "fmt"
  "strings"
  "strconv"
  "os"
  "time"
  "github.com/codegangsta/cli"
  "gopkg.in/vmihailenco/msgpack.v2"
  "silba/probes"
  "silba/pub"
)

func ProcessResults(res <-chan probes.Result, tags []string, zmqch chan<- []byte) {
  for {
    select {
    case r := <-res:
      s := r.GetTypeWithTags(tags)

      fmt.Printf("%v\n", s)
      b, err := msgpack.Marshal(&s)
	    if err != nil {
		     fmt.Println(err)
      }else{
        zmqch <- b
      }

      var v probes.Probe
	    err = msgpack.Unmarshal(b, &v)
	    if err != nil {
		    panic(err)
	     }

      fmt.Printf("Unpack: %#v\n", v)
    }
  }
}

var InitProbes =  func (c *cli.Context) {
  fmt.Printf(" >> Cmd Flags: [port: %s] %s %v [probes %s] [duration %d]\n",
            c.String("port"),
            c.GlobalFlagNames(),
            c.Args(),
            c.GlobalString("probes"),
            c.GlobalDuration("interval"))

  maxCnt, err := strconv.Atoi(c.GlobalString("count"))
  if err != nil {
    maxCnt = 5
  }
  cnt := 1

  // channels
//  tick := make(chan time.Time)
  quit := make(chan bool)
  res := make(chan probes.Result)

  fmt.Printf("zmq: %v", pub.ZeroMQPub)
  zmq := pub.ZeroMQPub.PublishData(quit)
  go ProcessResults(res, []string{"pt101"}, zmq)

  var SelProbes []chan time.Time

  collectFlag := c.GlobalString("probes")
  //fmt.Println(probes.AllProbes)
	for _, f := range strings.Split(collectFlag, ",") {
		p := strings.Trim(f, " ")

    if probe, found := probes.AllProbes[p]; found {
        SelProbes = append(SelProbes, probe.ReadData(res, quit))
    }else{
      fmt.Fprintf(os.Stderr, "Unknown probe option: '%s'\n", p)
    }

	}

  //fmt.Printf("CollectList %v\n", Probes)

  ticker := time.NewTicker(c.GlobalDuration("interval"))

  //show results immediately
  tnow := time.Now()
  for _, ch := range SelProbes {
    //fmt.Println("ticking ", v)
    ch <- tnow
  }

  fmt.Println("ticking stop")

  for {
    select {
    case t := <- ticker.C:
      switch {
        case maxCnt > 0 && cnt >= maxCnt:
          return
        case maxCnt == 0:
          cnt = 0
        default:
          cnt++
      }

      for _, ch := range SelProbes {
        ch <- t
      }

    }
  }
}
