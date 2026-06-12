# via

Run any Vercel AI Gateway model from your terminal. 

Built for the Built in London with Vercel and OpenAI event.

## Install

Install the latest published version:

```sh
go install github.com/bereilhp/via@latest
```

Or install from this checkout:

```sh
go install .
```

Go installs the binary into your Go bin directory. If `via` is not found after
installing, add that directory to your `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

To make that permanent for Bash:

```sh
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bash_profile
source ~/.bash_profile
```

Then run:

```sh
via "Why is the sky blue?"
```

Or build a local binary without installing:

```sh
go build -o via .
./via "Why is the sky blue?"
```

## Setup

Set your Vercel AI Gateway API key:

```sh
export AI_GATEWAY_API_KEY="your_key_here"
```

## Usage

```sh
via "Why is the sky blue?"
```

`via` prints the model response, followed by the model name and request latency.
Responses are requested in plain text for terminal display.

Override the model:

```sh
via -m openai/gpt-5-mini "Explain DNS in one paragraph"
```

The default model is:

```text
openai/gpt-4.1-mini
```

## Flags

```text
-m, --model <model>
    Override model.

-h, --help
    Show help.

-v, --version
    Show version.
```

## Notes

This MVP intentionally has no SDK, config file, streaming, chat history, TUI,
agents, or model list command.
