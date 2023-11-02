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
func wordSepNormalizeFunc(_ *pflag.FlagSet, name string) pflag.NormalizedName {

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

var FlagRemapFieldNames string

var FlagServiceName string
var FlagPort string
var FlagCert string
var FlagKey string
var FlagDest string
var FlagDestInsecure bool
var FlagTick bool
var FlagProxySuccess bool

//var FlagHeaderValue []string
//var FlagHeaderKey []string

func initFlags() {
	var err error

	hideFlags := make(map[string]string, 8)

	nFlags = pflag.NewFlagSet("default", pflag.ContinueOnError)
	nFlags.SetNormalizeFunc(wordSepNormalizeFunc)

	nFlags.BoolVarP(&FlagDebug, "debug", "d",
		false, "Enable additional informational and operational logging output for debug purposes")

	nFlags.BoolVarP(&FlagVerbose, "verbose", "v",
		false, "Supply additional run messages; use --debug for more information")

	nFlags.BoolVarP(&FlagHelp, "help", "h",
		false, "Display help message and usage information")

	nFlags.BoolVarP(&FlagQuiet, "quiet", "q",
		false, "Suppress output to stdout and stderr (output still goes to logfile)")

	// secret flags

	nFlags.StringVarP(&FlagOrganization, "organization", "",
		"P3IDTechnologies", "organization for email sends")
	hideFlags["FlagOrganization"] = "organization"

	// program flags

	nFlags.BoolVarP(&FlagProxySuccess, "proxy-success", "", false,
		"force all proxied xm2json requests to return an explicit success 200 status")

	nFlags.BoolVarP(&FlagTick, "tick", "", false, "enable a console tick every few seconds")

	nFlags.StringVarP(&FlagRemapFieldNames, "fieldNames", "", "",
		"Filename of conversion mapping, one pair per line, [oldName][newName], escape '[' and ']' by doubling them '[[' and ']]'. Case sensitive.")

	nFlags.BoolVarP(&FlagDestInsecure, "insecure", "", false,
		"Accesses the remote server without checking the remote "+
			"certificate's validity. THIS IS FOR TESTING PURPOSES ONLY. DO "+
			"NOT RUN WITH --insecure IN PRODUCTION.")

	nFlags.StringVarP(&FlagDest, "destination", "",
		"localhost",
		"destination for Xml2Json endpoint. "+
			"the value 'localhost' is a special value that "+
			"becomes \"https://localhost:<port>/reflect\" where"+
			"<port> is the port of this program.")

	/*
		nFlags.StringArrayVarP(&FlagHeaderKey, "header-key", "", []string{"AUTHORIZATION"},
			"Header Key (must be in same order as value)")

		nFlags.StringArrayVarP(&FlagHeaderValue, "header-value", "", []string{"bearer ****DuMmY*ToKeN****="},
			"Header Value(must be in same order as key)")
	*/

	nFlags.StringVarP(&FlagServiceName, "servername", "", "",
		"Name of service/FQDN \"microservice.example.com\" <<not fully tested>> ")

	nFlags.StringVarP(&FlagPort, "port", "",
		"9090", "port to listen on")

	nFlags.StringVarP(&FlagKey, "keyfile", "", "",
		"Key file for HTTPS service")

	nFlags.StringVarP(&FlagCert, "certfile", "", "",
		"Certificate file for HTTPS service")

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

	if FlagDest == "localhost" {
		var hmode string
		if misc.IsStringSet(&FlagCert) {
			hmode = "https"
		} else {
			hmode = "http"
		}
		err = nFlags.Set("destination",
			hmode+"://localhost:"+FlagPort+"/reflect")
		if nil != err {
			xLog.Printf("Could not reset command line option \"destination\" to \"%s\" because %s",
				"https://localhost:"+FlagPort+"/reflect", err.Error())
			myFatal()
		}
	}

	if FlagDebug && FlagVerbose {
		xLog.Println("\t\t/*** start program flags ***/\n")
		nFlags.VisitAll(logFlag)
		xLog.Println("\t\t/***   end program flags ***/")
	}
	/*
			if len(FlagHeaderKey) != len(FlagHeaderValue) {
				logPrintf("count of --header-key values (%d) does not equal count of --header-value (%d)",
					len(FlagHeaderKey), len(FlagHeaderValue))
				myFatal()
			}

			if FlagVerbose || FlagDebug {
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("\n%s\n\tCommand Line Headers\n", SEP))
				for ix := range FlagHeaderKey {
					sb.WriteString(fmt.Sprintf("[header %3d] %s=%s\n",
						ix, FlagHeaderKey[ix], FlagHeaderValue[ix]))
				}
				logPrintf(sb.String())
			}

		if FlagDestInsecure && !FlagDebug {
			xLog.Printf("--insecure cannot be used without --debug. DO NOT USE --insecure IN PRODUCTION.")
			myFatal()
		}
	*/
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
		_, _ = fmt.Fprintf(os.Stdout, "\t please see USAGE.MD for ")
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

	if misc.IsStringSet(&FlagRemapFieldNames) {
		FlagRemapMap = loadFieldTranslations(FlagRemapFieldNames)
	} else {
		FlagRemapMap = make(map[string]remapField)
	}

}

// logFlag -- This writes out to the logger the value of a
// particular flag. Called indirectly. `Write()` is used
// directly to prevent wierd interactions with backslash
// in filenames
func logFlag(flag *pflag.Flag) {
	var sb strings.Builder
	sb.WriteString(" flag ")
	sb.WriteString(flag.Name)
	sb.WriteString(" has value [")
	sb.Write([]byte(flag.Value.String()))
	sb.WriteString("] with default [")
	sb.Write([]byte(flag.DefValue))
	sb.WriteString("]\n")
	_, _ = xLog.Writer().Write([]byte(sb.String()))
}

// UsageMessage - describe capabilities and extended usage notes
func UsageMessage() {
	xLog.Printf("Please see README.MD for usage information")
}
