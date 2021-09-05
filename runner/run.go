package runner

import (
	"bufio"
	"fmt"
	"github.com/emirpasic/gods/sets/hashset"
	"io"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
)

func homeDir() string {
	current, err := user.Current()
	if nil == err {
		return current.HomeDir
	}
	return ""
}

func replaceHosts(newHosts map[string]string) error {
	hostPath := "/etc/hosts"
	file, err := os.OpenFile(hostPath, os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Cannot open /etc/hosts, err: [%v]", err)
		return err
	}
	defer file.Close()
	pos := int64(0)
	processHostSet := hashset.New()
	scanner := bufio.NewReader(file)
	for {
		bytes, _, err := scanner.ReadLine()
		if err == io.EOF {
			break
		}
		line := string(bytes)
		if line != "" && !strings.HasPrefix(line, "#") {
			hostItem := strings.Split(line, " ")
			if len(hostItem) == 2 {
				processHostSet.Add(hostItem[1])
				newIp := newHosts[hostItem[1]]
				if newIp != "" && newIp != "0.0.0.0" && newIp != hostItem[0] {
					fmt.Printf("begin replace %s from %s to %s\n", hostItem[1], hostItem[0], newIp)
					hostItem[0] = newIp
					line = strings.Join(hostItem, " ")
					file.WriteAt([]byte(line), pos)
				}
			}
		}
		pos += int64(len(line) + len("\n"))
	}

	for k, v := range newHosts {
		if !processHostSet.Contains(k) && v != "0.0.0.0" {
			fmt.Printf("begin add host item : %s %s\n", v, k)
			line := v + " " + k + "\n"
			file.WriteAt([]byte(line), pos)
			pos += int64(len(line))
		}
	}

	return nil
}

func Run() error {
	home := homeDir()
	domains, err := LoadFile(path.Join(home, ".github.txt"))
	if err != nil {
		return err
	}
	addrMap := ParseAddr(domains)
	return replaceHosts(addrMap)
}
