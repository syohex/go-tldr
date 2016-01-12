package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
)

const baseURL = `https://raw.githubusercontent.com/rprieto/tldr/master/pages/`

func tldlPage(query string) string {
	return fmt.Sprintf("%s/common/%s.md", baseURL, query)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: tldr query")
		os.Exit(1)
	}

	query := os.Args[1]
	url := tldlPage(query)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	if status == 404 {
		fmt.Fprintf(os.Stderr, "Not found '%s'\n", query)
		os.Exit(1)
	} else if status != 200 {
		fmt.Fprintf(os.Stderr, "Got %d status\n", status)
		os.Exit(1)
	}

	header := color.New(color.BgBlue, color.Bold, color.FgCyan)
	quote := color.New(color.FgWhite)
	list := color.New(color.FgGreen)
	code := color.New(color.BgHiBlack, color.FgWhite)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			fmt.Print("\n")
			continue
		}
		first := line[0]

		switch first {
		case '#':
			header.Print(strings.TrimLeft(line, `#`))
			fmt.Print("\n")
		case '>':
			fmt.Print("  ")
			quote.Println(strings.TrimLeft(line, `> `))
		case '-':
			list.Print(line)
			fmt.Print("\n")
		case '`':
			fmt.Print("  ")
			code.Print(strings.Trim(line, "`"))
			fmt.Print("\n")
		default:
			fmt.Print(line)
		}
	}

	fmt.Print("\n")

	if scanner.Err() != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
