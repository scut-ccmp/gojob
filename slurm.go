package gojob

import (
	"regexp"
	"strconv"
	"fmt"
)

type SlurmParameter struct {
	Name string
	NProc string
	NCom string
	Partion string
	Prepend string
	ExecCmd string
	Append string
}

type SlurmMgt struct {}

func NewSlurmMgt() SlurmMgt {
	return SlurmMgt{}
}


func (m SlurmMgt) parseJobState(s []byte) (state string, err error) {
	matched, err := regexp.Match(`Invalid job id specified`, s)
	if err != nil {
		return "", fmt.Errorf("parseJobState: invalid job id: %v", err)
	}
	if matched {
		return "NOJOBFOUND", nil
	}

  es := `(STATE)\s*(RUNNING|SUSPENDE|COMPLETING|PENDING)*`
	re := regexp.MustCompile(es)
	result := re.FindSubmatch(s)
	match := result[2]
	state = string(match)
	if state == "" {
		return "NOJOBFOUND", nil
	}
	return state, nil
}

func (m SlurmMgt) FindJobState(conn *Conn, id int) (string, error) {
	cmd := fmt.Sprintf(`squeue -j %d -o "%%.8T"`, id)
	sess, err := conn.Session()
	if err != nil {
		return "", fmt.Errorf("find job state: sess %v", err)
	}
	defer sess.Close()
	out, err := sess.Output(cmd)
	if err != nil {
		return "", fmt.Errorf("find job state: sess run %v", err)
	}
	s, err := m.parseJobState(out)
	if err != nil {
		return "", fmt.Errorf("FindJobState: parseJobState: %v", err)
	}
	return s, nil
}

func (m SlurmMgt) parseJobID(s []byte) (id int) {
  es := `^Submitted batch job (?P<jobID>\d+)`
	re := regexp.MustCompile(es)
	result := re.FindSubmatch(s)
	match := result[1]
	id, _ = strconv.Atoi(string(match))
	return id
}

func (m SlurmMgt) SubmitJob(conn *Conn, path string) (id int, err error) {
	cmd := "cd " + path + ";sbatch _job.sh"
	sess, err := conn.Session()
	if err != nil {
		return -1, fmt.Errorf("submit job: new sess: %v", err)
	}
	defer sess.Close()
	output, err := sess.Output(cmd)
	if err != nil {
		return -1, fmt.Errorf("submit job: run sess: %v", err)
	}
	id = m.parseJobID(output)
	return id, nil
}

func (m SlurmMgt) CheckDoneFunc(conn *Conn, id int) (done bool, err error) {
		state, err := m.FindJobState(conn, id)
		if err != nil {
			return true, fmt.Errorf("CheckDoneFunc: %v", err)
		}
		// fmt.Println(state)
		if state == "NOJOBFOUND" {
			return true, nil
		}
		return false, nil
}
