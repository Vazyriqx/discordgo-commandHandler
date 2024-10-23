package commandHandler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	flag "github.com/spf13/pflag"
)

// generate the commands help embed
func generateHelpEmbed(cmd *Command, groups []string) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	embed := &cmd.HelpEmbed

	if !cmd.HasCommandPermission(groups) {
		return
	}

	// set title if unset
	if embed.Title == "" {
		embed.Title = fmt.Sprintf("Help: %s ", cmd.Command)
	}

	// set description if unset
	if embed.Description == "" {
		embed.Description = cmd.LongDesc
	}

	// set default colour to magenta
	// later plan to make this set better
	if embed.Color == 0 {
		embed.Color = 0xFF00FF
	}

	embed.Fields = make([]*discordgo.MessageEmbedField, 0)

	addUse(embed, cmd)

	AddAliases(embed, cmd)

	addFlags(embed, cmd, cmd.flags.flags)

	addChildCommands(embed, cmd, groups)

}

func addUse(embed *discordgo.MessageEmbed, cmd *Command) {
	if cmd.Use != "" {
		usageField := discordgo.MessageEmbedField{
			Name:   "Usage",
			Value:  fmt.Sprintf("`%s`", cmd.Use),
			Inline: false,
		}
		embed.Fields = append(embed.Fields, &usageField)
	}
}

func AddAliases(embed *discordgo.MessageEmbed, cmd *Command) {
	if len(cmd.aliases) == 0 {
		return
	}
	var aliases strings.Builder
	for _, a := range cmd.aliases {
		aliases.WriteString(fmt.Sprintf("`%s` ", a))
	}
	addField(embed, "Aliases", aliases.String(), false)
}

func addFlags(embed *discordgo.MessageEmbed, cmd *Command, flags *flag.FlagSet) {
	if cmd.flags.flags == nil {
		return
	}
	var flagsStr strings.Builder
	var Usage strings.Builder
	flags.VisitAll(func(f *flag.Flag) {
		if f.Shorthand != "" {
			flagsStr.WriteString(fmt.Sprintf("`-%s, --%-15s %s`\n", f.Shorthand, f.Name, f.Usage))
		} else {
			flagsStr.WriteString(fmt.Sprintf("`    --%-15s %s`\n", f.Name, f.Usage))
		}
		Usage.WriteString(fmt.Sprintf("%-30s \n", f.Usage))
	})
	addField(embed, "Flags", flagsStr.String(), true)
}

func addChildCommands(embed *discordgo.MessageEmbed, cmd *Command, groups []string) {
	if len(cmd.commands) == 0 {
		return
	}
	var commands strings.Builder
	for _, child := range cmd.commands {
		if child.Hidden {
			continue
		} else if child.HasCommandPermission(groups) {
			commands.WriteString(fmt.Sprintf("`%s` - %s\n", child.Command, child.ShortDesc))
		}

	}
	addField(embed, "Commands", commands.String(), false)
}

func addField(embed *discordgo.MessageEmbed, name, value string, inline bool) {
	field := discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
	embed.Fields = append(embed.Fields, &field)
}
