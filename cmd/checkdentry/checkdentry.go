package main

import (
"fmt"
"os"
"bufio"
"flag"
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
    flag.IntVar(&threshold, "t", 100000000, "Warning threshold of dentry number")
}

func dentry_check(filename string) int {

	fd, err := os.Open(filename)
	
	if err != nil {
		fmt.Fprintf(os.Stdout, "read file %s failed %s\n", filename, err)
		return UNKNOWN
	}
	
	rd := bufio.NewReader(fd)

	var dentry_nr int64 = 0

	data, _, _ :=	rd.ReadLine()
	line := string(data)

	_, _ = fmt.Sscanf(line, "%d", &dentry_nr)

	
	fd.Close()

	if (dentry_nr  >= int64(threshold)) {
		fmt.Fprintf(os.Stdout, "Kernel dentry %d\n", dentry_nr)
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
	filename := procPath + "/sys/fs/dentry-state"

	ret := dentry_check(filename)
	os.Exit(ret)
}
