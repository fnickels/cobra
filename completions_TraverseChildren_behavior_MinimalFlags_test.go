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

func TestCompletionFlagTraverseChildrenBehaviorMinimalFlags(t *testing.T) {

	t.Parallel()

	getCmd := func(
		setRequired bool,
		mutuallyExclusive bool,
		oneRequired bool,
		requiredTogether bool,
	) *Command {

		rootCmd := &Command{
			Use:              "root",
			Run:              emptyRun,
			TraverseChildren: false,
		}
		childCmd := &Command{
			Use: "child",
			Run: emptyRun,
		}
		rootCmd.AddCommand(childCmd)
		childCmd2 := &Command{
			Use:       "secondchild",
			Run:       emptyRun,
			ValidArgs: []string{"arg1", "arg2", "arg3"},
		}
		rootCmd.AddCommand(childCmd2)

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

		childCmd2.PersistentFlags().Int("pchflag1", -1, "pchflag1")
		childCmd2.PersistentFlags().Int("pchflag2", -2, "pchflag2")
		childCmd2.PersistentFlags().Int("pchflag3", -3, "pchflag3")

		childCmd2.Flags().Bool("chflag1", false, "chflag1")
		childCmd2.Flags().Bool("chflag2", false, "chflag2")
		childCmd2.Flags().Bool("chflag3", false, "chflag3")

		if setRequired {
			rootCmd.MarkPersistentFlagRequired("pflag1")
			rootCmd.MarkFlagRequired("flag1")

			childCmd.MarkPersistentFlagRequired("pchflag1")
			childCmd.MarkFlagRequired("chflag1")

			childCmd2.MarkPersistentFlagRequired("pchflag1")
			childCmd2.MarkFlagRequired("chflag1")
		}

		if mutuallyExclusive {
			rootCmd.MarkFlagsMutuallyExclusive("pflag1", "pflag2")
			rootCmd.MarkFlagsMutuallyExclusive("flag1", "flag2")

			childCmd.MarkFlagsMutuallyExclusive("pchflag1", "pchflag2")
			childCmd.MarkFlagsMutuallyExclusive("chflag1", "chflag2")

			childCmd2.MarkFlagsMutuallyExclusive("pchflag1", "pchflag2")
			childCmd2.MarkFlagsMutuallyExclusive("chflag1", "chflag2")
		}

		if oneRequired {
			rootCmd.MarkFlagsOneRequired("pflag1", "pflag2")
			rootCmd.MarkFlagsOneRequired("flag1", "flag2")

			childCmd.MarkFlagsOneRequired("pchflag1", "pchflag2")
			childCmd.MarkFlagsOneRequired("chflag1", "chflag2")

			childCmd2.MarkFlagsOneRequired("pchflag1", "pchflag2")
			childCmd2.MarkFlagsOneRequired("chflag1", "chflag2")
		}

		if requiredTogether {
			rootCmd.MarkFlagsRequiredTogether("pflag1", "pflag3")
			rootCmd.MarkFlagsRequiredTogether("flag1", "flag3")

			childCmd.MarkFlagsRequiredTogether("pchflag1", "pchflag3")
			childCmd.MarkFlagsRequiredTogether("chflag1", "chflag3")

			childCmd2.MarkFlagsRequiredTogether("pchflag1", "pchflag3")
			childCmd2.MarkFlagsRequiredTogether("chflag1", "chflag3")
		}

		return rootCmd
	}

	testcases := []struct {
		name                                string
		input                               []string
		expectedNotRequired                 string
		expectedRequired                    string
		expectedMutual                      string
		expectedMutualRequired              string
		expectedOneRequired                 string
		expectedMutualOneRequired           string
		expectedRequiredTogether            string
		expectedRequiredTogetherAndRequired string
	}{
		{
			name:  "blank",
			input: []string{""},
			expectedNotRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"secondchild",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"secondchild",
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"secondchild",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"secondchild",
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedOneRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"secondchild",
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"secondchild",
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"secondchild",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				"help\tHelp about any command",
				"secondchild",
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
			expectedOneRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"child",
				"completion\tGenerate the autocompletion script for the specified shell",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
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
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{ //TODO why not show pflag1
				"--flag3\tflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--flag3\tflag3",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-' with persistent --pflag1",
			input: []string{"--pflag1", "1", "-"},
			expectedNotRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{ // why not show flag1?
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag3\tpflag3",
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
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--flag1\tflag1",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-' with persistent --pflag2",
			input: []string{"--pflag2", "2", "-"},
			expectedNotRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--flag3\tflag3",
				"--help\thelp for root",
				"-h\thelp for root",
				"--pflag1\tpflag1",
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
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
			expectedOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--flag1\tflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "incomplete flag '-' with persistent --pflag3",
			input: []string{"--pflag3", "3", "-"},
			expectedNotRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--flag1\tflag1",
				"--flag2\tflag2",
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"child",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"child",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"child",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				":0",
				"Completion ended with directive: ShellCompDirectiveDefault", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
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
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
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
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag3\tpflag3",
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
			expectedOneRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
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
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
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
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
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
			expectedOneRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--pflag3\tpflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag3\tpflag3",
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
			expectedOneRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
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
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
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
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--pflag1\tpflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with persistent child --pchflag1",
			input: []string{"child", "--pchflag1", "1", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
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
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with persistent child --pchflag2",
			input: []string{"child", "--pchflag2", "2", "-"},
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
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with persistent child --pchflag3 ",
			input: []string{"child", "--pchflag3", "3", "-"},
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
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with local child --chflag1",
			input: []string{"child", "--chflag1", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
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
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
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
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--chflag3\tchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag3\tchflag3",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with local child --chflag2",
			input: []string{"child", "--chflag2", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
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
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag3\tchflag3",
				"--help\thelp for child",
				"-h\thelp for child",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"--pchflag3\tpchflag3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "child incomplete flag '-' with local child --chflag3 ",
			input: []string{"child", "--chflag3", "-"},
			expectedNotRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--pflag3\tpflag3",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
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
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--chflag1\tchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
		{
			name:  "secondchild with local child --chflag3 ",
			input: []string{"secondchild", "--chflag3", ""},
			expectedNotRequired: strings.Join([]string{
				"arg1",
				"arg2",
				"arg3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				"arg1",
				"arg2",
				"arg3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutual: strings.Join([]string{
				"arg1",
				"arg2",
				"arg3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				"arg1",
				"arg2",
				"arg3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"arg1",
				"arg2",
				"arg3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedMutualOneRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--pflag2\tpflag2",
				"--chflag1\tchflag1",
				"--chflag2\tchflag2",
				"--pchflag1\tpchflag1",
				"--pchflag2\tpchflag2",
				"arg1",
				"arg2",
				"arg3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogether: strings.Join([]string{
				"--chflag1\tchflag1",
				"arg1",
				"arg2",
				"arg3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
			expectedRequiredTogetherAndRequired: strings.Join([]string{
				"--pflag1\tpflag1",
				"--chflag1\tchflag1",
				"--pchflag1\tpchflag1",
				"arg1",
				"arg2",
				"arg3",
				":4",
				"Completion ended with directive: ShellCompDirectiveNoFileComp", ""}, "\n"),
		},
	}

	for _, tc := range testcases {

		// required set to false
		t.Run(tc.name+"_Not_Required", func(t *testing.T) {
			t.Parallel()

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : no Flags Set Required")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(false, false, false, false), append([]string{ShellCompRequestCmd}, tc.input...)...)
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
			t.Parallel()

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag 1 Set Required")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(true, false, false, false), append([]string{ShellCompRequestCmd}, tc.input...)...)
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
			t.Parallel()

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag 1 & 2 Set Mutually Exclusive")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(false, true, false, false), append([]string{ShellCompRequestCmd}, tc.input...)...)
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
			t.Parallel()

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag 1 & 2 Set Mutually Exclusive & Flag 1 Set Required")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(true, true, false, false), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if output != tc.expectedMutualRequired {
				t.Logf("Expected:\n%+v", tc.expectedMutualRequired)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			}
		})

		// OneRequired set to true
		t.Run(tc.name+"s_OneRequired", func(t *testing.T) {
			t.Parallel()

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag 1 & 2 Set One Required")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(false, false, true, false), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if output != tc.expectedOneRequired {
				t.Logf("Expected:\n%+v", tc.expectedOneRequired)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			}
		})

		// mutual & OneRequired set to true
		t.Run(tc.name+"s_Mutual_OneRequired", func(t *testing.T) {
			t.Parallel()

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag 1 & 2 Set Mutually Exclusive & OneRequired")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(false, true, true, false), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if output != tc.expectedMutualOneRequired {
				t.Logf("Expected:\n%+v", tc.expectedMutualOneRequired)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			}
		})

		// Required Together set to true
		t.Run(tc.name+"s_RequiredTogether", func(t *testing.T) {
			t.Parallel()

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag 1 & 3 Set Required Together")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(false, false, false, true), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if output != tc.expectedRequiredTogether {
				t.Logf("Expected:\n%+v", tc.expectedRequiredTogether)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			}
		})

		// Required Together and Required set to true
		t.Run(tc.name+"s_RequiredTogether_Required", func(t *testing.T) {
			t.Parallel()

			t.Logf("Running test: %v", tc.name)
			t.Logf("Version     : Flag 1 & 3 Set Required Together & Flag 1 Set Required")
			t.Logf("Input       : %v", tc.input)

			output, err := executeCommand(getCmd(true, false, false, true), append([]string{ShellCompRequestCmd}, tc.input...)...)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if output != tc.expectedRequiredTogetherAndRequired {
				t.Logf("Expected:\n%+v", tc.expectedRequiredTogetherAndRequired)
				t.Logf("Got:\n%+v", output)
				t.Fail()
			}
		})
	}
}
