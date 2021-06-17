Enhancement: Support append() in cli.Command.Flags (Breaking Change)

Some oCIS extensions use multiple `[]cli.Flag` slices in `cli.Command{Flags: xxx}`
by appending them in place. The support for this to work is added in this PR. In order to 
return multiple flagset names, the type of ParsedCommand.Flags is changed
from string to []string. This breaks existing templates and the user needs to update them.

https://github.com/owncloud/flaex/pull/9
