package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type bgProcess struct {
	invocationCommand string
	id                int
}

type gitStatus struct {
	clean     bool
	ahead     bool
	unchanged bool
	branch    string
}

func getAllRegexMatchesFromString(input string, pattern string) [][]string {
	re := regexp.MustCompile(pattern)
	return re.FindAllStringSubmatch(input, -1)
}

func replaceSpecialChars(input string) string {
	matches := getAllRegexMatchesFromString(input, `~`)

	for _, match := range matches {
		switch match[0] {
		case "~":
			input = strings.Replace(input, match[0], "$HOME", -1)
		}
	}
	return input
}

func replaceEnvVars(input string) string {
	matches := getAllRegexMatchesFromString(input, `\$[A-Za-z0-9]+`)

	for _, match := range matches {
		value, isSet := os.LookupEnv(match[0][1:])
		if !isSet {
			continue
		}

		input = strings.Replace(input, match[0], value, 1)
	}

	return input
}

func execInput(input string, doneChannel chan int, id int) error {
	// Remove the newline
	input = strings.TrimSuffix(input, "\n")

	// Pass the arguments, split on space
	args := strings.Split(input, " ")

	// Check for built-in commands.
	switch args[0] {
	case "cd":
		// 'cd' to home dir with empty path not yet supported.
		if len(args) < 2 {
			return errors.New("path required")
		}
		// Change the directory and return the error.
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	default:
		// Prepare the command to execute
		cmd := exec.Command(args[0], args[1:]...)
		// Set the correct output devices
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		// Execute the command and return the error
		err := cmd.Run()
		if doneChannel != nil {
			doneChannel <- id
		}
		return err
	}
	return errors.New("Not able to execute command")
}

func getGitStatus() string {
	// Create an *exec.Cmd
	cmd := exec.Command("git", "status")

	// Combine stdout and stderr
	output, _ := cmd.CombinedOutput()
	output_str := string(output[:])

	result := ""
	if strings.Contains(output_str, "nothing added") {
		result = "-"
	} else if strings.Contains(output_str, "working directory clean") {
		result = "="
	} else if strings.Contains(output_str, "Your branch is ahead of") {
		result = "+"
	}
	return result
}

func printPrompt() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	gs := getGitStatus()

	fmt.Print("SSI: ", cwd + " " + gs, "\n > ")

	return
}

func main() {
	maxProcs := 2048
	reader := bufio.NewReader(os.Stdin)
	completedBackgroundJobs := make(chan int)
	backgroundJobs := make([]bgProcess, maxProcs)

	for {
		select {
		case id := <-completedBackgroundJobs:
			for i := range backgroundJobs {
				if backgroundJobs[i].id == id {
					cmd := strings.TrimSuffix(backgroundJobs[i].invocationCommand, "\n")
					fmt.Println(id, ": "+cmd+" has termintated.")
					// Found!
				}
			}
		default:
		}

		printPrompt()

		// Read the input
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// Sanitise input
		input = replaceEnvVars(replaceSpecialChars(input))

		// Check if this is a background process
		if len(input) > 2 && input[0:3] == "bg " {
			fmt.Println("Background process")
			id := rand.Intn(maxProcs)
			bgProc := bgProcess{input[3:], id}
			backgroundJobs[id] = bgProc
			go execInput(input[3:], completedBackgroundJobs, id)
			continue
		} else if len(input) > 6 && input[0:7] == "bglist " {
			fmt.Println(backgroundJobs)
		}

		// Handle the execution of the input.
		if err = execInput(input, nil, 0); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
