package main

import (
"fmt"
"os"
"io"
"bufio"
"flag"
"strings"
)

const (
    OK      = 0
    NONOK   = 1
    UNKNOWN = 2
)

var (
    procPath           string
    threshold   int
)

func init() {
    flag.StringVar(&procPath, "p", "/proc", "actual path of /proc")
    flag.IntVar(&threshold, "t", 1000, "Warning threshold of cgroup number")
}

func cgroup_check(filename string) int {

	fd, err := os.Open(filename)

	if err != nil {
		fmt.Fprintf(os.Stdout, "read file %s failed %s\n", filename, err)
		return UNKNOWN
	}

	rd := bufio.NewReader(fd)

	var cpu int = 0
	var mem int = 0
	name, hierarchy := "",0
	enabled := 0

	for {
		data, _, eof :=	rd.ReadLine()
		line := string(data)
		if eof == io.EOF {
			break
		}

		if cpu > 0 &&  mem > 0 {
			break
		}
		if strings.Contains(line, "memory") {
			_, _ = fmt.Sscanf(line, "%s %d %d %d", &name, &hierarchy, &mem, &enabled)
			continue
		}
		if strings.Contains(line, "cpuacct") {
			_, _ = fmt.Sscanf(line, "%s %d %d %d", &name, &hierarchy, &cpu, &enabled)
			continue
		}
	}
	fmt.Printf("mem %d cpu %d\n", mem, cpu)
	fd.Close()

	if (mem  >= threshold) || (cpu >= threshold) {
		fmt.Fprintf(os.Stdout, "Kernel cgroup mem:%d cpu:%d\n", mem, cpu)
		return NONOK
	}
	return OK
}

func main ()  {
	flag.Parse()
	if threshold <=0 {
        fmt.Fprintf(os.Stderr, "value of -t must Don't 0")
        os.Exit(UNKNOWN)
    }
	filename := procPath + "/cgroups"

	ret := cgroup_check(filename)
	os.Exit(ret)
}
