package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const syntaxError = 3
const runtimeError = 1

func usage(out io.Writer) {
	s := "Writes to a file or stdout in a pattern depending on several parameters.\n" +
		"By default, writes one 256k block to stdout\n" +
		"Each block is prefixed with the iteration number (starting at 0) and is filled with zeros\n" +
		" -f \tFile: file path to write to.\n" +
		" -s \tSize: How many bytes to write each iteration. Suffix with k or m for kilobytes or megabytes.\n" +
		" -c \tCount: How many writes to try before quitting.\n" +
		" -d \tDelay: How many seconds to delay between writes. Suffix with ms, m, or h.\n" +
		" -od\tOpen Delay: How many seconds to delay before opening the file. Suffix with ms, m, or h. Ignored\n" +
		"    \twithout -f.\n" +
		" -t \tTimeout: How many seconds (not counting Start Delay) to run before closing the file, unless Count is\n" +
		"	 \treached first. Suffix with ms, m, or h.\n" +
		" -sd\tStart Delay: How many seconds to delay before the first write. Suffix with ms, m, or h.\n" +
		" -l \tLog File: Filename to log to.\n" +
		" -rc\tReturn Code: Return code on successful exit. Will be overridden by any errors. Default is 0." +
		" -h \tHelp: Prints this text\n" +
		"Returns 0 on success, 1 if a runtime error is encountered, 3 if bad arguments are passed.\n"
	fmt.Fprintf(out, s)
	os.Exit(0)
}

func handleError(err error, out io.Writer, rc int) {
	if rc == syntaxError {
		fmt.Fprintf(out, "%v\nDo 'writer -h' for usage\n", err)
	} else {
		fmt.Fprintf(out, "%v\n", err)
	}
	os.Exit(rc)
}

