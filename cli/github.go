package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"utils"
)

type GithubCLI struct {
	Hosts      []Host
	ActiveHost *utils.Host
	Repos      []Repo
}

type Config struct {
	Github Host `yaml:"github.com"`
}

type Host struct {
	User        string `yaml:"user"`
	OauthToken  string `yaml:"oauth_token"`
	GitProtocol string `yaml:"git_protocol"`
	Host        string `yaml:"host"`
}

type Repo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (cli *GithubCLI) IsInstalled() {
	_, err := exec.Command("gh").Output()
	if err != nil {
		log.Fatalf("GitHub CLI is not installed. Install Instructions (https://github.com/cli/cli#installation) ")
	}
}

func (cli *GithubCLI) IsLoggedIn() {
	_, err := exec.Command("gh", "auth", "status").Output()
	if err != nil {
		log.Fatalf("You are not logged into any github hosts. Run \"gh auth login\" to authenticate\n")
	}
}

func (cli *GithubCLI) SetHost(host *utils.Host) {
	cli.ActiveHost = host

	cmd := exec.Command("gh", "auth", "login", "--with-token")

	var buf bytes.Buffer
	buf.Write([]byte(cli.ActiveHost.Token))

	cmd.Stdin = &buf
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		log.Fatalf("An error has occured: %v", err)
	}

}
func (cli *GithubCLI) SetGit() {
	cmd := exec.Command("gh", "auth", "setup-git")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error occured: %v", err)
	}
}

func (cli *GithubCLI) SetProtocol() {
	cmd := exec.Command("gh", "config", "set", "git_protocol", "ssh", "-h", "github.com")
	_, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error occured: %v", err)
	}
}

func (cli *GithubCLI) BackupRepositories() {
	cli.SetProtocol()
	cli.SetGit()

	count := len(cli.Repos)
	fmt.Println()
	fmt.Printf("Downloading [%v] Repositories\n", count)
	fmt.Println()

	for i, repo := range cli.Repos {
		fmt.Println(fmt.Sprintf("[%v/%v] Downloading %v [%v]", i+1, count, repo.Name, repo.Url))
		cmd := exec.Command("gh", "repo", "clone", repo.Url)
		output, err := cmd.StdoutPipe()
		if err := cmd.Start(); err != nil {
			log.Fatalf("Error occured: %v", err)
		}
		reader := bufio.NewReader(output)
		line, err := reader.ReadString('\n')
		for err == nil {
			_ = fmt.Sprintf("%v", line)
		}
		if err != nil {
			log.Fatalf("Error occured: %v", err)
		}
	}
}

func (cli *GithubCLI) GetRepositories() []Repo {
	cli.IsLoggedIn()

	cmd := exec.Command("gh", "repo", "list", "-L", "1000", "--json", "name,url")

	output, err := cmd.Output()

	if err != nil {
		log.Fatalf("Error occured: %v", err)
	}

	var repos []Repo

	err = json.Unmarshal(output, &repos)

	if err != nil {
		log.Fatalf("Error occured: %v", err)
	}

	cli.Repos = repos

	return repos
}

func (cli *GithubCLI) Setup() *Config {
	homeDir, _ := os.UserHomeDir()
	configPath := fmt.Sprintf("%s/.config/gh/hosts.yml", homeDir)
	file, err := os.OpenFile(configPath, os.O_RDONLY, 0600)

	if err != nil {
		log.Fatalf("Unable to initialize GitHub CLI host config. Message: %v\n", err)
	}

	conf := &Config{}

	config := yaml.NewDecoder(file)
	err = config.Decode(conf)

	if err != nil {
		log.Fatalf("%v", err)
	}

	cli.Hosts = []Host{
		{
			User:        conf.Github.User,
			OauthToken:  conf.Github.OauthToken,
			GitProtocol: conf.Github.GitProtocol,
			Host:        conf.Github.Host,
		},
	}

	return conf
}
