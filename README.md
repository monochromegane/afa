# AFA (AI for All)

$$\text{AI for All} \equiv \forall x \in X \, \exists \mathrm{AI}(x)$$

AFA is a terminal-friendly AI client that offers a rich UI for interactive chat and seamless integration with various commands through standard input and output.

## Usage

Run the interactive chat with:  

```sh
afa new
```  

Use a rich TUI viewer in chat mode with:  

```sh
afa new -V
```  

Start the interactive chat with extra information by running:  

```sh
afa new -p "What is this?" /path/to/file
```  

Continue from the last session with:  

```sh
afa resume
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
