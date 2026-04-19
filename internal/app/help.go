package app

import (
	"fmt"
	"io"
)

func printHelp(w io.Writer, version string) {
	fmt.Fprintf(w, `gentle-qa — Gentle-QA (%s)

USAGE
  gentle-qa                     Launch interactive TUI
  gentle-qa <command> [flags]

COMMANDS
  install      Configure AI coding agents on this machine
  uninstall    Remove Gentle-QA managed files from this machine
  sync         Sync agent configs and skills to current version
  update       Check for available updates
  upgrade      Apply updates to managed tools
  restore      Restore a config backup
  version      Print version

FLAGS
  --help, -h    Show this help

Run 'gentle-qa help' for this message.
Documentation: https://github.com/EduardoVeraE/Gentle-QA
`, version)
}
