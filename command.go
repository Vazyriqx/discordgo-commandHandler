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
	"sync"

	"github.com/bwmarrin/discordgo"
	flag "github.com/spf13/pflag"
)

type flagsMutex struct {
	mu    sync.Mutex
	flags *flag.FlagSet
}

type RunE func(cmd *Command, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error

type Command struct {
	// the bots prefix. This should only be set for the root command.
	Prefix byte

	// The command.
	Command string

	// Aliases for the command.
	aliases []string

	// Example usage for the command.
	Use string

	// short decription of the command.
	ShortDesc string

	// Full discription of the command.
	LongDesc string

	// number of args a message should have
	ExactArgs int

	// min https://discord.gg/dVFqNhXc
	MinArgs int

	// max args a command should have
	MaxArgs int

	// The embed shown to a user when --help or -h flag is added.
	// The Embed is generated automatically as long command Use and long desc are added.
	HelpEmbed discordgo.MessageEmbed

	// function to run prior to the command executing.
	PreRunE []RunE

	// command implemtation.
	RunE func(cmd *Command, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error

	// functions to run after the command has executed.
	PostRunE []RunE

	// whether the command should show as a child command in the help message.
	Hidden bool

	// The parent command of a sub command.
	parent *Command

	// child commands.
	commands []*Command

	// flags for the command.
	flags flagsMutex

	// Groups for limiting commands to users of a specific group
	Groups []string
}

// Adds one or more commands to the parent command
func (c *Command) AddCommand(cmds ...*Command) {
	for _, cmd := range cmds {
		if cmd == c {
			panic("Command cannot be a child of itself")
		}
		cmd.parent = c
		cmd.Flags().BoolP("help", "h", false, "Print the help dialog")
		c.commands = append(c.commands, cmd)
	}
}

// Returns []RunE. Used to set PreRunE and PostRunE
func AddRuns(postRuns ...RunE) []RunE {
	return postRuns
}

// Returns []string Used for settings permission groups for a command
func AddGroups(groups ...string) []string {
	return groups
}

// Returns []string Used for adding aliases to a command
func (c *Command) AddAliases(aliases ...string) {
	for _, command := range c.parent.commands {
		for _, a := range command.aliases {
			if contains(aliases, a) {
				panic(fmt.Sprintf("Alias %s redefined in command %s", a, command.Command))
			}
		}
	}
	c.aliases = aliases
}

// Helper function to determin if a string slice contains and item
func contains(slice []string, s string) bool {
	for _, a := range slice {
		if strings.EqualFold(a, s) {
			return true
		}
	}
	return false
}

// returns the flagset for the command
func (c *Command) Flags() *flag.FlagSet {
	if c.flags.flags == nil {
		c.flags.flags = flag.NewFlagSet(c.Command, flag.ContinueOnError)
	}
	return c.flags.flags
}
