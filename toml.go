package app

import (
	"fmt"
	"log"
	"regexp"

	"github.com/BurntSushi/toml"
)

type Repository struct {
	Repo           string `toml:"repo"`
	TargetFile     string `toml:"target_file"`
	TargetVariable string `toml:"target_variable"`
	Branch         string `toml:"branch"`
}

type Config struct {
	SourceRepoReleases string       `toml:"source_repo_releases"`
	TargetRepository   []Repository `toml:"target_repos"`
}

type netlifyConfig struct {
	Build   netlifyBuild `toml:"build"`
	Context netlifyBuild `toml:"context"`
}

type netlifyBuild struct {
	Command  string                  `toml:"command"`
	BuildEnv netlifyBuildEnvironment `toml:"environment"`
}
type netlifyBuildEnvironment struct {
	HugoVersion string `toml:"HUGO_VERSION"`
}

func ParseConfigFile(filepath string) (Config, error) {
	var conf Config
	if _, err := toml.DecodeFile(filepath, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func parseNetlifyConfFile(filepath string) (netlifyConfig, error) {
	var conf netlifyConfig
	if _, err := toml.DecodeFile(filepath, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func parseNetlifyConf(content string) (netlifyConfig, error) {
	var conf netlifyConfig
	if _, err := toml.Decode(content, &conf); err != nil {
		return conf, err
	}
	return conf, nil
}

func updateVersion(hugoVersion, deployContent string) string {

	regexp, err := regexp.Compile(`HUGO_VERSION = \"(\d+\.\d+\.\d+)\"`)
	if err != nil {
		log.Printf("Compile regexp error: %v", err)
	}
	replacement := fmt.Sprintf("HUGO_VERSION = \"%s\"", hugoVersion)

	return regexp.ReplaceAllString(deployContent, replacement)
}
