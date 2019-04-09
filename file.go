package gojob

import (
	"sync"
	"time"
	"os"
	"strconv"
	"path/filepath"
	"io/ioutil"
	"io"
	"path"

	"github.com/pkg/sftp"
)

var rand uint32
var randmu sync.Mutex

func reseed() uint32 {
	return uint32(time.Now().UnixNano() + int64(os.Getpid()))
}

func nextRandom() string {
	randmu.Lock()
	r := rand
	if r == 0 {
		r = reseed()
	}
	r = r*16644525 + 1013904223 // constants from Numerical Recipes
	rand = r
	randmu.Unlock()
	return strconv.Itoa(int(1e9 + r%1e9))[1:]
}

func TempDir(client *sftp.Client, dir, prefix string) (name string, err error) {
	if dir == "" {
		dir = os.TempDir()
	}
	nconflict := 0

	sshFxFailure := uint32(4)

	for i := 0; i < 10000; i++ {
		try := filepath.Join(dir, prefix+nextRandom())
		err = client.Mkdir(try)
		if status, ok := err.(*sftp.StatusError); ok {
			if status.Code == sshFxFailure {
				if nconflict++; nconflict > 10 {
					randmu.Lock()
					rand = reseed()
					randmu.Unlock()
				}
				continue
			}
			return "", err
		}
		name = try
		break
	}
	return name, nil
}

func SendFiles(client *sftp.Client, fromDir, toDir string) error {
	// loop over pwd files
	files, err := ioutil.ReadDir(fromDir)
	if err != nil {
		return err
	}
	Trace.Printf("Start send following files:\n")
	for _, f := range files {
		// create source file
		srcFile, err := os.Open(path.Join(fromDir, f.Name()))
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// create destination file
		dstFile, err := client.Create(path.Join(toDir, f.Name()))
		if err != nil {
			return err
		}
		defer dstFile.Close()
		// copy source file to destination file
		bytes, err := io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
		Trace.Printf("%s: %d bytes copied\n", f.Name(), bytes)
	}
	Info.Printf("Finished send files\n")
	return nil
}

func ReciveFiles(client *sftp.Client, fromDir, toDir string) error {
	// retrieve all files from the remote directory
	files, err := client.ReadDir(fromDir)
	if err != nil {
		return err
	}
	Trace.Printf("Start recive following files:\n")
	for _, f := range files {
		// create destination file
		var dstFile *os.File
		filename := path.Join(toDir, f.Name())
		dstFile, err = os.Create(filename)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		// open source file
		srcFile, err := client.Open(path.Join(fromDir, f.Name()))
		if err != nil {
			return err
		}

		// copy source file to destination file
		bytes, err := io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
		Trace.Printf("%d bytes copied\n", bytes)

		// flush in-memory copy
		err = dstFile.Sync()
		if err != nil {
			return err
		}
	}
	Info.Printf("Finished recive files\n")
	return nil
}
