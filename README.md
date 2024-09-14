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

### Chat using a Rich TUI

![Chat](examples/chat.gif)

### Command Suggestions using ZLE

![Command Suggestions](examples/command_suggestion.gif)

### Code Suggestions using Vim

![Code Suggestions](examples/code_suggestion.gif)

## Features

- Acts as a terminal-friendly AI command.
- Acts as a chat client with a rich terminal user interface (TUI).
- Supports contextual prompts for both system and user using templates.
- Accepts prompts, standard input, and file paths as context.
- Manages sessions, allowing for quick resumption via the `resume` sub-command.
- Supports structured output with a safely escaped JSON option, facilitating easy integration with other commands.
- The core application operates independently of third-party libraries.
- Supports `OpenAI` as an AI model (support for other AI models is planned for the future).

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

## Configuration

### Default options

The configuration file named `afa/option.json` should be located at the path specified by Go's [os.UserConfigDir](https://pkg.go.dev/os#UserConfigDir).

> On Unix systems, it returns `$XDG_CONFIG_HOME` as specified by [https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) if non-empty, else `$HOME/.config`. On Darwin, it returns `$HOME/Library/Application Support`. On Windows, it returns `%AppData%`. On Plan 9, it returns `$home/lib`.

### API Keys

The configuration file named `CONFIG_PATH/afa/secrets.json`.

### Tempates

AFA supports the use of template files, which can be placed in the `templates/{system,user}` directories with the `.tmpl` extension.
You can specify the template to use by providing the name without the extension using `-s` (for system templates) or `-u` (for user templates) options.

Templates allow you to dynamically insert information. You can utilize the following placeholders within your templates:

- `Message`: A string that can be replaced with a specific message by `-p` opition.
- `MessageStdin`: This placeholder can take input from the standard input as a message.
- `Files`: A collection of file objects, where each file has `Name` and `Content` members.

### Schemas

Similar to templates, schema files can be placed in the `schemas` directory and should have a `.json` extension.
You can specify the schema to use by providing the name without the `.json` extension using the `-j` option.

## Cache

### Sessions

The session files named `afa/sessions/SESSION_NAME.json` should be located at the path specified by Go's [os.UserCacheDir](https://pkg.go.dev/os#UserCacheDir).

> On Unix systems, it returns `$XDG_CACHE_HOME` as specified by [https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) if non-empty, else `$HOME/.cache`. On Darwin, it returns `$HOME/Library/Caches`. On Windows, it returns `%LocalAppData%`. On Plan 9, it returns `$home/lib/cache`.
