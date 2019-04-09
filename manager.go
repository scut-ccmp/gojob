package gojob

import (
	"time"
)

type Manager interface{
  FindJobState(c *Conn, id int) (state string, err error)
  SubmitJob(c *Conn, path string) (id int, err error)
  CheckDoneFunc(c *Conn, id int) (done bool, err error)
}

func JobManager(name string) Manager {
  table := map[string]Manager {
    "slurm": NewSlurmMgt(),
  }
  return table[name]
}

func Check(f func(*Conn, int) (bool, error), conn *Conn, id int) {
	timeout := time.After(10000 * time.Second)
	tick := time.Tick(5 * time.Second)

	for {
		select {
		case <- timeout:
			Error.Println("timeout!")
			return
		case <- tick:
			done, err := f(conn, id)
			if err != nil {
				c := conn.Connect.Conn
				// cInfo := c.ConnMetadata
				Error.Printf("Check process state error, conn: %v, id: %d", c, id)
			}
			if done {
				Info.Println("Remote Job DONE")
				return
			}
		}
	}
}
