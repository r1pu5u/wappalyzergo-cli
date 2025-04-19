# wappalyzergo-cli

`wappalyzergo-cli` is a command-line tool for detecting technologies used by websites. It leverages the high-performance [WappalyzerGo](https://github.com/projectdiscovery/wappalyzergo) library to analyze HTTP headers and HTML content, identifying technologies such as frameworks, platforms, and server software.

## Features

- Concurrent scanning with configurable thread count
- Custom User-Agent support
- JSON output for easy integration
- Built-in Wappalyzer fingerprinting
- Progress logging for each processed URL

## Installation

Ensure you have Go installed (version 1.16 or later).

```sh
go install github.com/yourname/wappalyzergo-cli@latest
```


Replace `yourname` with your actual GitHub username.

## Usage

```sh
wappalyzergo-cli -l <list_file> -t <threads> -o <output_file>
```


### Parameters

- `-l`: Path to the file containing a list of URLs (one per line).
- `-t`: Number of concurrent threads (default is 5).
- `-o`: Output file path for the JSON results.

### Example

```sh
wappalyzergo-cli -l urls.txt -t 10 -o results.json
```


## Output

The output is a JSON file mapping each URL to its detected technologies.

Example:

```json
{
  "https://example.com": {
    "Amazon EC2": {},
    "Cloudflare": {},
    "Nginx": {}
  },
  "https://anotherexample.com": {
    "Apache": {},
    "PHP": {}
  }
}
```

