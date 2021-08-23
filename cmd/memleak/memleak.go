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
    percentThreshold   int
)

func init() {
    flag.StringVar(&procPath, "p", "/proc", "actual path of /proc")
    flag.IntVar(&percentThreshold, "t", 15, "Warning threshold of percentage of memory use")
}

func memleak_check(filename string) int {

	fd, err := os.Open(filename)

	if err != nil {
		fmt.Fprintf(os.Stdout, "read file %s failed %s\n", filename, err)
		return UNKNOWN
	}

	rd := bufio.NewReader(fd)

	var unrecSlab int64 = 0
	var memTotal int64 = 0
	name, kb := "", ""

	for {
		data, _, eof :=	rd.ReadLine()
		line := string(data)
		if eof == io.EOF {
			break
		}

		if unrecSlab > 0 &&  memTotal > 0 {
			break
		}
		if strings.Contains(line, "SUnreclaim") {
			_, _ = fmt.Sscanf(line, "%s %d %s", &name, &unrecSlab, &kb)
			continue
		}
		if strings.Contains(line, "MemTotal") {
			_, _ = fmt.Sscanf(line, "%s %d %s", &name, &memTotal, &kb)
			continue
		}
	}

	fd.Close()

	/*SUnreclaim >= memTotal*15% and SUnreclaim >= 800M */
	if (unrecSlab * 100  >= memTotal * int64(percentThreshold)) && (unrecSlab >= 800 * 1024) {
		fmt.Fprintf(os.Stdout, "Kernel memleak SUnreclaim:%d MemTotal:%d\n", unrecSlab, memTotal)
		return NONOK
	}
	return OK
}

func main ()  {
	flag.Parse()
	if percentThreshold >= 100 || percentThreshold <=0 {
        fmt.Fprintf(os.Stderr, "value of -t must between 0 and 100")
        os.Exit(UNKNOWN)
    }
	filename := procPath + "/meminfo"

	ret := memleak_check(filename)
	os.Exit(ret)
}
