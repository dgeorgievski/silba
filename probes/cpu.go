package probes

import (
  "fmt"
  "time"
  "github.com/cloudfoundry/gosigar"
)

type Cpu struct {
  Probe
}

func (c Cpu) String() string {
  var cpu sigar.Cpu
  cpu = c.data.(sigar.Cpu)
  return fmt.Sprintf("Cpu[etime: %v, tags: %v, User: %d Nice: %d Sys: %d Idle: %d Wait: %d Irq: %d SoftIrq: %d Stolen: %d]",
                          c.etime,
                          c.tags,
                          cpu.User,
                          cpu.Nice,
                          cpu.Sys,
                          cpu.Idle,
                          cpu.Wait,
                          cpu.Irq,
                          cpu.SoftIrq,
                          cpu.Stolen)
}

func (c Cpu) ReadData(r chan<- Result, quit <-chan bool) chan time.Time {
  var cpuUsage sigar.Cpu
  tchan := make(chan time.Time)

  go func() {
      for {
        select {
        case tick := <- tchan:
            previousCpuUsage := cpuUsage
            cpuUsage.Get()
            cpuDelta := cpuUsage.Delta(previousCpuUsage)
            total := cpuDelta.Total()

            var cpuRes sigar.Cpu

            if total > 0 {
              cpuRes.User  = 100 * cpuDelta.User/total
              cpuRes.Nice   = 100 * cpuDelta.Nice/total
              cpuRes.Sys   = 100 * cpuDelta.Sys/total
              cpuRes.Idle  = 100 * cpuDelta.Idle/total
              cpuRes.Wait  = 100 * cpuDelta.Wait/total
              cpuRes.Irq  = 100 * cpuDelta.Irq/total
              cpuRes.SoftIrq  = 100 * cpuDelta.SoftIrq/total
              cpuRes.Stolen  = 100 * cpuDelta.Stolen/total
            }

            //fmt.Println("\t CPU total: ", total, puser)
            c.data = cpuRes
            ncpu := Result{Cpu{Probe{tick,
                        c.tags,
                        c.data}}}
            r <- ncpu
        case <-quit:
            return
        }
      }
    }()

  return tchan
}

func (c Cpu) GetTypeWithTags(t []string) Probe {
  var cpu sigar.Cpu
  cpu = c.data.(sigar.Cpu)
  return Probe{
    etime: c.etime,
    tags: append(c.tags, t...),
    data: sigar.Cpu{
            User: cpu.User,
            Nice: cpu.Nice,
            Sys: cpu.Sys,
            Idle: cpu.Idle,
            Wait: cpu.Wait,
            Irq: cpu.Irq,
            SoftIrq: cpu.SoftIrq,
            Stolen: cpu.Stolen,
          },
    }
}

func init() {
  var cpu sigar.Cpu
  AllProbes["cpu"] = Cpu{Probe{time.Now(), []string{"cpu"}, cpu}}
}
