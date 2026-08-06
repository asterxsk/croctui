package croc

import (
	"regexp"
	"strconv"
	"strings"
)

type EventKind int

const (
	EventStatus EventKind = iota
	EventCode
	EventFile
	EventProgress
	EventError
)

type Event struct {
	Kind    EventKind
	Code    string
	File    string
	Rate    string
	Text    string
	Percent float64 // -1 when unknown
}

var (
	codeRe    = regexp.MustCompile(`Code is:\s+(\S+)`)
	fileRe    = regexp.MustCompile(`(?i)(Sending|Receiving)\s+'([^']+)'`)
	errorRe   = regexp.MustCompile(`(?i)(could not connect|bad response|error:|connection failed|transfer failed)`)
	percentRe = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*%`)
	rateRe    = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*[KMGT]?B/s`)
)

// ParseLine parses a single line/fragment of croc output into an Event.
func ParseLine(line string) Event {
	if code := codeRe.FindStringSubmatch(line); code != nil {
		return Event{Kind: EventCode, Code: code[1], Text: line, Percent: -1}
	}
	if f := fileRe.FindStringSubmatch(line); f != nil {
		return Event{Kind: EventFile, File: f[2], Text: line, Percent: -1}
	}
	if errorRe.MatchString(line) {
		return Event{Kind: EventError, Text: line, Percent: -1}
	}
	if p := percentRe.FindStringSubmatch(line); p != nil {
		pct, _ := strconv.ParseFloat(p[1], 64)
		rate := ""
		if r := rateRe.FindString(line); r != "" {
			rate = r
		}
		return Event{Kind: EventProgress, Percent: pct, Rate: rate, Text: line}
	}
	return Event{Kind: EventStatus, Text: line, Percent: -1}
}

// ParseChunk splits raw output on \r and \n, drops blank/whitespace
// fragments, and returns an Event per fragment.
func ParseChunk(chunk string) []Event {
	var out []Event
	for _, frag := range strings.FieldsFunc(chunk, func(r rune) bool { return r == '\r' || r == '\n' }) {
		if strings.TrimSpace(frag) == "" {
			continue
		}
		out = append(out, ParseLine(frag))
	}
	return out
}