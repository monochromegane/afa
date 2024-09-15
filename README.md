# AFA (AI for All)

$$\text{AI for All} := \forall x \in X \, \exists \mathrm{AI}(x)$$

AFA is a terminal-friendly AI command. It enables new behaviors through ad-hoc prompts without requiring programming implementation.

AFA processes text streams as its input and output. Text is a flexible and universal interface, making it easy to interact across various environments and tools.
Additionally, it integrates with rich TUI tools to provide an interactive chat experience in the terminal.

With AFA, let's collaborate with both existing and unknown commands in line with the [UNIX philosophy](https://en.wikipedia.org/wiki/Unix_philosophy).

- *Prompt* programs that do one thing and do it well.
- *Prompt* programs to work together.
- *Prompt* programs to handle text streams, because that is a universal interface.

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

### Default Options

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

## Practical Examples

### Command Suggestions using ZLE

Firstly, we prepare a wrapper for suggestions named `afa_command_suggestion.zsh`, as shown below:

```zsh
#!/bin/zsh

prompt_=""

while getopts p: OPT
do
  case $OPT in
    "p" ) prompt_="$OPTARG";;
    *) echo "Error: Invalid option." >&2; exit 1;;
  esac
done

shift `expr $OPTIND - 1`

suggested_command=$(afa new -script -Q -j command_suggestion -u command_suggestion -p "$prompt_")
if [ $? -ne 0 ]; then
  echo "Error: Failed to generate suggested command." >&2
  exit 1
fi

command_new=$(printf "%s" "$suggested_command" | jq '. | fromjson' | jq -r '.suggested_command')
if [ $? -ne 0 ]; then
  echo "Error: No suggested command received." >&2
  exit 1
fi

echo "$command_new"
```

And we also prepare user prompt and schema files.

`CONFIG_PATH/afa/templates/user/command_suggestion.json`

```json
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
```

`CONFIG_PATH/afa/schemas/command_suggestion.tmpl`

~~~markdown
You are an assistant supporting operations in the terminal. Please suggest commands based on the following requirements.

## Objective

{{ .Message }}

## Requirements

- Commands executable in macOS's zsh.
- Suggest command only, without explanations or comments.
- The output should follow the provided json_schema.

## json_schema

```json
{
  "suggested_command": "<Final suggested command>"
}
```
~~~

Next, we prepare a function for ZLE:

```zsh
function _afa-suggest-command() {
  local command=$(afa_command_suggestion.zsh -p "$BUFFER")
  if [ -n "$command" ]; then
    BUFFER="$command"
  fi
  CURSOR=$#BUFFER
  zle reset-prompt
}
```

Finally, we add a function setting and key bindings to the `.zshrc` as follows:

```zsh
autoload -Uz _afa-suggest-command

zle -N afa-suggest-command _afa-suggest-command
bindkey '^G^K' afa-suggest-command
```

### Code Suggestions using Vim

First, we prepare a wrapper for suggestions named `afa_code_suggestion.zsh`, as shown below:

```zsh
#!/bin/zsh

file_path=""
file_type=""
prompt_=""

while getopts f:e:p: OPT
do
  case $OPT in
    "f" ) file_path="$OPTARG";;
    "e" ) file_type="$OPTARG";;
    "p" ) prompt_="$OPTARG";;
    *) echo "Error: Invalid option." >&2; exit 1;;
  esac
done

shift `expr $OPTIND - 1`

code_org=$(cat)
suggested_code=$(echo "$code_org" | afa new -script -Q -j code_suggestion -u code_suggestion -p "$prompt_" <(echo "$file_path") <(echo "$file_type") "$@")
if [ $? -ne 0 ]; then
  echo "Error: Failed to generate suggested code." >&2
  echo "$code_org"
  exit 1
fi

code_new=$(printf "%s" "$suggested_code" | jq '. | fromjson' | jq -r '.suggested_code')
if [ $? -ne 0 ]; then
  echo "Error: No suggested code received." >&2
  echo "$code_org"
  exit 1
fi

echo "$code_new"
```

And we also prepare user prompt and schema files.

`CONFIG_PATH/afa/templates/user/code_suggestion.json`

```json
{
  "type": "object",
  "properties": {
    "suggested_code": {
      "type": "string"
    }
  },
  "additionalProperties": false,
  "required": [
    "suggested_code"
  ]
}
```

`CONFIG_PATH/afa/schemas/code_suggestion.tmpl`

~~~markdown
You are an assistant supporting coding tasks. Please suggest code modifications or generate new code based on the following requirements.

{{- if .MessageStdin }}

## Current Code

- File: {{ (index .Files 0).Content -}}
- Language: {{ (index .Files 1).Content -}}

```
{{ .MessageStdin }}```
{{ end -}}

{{- if ge (len .Files) 3 }}
## Content of Related Files
{{ range $i, $f := .Files }}
{{- if ge $i 2 }}
- File: {{ $f.Name }}
```
{{ $f.Content }}```
{{ end -}}
{{ end }}
{{ end -}}

## Objective

{{ .Message }}

## Requirements

- Maintain the current structure of the existing functions and classes as much as possible.
- Necessary libraries are already installed, so installation steps are not required.
- Suggest code only, without explanations or comments.
- The output should follow the provided json_schema.

## json_schema

```json
{
  "suggested_code": "<Final suggested code>"
}
```
~~~

Finally, we add a command for Vim:

