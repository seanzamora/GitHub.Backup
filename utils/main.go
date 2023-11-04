package utils

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

var boldBlue = color.New(color.Bold, color.FgBlue)
var bold = color.New(color.Bold)
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

type Config struct {
	Hosts map[string][]Host
}

type Host struct {
	User  string
	Token string
}

func (config Config) SelectHost() int {
	fmt.Println()
	fmt.Println("Select account you would like to backup.")
	fmt.Println()

	count := 0
	for key, hosts := range config.Hosts {
		for _, host := range hosts {
			count++

			fmt.Printf(
				"%v %v %v\n",
				bold.Sprintf("%v.", count),
				bold.Sprintf(host.User),
				boldBlue.Sprintf("[%s]", key),
			)
		}
	}

	if count == 0 {
		log.Fatalf("No account have been added to Git.Backup. Try authenticating again.")
	}

	buf := bufio.NewReader(os.Stdin)

	fmt.Println()

	proceed := false

	for {
		fmt.Print(fmt.Sprintf("Select an account: [1-%v] ", count))

		input, _ := buf.ReadString('\n')

		index, err := strconv.Atoi(strings.Trim(input, "\n"))

		if err != nil {
			fmt.Println("Invalid input entry, try again.")
			continue
		}

		if index < 1 || index > count {
			fmt.Println("Invalid input range, try again.")
			continue
		}

		proceed = true

		if proceed {
			return index
		}
	}

	return -1
}

func (config Config) GetHost(index int) (*Host, error) {
	count := 0
	for _, hosts := range config.Hosts {
		for _, host := range hosts {
			count++
			if count == index {
				return &host, nil
			}
		}
	}
	return nil, fmt.Errorf("invalid selection")
}

func (config Config) HostExists(host string, user string) bool {
	for _, host := range config.Hosts[host] {
		if host.User == user {
			return true
		}
	}
	return false
}

func (config Config) AddHost(host string, user string, token string) {
	homeDir, _ := os.UserHomeDir()
	configPath := fmt.Sprintf("%s/.config/git-backup/hosts.conf", homeDir)

	file, _ := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)

	line := fmt.Sprintf("%s:%s|%s:%s\n", "host", host, user, token)

	_, err := file.Write([]byte(line))

	if err != nil {
		log.Fatalf("Unable to save host to config. Message: %v", err)
	}
}

func GenerateHash(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func GetHash() []byte {
	homeDir, _ := os.UserHomeDir()
	configDir := fmt.Sprintf("%s/.config/git-backup", homeDir)
	hash := fmt.Sprintf("%s/hash", configDir)

	file, err := os.ReadFile(hash)

	if err != nil {
		log.Fatalf("%v/n", err)
	}

	return file
}

func GetConfig() Config {
	config := Config{}
	config.Hosts = make(map[string][]Host)

	homeDir, _ := os.UserHomeDir()
	configDir := fmt.Sprintf("%s/.config/git-backup", homeDir)
	hosts := fmt.Sprintf("%s/hosts.conf", configDir)

	file, err := os.ReadFile(hosts)

	if err != nil {
		log.Fatalf("%v/n", err)
	}

	configRaw := string(file)
	rows := strings.Split(configRaw, "\n")

	for _, row := range rows {
		var host *string
		for i, col := range strings.Split(row, "|") {
			keyVal := strings.Split(col, ":")
			if len(keyVal) == 2 {
				key := &keyVal[0]
				val := &keyVal[1]
				if i == 0 && *key == "host" {
					host = val
				} else {
					config.Hosts[*host] = append(config.Hosts[*host], Host{*key, *val})
				}
			}
		}
	}

	return config
}
