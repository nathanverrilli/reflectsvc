package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflectsvc/misc"
	"runtime"
	"strings"
	"sync"
)

var xLogFile *os.File
var xLogBuffer *bufio.Writer
var xLog log.Logger

// flushLog just flushes the log to disk
var flmx sync.Mutex

func flushLog() {
	flmx.Lock()
	defer flmx.Unlock()
	if nil != xLogBuffer {
		err := xLogBuffer.Flush()
		if nil != err {
			safeLogPrintf("huh? could not flush xLogBuffer because %s", err.Error())
		}
	}
}

var clmx sync.Mutex

// closeLog shuts the logging service down
// cleanly, flushing buffers (and thus
// preserving the most likely error of
// interest)
func closeLog() {
	clmx.Lock()
	lpmx.Lock()
	var err01, err02 error
	if nil != xLogBuffer {
		flushLog()
		xLogBuffer = nil
	}
	if nil != xLogFile {
		err01 = xLogFile.Close()
		xLogFile = nil
	}
	lpmx.Unlock()
	clmx.Unlock()
	err := misc.ConcatenateErrors(err01, err02)
	if nil != err {
		safeLogPrintf(err.Error())
	}

}

// initLog starts up a logging service to logfile and
// stderr. FlagQuiet (--quiet) will shut off a great deal of
// console output, but the log is never suppressed.
func initLog(lfName string) {
	var err error
	var logWriters = make([]io.Writer, 0, 2)
	xLogFile, err = os.OpenFile(lfName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if nil != err {
		safeLogPrintf("error opening log file %s because %s",
			lfName, err.Error())
	}

	xLogBuffer = bufio.NewWriter(xLogFile)
	logWriters = append(logWriters, os.Stderr)
	logWriters = append(logWriters, xLogBuffer)
	xLog.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
	xLog.SetOutput(io.MultiWriter(logWriters...))

	logPath, err := filepath.Abs(xLogFile.Name())
	if nil != err {
		logPrintf("huh? could not resolve logfilename %s because %s",
			xLogFile.Name(), err.Error())
		myFatal()
	}
	logPrintf("Logfile set to %s", logPath)
}

// myFatal is meant to close the program, and close the
// log files properly. Go doesn't support optional arguments,
// but variadic arguments allow finessing this. myFatal() gets
// a default RC of -1, and that's overridden by the first int
// in the slice of integers argument (which is present
// even if the length is 0).
//
// At some point, might create a more
// thorough at-close routine and register closing the file
// and log as part of the things to do 'at close'.
func myFatal(rcList ...int) {
	rc := -1
	if len(rcList) > 0 {
		rc = rcList[0]
	}

	// if this is an expected exit, and FlagQuiet is set,
	// this doesn't need to be logged
	if rc != 0 {
		_, srcFile, srcLine, ok := runtime.Caller(1)
		if ok {
			srcFile = filepath.Base(srcFile)
			logPrintf("\n\t\t/*** myFatal called ***/\n"+
				"\tfrom file:line %12s:%04d\n"+
				"\t\t/*** myFatal ended ***/", srcFile, srcLine)
		} else {
			logPrintf("\n\t\t/*** myFatal called ***/\n" +
				"\tbut could not get stack information for caller\n" +
				"\t\t/*** myFatal ended ***/")
		}
	}
	closeLog()
	os.Exit(rc)
}

var lpmx sync.Mutex

func logPrintf(format string, a ...any) {
	lpmx.Lock()
	defer lpmx.Unlock()
	safeLogPrintf(format, a...)
}

// safeLogPrintf may be called in lieu of xLog.Printf() if there
// is a possibility the log may not be open. If the log is
// available, well and good. Otherwise, print the message to
// STDERR.
func safeLogPrintf(format string, a ...any) {

	if nil != xLogBuffer && nil != xLogFile {
		xLog.Printf(format, a...)
	} else {
		_, _ = fmt.Fprintf(os.Stderr,
			"\n\tSAFELOG\n"+format+"\n",
			a...)
	}
}

func debugMapStringString(params map[string]string) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("\n\tgot map[string]string size %d\n", len(params)))
	for k, v := range params {
		sb.WriteString(fmt.Sprintf("\t[ %-20s ][ %-20s ]\n", k, v))
	}
}
