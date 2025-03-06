// Copyright 2013-2023 The Cobra Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cobra

import (
	"strings"
	"testing"
)

func TestCommingleArgsAndFlags(t *testing.T) {
	rootCmd := &Command{
		Use: "root",
		Run: emptyRun,
		CompletionOptions: CompletionOptions{
			CommingleArgsAndFlags: true,
		},
	}
	childCmd := &Command{
		Use:       "childCmd",
		Short:     "first command",
		Version:   "1.2.3",
		ValidArgs: []string{"arg1", "arg2"},
		Run:       emptyRun,
	}
	childCmd2 := &Command{
		Use:       "childCmd2",
		Short:     "second command",
		Version:   "1.2.3",
		ValidArgs: []string{"arg1", "arg2"},
		Run:       emptyRun,
		CompletionOptions: CompletionOptions{
			CommingleArgsAndFlags: true,
		},
	}
	rootCmd.AddCommand(childCmd)
	rootCmd.AddCommand(childCmd2)

	rootCmd.Flags().IntP("first", "f", -1, "first flag\nlonger description for flag")
	rootCmd.PersistentFlags().BoolP("second", "s", false, "second flag")
	childCmd.Flags().String("subFlag", "", "sub flag")
	childCmd2.Flags().String("subFlag2", "", "sub flag2")

	// Test that flag names are not shown if the user has not given the '-' prefix
	output, err := executeCommand(rootCmd, ShellCompRequestCmd, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"childCmd\tfirst command",
		"childCmd2\tsecond command",
		"completion\tGenerate the autocompletion script for the specified shell",
		"help\tHelp about any command",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--first\tfirst flag",
		"-f\tfirst flag",
		"--help\thelp for root",
		"-h\thelp for root",
		"--second\tsecond flag",
		"-s\tsecond flag",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed when a prefix is given
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "--f")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--first\tfirst flag",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed in a sub-cmd
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "childCmd", "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--second\tsecond flag",
		"-s\tsecond flag",
		"--help\thelp for childCmd",
		"-h\thelp for childCmd",
		"--subFlag\tsub flag",
		"--version\tversion for childCmd",
		"-v\tversion for childCmd",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that flag names are completed in a sub-cmd
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "childCmd2", "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--second\tsecond flag",
		"-s\tsecond flag",
		"--help\thelp for childCmd2",
		"-h\thelp for childCmd2",
		"--subFlag2\tsub flag2",
		"--version\tversion for childCmd2",
		"-v\tversion for childCmd2",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	// Test that arg names are completed in a sub-cmd
	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "childCmd", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"arg1",
		"arg2",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "childCmd2", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"arg1",
		"arg2",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}
}
