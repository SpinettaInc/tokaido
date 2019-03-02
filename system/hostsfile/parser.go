// Originally by @lextoumbourou https://github.com/lextoumbourou/goodhosts
// Minor modifications in place to suit Tokaido

package hostsfile

import (
	"github.com/ironstar-io/tokaido/utils"

	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const commentChar string = "#"

// HostsLine - Represents a single line in the hosts file.
type HostsLine struct {
	IP    string
	Hosts []string
	Raw   string
	Err   error
}

// IsComment - Return ```true``` if the line is a comment.
func (l HostsLine) IsComment() bool {
	trimLine := strings.TrimSpace(l.Raw)
	isComment := strings.HasPrefix(trimLine, commentChar)
	return isComment
}

// NewHostsLine - Return a new instance of ```HostsLine```.
func NewHostsLine(raw string) HostsLine {
	fields := strings.Fields(raw)
	if len(fields) == 0 {
		return HostsLine{Raw: raw}
	}

	output := HostsLine{Raw: raw}
	if !output.IsComment() {
		rawIP := fields[0]
		if net.ParseIP(rawIP) == nil {
			output.Err = fmt.Errorf("Bad hosts line: %q", raw)
		}

		output.IP = rawIP
		output.Hosts = fields[1:]
	}

	return output
}

// Hosts - Represents a hosts file.
type Hosts struct {
	Path  string
	Lines []HostsLine
}

// IsWritable - Return ```true``` if hosts file is writable.
func (h *Hosts) IsWritable() bool {
	_, err := os.OpenFile(h.Path, os.O_WRONLY, 0660)
	if err != nil {
		return false
	}

	return true
}

// Load the hosts file into ```l.Lines```.
// ```Load()``` is called by ```NewHosts()``` and ```Hosts.Flush()``` so you
// generally you won't need to call this yourself.
func (h *Hosts) Load() error {
	var lines []HostsLine

	file, err := os.Open(h.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := NewHostsLine(scanner.Text())
		if err != nil {
			return err
		}

		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	h.Lines = lines

	return nil
}

// Flush any changes made to hosts file.
func (h Hosts) Flush() error {
	file, err := os.Create(h.Path)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)

	for _, line := range h.Lines {
		fmt.Fprintf(w, "%s%s", line.Raw, eol)
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return h.Load()
}

// WriteElevated - Request a password then write the new hostsfile
func (h *Hosts) WriteElevated() error {
	var s string
	for _, line := range h.Lines {
		s = s + line.Raw + "\n"
	}
	utils.BashStringCmd("echo \"" + s + "\" | sudo tee " + h.Path)

	return nil
}

// Add an entry to the hosts file.
func (h *Hosts) Add(ip string, hosts ...string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("%q is an invalid IP address", ip)
	}

	position := h.getIPPosition(ip)
	if position == -1 {
		endLine := NewHostsLine(buildRawLine(ip, hosts))
		// Ip line is not in file, so we just append our new line.
		h.Lines = append(h.Lines, endLine)
	} else {
		// Otherwise, we replace the line in the correct position
		newHosts := h.Lines[position].Hosts
		for _, addHost := range hosts {
			if itemInSlice(addHost, newHosts) {
				continue
			}

			newHosts = append(newHosts, addHost)
		}
		endLine := NewHostsLine(buildRawLine(ip, newHosts))
		h.Lines[position] = endLine
	}

	return nil
}

// Has - Return a bool if ip/host combo in hosts file.
func (h Hosts) Has(ip string, host string) bool {
	pos := h.getHostPosition(ip, host)

	return pos != -1
}

// Remove an entry from the hosts file.
func (h *Hosts) Remove(ip string, hosts ...string) error {
	var outputLines []HostsLine

	if net.ParseIP(ip) == nil {
		return fmt.Errorf("%q is an invalid IP address", ip)
	}

	for _, line := range h.Lines {

		// Bad lines or comments just get readded.
		if line.Err != nil || line.IsComment() || line.IP != ip {
			outputLines = append(outputLines, line)
			continue
		}

		var newHosts []string
		for _, checkHost := range line.Hosts {
			if !itemInSlice(checkHost, hosts) {
				newHosts = append(newHosts, checkHost)
			}
		}

		// If hosts is empty, skip the line completely.
		if len(newHosts) > 0 {
			newLineRaw := line.IP

			for _, host := range newHosts {
				newLineRaw = fmt.Sprintf("%s %s", newLineRaw, host)
			}
			newLine := NewHostsLine(newLineRaw)
			outputLines = append(outputLines, newLine)
		}
	}

	h.Lines = outputLines
	return nil
}

func (h Hosts) getHostPosition(ip string, host string) int {
	for i := range h.Lines {
		line := h.Lines[i]
		if !line.IsComment() && line.Raw != "" {
			if ip == line.IP && itemInSlice(host, line.Hosts) {
				return i
			}
		}
	}

	return -1
}

func (h Hosts) getIPPosition(ip string) int {
	for i := range h.Lines {
		line := h.Lines[i]
		if !line.IsComment() && line.Raw != "" {
			if line.IP == ip {
				return i
			}
		}
	}

	return -1
}

// NewHosts - Return a new instance of ``Hosts``.
func NewHosts() (Hosts, error) {
	osHostsFilePath := os.ExpandEnv(filepath.FromSlash(hostsFilePath))

	hosts := Hosts{Path: osHostsFilePath}

	err := hosts.Load()
	if err != nil {
		return hosts, err
	}

	return hosts, nil
}

func itemInSlice(item string, list []string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}

	return false
}

func buildRawLine(ip string, hosts []string) string {
	output := ip
	for _, host := range hosts {
		output = fmt.Sprintf("%s %s", output, host)
	}

	return output
}
