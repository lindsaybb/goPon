package goPon

import (
	//"bytes"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/secsy/goftp"
)

var buf bytes.Buffer

var LocalOnlyLogs = []string{
	"chrony",
	"license",
	"onu-default",
	"dtrc",
	"ipsgnd.store.dat",
	"ipsg.store.dat",
	"ipsgv6.store.dat",
	"logs",
	"lost+found",
	"lpr.log",
	"mail.err",
	"mail.info",
	"mail.log",
	"mail.warn",
	"news.crit",
	"news.err",
	"news.notice",
	"rabbitmq",
	"samba",
	"sa",
	"tracer",
	"watchdog",
	"xferlog",
}

var OltDirs = []string{
	"/config",
	"/config/onu-default",
	"/log",
	"/bin",
	"/src",
}

/*
goftp: 0.103 #1 sending command FEAT
goftp: 0.104 #1 got 211-Features:
 EPRT
 EPSV
 MDTM
 PASV
 REST STREAM
 SIZE
 TVFS
 UTF8
End

*/

// NewFtpClient wraps the goftp Client config with some path and directory setup
func NewFtpClient(ip, auth string) (*goftp.Client, error) {
	un, pw := parseAuth(auth)
	config := goftp.Config{
		User:               un,
		Password:           pw,
		ConnectionsPerHost: 5,
		Timeout:            3 * time.Second,
		Logger:             &buf, // change this to 'os.Stdout' for debugging
		ActiveTransfers:    false,
	}

	cl, err := goftp.DialConfig(config, ip)
	return cl, err
}

// GetOltLogs is an extension on the goftp.Client object that handles Olt-specific details for retrieving the logs
func GetOltLogs(cl *goftp.Client, path string, tail bool) (int, error) {
	files, err := cl.ReadDir("/log")
	if err != nil {
		return 0, err
	}
	newDir, err := dirInit(path)
	if err != nil {
		return 0, err
	}
	var gotFiles int
	for _, f := range files {
		if !localOnlyLog(f.Name()) {
			locFilePath := filepath.Join(newDir, f.Name())
			remFilePath := fmt.Sprintf("/log/%s", f.Name())
			outFile, err := os.Create(locFilePath)
			if err != nil {
				return gotFiles, err
			}
			err = cl.Retrieve(remFilePath, outFile)
			if err != nil {
				return gotFiles, err
			}
			gotFiles++
			if tail {
				fs, tailBuf, err := tailLogfile(outFile)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("Retrieved %d bytes from %s\n%s\n", fs, outFile.Name(), string(tailBuf))
			} else {
				fmt.Printf("Retrieved: %s\n", outFile.Name())
			}
		}
	}
	return gotFiles, nil
}

func PutOltConfig(cl *goftp.Client, file *os.File) error {
	return cl.Store(fmt.Sprintf("/config/%s", filepath.Base(file.Name())), file)
}

func PutInnboxConfig(cl *goftp.Client, file *os.File) error {
	return cl.Store(fmt.Sprintf("/config/onu-default/%s", filepath.Base(file.Name())), file)
}

func DeleteOltConfig(cl *goftp.Client, file string) error {
	return cl.Delete(fmt.Sprintf("/config/%s", file))
}

func DeleteInnboxConfig(cl *goftp.Client, file string) error {
	return cl.Delete(fmt.Sprintf("/config/onu-default/%s", file))
}

func dirInit(path string) (string, error) {
	ds := generateDatestamp()
	newDir := filepath.Join(path, ds, "log")
	err := os.MkdirAll(newDir, 0755)
	return newDir, err
}

func localOnlyLog(file string) bool {
	for _, v := range LocalOnlyLogs {
		if v == file {
			return true
		}
	}
	return false
}

func generateDatestamp() string {
	t := time.Now()
	tf := strings.Replace(t.Format(time.RFC3339), ":", "-", -1)
	return tf[:len(tf)-6]
}

func tailLogfile(file *os.File) (int64, []byte, error) {
	buf := make([]byte, 8000)
	stat, err := os.Stat(file.Name())
	if err != nil {
		return int64(0), buf, err
	}
	fileSize := stat.Size()
	_, err = file.Read(buf)
	return fileSize, buf, err

}
