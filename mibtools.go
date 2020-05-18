package mibtools

// [2020-05-18] 1505 
// S:\MiGo\mibGOPATH\src\mimodules\callmibtools\mibtools\mibtools.go

// Usage:
// - Set path to config file
// - Populate config object with values from config file
// - Set path to Logs
// - Get config key value 
// - Add temporary key to config object
// - Write to Log

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	"time"
	"errors"
	//"strings"
)

// Version current version
func Version() string {
	return "mibtools: MAY-18-2020"
} 

// ConfigEntry key/value pair model
type ConfigEntry struct {
	Key string `json:"Key"`
	Value string `json:"Value"`
}

// appConfigSlice slice of ConfigEntry key/value pair model
type appConfigSlice []ConfigEntry

// AppConfig app config
var AppConfig *appConfigSlice = new(appConfigSlice) //config object

// Item single config key
var Item *ConfigEntry

// DefaultPathToConfig Default Path To Config 
var DefaultPathToConfig string = "./config/app.config.json"

// DefaultPathToLogKey Default Path To Log Key
var DefaultPathToLogKey string = "PathToLog"  // config key for path to log

var pathToConfig string = ""
var pathToLog string = ""


// SetPathToLog set internal path to Log 
func SetPathToLog (path string) {
	pathToLog = path
}

// getPathToLog private method
func getPathToLog() string {
	return pathToLog
}

// SetPathToConfig Set path to config file
func SetPathToConfig(path string) {
	pathToConfig = path
}

// getPathToConfig private method
func getPathToConfig() string {
	return pathToConfig
}

//______________________________________________________ ErrorLevel

// ErrorLevel error code
type ErrorLevel int
// Error levels names
const (
  OK ErrorLevel = 0 + iota
	INFO
	WARNING
	RETRY
	ERROR
  EXCEPTION
)
// ErrorLevelDescription description of error level
var ErrorLevelDescription = [...]string{
  "OK",
  "Information",
  "Warning",
  "Network error or Timeout",
  "Error",
  "System Exception",
}
func (errc ErrorLevel) String() string { 
	return ErrorLevelDescription[errc] 
}

//______________________________________________________ Error structure

// ErrorItem my error item interface 
type ErrorItem struct {
	FunctionName string
	Code  ErrorLevel
	Details string
}
func (e *ErrorItem) Error() string {
	var returnedValue string
	returnedValue = e.FunctionName + ". " + e.Code.String() + ". " + e.Details
  return returnedValue
}

//___________________________________________________________



// InitTools Initialize related objects
func InitTools(pathToConfig string) bool {
	var returnedValue bool = false
	var actualpathToConfig, actualpathToLog string

	if _, err := os.Stat(pathToConfig); err == nil {
		actualpathToConfig = pathToConfig
	} else {
		fmt.Println("Invalid pathToConfig")
		return returnedValue
	}

	// Set path to config
	SetPathToConfig(actualpathToConfig)

	// Read and set path to log
	actualpathToLog = ReadConfigKey(DefaultPathToLogKey)
	if actualpathToLog == "" {
		fmt.Println("Path to Log is not set")
		return returnedValue
	}
	// Does Folder Exist
	//_ => folderInfo
	_, err := os.Stat(actualpathToLog)
	if os.IsNotExist(err) {
		fmt.Println("Invalid path to Log: " + actualpathToLog)
		return returnedValue
	}

	SetPathToLog(actualpathToLog)
	
	err2 := PopulateAppConfig(AppConfig, actualpathToConfig)
	if err2 != nil {
		fmt.Println(err2)
		return returnedValue
	} 
	returnedValue = true
	
	return returnedValue
}


// ReadConfigKey return value by key name provided, from config file at private pathToConfig
func ReadConfigKey(key string) string {
	var returnedValue string = ""
	data, err := ioutil.ReadFile(pathToConfig)
	if err != nil {
			fmt.Println("File reading error", err, pathToConfig)
			return returnedValue
	}
	var target appConfigSlice
	if err := json.Unmarshal([]byte(data), &target); err != nil {
		fmt.Println("Error:", err)
		return returnedValue
	}
	for i := 0; i < len(target); i++ {
		if target[i].Key == key {
			returnedValue = target[i].Value
		}
	}
	return returnedValue
}

