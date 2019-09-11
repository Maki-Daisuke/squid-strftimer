package main

import (
	"strconv"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/leekchan/timeutil"
	"github.com/mattn/go-forlines"
)

const defaultFormat = `[%Y-%m-%dT%H:%M:%S.%fZ]`

func formatTime(t time.Time) string {
	return timeutil.Strftime(&t, defaultFormat)
}

var reTimestamp = regexp.MustCompile(`^(\d+)\.(\d+)`)

func main(){
	forlines.Must(os.Stdin, func(line string)error{
		match := reTimestamp.FindStringSubmatchIndex(line)
		if match == nil {
			fmt.Println(line)
			return nil
		}
		sec, err := strconv.ParseInt(line[match[2]:match[3]], 10, 64)
		if err != nil {
			return err
		}
		msec, err := strconv.ParseInt(line[match[4]:match[5]], 10, 64)
		if err != nil {
			return err
		}
		t := time.Unix(sec, int64(time.Duration(msec) * time.Millisecond / time.Nanosecond))
		fmt.Printf("%s%s\n", formatTime(t), line[match[1]:])
		return nil
	})
}
