package myiniconf

import "fmt"
import "io/ioutil"
import "strings"
import "reflect"
import "strconv"

type myiniconf interface{}

var confMaps map[string](map[string]string)
var confContent string

func init() {
	fmt.Println("This is a test of ini conf file parser.")
	confMaps = make(map[string](map[string]string))
}

func showAll() {
	fmt.Println("-------------all varibles have been loaded.-------------")
	var col int = 1
	for k, vs := range confMaps {
		for i, v := range vs {
			fmt.Printf("%d : %s.%s = %s\n", col, k, i, v)
			col++
		}
	}
	fmt.Println("--------------------------------------------------------")
}

func LoadConf(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Open file failed, err: %v \n", err)
		return
	}
	confContent = string(content)
	iniConf(string(content))
}

func iniConf(content string) {
	lines := strings.Split(content, "\n")
	var key string
	var startFlag bool
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line[0] == ';' || line[0] == '#' {
			continue
		}
		if line[0] == '[' && line[len(line)-1] == ']' {
			startFlag = true
			key = line[1 : len(line)-1]
			confMaps[key] = make(map[string]string)
		}
		if !startFlag {
			continue
		}
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		confMaps[key][strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	showAll()
}

func checkStruct(str interface{}) bool {
	return true
}

func Parser(section string, str interface{}) {
	if _, ok := confMaps[section]; !ok {
		return
	}
	if !checkStruct(str) {
		return
	}

	d := reflect.ValueOf(str)
	t := reflect.TypeOf(str)

	for k, v := range confMaps[section] {
		for idx := 0; idx < t.Elem().NumField(); idx++ {
			field := t.Elem().Field(idx)
			if field.Tag.Get("ini") == k {
				valueName := t.Elem().Field(idx).Name
				switch field.Type.Kind() {
				case reflect.String:
					d.Elem().FieldByName(valueName).SetString(v)
				case reflect.Int:
					value, _ := strconv.Atoi(v)
					d.Elem().FieldByName(valueName).SetInt(int64(value))
				case reflect.Bool:
					value, _ := strconv.ParseBool(v)
					d.Elem().FieldByName(valueName).SetBool(value)
				default:
					fmt.Println("Dont Know which type it is. ", field.Type.Kind())
				}
			}
		}
	}
}
