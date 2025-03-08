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
	childCmd3 := &Command{
		Use:       "childCmd3",
		Short:     "third command - no modifiers",
		ValidArgs: []string{"arg1", "arg2"},
		Run:       emptyRun,
	}
	childCmd4 := &Command{
		Use:       "childCmd4",
		Short:     "fourth command - Show all flags",
		ValidArgs: []string{"arg1", "arg2"},
		Run:       emptyRun,
		CompletionOptions: CompletionOptions{
			ShowAllFlags: true,
		},
	}
	childCmd5 := &Command{
		Use:       "childCmd5",
		Short:     "fifth command - commingle args and flags",
		ValidArgs: []string{"arg1", "arg2"},
		Run:       emptyRun,
		CompletionOptions: CompletionOptions{
			CommingleArgsAndFlags: true,
		},
	}
	childCmd6 := &Command{
		Use:       "childCmd6",
		Short:     "sixth command - both",
		ValidArgs: []string{"arg1", "arg2"},
		Run:       emptyRun,
		CompletionOptions: CompletionOptions{
			ShowAllFlags:          true,
			CommingleArgsAndFlags: true,
		},
	}

	rootCmd.AddCommand(childCmd)
	rootCmd.AddCommand(childCmd2)
	rootCmd.AddCommand(childCmd3)
	rootCmd.AddCommand(childCmd4)
	rootCmd.AddCommand(childCmd5)
	rootCmd.AddCommand(childCmd6)

	rootCmd.Flags().IntP("first", "f", -1, "first flag\nlonger description for flag")
	rootCmd.PersistentFlags().BoolP("second", "s", false, "second flag")
	rootCmd.PersistentFlags().BoolP("third", "t", false, "third flag required")
	rootCmd.Flags().IntP("fourth", "4", -1, "fourth flag\nfourth flag required")
	rootCmd.MarkFlagRequired("third")
	rootCmd.MarkFlagRequired("fourth")

	childCmd.Flags().String("subFlag", "", "sub flag")
	childCmd2.Flags().String("subFlag2", "", "sub flag2")
	childCmd3.Flags().String("notrequired", "n", "sub flag3 not required")
	childCmd3.Flags().String("required", "r", "sub flag3 required")
	childCmd3.MarkFlagRequired("required")
	childCmd4.Flags().String("notrequired", "n", "sub flag4 not required")
	childCmd4.Flags().String("required", "r", "sub flag4 required")
	childCmd4.MarkFlagRequired("required")
	childCmd5.Flags().String("notrequired", "n", "sub flag5 not required")
	childCmd5.Flags().String("required", "r", "sub flag5 required")
	childCmd5.MarkFlagRequired("required")
	childCmd6.Flags().String("notrequired", "n", "sub flag6 not required")
	childCmd6.Flags().String("required", "r", "sub flag6 required")
	childCmd6.MarkFlagRequired("required")

	// Test that flag names are not shown if the user has not given the '-' prefix
	output, err := executeCommand(rootCmd, ShellCompRequestCmd, "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := strings.Join([]string{
		"childCmd\tfirst command",
		"childCmd2\tsecond command",
		"childCmd3\tthird command - no modifiers",
		"childCmd4\tfourth command - Show all flags",
		"childCmd5\tfifth command - commingle args and flags",
		"childCmd6\tsixth command - both",
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
		"--third\tthird flag required",
		"-t\tthird flag required",
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
		"--third\tthird flag required",
		"-t\tthird flag required",
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
		"--third\tthird flag required",
		"-t\tthird flag required",
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

	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "childCmd3", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--required\tsub flag3 required",
		"arg1",
		"arg2",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "childCmd3", "-")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected = strings.Join([]string{
		"--required\tsub flag3 required",
		":4",
		"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n")

	if output != expected {
		t.Errorf("expected: %q, got: %q", expected, output)
	}

	output, err = executeCommand(rootCmd, ShellCompRequestCmd, "childCmd3", "a")
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

func TestCompletionFlagBehaviorNormal(t *testing.T) {

	getCmd := func(setRequired bool, mutuallyExclusive bool) *Command {
		rootCmd := &Command{
			Use: "root",
			Run: emptyRun,
		}
		childCmd := &Command{
			Use: "child",
			Run: emptyRun,
		}
		rootCmd.AddCommand(childCmd)

		rootCmd.PersistentFlags().Int("pflag1", -1, "pflag1")
		rootCmd.PersistentFlags().Int("pflag2", -2, "pflag2")
		rootCmd.PersistentFlags().Int("pflag3", -3, "pflag3")
		rootCmd.Flags().Bool("flag1", false, "flag1")
		rootCmd.Flags().Bool("flag2", false, "flag2")
		rootCmd.Flags().Bool("flag3", false, "flag3")

		childCmd.PersistentFlags().Int("pchflag1", -1, "pchflag1")
		childCmd.PersistentFlags().Int("pchflag2", -2, "pchflag2")
		childCmd.PersistentFlags().Int("pchflag3", -3, "pchflag3")

		childCmd.Flags().Bool("chflag1", false, "chflag1")
		childCmd.Flags().Bool("chflag2", false, "chflag2")
		childCmd.Flags().Bool("chflag3", false, "chflag3")

		if setRequired {
			rootCmd.MarkPersistentFlagRequired("pflag1")
			rootCmd.MarkFlagRequired("flag1")

			childCmd.MarkPersistentFlagRequired("pchflag1")
			childCmd.MarkFlagRequired("chflag1")
		}
		if mutuallyExclusive {
			rootCmd.MarkFlagsMutuallyExclusive("pflag1", "pflag2")
			rootCmd.MarkFlagsMutuallyExclusive("flag1", "flag2")

			childCmd.MarkFlagsMutuallyExclusive("pchflag1", "pchflag2")
			childCmd.MarkFlagsMutuallyExclusive("chflag1", "chflag2")
		}
		//rootCmd.MarkFlagsMutuallyExclusive()

		//rootCmd.MarkFlagsOneRequired()
		//rootCmd.MarkFlagsRequiredTogether("flag1", "flag2")

		return rootCmd
	}

	testcases := []struct {
		name                   string
		input                  []string
		expectedNotRequired    string
		expectedRequired       string
		expectedMutual         string
		expectedMutualRequired string
	}{
		{
			name:  "blank",
			input: []string{""},
			expectedNotRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete subcommand 'c'",
			input: []string{"c"},
			expectedNotRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-'",
			input: []string{"-"},
			expectedNotRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '--'",
			input: []string{"--"},
			expectedNotRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-' with local --flag1",
			input: []string{"--flag1", "-"},
			expectedNotRequired: strings.Join([]string{
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-' with persistent --pflag1",
			input: []string{"--pflag1", "1", "-"},
			expectedNotRequired: strings.Join([]string{ //TODO this does not seem right
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--flag1\tflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--flag1\tflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-' with local --flag2",
			input: []string{"--flag2", "-"},
			expectedNotRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-' with persistent --pflag2",
			input: []string{"--pflag2", "2", "-"},
			expectedNotRequired: strings.Join([]string{ //TODO this does not seem right
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--flag1\tflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-' with local --flag3",
			input: []string{"--flag3", "-"},
			expectedNotRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-' with persistent --pflag3",
			input: []string{"--pflag3", "3", "-"},
			expectedNotRequired: strings.Join([]string{ //TODO this does not seem right
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child command incomplete",
			input: []string{"child"},
			expectedNotRequired: strings.Join([]string{
				"child",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"child",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"child",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"child",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child command",
			input: []string{"child", ""},
			expectedNotRequired: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
		},
		{
			name:  "child incomplete subsommand 'c'",
			input: []string{"child", "c"},
			expectedNotRequired: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-'",
			input: []string{"child", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '--'",
			input: []string{"child", "--"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with persistent --pflag1 post child",
			input: []string{"child", "--pflag1", "1", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with persistent --pflag2 post child",
			input: []string{"child", "--pflag2", "2", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with persistent --pflag3 post child",
			input: []string{"child", "--pflag3", "3", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with persistent --pflag1 pre child",
			input: []string{"--pflag1", "1", "child", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with persistent --pflag2 pre child",
			input: []string{"--pflag2", "2", "child", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with persistent --pflag3 pre child",
			input: []string{"--pflag3", "3", "child", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
	}

	for _, tc := range testcases {

		// required set to false
		t.Run(tc.name+"_Not_Required", func(t *testing.T) {

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : no Flags Set Required")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(false, false), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if output != tc.expectedNotRequired {
				t.Logf("Expected:\n%+v", tc.expectedNotRequired)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			}

		})

		// required set to true
		t.Run(tc.name+"s_Required", func(t *testing.T) {

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag1 Set Required")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(true, false), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if output != tc.expectedRequired {
				t.Logf("Expected:\n%+v", tc.expectedRequired)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			}
		})

		// mutual set to true
		t.Run(tc.name+"s_Mutual", func(t *testing.T) {

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag1 Set Mutually Exclusive")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(false, true), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if output != tc.expectedMutual {
				t.Logf("Expected:\n%+v", tc.expectedMutual)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			}
		})

		// mutual & Required set to true
		t.Run(tc.name+"s_Mutual_Required", func(t *testing.T) {

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag1 Set Mutually Exclusive & Required")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(true, true), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if output != tc.expectedMutualRequired {
				t.Logf("Expected:\n%+v", tc.expectedMutualRequired)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			}
		})
	}
}
