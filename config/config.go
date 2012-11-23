package config

import (
	"TOGY/util"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const timeFormat = "15:04"
const dateFormat = "2006-1-2"

//Cofiguration as unmarshaled from JSON.
//Times are specified with minute precision in a format like this: 15:04
//Dates should have this format: 2010-3-12
//These two intermediate structs could be replaced by a map[string]interface{},
//but I tried to do it and the simplification that it affords 
//would be outweighed by the need to use type assertions everywhere
//and some strange type errors.
type localConfig struct {
	PowerPoint  string
	UpdateURL   string
	LogFile     string
	Name        string
	CentralPath string
}

type centralConfig struct {
	StandardTimeSettings map[string]string
	OverrideDays         map[string]map[string]string
	OverrideOn           bool
	OverrideOff          bool
	UpdateInterval       int
}

//The real configuration struct.
type Config struct {
	Presentation         string
	UpdatePath           string
	StandardTimeSettings TimeConfig
	OverrideDays         map[time.Time]TimeConfig
	OverrideOn           bool
	OverrideOff          bool
	PowerPoint           string
	UpdateURL            string
	UpdateInterval       int
	Log                  *log.Logger
	Name                 string
	CentralPath          string
}

//Struct representing time when the TV should be running.
type TimeConfig struct {
	TurnOn  time.Time
	TurnOff time.Time
}

//Loads configuration file from the specified path.
func getLocal(path string) (l localConfig, err error) {
	err = getJSONFile(path, &l)
	return
}

func (l localConfig) GetCentral() (c centralConfig, err error) {
	err = getJSONFile(l.CentralPath, &c)
	return
}

func Get(path string) (conf Config, err error) {
	l, err := getLocal(path)
	if err != nil {
		return
	}
	c, err := l.GetCentral()
	if err != nil {
		return
	}
	conf, err = joinConfigs(l, c)
	return
}

//Converts jsonConfig to Config.
func joinConfigs(l localConfig, c centralConfig) (conf Config, err error) {
	conf.StandardTimeSettings, err = makeTimeConfig(c.StandardTimeSettings)
	conf.OverrideDays = make(map[time.Time]TimeConfig)
	if err != nil {
		return
	}
	for k, v := range c.OverrideDays {
		var key time.Time
		key, err = time.Parse(dateFormat, k)
		if err != nil {
			return
		}

		key = util.NormalizeDate(key)

		conf.OverrideDays[key], err = makeTimeConfig(v)
		if err != nil {
			return
		}
	}

	logOut, err := os.Create(l.LogFile)
	if err != nil {
		logOut = os.Stderr
	}

	conf.Log = log.New(logOut, "", log.LstdFlags)

	conf.OverrideOn = c.OverrideOn
	conf.OverrideOff = c.OverrideOff
	conf.PowerPoint = l.PowerPoint
	conf.UpdateURL = l.UpdateURL
	conf.UpdateInterval = c.UpdateInterval
	conf.Name = l.Name
	conf.CentralPath = l.CentralPath
	return
}

//Converts map of strings to strings (formatted as time) to a TimeConfig struct.
func makeTimeConfig(times map[string]string) (tc TimeConfig, err error) {
	on, err := time.Parse(timeFormat, times["TurnOn"])
	if err != nil {
		return
	}

	off, err := time.Parse(timeFormat, times["TurnOff"])
	if err != nil {
		return
	}

	tc.TurnOn = util.NormalizeTime(on)
	tc.TurnOff = util.NormalizeTime(off)
	return
}

func getJSONFile(path string, d interface{}) (err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, d)
	if err != nil {
		return
	}
	return
}
