package probes

import (
  "fmt"
  "time"
  "github.com/cloudfoundry/gosigar"
)

type Mem struct {
  Probe
}

func (m Mem) String() string {
  var mem  sigar.Mem
  mem = m.data.(sigar.Mem)
  return fmt.Sprintf("Mem[etime: %v, tags: %v, Total: %d Used: %d Free: %d ActualFree: %d ActualUsed: %d]",
                          m.etime,
                          m.tags,
                          mem.Total,
                          mem.Used,
                          mem.Free,
                          mem.ActualFree,
                          mem.ActualUsed)
}

func (m Mem) ReadData(r chan<- Result, quit <-chan bool) chan time.Time {
  var mem sigar.Mem
  tchan := make(chan time.Time)

  go func() {
    for {
      select {
      case tick := <- tchan:
          //fmt.Println("total: ", total, puser)
          mem.Get()
          m.data = mem
          r <- Result{Mem{Probe{tick, m.tags, m.data}}}

      case <-quit:
          return
      }
    }
  }()

  return tchan
}

func (m Mem) GetTypeWithTags(t []string) Probe {
  var mem sigar.Mem
  mem = m.data.(sigar.Mem)

  return Probe{
    etime: m.etime,
    tags: append(m.tags, t...),
    data:  sigar.Mem{
            Total: mem.Total,
            Used: mem.Used,
            Free: mem.Free,
            ActualFree: mem.ActualFree,
            ActualUsed: mem.ActualUsed,
          },
  }
}

func init() {
  var mem sigar.Mem
  AllProbes["mem"] = Mem{Probe{time.Now(),
                                []string{"mem"},
                                mem}}

}
