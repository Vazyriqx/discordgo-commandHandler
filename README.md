# Command Handler

This project provides a customizable command handler for Discord bots using the `discordgo` library. It allows for easy command creation, aliasing, permission grouping, and automatic help message generation.

## Features

- **Command Execution**: Supports executing commands with predefined arguments, aliases, and permission groups.
- **Help Embed Generation**: Automatically generates help messages for commands.
- **Flag Parsing**: Uses `pflag` for handling command-specific flags CLI style.
- **Command Hooks**: Supports pre-run and post-run command hooks.
- **Permission Groups**: Command permission groups only allow members with the correct permission to execute and see commands in the help embed.
- **Hidden Commands**: Commands can be hidden from users in the help embed while still allowing execution of the command.

There is no support for slash commands. 

## Usage

### Command Structure

Each command can have:
- A prefix (root-level only).
- A command name and optional aliases.
- A short and long description for help generation.
- Pre-run and post-run hooks.
- Permission groups to restrict access.
- Flag parsing using `pflag`.

### Example command

```go
cmd := &Command{
    Command:   "example",
    Use:       "example [args]",
    ShortDesc: "This is an example command.",
    LongDesc:  "This is a more detailed description of the example command.",
    PreRunE:   commandHandler.AddRuns(somefunc1, sumfunc2),
    RunE: func(cmd *Command, args []string, s *discordgo.Session, m *discordgo.MessageCreate) error {
        s.ChannelMessageSend(m.ChannelID, "Example command executed!")
        return nil
    },
    PostRunE:  commandHandler.AddRuns(somefunc3),
}
```
# Disclaimer
This code is in very early stages of development. Things are likely to change. This command handler is also not designed to be used for larger bots and does not provide some of the customisations you would want.
