package main

import (
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

const exe = "/master"
const cf = "/master.ini"

var cfgVar = regexp.MustCompile(`(?s)\AMASIF_MASTER_(\w+)_(\w+)=(.*)\z`)

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		cmd := exec.Command(exe)
		cmd.Stdout = os.Stderr
		cmd.Stderr = os.Stderr

		if errRn := cmd.Run(); errRn != nil {
			if errEx, ok := errRn.(*exec.ExitError); !(ok && errEx.ExitCode() > -1) {
				log.Fatal(errRn.Error())
			}
		}
	}

	{
		cfg := ini.Empty()
		log.Info("Translating env vars like MASIF_MASTER_*_* to .ini config")

		for _, ev := range os.Environ() {
			if match := cfgVar.FindStringSubmatch(ev); match != nil {
				_, errNK := cfg.Section(strings.ToLower(match[1])).NewKey(strings.ToLower(match[2]), match[3])
				if errNK != nil {
					log.Fatal(errNK.Error())
				}
			}
		}

		if errST := cfg.SaveTo(cf); errST != nil {
			log.Fatal(errST.Error())
		}
	}

	log.Info("Starting actual daemon via exec(3)")

	errEx := syscall.Exec(exe, []string{exe, "--config", cf}, os.Environ())
	if errEx == nil {
		errEx = syscall.Errno(0)
	}

	log.Fatal(errEx.Error())
}