// ReadConfigKeyWithPath return value by key name provided, from config file at path provided
func ReadConfigKeyWithPath(path string, key string) string {
	var returnedValue string = ""
	data, err := ioutil.ReadFile(path)
	if err != nil {
			fmt.Println("File reading error", err, path)
			return returnedValue
	}
	var target appConfigSlice
	if err := json.Unmarshal([]byte(data), &target); err != nil {
		fmt.Println("Error:", err)
		return returnedValue
	}
	for i := 0; i < len(target); i++ {
		if target[i].Key == key {
			returnedValue = target[i].Value
		}
	}
	return returnedValue
}

// GetConfigKeyValue return value by key name provided, from AppConfig slice
func GetConfigKeyValue(key string, app *appConfigSlice) string {
	var returnedValue string
	returnedValue = "NONE"
   for _, currItem := range *app {
		if currItem.Key == key {
			returnedValue = currItem.Value
			break
		}
	}
 	return returnedValue
}

// PopulateAppConfig populates AppConfig slice from config file at private member pathToConfig.
// If pathToConfig is not previously set - optionalpath is used, if it is valid.
func PopulateAppConfig(app *appConfigSlice, optionalpath string) (reterr error) {
	reterr = nil
	var item *ConfigEntry
	var actualpath string
	var errormessage string
	
	if pathToConfig == "" {
		if optionalpath == "" {
			reterr = errors.New("Both pathToConfig and optionalpath are empty")
			return reterr
		}
		actualpath = optionalpath
	} else {
		actualpath = pathToConfig
	}
	data, err := ioutil.ReadFile(actualpath)
	if err != nil {
		errormessage = fmt.Sprintf("File reading error. %s %s", err, actualpath)
		reterr = errors.New(errormessage)
		return reterr
	}
	var target appConfigSlice
	if err := json.Unmarshal([]byte(data), &target); err != nil {
		errormessage = fmt.Sprintf("Error. %s", err)
		reterr = errors.New(errormessage)
		return reterr
	}
	for i := 0; i < len(target); i++ {
		item = new(ConfigEntry)
	  item.Key = target[i].Key
	  item.Value = target[i].Value
		*app = append(*app, *item)
	}
	return reterr
}

// PutConfigKeyValue add key/value pair to AppConfig slice
func PutConfigKeyValue(key string, value string, app *appConfigSlice) {
	var item *ConfigEntry
	item = new(ConfigEntry)
	item.Key = key
	item.Value = value
	if app == nil {
		app = new(appConfigSlice)
	}
	*app = append(*app, *item)
}

// LogError report error details
func LogError(e ErrorItem) {
	var message string
	message = e.Code.String() + ". " + e.Details
	WriteToLog(e.FunctionName, message)
}

// WriteToLog write header and message to Log file at private member pathToLog
func WriteToLog(header, message string) (reterr error) {
  reterr = nil
  if pathToLog == "" {
    reterr = errors.New("PathToLog is not set")
		return reterr
	}
	var prefix, layout, filename, fullpath string
	layout = "2006-01-02 15:04:05"
	current := time.Now()
	prefix = current.Format(layout) + ": " + header + "\n"  //time.RFC1123 - closest
	filename = "LOG_" + prefix[0:10] + ".txt"
	fullpath = pathToLog + "\\" + filename

	outfile, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		reterr = fmt.Errorf("PathToLog is not set %s", err)
		return reterr
	}

	ln, err := outfile.WriteString(prefix)
	if err != nil || ln == 0 {
		outfile.Close()
		reterr = fmt.Errorf("Unable to write to file %s. %s", pathToLog, err)
		return reterr
	}
	outfile.WriteString(message + "\n")

	err = outfile.Close()
	if err != nil {
		reterr = fmt.Errorf("Unable to close file: %s. %s", pathToLog, err)	
		return reterr
	}
  return reterr
}

// WriteToLogWithPath write header and message to Log file at path provided
func WriteToLogWithPath(path, header, message string) (reterr error) {
  reterr = nil
	var prefix, layout, filename, fullpath string
	layout = "2006-01-02 15:04:05"
	current := time.Now()
	prefix = current.Format(layout) + ": " + header + "\n"  //time.RFC1123 - closest format
	filename = "LOG_" + prefix[0:10] + ".txt"
	
	fullpath = path + "\\" + filename

	outfile, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		reterr = fmt.Errorf("PathToLog is not set %s", err)
		return reterr
	}

	l, err := outfile.WriteString(prefix)
	if err != nil || l == 0 {
		outfile.Close()
		reterr = fmt.Errorf("Unable to write to file %s. %s", pathToLog, err)
		return reterr
	}
	outfile.WriteString(message + "\n")

	err = outfile.Close()
	if err != nil {
		reterr = fmt.Errorf("Unable to close file: %s. %s", pathToLog, err)	
		return reterr
	}
  return reterr
}
