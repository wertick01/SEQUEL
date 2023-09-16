package config

import (
	"log"
	"os/exec"
)

type Config struct {
	System *SystemInfo
}

type SystemInfo struct {
	OperationSystem string
	JavaInfo        string
}

func GetConfig() *Config {
	out, err := exec.Command("java", "-version").Output()
	log.Println(out)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Command Successfully Executed")
	output := string(out[:])
	log.Println(output)
	// log.Println(string(resp), err)
	return &Config{}
}
