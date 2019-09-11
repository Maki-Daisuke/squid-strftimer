package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/leekchan/timeutil"
	"github.com/mattn/go-forlines"
)

const defaultPort = 36059
const defaultFormat = `%Y-%m-%dT%H:%M:%S.%fZ`

func formatTime(t time.Time) string {
	return timeutil.Strftime(&t, format)
}

var port int
var format string
var writer io.WriteCloser

func init() {
	flag.IntVar(&port, "port", defaultPort, "TCP port to listen")
	flag.StringVar(&format, "format", defaultFormat, "Strftime-style format")
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		writer = os.Stdout
	} else if flag.NArg() != 1 {
		flag.CommandLine.Usage()
		os.Exit(1)
	} else {
		var err error
		writer, err = os.OpenFile(flag.Arg(0), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("Can't open `%s': %s", flag.Arg(0), err)
			os.Exit(1)
		}
		defer writer.Close()
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Printf("Can't listen port %d: %s", port, err)
		os.Exit(1)
	}
	log.Printf("Listening port %d...", port)

	// Monitor signal and cancel the context when SIGINT or SIGTERM comes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		listener.Close()
		cancel()
	}()

	wg := &sync.WaitGroup{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			if e, ok := err.(*net.OpError); ok {
				if e.Err.Error() == "use of closed network connection" {
					// The passive socket is closed because of a signal, and so, end the loop.
					break
				}
			}
			log.Panic(err)
		}
		log.Print("Accepted a connection")
		wg.Add(1)
		go handleConnection(ctx, wg, conn)
	}
	log.Print("Shutting down... ")
	wg.Wait()
	log.Print("Stopped")
}

func handleConnection(ctx context.Context, wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()

	c, cancel := context.WithCancel(ctx)
	go func() {
		<-c.Done() // Wait until either of handleConnection or main indicates to terminate this context.
		conn.Close()
	}()
	defer cancel() // This ensures that the above goroutine closes conn after the following loop ends.

	forlines.Must(conn, func(line string) error {
		return handleLine(line)
	})
	log.Print("Disconnected from ", conn.RemoteAddr())
}

var reTimestamp = regexp.MustCompile(`^(\d+)\.(\d+)`)

func handleLine(line string) error {
	match := reTimestamp.FindStringSubmatchIndex(line)
	if match == nil {
		fmt.Fprintln(writer, line)
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
	t := time.Unix(sec, int64(time.Duration(msec)*time.Millisecond/time.Nanosecond))
	fmt.Fprintf(writer, "%s%s\n", formatTime(t), line[match[1]:])
	return nil
}
