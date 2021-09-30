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
	s := "Reads data from a file or stdin and writes it to a file or stdout in a pattern depending on several parameters.\n" +
		"By default, reads from stdin and writes to stdout in 256k blocks until EOF\n" +
		" -i \tInput file: file path to read from.\n" +
		" -o \tOutput file: file path to write to.\n" +
		" -s \tSize: How many bytes to attempt to read and write each iteration. Suffix with k or m for kilobytes or\n" +
		"    \tmegabytes.\n" +
		" -c \tCount: How many iterations to try before quitting, unless EOF is reached first.\n" +
		" -d \tDelay: How many seconds to delay between iterations. Suffix with ms, m, or h.\n" +
		" -od\tOpen Delay: How many seconds to delay before opening the files. Suffix with ms, m, or h. Ignored without\n" +
		"    \t-o or -i.\n" +
		" -t \tTimeout: How many seconds (not counting Start Delay) to run before quitting, unless Count is reached\n" +
		"	 \tfirst. Suffix with ms, m, or h.\n" +
		" -sd\tStart Delay: How many seconds to delay before beginning to read and write. Suffix with ms, m, or h.\n" +
		" -l \tLog File: Filename to log to.\n" +
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
	var inFile, outFile string
	var size = 256 * 1024
	var count int
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
		case "-i":
			arg += 1
			skip = true
			inFile = os.Args[arg]
		case "-o":
			arg += 1
			skip = true
			outFile = os.Args[arg]
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
		case "-h":
			usage(log)
		default:
			handleError(fmt.Errorf("Invalid argument '%s'\n", os.Args[arg]), log, syntaxError)
		}
	}

	buf := make([]byte, 0, size)
	var bytesIn, bytesOut int
	var output = os.Stdout
	var input = os.Stdin
	var err error
	if inFile != "" || outFile != "" {
		time.Sleep(openDelay)
		if outFile != "" {
			output, err = os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE, 0755)
			if err != nil {
				handleError(err, log, runtimeError)
			}
		}
		if inFile != "" {
			input, err = os.OpenFile(inFile, os.O_RDONLY|os.O_CREATE, 0755)
			if err != nil {
				handleError(err, log, runtimeError)
			}
		}
	}
	time.Sleep(startDelay)
	start := time.Now()
	var runtime time.Duration
	var eof bool
	x := 0
	for {
		runtime = time.Since(start)
		if runtime >= timeout && timeout != 0 {
			break
		}
		b, err := input.Read(buf)
		bytesIn += b
		if err != nil && err != io.EOF {
			fmt.Fprintf(log, "Error encountered while reading: %v\n")
			rc = runtimeError
			break
		} else if err == io.EOF {
			eof = true
		}
		b, err = output.Write(buf)
		bytesOut += b
		if err != nil {
			fmt.Fprintf(log, "Error encountered while writing: %v\n")
			rc = runtimeError
			break
		}
		if eof {
			break
		}
		if count == 0 {
			continue
		}
		x += 1
		if x == count {
			break
		}
		time.Sleep(delay)
	}
	var bytesRead string
	if bytesIn > 1024*1024*10 {
		bytesRead = fmt.Sprintf("%d megabytes", bytesIn/1024/1024)
	} else if bytesIn > 1024*10 {
		bytesRead = fmt.Sprintf("%d kilobytes", bytesIn/1024)
	} else {
		bytesRead = fmt.Sprintf("%d bytes", bytesIn)
	}
	var bytesWritten string
	if bytesOut > 1024*1024*10 {
		bytesWritten = fmt.Sprintf("%d megabytes", bytesOut/1024/1024)
	} else if bytesOut > 1024*10 {
		bytesWritten = fmt.Sprintf("%d kilobytes", bytesOut/1024)
	} else {
		bytesWritten = fmt.Sprintf("%d bytes", bytesOut)
	}
	fmt.Fprintf(log, "Read %s and wrote %s in %s\n", bytesRead, bytesWritten, runtime.String())
	os.Exit(rc)
}
