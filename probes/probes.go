package probes

import (
  "fmt"
  "time"
  "gopkg.in/vmihailenco/msgpack.v2"
)

var (
	_ msgpack.CustomEncoder = &Probe{}
	_ msgpack.CustomDecoder = &Probe{}
)

type Metric interface {
  ReadData(chan<- Result, <-chan bool) chan time.Time
  GetTypeWithTags([]string) Probe
}

type Result struct {
  Metric
}

type Data interface{}

type Probe struct {
  etime  time.Time
  tags []string
  data Data
}

func (p Probe) String() string {
    return fmt.Sprintf("Result<time: %s, tags: %v, Metric: %v>\n",
        p.etime.Format("03:04:05 PM"), p.tags, p.data)
}

func (p Probe) Get() (time.Time, []string, Data) {
    return p.etime, p.tags, p.data
}

func (p *Probe) EncodeMsgpack(enc *msgpack.Encoder) error {
	return enc.Encode(p.etime, p.tags, p.data)
}

func (p *Probe) DecodeMsgpack(dec *msgpack.Decoder) error {
	return dec.Decode(&p.etime, &p.tags, &p.data)
}

func NewProbe(et time.Time, t []string, d Data) Probe {
  return Probe {
      etime: et,
      tags: t,
      data: d,
  }
}

var ProbeTags []string
var AllProbes = make(map[string]Metric, 10)
