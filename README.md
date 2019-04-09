# Gojob

A library for build cmdline tools that run remote servers' jobs over job scheduler locally.

Now Gojob support `slurm` job scheduler.

## Install

The library can be install by:
```sh
$ go get github.com/scut-ccmp/gojob
```

I also provide a simple executable binary which can run jobs over slurm scheduler.

You can download binary in here. [slurm_linux_x64]()

## Examples

### Vasp over slurm

Download the binary executable ([slurm_linux_x64]())
 or build the following go [source code](https://github.com/scut-ccmp/gojob/blob/master/examples/slurm.go) code.

Put the following json config file in `$HOME/.config/gojob/`.
Make sure you can connect server.
And create the tempDir in remote server.
```json
{
        "server": {
                "host": "00.0.0.1",
                "port": "22",
                "user": "foo",
                "password": "password"
        },
        "file": {
                "tempDir": "/home/foo/scratch",
                "dirPrefix": "tmp"
        },
        "job": {
                "nproc": "12",
                "ncom": "1",
                "partion": "pri_q",
                "prepend": "module load vasp/5.4.4-impi-mkl",
                "exec": "mpirun -n ${SLURM_NPROCS} vasp_std"
        }

}
```

Put the necessary files of VASP (INCAR, POSCAR, POTCAR, KPOINTS) in your local machine.

By running the executable, it will create the `_job.sh` file and the upload all files in current directory to remote server.
After submitting the job over slurm scheduler, program wait in cmd until the remote job is finish.
Program then download the remote output files to local machine.
