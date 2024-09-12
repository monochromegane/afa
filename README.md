# AFA (AI for All)

$$\text{AI for All} \equiv \forall x \in X \, \exists \mathrm{AI}(x)$$

AFA is a terminal-friendly AI command. It enables new behaviors through ad-hoc prompts without requiring programming implementation.

AFA processes text streams as its input and output. Text is a flexible and universal interface, making it easy to interact across various environments and tools.
Additionally, it integrates with rich TUI tools to provide an interactive chat experience in the terminal.

With AFA, let's collaborate with both existing and unknown commands in line with the UNIX philosophy.

- Prompt programs that do one thing and do it well.
- Prompt programs to work together.
- Prompt programs to handle text streams, because that is a universal interface.

## Demo

### Error Message Explanations

### Command Suggestions

### Code Suggestions for Vim

### Git Commit Message Suggestions

### GitHub Pull Request Content Suggestions

### Interactive Chat with Rich TUI

## Features

- Acts as a terminal-friendly AI command.
- Functions as a chat client with a rich TUI.
- Supports contextual system and user prompts using templates.
- Accepts prompts, standard input, and file paths as context.
- Manages sessions, allowing quick resumption through the `resume` sub-command.
- Supports structured output with a safely escaped JSON option, facilitating easy integration with other commands.

## Usage

Run the interactive chat with:

```sh
afa new
```

Use a rich TUI viewer in chat mode with:

```sh
afa new -V
```

Start the interactive chat with additional information by executing:

```sh
echo $ERROR_MESSAGE | afa new -p "What is happening?" /path/to/file1 /path/to/file2
# Please be cautious; when standard input is provided, interactive mode is disabled.
# Consider using process substitution.
#=> afa new -p "What is happening?" /path/to/file1 /path/to/file2 <( echo $ERROR_MESSAGE )
```

Continue from the last session with:

```sh
afa resume
```

Continue from a specified session with:

```sh
# The command `afa list` displays past sessions.
afa source -l SESSION_NAME
```

Specify the user prompt with:

```sh
# `Message`, `MessageStdin`, and `Files` that include `File` with `Name` and `Content` members can be used in the template file.
echo "Please explain the following.\n{{ (index .Files 0).Content }}" > CONFIG_PATH/templates/user/explain.tmpl
afa -u explain /path/to/file
```

Specify the schema for structured output with:

```sh
cat <<< EOS > CONFIG_PATH/schemas/command_suggestion.json
{
  "type": "object",
  "properties": {
    "suggested_command": {
      "type": "string"
    }
  },
  "additionalProperties": false,
  "required": [
    "suggested_command"
  ]
}
EOS

P="List Go files from the directory named 'internal' and print the first line of each file."
afa new -script -Q -j command_suggestion -p $P | jq '. | fromjson' | jq -r '.suggested_command'
#=> find internal -name '*.go' -exec head -n 1 {} \;
```

## Installation

Follow these steps to install the tool and viewer:

```sh
# Install the core application
go install github.com/monochromegane/afa@latest

# Install the TUI viewer
go install github.com/monochromegane/afa-tui@latest
```

Initialize the setup:

```sh
afa init
```
