package croc

import "testing"

func TestParseCode(t *testing.T) {
	ev := ParseLine("Code is: smoketest")
	if ev.Kind != EventCode || ev.Code != "smoketest" {
		t.Fatalf("got %+v, want EventCode{Code: smoketest}", ev)
	}
}

func TestParseFile(t *testing.T) {
	ev := ParseLine("Sending 'crocgui-probe.txt' (12 B)")
	if ev.Kind != EventFile || ev.File != "crocgui-probe.txt" {
		t.Fatalf("got %+v, want EventFile{crocgui-probe.txt}", ev)
	}
}

func TestParseProgress(t *testing.T) {
	ev := ParseLine("Sending 0.5 MB (25%) [===>      ] 1.2 MB/s")
	if ev.Kind != EventProgress {
		t.Fatalf("got kind %v, want EventProgress", ev.Kind)
	}
	if ev.Percent != 25 {
		t.Fatalf("got percent %v, want 25", ev.Percent)
	}
	if ev.Rate != "1.2 MB/s" {
		t.Fatalf("got rate %q, want %q", ev.Rate, "1.2 MB/s")
	}
}

func TestParseError(t *testing.T) {
	ev := ParseLine("could not connect to 142.132.189.179:9009: bad response")
	if ev.Kind != EventError {
		t.Fatalf("got kind %v, want EventError", ev.Kind)
	}
}

func TestParseUnknownIsStatus(t *testing.T) {
	ev := ParseLine("On the other computer run:")
	if ev.Kind != EventStatus {
		t.Fatalf("got kind %v, want EventStatus", ev.Kind)
	}
}

func TestParseChunkRealProbeOutput(t *testing.T) {
	// Raw bytes captured from a real `croc --yes send` on 10.6.0.
	chunk := "\r                                \rSending 0 files (12 B)\r                                \rSending 'crocgui-probe.txt' (12 B)\nCode is: smoketest\n"
	evs := ParseChunk(chunk)
	found := map[EventKind]bool{}
	for _, ev := range evs {
		found[ev.Kind] = true
	}
	if !found[EventCode] {
		t.Fatalf("expected an EventCode in %+v", evs)
	}
	if !found[EventFile] {
		t.Fatalf("expected an EventFile in %+v", evs)
	}
}

func TestParseChunkDropsBlanks(t *testing.T) {
	evs := ParseChunk("\r  \r\n\r\n")
	if len(evs) != 0 {
		t.Fatalf("expected 0 events, got %d: %+v", len(evs), evs)
	}
}

func TestSendArgs(t *testing.T) {
	got := SendArgs([]string{"a.txt", "b.txt"})
	want := []string{"--yes", "send", "a.txt", "b.txt"}
	if len(got) != len(want) {
		t.Fatalf("SendArgs = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("SendArgs = %v, want %v", got, want)
		}
	}
}

func TestReceiveArgs(t *testing.T) {
	got := ReceiveArgs("mycode", `C:\Users\me\Downloads`)
	want := []string{"--yes", "--rename", "mycode", "--out", `C:\Users\me\Downloads`}
	if len(got) != len(want) {
		t.Fatalf("ReceiveArgs = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ReceiveArgs = %v, want %v", got, want)
		}
	}
}