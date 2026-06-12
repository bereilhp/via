package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultModel = "openai/gpt-4.1-mini"
	apiURL       = "https://ai-gateway.vercel.sh/v1/chat/completions"
	appVersion   = "0.1.0"
	systemPrompt = "You are running in a terminal. Respond in plain text only. Do not use Markdown formatting."
)

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []choice `json:"choices"`
}

type choice struct {
	Message *responseMessage `json:"message"`
}

type responseMessage struct {
	Content *string `json:"content"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	var model string
	var help bool
	var version bool

	flags := flag.NewFlagSet("via", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.StringVar(&model, "m", defaultModel, "model to use")
	flags.StringVar(&model, "model", defaultModel, "model to use")
	flags.BoolVar(&help, "h", false, "show help")
	flags.BoolVar(&help, "help", false, "show help")
	flags.BoolVar(&version, "v", false, "show version")
	flags.BoolVar(&version, "version", false, "show version")

	if err := flags.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		printUsage(stderr)
		return 1
	}

	if help {
		printUsage(stdout)
		return 0
	}

	if version {
		fmt.Fprintf(stdout, "via %s\n", appVersion)
		return 0
	}

	prompt := strings.Join(flags.Args(), " ")
	if prompt == "" {
		printUsage(stderr)
		return 1
	}

	apiKey := os.Getenv("AI_GATEWAY_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(stderr, `Missing AI_GATEWAY_API_KEY.
Set it with:
  export AI_GATEWAY_API_KEY="your_key_here"`)
		return 1
	}

	reqBody, err := json.Marshal(chatRequest{
		Model: model,
		Messages: []message{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	})
	if err != nil {
		fmt.Fprintf(stderr, "Failed to build request: %v\n", err)
		return 1
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(reqBody))
	if err != nil {
		fmt.Fprintf(stderr, "Failed to build request: %v\n", err)
		return 1
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	latency := time.Since(start)
	if err != nil {
		fmt.Fprintf(stderr, "Request failed: %v\n", err)
		return 1
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(stderr, "Request failed: %v\n", err)
		return 1
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fmt.Fprintf(stderr, "Request failed: HTTP %d\n%s", resp.StatusCode, string(respBody))
		if !strings.HasSuffix(string(respBody), "\n") {
			fmt.Fprintln(stderr)
		}
		return 1
	}

	var parsed chatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		fmt.Fprintf(stderr, "Failed to parse response: %v\n", err)
		return 1
	}

	if len(parsed.Choices) == 0 || parsed.Choices[0].Message == nil || parsed.Choices[0].Message.Content == nil {
		fmt.Fprintln(stderr, "No message content found in response.")
		return 1
	}

	fmt.Fprint(stdout, *parsed.Choices[0].Message.Content)
	if !strings.HasSuffix(*parsed.Choices[0].Message.Content, "\n") {
		fmt.Fprintln(stdout)
	}
	fmt.Fprintf(stdout, "\nmodel: %s\nlatency: %s\n", model, latency.Round(time.Millisecond))

	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Run any Vercel AI Gateway model from your terminal.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  via [flags] <prompt>")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, `  via "Why is the sky blue?"`)
	fmt.Fprintln(w, `  via -m openai/gpt-5-mini "Explain DNS in one paragraph"`)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintf(w, "  -m, --model <model>\n      Override model. Default: %s\n", defaultModel)
	fmt.Fprintln(w, "  -h, --help")
	fmt.Fprintln(w, "      Show help.")
	fmt.Fprintln(w, "  -v, --version")
	fmt.Fprintln(w, "      Show version.")
}
