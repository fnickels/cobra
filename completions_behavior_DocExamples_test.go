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
			Run:       emptyRun,
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
				t.Logf("Expected Len(): %v", len(tc.expectedOutput))
				t.Logf("Got      Len(): %v", len(output))

				t.Logf("Expected:\n%+v", tc.expectedOutput)
				t.Logf("Got:\n%+v", output)

				//	min := len(tc.expectedOutput)
				//	if min > len(output) {
				//		min = len(output)
				//	}
				//	//if min > 30 {
				//	//	min = 30
				//	//}
				//
				//	for i := 0; i < min; i++ {
				//		t.Logf("(%2v) Expected %v -> Got %v", i, tc.expectedOutput[i], output[i])
				//	}

				t.Fail()
			}
		})
	}
}

//
//
//  func TestCompletionFlagBehaviorDocExamples(t *testing.T) {
//
//
//  	getCmd := func() *Command {
//
//  		rootCmd := &Command{
//  			Use:       "root",
//  			Run:       emptyRun,
//  			ValidArgs: []string{"rarg1", "rarg2"},
//  			CompletionBehaviors: &CompletionBehaviors{
//  				FlagVerbosity: MoreVerboseFlags,
//  			},
//  		}
//  		childCmd := &Command{
//  			Use: "child",
//  			Run: emptyRun,
//  		}
//  		rootCmd.AddCommand(childCmd)
//  		childCmd2 := &Command{
//  			Use:       "secondchild",
//  			Run:       emptyRun,
//  			ValidArgs: []string{"arg1", "arg2", "arg3"},
//  		}
//  		rootCmd.AddCommand(childCmd2)
//
//  		rootCmd.PersistentFlags().Int("pflag1", -1, "pflag1")
//  		rootCmd.PersistentFlags().Int("pflag2", -2, "pflag2")
//  		rootCmd.PersistentFlags().Int("pflag3", -3, "pflag3")
//  		rootCmd.Flags().Bool("flag1", false, "flag1")
//  		rootCmd.Flags().Bool("flag2", false, "flag2")
//  		rootCmd.Flags().Bool("flag3", false, "flag3")
//
//  		childCmd.PersistentFlags().Int("pchflag1", -1, "pchflag1")
//  		childCmd.PersistentFlags().Int("pchflag2", -2, "pchflag2")
//  		childCmd.PersistentFlags().Int("pchflag3", -3, "pchflag3")
//
//  		childCmd.Flags().Bool("chflag1", false, "chflag1")
//  		childCmd.Flags().Bool("chflag2", false, "chflag2")
//  		childCmd.Flags().Bool("chflag3", false, "chflag3")
//
//  		childCmd2.PersistentFlags().Int("pchflag1", -1, "pchflag1")
//  		childCmd2.PersistentFlags().Int("pchflag2", -2, "pchflag2")
//  		childCmd2.PersistentFlags().Int("pchflag3", -3, "pchflag3")
//
//  		childCmd2.Flags().Bool("chflag1", false, "chflag1")
//  		childCmd2.Flags().Bool("chflag2", false, "chflag2")
//  		childCmd2.Flags().Bool("chflag3", false, "chflag3")
//
//  		if setRequired {
//  			rootCmd.MarkPersistentFlagRequired("pflag1")
//  			rootCmd.MarkFlagRequired("flag1")
//
//  			childCmd.MarkPersistentFlagRequired("pchflag1")
//  			childCmd.MarkFlagRequired("chflag1")
//
//  			childCmd2.MarkPersistentFlagRequired("pchflag1")
//  			childCmd2.MarkFlagRequired("chflag1")
//  		}
//
//  		if mutuallyExclusive {
//  			rootCmd.MarkFlagsMutuallyExclusive("pflag1", "pflag2")
//  			rootCmd.MarkFlagsMutuallyExclusive("flag1", "flag2")
//
//  			childCmd.MarkFlagsMutuallyExclusive("pchflag1", "pchflag2")
//  			childCmd.MarkFlagsMutuallyExclusive("chflag1", "chflag2")
//
//  			childCmd2.MarkFlagsMutuallyExclusive("pchflag1", "pchflag2")
//  			childCmd2.MarkFlagsMutuallyExclusive("chflag1", "chflag2")
//  		}
//
//  		if oneRequired {
//  			rootCmd.MarkFlagsOneRequired("pflag1", "pflag2")
//  			rootCmd.MarkFlagsOneRequired("flag1", "flag2")
//
//  			childCmd.MarkFlagsOneRequired("pchflag1", "pchflag2")
//  			childCmd.MarkFlagsOneRequired("chflag1", "chflag2")
//
//  			childCmd2.MarkFlagsOneRequired("pchflag1", "pchflag2")
//  			childCmd2.MarkFlagsOneRequired("chflag1", "chflag2")
//  		}
//
//  		if requiredTogether {
//  			rootCmd.MarkFlagsRequiredTogether("pflag1", "pflag3")
//  			rootCmd.MarkFlagsRequiredTogether("flag1", "flag3")
//
//  			childCmd.MarkFlagsRequiredTogether("pchflag1", "pchflag3")
//  			childCmd.MarkFlagsRequiredTogether("chflag1", "chflag3")
//
//  			childCmd2.MarkFlagsRequiredTogether("pchflag1", "pchflag3")
//  			childCmd2.MarkFlagsRequiredTogether("chflag1", "chflag3")
//  		}
//
//  		return rootCmd
//  	}
//
//  	testcases := []struct {
//  		name           string
//  		input          []string
//  		expectedOutput string
//  	}{
//  		{
//  			name:  "blank",
//  			input: []string{""},
//  			expectedOutput: strings.Join([]string{
//  				"child",
//  				"completion\tGenerate the autocompletion script for the specified shell",
//  				"help\tHelp about any command",
//  				"secondchild",
//  				"rarg1",
//  				"rarg2",
//  				":4",
//  				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
//  		},	}
//
//  	for _, tc := range testcases {
//
//  		// required set to false
//  		t.Run(tc.name+"_Not_Required", func(t *testing.T) {
//
//  			t.Logf("Running test: %v", tc.name)
//  			t.Logf("Input       : %v", tc.input)
//
//  			output, err := executeCommand(getCmd(), append([]string{ShellCompRequestCmd}, tc.input...)...)
//  			if err != nil {
//  				t.Logf("Unexpected error: %v", err)
//  				t.Fail()
//  			}
//
//  			if output != tc.expectedOutput {
//  				t.Logf("Expected:\n%+v", tc.expectedOutput)
//  				t.Logf("Got:\n%+v", output)
//  				t.Fail()
//  			}
//  		})
//  }
