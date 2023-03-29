package main

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	markdown "github.com/MichaelMure/go-term-markdown"
	"github.com/spf13/pflag"
	"os"
	"path/filepath"
	"reflectsvc/misc"
	"runtime/debug"
	"strconv"
	"strings"
)

//go:embed all:Resources
var efs embed.FS

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
var FlagRemapMap map[string]string
var FlagServiceName string
var FlagPort string
var FlagCert string
var FlagKey string
var FlagDest string
var FlagDestInsecure bool
var FlagHeaderValue []string
var FlagHeaderKey []string

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

	nFlags.StringArrayVarP(&FlagHeaderKey, "header-key", "", []string{"AUTHORIZATION"},
		"Header Key (must be in same order as value)")

	nFlags.StringArrayVarP(&FlagHeaderValue, "header-value", "", []string{"bearer ****DuMmY*ToKeN****="},
		"Header Value(must be in same order as key)")

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
		err = nFlags.Set("destination",
			"https://localhost:"+FlagPort+"/reflect")
		if nil != err {
			logPrintf("Could not reset command line option \"destination\" to \"%s\" because %s",
				"https://localhost:"+FlagPort+"/reflect", err.Error())
			myFatal()
		}
	}

	if FlagDebug && FlagVerbose {
		xLog.Println("\n\t\t/*** program flags ***/\n" +
			"\tplease note that the double backslash is " +
			"an artifact to prevent Windows from corrupting the " +
			"display output. The actual string only has one backslash.\n")
		nFlags.VisitAll(logFlag)
		xLog.Println("\t\t/*** end program flags ***/")
	}

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
		logPrintf("--insecure cannot be used without --debug. DO NOT USE --insecure IN PRODUCTION.")
		myFatal()
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

	if misc.IsStringSet(&FlagRemapFieldNames) {
		FlagRemapMap = loadRemapMap(FlagRemapFieldNames)
	} else {
		FlagRemapMap = make(map[string]string, 0)
	}

}

func loadRemapMap(fn string) map[string]string {
	remap := make(map[string]string, 64)

	f, err := os.Open(fn)
	if nil != err {
		xLog.Printf("Could not open field name conversion file because %s", err.Error())
		myFatal()
	}
	defer misc.DeferError(f.Close)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if "" == line {
			continue
		}
		key, val, err := parseTokens(line)
		if nil != err {
			xLog.Printf("Could not parse line because %s", err.Error())
			continue
		}
		if misc.IsStringSet(&key) && misc.IsStringSet(&val) {
			remap[key] = val
		}
	}
	return remap
}

func parseTokens(line string) (key string, val string, err error) {
	var sbKey, sbVal strings.Builder
	runes := []rune(line)
	err = errors.New("token remap line has bad format: { " + line + " }")

	if len(runes) <= 0 {
		return "", "", nil
	}

	// start a comment, or begin reading key token
	switch runes[0] {
	case '[':
		break
	case '#': // comment
		return "", "", nil
	default:
		return "", "", err
	}
	var ix int
	for ix = 1; ix < len(runes) && ']' != runes[ix]; ix++ {
		sbKey.WriteRune(runes[ix])
	}
	if ix+2 >= len(runes) || runes[ix] != ']' || runes[ix+1] != '[' {
		return "", "", err
	}
	for ix += 2; ix < len(runes) && ']' != runes[ix]; ix++ {
		sbVal.WriteRune(runes[ix])
	}
	if ']' != runes[ix] || (len(runes)-1) != ix {
		return "", "", err
	}
	return sbKey.String(), sbVal.String(), nil
}

func logFlag(flag *pflag.Flag) {
	logPrintf(" flag \"%s\" has value \"%s\" with default %s",
		flag.Name, misc.WinSep(flag.Value.String()), misc.WinSep(flag.DefValue))
}

// UsageMessage - describe capabilities and extended usage notes
func UsageMessage() {
	src, err := efs.ReadFile("Resources/USAGE.MD")
	if nil != err {
		logPrintf("Could not open embedded resource USAGE.MD because %s", err.Error())
		myFatal()
	}
	result := markdown.Render(string(src), 80, 5)
	s2 := string(result)
	fmt.Println(s2)
}
