package config

import "os"

type Config struct {
	NotionApiKey     string
	NotionDatabaseId string
	GithubUser       string
}

func GetConfig() *Config {
	return &Config{
		NotionApiKey:     os.Getenv("NOTION_API_KEY"),
		NotionDatabaseId: os.Getenv("NOTION_DATABASE_ID"),
		GithubUser:       os.Getenv("GITHUB_USER"),
	}
}
