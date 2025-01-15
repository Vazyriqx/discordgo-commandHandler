// Copyright 2024 Vazyriqx
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

package commandHandler

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	flag "github.com/spf13/pflag"
)

// Entry function for the bot to process a command
func (c *Command) ProcessMessage(s *discordgo.Session, m *discordgo.MessageCreate, groups []string) error {
	if len(m.Content) <= 1 || m.Content[0] != c.Prefix {
		// Check if it is potentially a valid command.
		return nil
	}
	// this is simply so iOS users don't run into issues with -- being autoreplaced with —
	fixedIosString := strings.Replace(m.Content, "—", "--", -1)
	args := strings.Fields(fixedIosString)
	args[0] = args[0][1:]
	return c.execute(args, s, m, groups, nil)
}

// Executes the commands. You should be using processess message in most cases
func (c *Command) execute(args []string, s *discordgo.Session, m *discordgo.MessageCreate, groups []string, flags *flag.FlagSet) error {
	defer c.resetFlags() // Make sure flags are reset when leaving a command

	cmdName := args[0]
	args = args[1:]

	// TO BE IMPROVED
	if cmdName == "help" {
		generateHelpEmbed(c, groups)
		s.ChannelMessageSendEmbed(m.ChannelID, &c.HelpEmbed)
		return nil
	}

	for _, cmd := range c.commands {
		if !isCommand(cmd, cmdName) || !cmd.HasCommandPermission(groups) {
			continue
		}

		cmd.flags.mu.Lock()
		defer cmd.flags.mu.Unlock()

		// Innitliase flags if not already innitialised
		if flags == nil {
			flags = flag.NewFlagSet(cmd.Command, flag.ContinueOnError)
			if cmd.flags.flags != nil {
				flags.AddFlagSet(cmd.flags.flags)
			}
			flags.Parse(args)
			args = flags.Args()
		}

		// Check for child commands
		if len(args) > 0 && len(cmd.commands) > 0 {
			for _, command := range cmd.commands {
				if isCommand(command, args[0]) {
					return cmd.execute(args, s, m, groups, flags)
				}
			}
		}

		defer cmd.resetFlags()

		// Print Help diaglog
		if helpflag, _ := flags.GetBool("help"); helpflag {
			generateHelpEmbed(cmd, groups)
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, &cmd.HelpEmbed)
			return err
		}

		if err := validateArgs(cmd, args); err != nil {
			return err
		}

		if err := runHooks(cmd.PreRunE, cmd, args, s, m); err != nil {
			return err
		}
		if err := cmd.RunE(cmd, args, s, m); err != nil {
			return err
		}

		return runHooks(cmd.PostRunE, cmd, args, s, m)
	}

	return nil
}

func isCommand(cmd *Command, cmdName string) bool {
	if strings.EqualFold(cmd.Command, cmdName) || contains(cmd.aliases, cmdName) {
		return true
	}
	return false
}

func runHooks(hooks []RunE, cmd *Command, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
	for _, hook := range hooks {
		if err := hook(cmd, args, s, m); err != nil {
			return err
		}
	}
	return nil
}

func validateArgs(cmd *Command, args []string) error {
	if cmd.MinArgs != 0 && len(args) < cmd.MinArgs {
		return fmt.Errorf("too few arguments")
	}
	if cmd.MaxArgs != 0 && len(args) > cmd.MaxArgs {
		return fmt.Errorf("too many arguments")
	}
	if cmd.ExactArgs != 0 && len(args) != cmd.ExactArgs {
		return fmt.Errorf("incorrect number of arguments")
	}
	return nil
}

func (c *Command) resetFlags() {
	if c.flags.flags != nil {
		c.flags.flags.VisitAll(func(f *flag.Flag) {
			f.Value.Set(f.DefValue)
		})
	}
}
