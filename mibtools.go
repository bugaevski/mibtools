package mibtools

// [2020-05-12] 1535
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
	return "mibtools: MAY-12-2020"
} 

var pathToConfig string = ""
var pathToLog string = ""

// ConfigEntry key/value pair model
type ConfigEntry struct {
	Key string `json:"Key"`
	Value string `json:"Value"`
}

// AppConfig slice of ConfigEntry key/value pair model
type AppConfig []ConfigEntry

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

// ReadConfigKey return value by key name provided, from config file at private pathToConfig
func ReadConfigKey(key string) string {
	var returnedValue string = ""
	data, err := ioutil.ReadFile(pathToConfig)
	if err != nil {
			fmt.Println("File reading error", err, pathToConfig)
			return returnedValue
	}
	var target AppConfig
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
	var target AppConfig
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
func GetConfigKeyValue(key string, app *AppConfig) string {
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
func PopulateAppConfig(app *AppConfig, optionalpath string) (reterr error) {
	reterr = nil
	var item *ConfigEntry
	var actualpath string
	var errormessage string
	
	if pathToConfig == "" {
		if optionalpath == "" {
			reterr = errors.New("Both pathToConfig and optionalpath are empty")
			return reterr
		}else {
		  actualpath = optionalpath
		}
	} else {
		actualpath = pathToConfig
	}
	data, err := ioutil.ReadFile(actualpath)
	if err != nil {
		errormessage = fmt.Sprintf("File reading error. %s %s", err, actualpath)
		reterr = errors.New(errormessage)
		return reterr
	}
	var target AppConfig
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
func PutConfigKeyValue(key string, value string, app *AppConfig) {
	var item *ConfigEntry
	item = new(ConfigEntry)
	item.Key = key
	item.Value = value
	*app = append(*app, *item)
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
