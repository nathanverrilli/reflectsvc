package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"reflectsvc/misc"
	"runtime/debug"
	"strconv"
	"strings"
)

// wordSepNormalizeFunc all options are lowercase, so
// ... lowercase they shall be
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	/**********************************************
	from := []string{"-", "_"}
	to := "."
	if nil != f {
		for _, sep := range from {
			name = strings.Replace(name, sep, to, -1)
		}
	}
	********************************************/
	return pflag.NormalizedName(strings.ToLower(name))
}

var nFlags *pflag.FlagSet

/* secret flags */
var FlagOrganization string

/* standard flags */

var FlagHelp bool
var FlagQuiet bool
var FlagVerbose bool
var FlagDebug bool

/* program specific flags */
var FlagPort string
var FlagLogAppend bool

func initFlags() {
	var err error

	hideFlags := make(map[string]string, 8)

	nFlags = pflag.NewFlagSet("default", pflag.ContinueOnError)
	nFlags.SetNormalizeFunc(wordSepNormalizeFunc)

	nFlags.BoolVarP(&FlagDebug, "debug", "d",
		true, "Enable additional informational and operational logging output for debug purposes")

	nFlags.BoolVarP(&FlagVerbose, "verbose", "v",
		true, "Supply additional run messages; use --debug for more information")

	nFlags.BoolVarP(&FlagHelp, "help", "h",
		false, "Display help message and usage information")

	nFlags.BoolVarP(&FlagQuiet, "quiet", "q",
		false, "Suppress output to stdout and stderr (output still goes to logfile)")

	// secret flags

	nFlags.StringVarP(&FlagOrganization, "organization", "",
		"P3IDTechnologies", "organization for email sends")
	hideFlags["FlagOrganization"] = "organization"

	// program flags
	nFlags.StringVarP(&FlagPort, "port", "",
		"9090", "port to listen on")

	for flagName, optName := range hideFlags {
		err = nFlags.MarkHidden(optName)
		if nil != err {
			xLog.Printf("could not mark option %s as %s hidden because %s\n",
				optName, flagName, err.Error())
			myFatal()
		}
	}

	// Fetch and load the program flags
	err = nFlags.Parse(os.Args[1:])
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "\n%s\n", nFlags.FlagUsagesWrapped(75))
		xLog.Fatalf("\nerror parsing flags because: %s\n%s %s\n%s\n\t%v\n",
			err.Error(),
			"  common issue: 2 hyphens for long-form arguments,",
			"  1 hyphen for short-form argument",
			"  Program arguments are: ",
			os.Args)
	}

	// do quietness setup first
	// only write to logfile not stderr
	// for debug and verbose messages
	if FlagQuiet {
		xLog.SetOutput(xLogBuffer)
		// messages only to logfile, not stderr
	}

	if FlagDebug && FlagVerbose {
		xLog.Println("\n\t\t/*** program flags ***/\n" +
			"\tplease note that the double backslash is " +
			"an artifact to prevent Windows from corrupting the " +
			"display output. The actual string only has one backslash.\n")
		nFlags.VisitAll(logFlag)
		xLog.Println("\t\t/*** end program flags ***/")
	}

	// next simplest
	if FlagHelp {
		var err1, err2 error
		_, thisCmd := filepath.Split(os.Args[0])
		_, err1 = fmt.Fprint(os.Stdout, "\n", "usage for ", thisCmd, ":\n")
		_, err2 = fmt.Fprintf(os.Stdout, "%s\n", nFlags.FlagUsagesWrapped(75))
		if nil != err1 || nil != err2 {
			xLog.Printf("huh? can't write to os.stdout because\n%s",
				misc.ConcatenateErrors(err1, err2).Error())
		}
		UsageMessage()
		myFatal(0)
	}

	if FlagVerbose {
		errMsg := ""
		user, host, err := misc.UserHostInfo()
		if nil != err || nil == misc.SafeString(&errMsg) {
			errMsg = " (encountered error " + err.Error() + ")"
		}
		xLog.Printf("Verbose mode active (all debug and informative messages) for %s@%s%s",
			user, host, errMsg)
	}

	if FlagDebug && FlagVerbose {
		_, exeName := filepath.Split(os.Args[0])
		exeName = strings.TrimSuffix(exeName, filepath.Ext(exeName))
		bi, ok := debug.ReadBuildInfo()
		if !ok {
			xLog.Printf("huh? Could not read build information for %s "+
				"-- perhaps compiled without module support?", exeName)
		} else {
			xLog.Printf("\n***** %s BuildInfo: *****\n%s\n%s\n",
				exeName, bi.String(), strings.Repeat("*", 22+len(exeName)))
		}

	}

	// sanity check for FlagPort
	portNumber, err := strconv.Atoi(FlagPort)
	if nil != err {
		xLog.Printf("Got bad value for --port: %s (must be an integer) [error: %s]",
			FlagPort, err.Error())
		myFatal()
	}
	if FlagVerbose || FlagDebug {
		xLog.Printf("Listening on port %d", portNumber)
	}

}

func logFlag(flag *pflag.Flag) {
	xLog.Printf(" flag %s has value %s with default %s",
		flag.Name, misc.WinSep(flag.Value.String()), misc.WinSep(flag.DefValue))
}

// UsageMessage - describe capabilities and extended usage notes
func UsageMessage() {
	var sb strings.Builder
	sb.WriteString(" cpAuthOrg: \n")
	sb.WriteString(" --\n")
	sb.WriteString("Useful information goes here\n")

}
