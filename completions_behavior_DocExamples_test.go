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
	"sort"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func TestCompletionFlagBehaviorDocExamplesValidFlagsFunction(t *testing.T) {

	getCmd := func() *Command {

		cmd := &Command{
			Use:       "root",
			Run:       func(*Command, []string) {},
			ValidArgs: []string{"one", "two"},
			ValidFlagsFunction: func(cmd *Command, args []string, toComplete string) ([]Completion, ShellCompDirective) {
				list := []string{}

				cmd.Flags().VisitAll(func(flag *pflag.Flag) {

					// logic for selecting which flags to exclude
					if flag.Name == "flag1" {
						for _, v := range args {
							if v == "one" {
								// exclude 'flag1' when "one" exists on the command line
								return
							}
						}
					}

					if flag.Name == "flag2" {
						for _, v := range args {
							if v == "two" {
								// exclude 'flag2' when "two" exists on the command line
								return
							}
						}
					}

					// exclude hidden flags and flags already present on the command line
					if flag.Hidden || flag.Changed {
						return
					}

					if flag.Name != "" {
						item := "--" + flag.Name
						if strings.HasPrefix(item, toComplete) {
							list = append(list, CompletionWithDesc(item, flag.Usage))
						}
					}

					if flag.Shorthand != "" {
						item := "-" + flag.Shorthand
						if strings.HasPrefix(item, toComplete) {
							list = append(list, CompletionWithDesc(item, flag.Usage))
						}
					}
				})

				sort.Strings(list)

				return list, ShellCompDirectiveNoFileComp
			},
		}
		cmd.Flags().BoolP("flag1", "a", false, "flag1 description")
		cmd.Flags().BoolP("flag2", "b", false, "flag2 description")
		cmd.Flags().BoolP("flag3", "c", false, "flag3 description")

		return cmd

	}

	testcases := []struct {
		name           string
		input          []string
		expectedOutput string
	}{
		{
			name:  "blank",
			input: []string{""},
			expectedOutput: strings.Join([]string{
				"--flag1\tflag1 description",
				"--flag2\tflag2 description",
				"--flag3\tflag3 description",
				"--help\thelp for root",
				"-a\tflag1 description",
				"-b\tflag2 description",
				"-c\tflag3 description",
				"-h\thelp for root",
				"one",
				"two",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "Arg one - incomplete flag",
			input: []string{"one", "-"},
			expectedOutput: strings.Join([]string{
				"--flag2\tflag2 description",
				"--flag3\tflag3 description",
				"--help\thelp for root",
				"-b\tflag2 description",
				"-c\tflag3 description",
				"-h\thelp for root",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "Arg two - incomplete flag",
			input: []string{"two", "-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tflag1 description",
				"--flag3\tflag3 description",
				"--help\thelp for root",
				"-a\tflag1 description",
				"-c\tflag3 description",
				"-h\thelp for root",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "Arg one & two - incomplete flag",
			input: []string{"one", "two", "-"},
			expectedOutput: strings.Join([]string{
				"--flag3\tflag3 description",
				"--help\thelp for root",
				"-c\tflag3 description",
				"-h\thelp for root",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "flag3",
			input: []string{"--flag3", "-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tflag1 description",
				"--flag2\tflag2 description",
				"--help\thelp for root",
				"-a\tflag1 description",
				"-b\tflag2 description",
				"-h\thelp for root",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "flag3 blank",
			input: []string{"--flag3", ""},
			expectedOutput: strings.Join([]string{
				"--flag1\tflag1 description",
				"--flag2\tflag2 description",
				"--help\thelp for root",
				"-a\tflag1 description",
				"-b\tflag2 description",
				"-h\thelp for root",
				"one",
				"two",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "flag3 short",
			input: []string{"-c", "-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tflag1 description",
				"--flag2\tflag2 description",
				"--help\thelp for root",
				"-a\tflag1 description",
				"-b\tflag2 description",
				"-h\thelp for root",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "flag3 short blank",
			input: []string{"-c", ""},
			expectedOutput: strings.Join([]string{
				"--flag1\tflag1 description",
				"--flag2\tflag2 description",
				"--help\thelp for root",
				"-a\tflag1 description",
				"-b\tflag2 description",
				"-h\thelp for root",
				"one",
				"two",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "flag prefix",
			input: []string{"two", "--fl"},
			expectedOutput: strings.Join([]string{
				"--flag1\tflag1 description",
				"--flag3\tflag3 description",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
	}

	for _, tc := range testcases {

		// required set to false
		t.Run(tc.name, func(t *testing.T) {

			t.Logf("Running test: %v", tc.name)
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Logf("Unexpected error: %v", err)
				t.Fail()
			}

			if output != tc.expectedOutput {
				t.Logf("Expected:\n%+v", tc.expectedOutput)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			} else {
				// display actual output when successful.  Use with 'go test -v'
				t.Logf("\n%+v", output)
			}
		})
	}
}

func TestCompletionFlagBehaviorDocExamplesVerbosity(t *testing.T) {

	getCmd := func() *Command {

		rootCmd := &Command{
			Use:       "root",
			Run:       func(*Command, []string) {},
			ValidArgs: []string{"one", "two"},
		}
		child1Cmd := &Command{
			Use:       "child-min",
			Run:       func(*Command, []string) {},
			ValidArgs: []string{"min", "max"},
			CompletionBehaviors: &CompletionBehaviors{
				FlagVerbosity: MinimalFlags,
			},
		}
		rootCmd.AddCommand(child1Cmd)
		child2Cmd := &Command{
			Use:       "child-more",
			Run:       func(*Command, []string) {},
			ValidArgs: []string{"more", "less"},
			CompletionBehaviors: &CompletionBehaviors{
				FlagVerbosity: MoreVerboseFlags,
			},
		}
		rootCmd.AddCommand(child2Cmd)
		child3Cmd := &Command{
			Use:       "child-all",
			Run:       func(*Command, []string) {},
			ValidArgs: []string{"all", "none"},
			CompletionBehaviors: &CompletionBehaviors{
				FlagVerbosity: AllFlags,
			},
		}
		rootCmd.AddCommand(child3Cmd)

		rootCmd.PersistentFlags().BoolP("flag1", "1", false, "persistent flag 1 description")
		rootCmd.PersistentFlags().BoolP("flag2", "2", false, "persistent flag 2 description")

		rootCmd.Flags().BoolP("flagA", "a", false, "flag A description")
		rootCmd.Flags().BoolP("flagB", "b", false, "flag B description")

		child1Cmd.Flags().BoolP("flagC", "c", false, "flag C description")
		child1Cmd.Flags().BoolP("flagD", "d", false, "flag D description")

		child2Cmd.Flags().BoolP("flagE", "e", false, "flag E description")
		child2Cmd.Flags().BoolP("flagF", "f", false, "flag F description")

		child3Cmd.Flags().BoolP("flagG", "g", false, "flag G description")
		child3Cmd.Flags().BoolP("flagI", "i", false, "flag I description")

		return rootCmd
	}

	testcases := []struct {
		name           string
		input          []string
		expectedOutput string
	}{
		{
			name:  "blank",
			input: []string{""},
			expectedOutput: strings.Join([]string{
				"child-all",
				"child-min",
				"child-more",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"one",
				"two",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag",
			input: []string{"-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flag2\tpersistent flag 2 description",
				"-2\tpersistent flag 2 description",
				"--flagA\tflag A description",
				"-a\tflag A description",
				"--flagB\tflag B description",
				"-b\tflag B description",
				"--help\thelp for root",
				"-h\thelp for root",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "Min Verbosity - blank",
			input: []string{"child-min", ""},
			expectedOutput: strings.Join([]string{
				"min",
				"max",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "Min Verbosity - incomplete flag",
			input: []string{"child-min", "-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flag2\tpersistent flag 2 description",
				"-2\tpersistent flag 2 description",
				"--flagC\tflag C description",
				"-c\tflag C description",
				"--flagD\tflag D description",
				"-d\tflag D description",
				"--help\thelp for child-min",
				"-h\thelp for child-min",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "More Verbosity - blank",
			input: []string{"child-more", ""},
			expectedOutput: strings.Join([]string{
				"more",
				"less",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "More Verbosity - incomplete flag",
			input: []string{"child-more", "-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flag2\tpersistent flag 2 description",
				"-2\tpersistent flag 2 description",
				"--flagE\tflag E description",
				"-e\tflag E description",
				"--flagF\tflag F description",
				"-f\tflag F description",
				"--help\thelp for child-more",
				"-h\thelp for child-more",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "All Verbosity - blank",
			input: []string{"child-all", ""},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flag2\tpersistent flag 2 description",
				"-2\tpersistent flag 2 description",
				"--flagG\tflag G description",
				"-g\tflag G description",
				"--flagI\tflag I description",
				"-i\tflag I description",
				"--help\thelp for child-all",
				"-h\thelp for child-all",
				"all",
				"none",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "All Verbosity - incomplete flag",
			input: []string{"child-all", "-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flag2\tpersistent flag 2 description",
				"-2\tpersistent flag 2 description",
				"--flagG\tflag G description",
				"-g\tflag G description",
				"--flagI\tflag I description",
				"-i\tflag I description",
				"--help\thelp for child-all",
				"-h\thelp for child-all",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
	}

	for _, tc := range testcases {

		// required set to false
		t.Run(tc.name, func(t *testing.T) {

			t.Logf("Running test: %v", tc.name)
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Logf("Unexpected error: %v", err)
				t.Fail()
			}

			if output != tc.expectedOutput {
				t.Logf("Expected:\n%+v", tc.expectedOutput)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			} else {
				// display actual output when successful.  Use with 'go test -v'
				t.Logf("\n%+v", output)
			}
		})
	}
}

func TestCompletionFlagBehaviorDocExamplesVerbosityAndRequired(t *testing.T) {

	getCmd := func() *Command {

		rootCmd := &Command{
			Use:       "root",
			Run:       func(*Command, []string) {},
			ValidArgs: []string{"one", "two"},
		}
		child1Cmd := &Command{
			Use:       "child-min",
			Run:       func(*Command, []string) {},
			ValidArgs: []string{"min", "max"},
			CompletionBehaviors: &CompletionBehaviors{
				FlagVerbosity: MinimalFlags,
			},
		}
		rootCmd.AddCommand(child1Cmd)
		child2Cmd := &Command{
			Use:       "child-more",
			Run:       func(*Command, []string) {},
			ValidArgs: []string{"more", "less"},
			CompletionBehaviors: &CompletionBehaviors{
				FlagVerbosity: MoreVerboseFlags,
			},
		}
		rootCmd.AddCommand(child2Cmd)
		child3Cmd := &Command{
			Use:       "child-all",
			Run:       func(*Command, []string) {},
			ValidArgs: []string{"all", "none"},
			CompletionBehaviors: &CompletionBehaviors{
				FlagVerbosity: AllFlags,
			},
		}
		rootCmd.AddCommand(child3Cmd)

		rootCmd.PersistentFlags().BoolP("flag1", "1", false, "persistent flag 1 description")
		rootCmd.PersistentFlags().BoolP("flag2", "2", false, "persistent flag 2 description")
		rootCmd.MarkPersistentFlagRequired("flag1")

		rootCmd.Flags().BoolP("flagA", "a", false, "flag A description")
		rootCmd.Flags().BoolP("flagB", "b", false, "flag B description")
		rootCmd.MarkFlagRequired("flagA")

		child1Cmd.Flags().BoolP("flagC", "c", false, "flag C description")
		child1Cmd.Flags().BoolP("flagD", "d", false, "flag D description")
		child1Cmd.MarkFlagRequired("flagC")

		child2Cmd.Flags().BoolP("flagE", "e", false, "flag E description")
		child2Cmd.Flags().BoolP("flagF", "f", false, "flag F description")
		child2Cmd.MarkFlagRequired("flagE")

		child3Cmd.Flags().BoolP("flagG", "g", false, "flag G description")
		child3Cmd.Flags().BoolP("flagI", "i", false, "flag I description")
		child3Cmd.MarkFlagRequired("flagG")

		return rootCmd
	}

	testcases := []struct {
		name           string
		input          []string
		expectedOutput string
	}{
		{
			name:  "blank",
			input: []string{""},
			expectedOutput: strings.Join([]string{
				"child-all",
				"child-min",
				"child-more",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flagA\tflag A description",
				"-a\tflag A description",
				"one",
				"two",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag",
			input: []string{"-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flagA\tflag A description",
				"-a\tflag A description",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "Min Verbosity - blank",
			input: []string{"child-min", ""},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flagC\tflag C description",
				"-c\tflag C description",
				"min",
				"max",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "Min Verbosity - incomplete flag",
			input: []string{"child-min", "-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flagC\tflag C description",
				"-c\tflag C description",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "More Verbosity - blank",
			input: []string{"child-more", ""},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flagE\tflag E description",
				"-e\tflag E description",
				"more",
				"less",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "More Verbosity - incomplete flag",
			input: []string{"child-more", "-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flag2\tpersistent flag 2 description",
				"-2\tpersistent flag 2 description",
				"--flagE\tflag E description",
				"-e\tflag E description",
				"--flagF\tflag F description",
				"-f\tflag F description",
				"--help\thelp for child-more",
				"-h\thelp for child-more",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "All Verbosity - blank",
			input: []string{"child-all", ""},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flag2\tpersistent flag 2 description",
				"-2\tpersistent flag 2 description",
				"--flagG\tflag G description",
				"-g\tflag G description",
				"--flagI\tflag I description",
				"-i\tflag I description",
				"--help\thelp for child-all",
				"-h\thelp for child-all",
				"all",
				"none",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "All Verbosity - incomplete flag",
			input: []string{"child-all", "-"},
			expectedOutput: strings.Join([]string{
				"--flag1\tpersistent flag 1 description",
				"-1\tpersistent flag 1 description",
				"--flag2\tpersistent flag 2 description",
				"-2\tpersistent flag 2 description",
				"--flagG\tflag G description",
				"-g\tflag G description",
				"--flagI\tflag I description",
				"-i\tflag I description",
				"--help\thelp for child-all",
				"-h\thelp for child-all",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
	}

	for _, tc := range testcases {

		// required set to false
		t.Run(tc.name, func(t *testing.T) {

			t.Logf("Running test: %v", tc.name)
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Logf("Unexpected error: %v", err)
				t.Fail()
			}

			if output != tc.expectedOutput {
				t.Logf("Expected:\n%+v", tc.expectedOutput)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			} else {
				// display actual output when successful.  Use with 'go test -v'
				t.Logf("\n%+v", output)
			}
		})
	}
}
