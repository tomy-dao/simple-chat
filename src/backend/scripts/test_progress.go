package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	Green  = "\033[0;32m"
	Red    = "\033[0;31m"
	Yellow = "\033[1;33m"
	Blue   = "\033[0;34m"
	Cyan   = "\033[0;36m"
	Reset  = "\033[0m"
)

func main() {
	fmt.Printf("%s========================================%s\n", Blue, Reset)
	fmt.Printf("%s   Running Tests with Progress%s\n", Blue, Reset)
	fmt.Printf("%s========================================%s\n\n", Blue, Reset)

	// Get all test packages
	cmd := exec.Command("go", "list", "./test/...")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("%sError: %v%s\n", Red, err, Reset)
		os.Exit(1)
	}

	packages := strings.Fields(string(output))
	total := len(packages)

	fmt.Printf("%sFound %d test packages%s\n\n", Yellow, total, Reset)

	passed := 0
	failed := 0

	// Run tests sequentially for cleaner output
	for i, pkg := range packages {
		fmt.Printf("%s[%d/%d]%s Testing: %s\n", Blue, i+1, total, Reset, pkg)

		// Run tests
		cmd := exec.Command("go", "test", "-v", pkg)
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		if err := cmd.Start(); err != nil {
			fmt.Printf("  %s✗ Error starting test%s\n", Red, Reset)
			failed++
			fmt.Println()
			continue
		}

		// Count PASS and FAIL
		passCount := 0
		failCount := 0
		var outputLines []string

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			outputLines = append(outputLines, line)
			if strings.Contains(line, "--- PASS:") {
				passCount++
			}
			if strings.Contains(line, "--- FAIL:") {
				failCount++
			}
		}

		scanner = bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			outputLines = append(outputLines, line)
		}

		cmd.Wait()
		exitCode := cmd.ProcessState.ExitCode()

		if exitCode == 0 {
			fmt.Printf("  %s✓ PASS%s - %d test(s) passed\n", Green, Reset, passCount)
			passed++
		} else {
			fmt.Printf("  %s✗ FAIL%s - %d test(s) failed\n", Red, Reset, failCount)
			// Show first few errors
			errorCount := 0
			for _, line := range outputLines {
				if (strings.Contains(line, "FAIL") || strings.Contains(line, "Error") || strings.Contains(line, "panic")) && errorCount < 3 {
					fmt.Printf("    %s%s%s\n", Red, line, Reset)
					errorCount++
				}
			}
			failed++
		}
		fmt.Println()
	}

	fmt.Printf("%s========================================%s\n", Blue, Reset)
	fmt.Printf("%s   Test Summary%s\n", Blue, Reset)
	fmt.Printf("%s========================================%s\n", Blue, Reset)
	fmt.Printf("%sTotal Packages: %d%s\n", Cyan, total, Reset)
	fmt.Printf("%sPassed: %d%s\n", Green, passed, Reset)
	if failed > 0 {
		fmt.Printf("%sFailed: %d%s\n", Red, failed, Reset)
	} else {
		fmt.Printf("%sFailed: %d%s\n", Green, failed, Reset)
	}

	if failed > 0 {
		os.Exit(1)
	}
}

