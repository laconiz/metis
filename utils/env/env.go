package env

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// 从参数中获取合法的IP地址及端口
func Address(name string, defaults string) (string, bool) {
	value, ok := argument(name)
	if !ok {
		value = defaults
	}
	array := strings.Split(value, ":")
	if len(array) != 2 {
		return "", false
	}
	if ip := net.ParseIP(array[0]); ip == nil {
		return "", false
	}
	port, err := strconv.Atoi(array[1])
	if err != nil || port < 0 || port > 65535 {
		return "", false
	}
	return value, true
}

func String(name string, defaults string) string {
	if str, ok := argument(name); ok {
		return str
	}
	return defaults
}

func Bool(name string) bool {
	value, ok := argument(name)
	if !ok {
		return false
	}
	if value == "" {
		return true
	}
	boolean, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return boolean
}

func argument(name string) (string, bool) {
	for index, argument := range arguments {
		// -name | --name
		if argument == "-"+name || argument == "--"+name {
			if index+1 >= len(arguments) {
				return "", true
			}
			value := arguments[index+1]
			if strings.HasPrefix(value, "-") {
				return "", true
			}
			if strings.HasPrefix(value, "--") {
				return "", true
			}
			return value, true
		}
		// -name=
		if strings.HasPrefix(argument, fmt.Sprintf("-%s=", name)) {
			return strings.SplitN(argument, "=", 2)[1], true
		}
		// --name=
		if strings.HasPrefix(argument, "--"+name+"=") {
			return strings.SplitN(argument, "=", 2)[1], true
		}
	}
	// 没有查找到参数
	return "", false
}

var arguments = os.Args