func main() {
	var fileName string
	var size = 256 * 1024
	var count = 1
	var delay, openDelay, startDelay, timeout time.Duration
	var log = io.Discard
	var rc int

	var skip = true
	for arg := range os.Args {
		if skip {
			skip = false
			continue
		}
		switch os.Args[arg] {
		case "-f":
			arg += 1
			skip = true
			fileName = os.Args[arg]
		case "-l":
			arg += 1
			skip = true
			l, err := os.Create(os.Args[arg])
			if err != nil {
				handleError(fmt.Errorf("could not open log file: %v", err), log, syntaxError)
			}
			log = l
		case "-s":
			arg += 1
			skip = true
			var mult = 1
			var s string
			if strings.HasSuffix(os.Args[arg], "k") {
				mult = 1024
				s = strings.TrimSuffix(os.Args[arg], "k")
			} else if strings.HasSuffix(os.Args[arg], "m") {
				mult = 1024 * 1024
				s = strings.TrimSuffix(os.Args[arg], "m")
			} else {
				s = os.Args[arg]
			}
			sz, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				handleError(fmt.Errorf("invalid argument for size '%s': %v", s, err), log, syntaxError)
			}
			size = int(sz) * mult
		case "-c":
			arg += 1
			skip = true
			var err error
			c, err := strconv.ParseInt(os.Args[arg], 10, 32)
			if err != nil {
				handleError(fmt.Errorf("invalid argument for size '%s': %v", os.Args[arg], err), log, syntaxError)
			}
			count = int(c)
		case "-d":
			arg += 1
			skip = true
			var mult = time.Second
			var s string
			if strings.HasSuffix(os.Args[arg], "ms") {
				mult = time.Millisecond
				s = strings.TrimSuffix(os.Args[arg], "ms")
			} else if strings.HasSuffix(os.Args[arg], "m") {
				mult = time.Minute
				s = strings.TrimSuffix(os.Args[arg], "m")
			} else if strings.HasSuffix(os.Args[arg], "h") {
				mult = time.Hour
				s = strings.TrimSuffix(os.Args[arg], "h")
			} else {
				s = os.Args[arg]
			}
			t, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				handleError(fmt.Errorf("invalid argument for delay '%s': %v", s, err), log, syntaxError)
			}
			delay = time.Duration(t) * mult
		case "-od":
			arg += 1
			skip = true
			var mult = time.Second
			var s string
			if strings.HasSuffix(os.Args[arg], "ms") {
				mult = time.Millisecond
				s = strings.TrimSuffix(os.Args[arg], "ms")
			} else if strings.HasSuffix(os.Args[arg], "m") {
				mult = time.Minute
				s = strings.TrimSuffix(os.Args[arg], "m")
			} else if strings.HasSuffix(os.Args[arg], "h") {
				mult = time.Hour
				s = strings.TrimSuffix(os.Args[arg], "h")
			} else {
				s = os.Args[arg]
			}
			t, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				handleError(fmt.Errorf("invalid argument for open delay '%s': %v", s, err), log, syntaxError)
			}
			openDelay = time.Duration(t) * mult
		case "-t":
			arg += 1
			skip = true
			var mult = time.Second
			var s string
			if strings.HasSuffix(os.Args[arg], "ms") {
				mult = time.Millisecond
				s = strings.TrimSuffix(os.Args[arg], "ms")
			} else if strings.HasSuffix(os.Args[arg], "m") {
				mult = time.Minute
				s = strings.TrimSuffix(os.Args[arg], "m")
			} else if strings.HasSuffix(os.Args[arg], "h") {
				mult = time.Hour
				s = strings.TrimSuffix(os.Args[arg], "h")
			} else {
				s = os.Args[arg]
			}
			t, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				handleError(fmt.Errorf("invalid argument for timeout '%s': %v", s, err), log, syntaxError)
			}
			timeout = time.Duration(t) * mult
		case "-sd":
			arg += 1
			skip = true
			var mult = time.Second
			var s string
			if strings.HasSuffix(os.Args[arg], "ms") {
				mult = time.Millisecond
				s = strings.TrimSuffix(os.Args[arg], "ms")
			} else if strings.HasSuffix(os.Args[arg], "m") {
				mult = time.Minute
				s = strings.TrimSuffix(os.Args[arg], "m")
			} else if strings.HasSuffix(os.Args[arg], "h") {
				mult = time.Hour
				s = strings.TrimSuffix(os.Args[arg], "h")
			} else {
				s = os.Args[arg]
			}
			t, err := strconv.ParseInt(s, 10, 32)
			if err != nil {
				handleError(fmt.Errorf("invalid argument for start delay '%s': %v", s, err), log, syntaxError)
			}
			startDelay = time.Duration(t) * mult
		case "-rc":
			arg += 1
			skip = true
			returnCode, err := strconv.Atoi(os.Args[arg])
			if err != nil {
				handleError(fmt.Errorf("Invalid return code '%s': %v", os.Args[arg], err), log, syntaxError)
			}
			rc = returnCode
		case "-h":
			usage(log)
		default:
			handleError(fmt.Errorf("Invalid argument '%s'\n", os.Args[arg]), log, syntaxError)
		}
	}

	buf := make([]byte, size)
	var bytes int
	var output io.Writer
	var err error
	if fileName == "" {
		output = os.Stdout
	} else {
		time.Sleep(openDelay)
		output, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			handleError(err, log, runtimeError)
		}
	}
	time.Sleep(startDelay)
	start := time.Now()
	var runtime time.Duration
	for itr := 0; itr != count; itr += 1 {
		runtime = time.Since(start)
		if runtime >= timeout && timeout != 0 {
			break
		}
		b, err := output.Write(buf)
		bytes += b
		if err != nil {
			fmt.Fprintf(log, "Error encountered while writing: %v\n")
			rc = runtimeError
			break
		}
		time.Sleep(delay)
	}
	var bytesWritten string
	if bytes > 1024*1024*10 {
		bytesWritten = fmt.Sprintf("%d megabytes", bytes/1024/1024)
	} else if bytes > 1024*10 {
		bytesWritten = fmt.Sprintf("%d kilobytes", bytes/1024)
	} else {
		bytesWritten = fmt.Sprintf("%d bytes", bytes)
	}
	fmt.Fprintf(log, "Wrote %s in %s\n", bytesWritten, runtime.String())
	os.Exit(rc)
}
