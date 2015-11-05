package probes

import (
  "fmt"
  "time"
  "github.com/cloudfoundry/gosigar"
)

type Uptime struct {
  Probe
}

func (m Uptime) String() string {
  var upt sigar.Uptime
  upt = m.data.(sigar.Uptime)
  return fmt.Sprintf("Uptime[etime: %v, tags: %v, Length: %f]",
                                    m.etime,
                                    m.tags,
                                    upt.Length)
}

func (m Uptime) ReadData(r chan<- Result, quit <-chan bool) chan time.Time {
  var uptime sigar.Uptime
  tchan := make(chan time.Time)

  go func() {
    for {
      select {
      case tick := <- tchan:
          uptime.Get()
          m.data = uptime
          r <- Result{Uptime{Probe{ tick,
                      m.tags,
                      m.data}}}

      case <-quit:
          return
      }
    }
  }()

  return tchan
}

func (m Uptime) GetTypeWithTags(t []string) Probe {
  var upt sigar.Uptime
  upt = m.data.(sigar.Uptime)
  return Probe{
      etime: m.etime,
      tags: append(m.tags, t...),
      data:  sigar.Uptime{
                Length: upt.Length,
              },
      }
}

func init() {
  var upt sigar.Uptime
  AllProbes["uptime"] = Uptime{Probe{time.Now(),
                                    []string{"uptime"},
                                    upt}}
}