```vimrc
set splitright
command DiffOrig vert new | set bt=nofile | r ++edit # | 0d_
      \ | diffthis | wincmd p | diffthis

command -nargs=* -range=% -complete=file Afa <line1>,<line2>call AfaFn(<f-args>)
function AfaFn(...) range
  let user_input = input("Enter prompt: ")
  redraw
  let cmd = a:firstline . ',' . a:lastline . '! afa_code_suggestion.zsh -f % -e %:e -p ' . shellescape(user_input) . ' ' . join(a:000, ' ')
  execute cmd
endfunction
```

### Error Message Explanations using TMUX Capture Panel Function

We prepare a Zsh function to capture the output of the last command:

```zsh
function _afa-capture() {
    local start_line=-100
    local last_command=$(fc -l -n -1)
    local capture=$(tmux capture-pane -S "$start_line" -p)

    if match=$(echo "$capture" | grep -F -n "$last_command"); then
      local last_line_num=$(echo "$match" | tail -n 1 | cut -d":" -f1)
      local result=$(echo "$capture" | awk -v start="$last_line_num" 'NR >= start')
      echo $result | afa new -u explain
    else
      echo "\"$last_command\" not found in the capture panel."
    fi

    zle reset-prompt
}
```

Here is an example of a general `explain` user prompt template:

~~~markdown
Please explain the following commands and their results, as well as the content of the provided files. Additionally, provide solutions if necessary.

{{ .Message }}
{{ if .MessageStdin }}
```
{{ .MessageStdin }}```
{{- end }}
{{ range .Files }}
- File: {{ .Name }}
```
{{ .Content }}```
{{ end -}}
~~~

And we add a function setting and key bindings to the `.zshrc` as follows:

```zsh
autoload -Uz _afa-capture

zle -N afa-capture _afa-capture
bindkey '^G^E' afa-capture
```

### GitHub Pull Request Content Suggestions using ZLE and gh Command

Firstly, we prepare a wrapper for suggestions named `afa_github_pull_request.zsh`, as shown below:

```zsh
#!/bin/zsh

prompt_=""

while getopts p: OPT
do
  case $OPT in
    "p" ) prompt_="$OPTARG";;
    *) echo "Error: Invalid option." >&2; exit 1;;
  esac
done

shift `expr $OPTIND - 1`

pull_request_template=".github/pull_request_template.md"

current_branch=$(git branch --show-current)
git fetch --quiet
if ! git diff --quiet HEAD origin/"$current_branch"; then
  echo "You have unsynced changes. Please push to the remote." >&2
  exit 1
fi

pull_request=$(afa new -script -Q -u github_pull_request -j github_pull_request -p "$prompt_" <( git diff --no-ext-diff origin ) <( git log --format="- %s" --no-merges origin..HEAD ))
if [ $? -ne 0 ]; then
  echo "Error: Failed to generate suggested github pull request." >&2
  exit 1
fi

title=$(printf "%s" "$pull_request" | jq '. | fromjson' | jq -r '.title_for_github_pull_request')
if [ $? -ne 0 ]; then
  echo "Error: No suggested github pull request received." >&2
  exit 1
fi

body=$(printf "%s" "$pull_request" | jq '. | fromjson' | jq -r '.body_for_github_pull_request')
if [ $? -ne 0 ]; then
  echo "Error: No suggested github pull request received." >&2
  exit 1
fi

if [[ -f "$pull_request_template" ]]; then
  template_body=$(<"$pull_request_template")
  body_with_template=$(printf "%s\n\n---\n%s" "$body" "$template_body")
else
  body_with_template=$body
fi

gh pr create --web --title="$title" --body="$body_with_template"
```

And we also prepare user prompt and schema files.

`CONFIG_PATH/afa/templates/user/github_pull_request.json`

```json
{
  "type": "object",
  "properties": {
    "title_for_github_pull_request": {
      "type": "string"
    },
    "body_for_github_pull_request": {
      "type": "string"
    }
  },
  "additionalProperties": false,
  "required": [
    "title_for_github_pull_request",
    "body_for_github_pull_request"
  ]
}
```

`CONFIG_PATH/afa/schemas/github_pull_request.tmpl`

~~~markdown
Based on the following information, please propose a title and body for a GitHub pull request.

# Summary of Changes (from git diff information):

{{ (index .Files 0).Content }}
# Related Commit Messages (from git log information):

{{ (index .Files 1).Content }}

{{- if ne .Message "" }}
## Background of Changes

{{ .Message }}
{{ end }}
## Requirements:

- Use Markdown format.
- Focus on the purpose and background rather than the details of the changes.
- Propose only the title and body without extra explanations or comments.
- Start by "# Automatically Generated Pull Request Description".
- And then:
- "## Summary"
- "## Changes", List multiple key changes. Provide a brief explanation for each key change.
- Write Outlines only.
~~~

Next, we prepare a function for ZLE:

```zsh
function _afa-github-pull-request() {
  afa_github_pull_request.zsh -p "$BUFFER"
  BUFFER=""
  CURSOR=$#BUFFER
  zle reset-prompt
}
```

Finally, we add a function setting and key bindings to the `.zshrc` as follows:

```zsh
autoload -Uz _afa-github-pull-request

zle -N afa-github-pull-request _afa-github-pull-request
bindkey '^G^P' afa-github-pull-request
```

## License

[MIT](https://github.com/monochromegane/afa/blob/master/LICENSE)

## Author

[monochromegane](https://github.com/monochromegane)
