package config

import "os"

func DBUser() string {
	return os.Getenv("DB_USER")
}

func DBPassword() string {
	return os.Getenv("DB_PASSWORD")
}

func DBHost() string {
	return os.Getenv("DB_HOST")
}

func DBPort() string {
	return os.Getenv("DB_PORT")
}

func DBName() string {
	return os.Getenv("DB_NAME")
}
