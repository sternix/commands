package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"syscall"
)

const (
	MFLAG = 1 << iota
	NFLAG
	PFLAG
	RFLAG
	SFLAG
	VFLAG
	IFLAG
	KFLAG
)

func main() {
	allFlag := flag.Bool("a", false, "")
	identFlag := flag.Bool("i", false, "")
	platformFlag := flag.Bool("m", false, "")
	hostnameFlag := flag.Bool("n", false, "")
	osnameFlag := flag.Bool("o", false, "")
	archFlag := flag.Bool("p", false, "")
	releaseFlag := flag.Bool("r", false, "")
	sysnameFlag := flag.Bool("s", false, "")
	versionFlag := flag.Bool("v", false, "")
	kernversFlag := flag.Bool("K", false, "")
	helpFlag := flag.Bool("h", false, "")

	flag.Parse()

	var flags int = 0

	if *allFlag {
		flags |= (MFLAG | NFLAG | RFLAG | SFLAG | VFLAG)
	}

	if *identFlag {
		flags |= IFLAG
	}

	if *osnameFlag || *sysnameFlag {
		flags |= SFLAG
	}

	if *platformFlag {
		flags |= MFLAG
	}

	if *hostnameFlag {
		flags |= NFLAG
	}

	if *archFlag {
		flags |= PFLAG
	}

	if *releaseFlag {
		flags |= RFLAG
	}

	if *versionFlag {
		flags |= VFLAG
	}

	if *kernversFlag {
		flags |= KFLAG
	}

	if *helpFlag {
		usage()
	}

	if flags == 0 {
		flags |= SFLAG
	}

	printUname(flags)
}

var isPrinted bool

func printFlag(flags, flg int, envVar, sysctlName string, fn func(string, string) string) {
	if flags&flg == flg {
		if isPrinted {
			fmt.Print(" ")
		}

		fmt.Print(fn(envVar, sysctlName))
		isPrinted = true
	}
}

func printUname(flags int) {
	printFlag(flags, SFLAG, "s", "kern.ostype", getValue)
	printFlag(flags, NFLAG, "n", "kern.hostname", getValue)
	printFlag(flags, RFLAG, "r", "kern.osrelease", getValue)
	printFlag(flags, VFLAG, "v", "kern.version", getVersion)
	printFlag(flags, MFLAG, "m", "hw.machine", getValue)
	printFlag(flags, PFLAG, "p", "hw.machine_arch", getValue)
	printFlag(flags, IFLAG, "i", "kern.ident", getValue)
	printFlag(flags, KFLAG, "K", "kern.osreldate", getIntValue)

	fmt.Println()
}

func isEnvSet(flg string) (bool, string) {
	if val := os.Getenv("UNAME_" + flg); val != "" {
		return true, val
	} else {
		return false, ""
	}
}

func getSysctl(sysctlName string) (val string) {
	var err error
	if val, err = syscall.Sysctl(sysctlName); err != nil {
		fmt.Fprintf(os.Stderr, "%s - %v\n", sysctlName, err)
		os.Exit(1)
	}

	return val
}

func getSysctlUint32AsString(sysctlName string) (val string) {
	if ret, err := syscall.SysctlUint32(sysctlName); err != nil {
		fmt.Fprintf(os.Stderr, "%s - %v", sysctlName, err)
		os.Exit(1)
	} else {
		val = strconv.FormatUint(uint64(ret), 10)
	}

	return
}

func getValue(envVar, sysctlName string) string {
	if yes, envVal := isEnvSet(envVar); yes {
		return envVal
	} else {
		return getSysctl(sysctlName)
	}
}

func getIntValue(envVar, sysctlName string) string {
	if yes, envVal := isEnvSet(envVar); yes {
		return envVal
	} else {
		return getSysctlUint32AsString(sysctlName)
	}
}

func getVersion(envVar, sysctlName string) string {
	if yes, envVal := isEnvSet(envVar); yes {
		return envVal
	} else {
		version := []byte(getSysctl(sysctlName))
		for i, r := range version {
			if r == '\n' || r == '\t' {
				version[i] = ' '
			}
		}

		return string(version)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: uname [-aiKmnoprsUv]\n")
	os.Exit(1)
}
