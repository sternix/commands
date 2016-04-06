package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

import (
	"github.com/sternix/commands/lib/sysctl"
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

var (
	allFlag      = flag.Bool("a", false, "Behave as though the options -m, -n, -r, -s, and -v were specified.")
	identFlag    = flag.Bool("i", false, "Write the kernel ident to standard output.")
	platformFlag = flag.Bool("m", false, "Write the type of the current hardware platform to standard output.")
	hostnameFlag = flag.Bool("n", false, "Write the name of the system to standard output.")
	osnameFlag   = flag.Bool("o", false, "This is a synonym for the -s option, for compatibility with other systems.")
	archFlag     = flag.Bool("p", false, "Write the type of the machine processor architecture to standard output.")
	releaseFlag  = flag.Bool("r", false, "Write the current release level of the operating system to standard output.")
	sysnameFlag  = flag.Bool("s", false, "Write the name of the operating system implementation to standard output.")
	versionFlag  = flag.Bool("v", false, "Write the version level of this release of the operating system to standard output.")
	kernversFlag = flag.Bool("K", false, "Write the FreeBSD version of the kernel.")
)

func init() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var flags int

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

func getSysctl(sysctlName string) string {
	if val, err := sysctl.ByName(sysctlName); err != nil {
		log.Fatalf("%s - %v\n", sysctlName, err)
		return "" //NOTREACHED
	} else {
		return val
	}
}

func getSysctlUint32AsString(sysctlName string) string {
	if ret, err := sysctl.Uint32(sysctlName); err != nil {
		log.Fatalf("%s - %v", sysctlName, err)
		return "" //NOTREACHED
	} else {
		return strconv.FormatUint(uint64(ret), 10)
	}
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
	log.Fatalln("usage: uname [-aiKmnoprsUv]")
}
