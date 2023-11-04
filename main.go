package main

import (
	"cli"
	"fmt"
	"log"
	"os"
	"utils"
)

func main() {
	_ = os.Args

	gh := &cli.GithubCLI{}

	config := Setup(gh)

	hostIdx := config.SelectHost()

	host, err := config.GetHost(int(hostIdx))

	gh.SetHost(host)

	if err != nil {
		log.Fatalf("%v\n", err)
		return
	}

	gh.GetRepositories()

	gh.BackupRepositories()
}

func Setup(gh *cli.GithubCLI) utils.Config {

	gh.IsInstalled()
	gh.IsLoggedIn()
	conf := gh.Setup()

	homeDir, _ := os.UserHomeDir()
	configDir := fmt.Sprintf("%s/.config/git-backup", homeDir)
	configFile := fmt.Sprintf("%s/hosts.conf", configDir)
	hash := fmt.Sprintf("%s/hash", configDir)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.Mkdir(configDir, 0755)
		if err != nil {
			log.Fatalf("Unable to create config directory. Message: %v/n", err)
		}
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		err := os.WriteFile(configFile, []byte(""), 0755)
		if err != nil {
			log.Fatalf("Unable to create config. Message: %v/n", err)
		}
	}

	if _, err := os.Stat(hash); os.IsNotExist(err) {
		hashStr := utils.GenerateHash(150)
		err := os.WriteFile(hash, []byte(hashStr), 0755)
		if err != nil {
			log.Fatalf("Unable to create config. Message: %v/n", err)
		}
	}

	config := utils.GetConfig()

	if _, ok := config.Hosts["github.com"]; !ok {
		config.Hosts["github.com"] = []utils.Host{}
	}

	host := utils.Host{User: conf.Github.User, Token: conf.Github.OauthToken}

	if !config.HostExists("github.com", host.User) {
		config.AddHost("github.com", host.User, host.Token)
		config.Hosts["github.com"] = append(config.Hosts["github.com"], host)
	}

	return config
}
