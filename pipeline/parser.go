package pipeline

import (
	"bufio"
	"os"
	"strings"
)

func ParsePipelineFile(path string) (Pipeline, error) {
	pipeline := Pipeline{}

	file, err := os.Open(path)
	if err != nil {
		return Pipeline{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanRepo := false
	scanBuild := false
	scanDeploy := false
	for scanner.Scan() {
		if scanRepo {
			repo := scanner.Text()
			url := strings.Split(repo, ":")[1]
			pipeline.repo = url
			println(url)
		}

		if scanBuild {

		}

		if scanDeploy {

		}

		switch scanner.Text() {
		case "- INFO -":
			scanRepo = true
		case "- BUILD PHASE -":
			scanBuild = true
		case "- DEPLOY PHASE -":
			scanDeploy = true
		}

	}

	if err := scanner.Err(); err != nil {
		return Pipeline{}, err
	}
	return Pipeline{}, err
}
