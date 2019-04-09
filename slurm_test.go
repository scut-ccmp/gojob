package gojob

import (
  "testing"
)

func TestparsegojobID(t *testing.T) {
  tests := []struct {
    in []byte
    expected int
  }{
    {[]byte(`Submitted batch job 54264`), 54264},
    {[]byte(`Submitted batch job 99999`), 99999},
  }

  m := NewSlurmMgt()
  for _, test := range tests {
    if got := m.parseJobID(test.in); got != test.expected {
      t.Errorf("FindJobID(%s) got %d", string(test.in), got)
    }
  }
}

func TestParseJobState(t *testing.T) {
  tests := []struct {
    in []byte
    expected string
  }{
    {[]byte(`STATE COMPLETING`), "COMPLETING"},
    {[]byte(`STATE RUNNING`), "RUNNING"},
    {[]byte(`foobar STATE RUNNING`), "RUNNING"},
    {[]byte(`STATE`), "NOJOBFOUND"},
    {[]byte(`slurm_load_jobs error: Invalid job id specified`), "NOJOBFOUND"},
  }

  m := NewSlurmMgt()
  for _, test := range tests {
    if got, _ := m.parseJobState(test.in); got != test.expected {
      t.Errorf("JobState(%s) got %s", string(test.in), got)
    }
  }
}
