package gojob

const SlurmTmpl = `#!/bin/bash -l
#SBATCH -J {{.Name}}
#SBATCH -p {{.Partion}}
#SBATCH -N {{.NCom}}
#SBATCH -n {{.NProc}}
#SBATCH -o _stdout.out
#SBATCH -e _stderr.out

{{.Prepend}}

{{.ExecCmd}}

{{.Append}}
`
