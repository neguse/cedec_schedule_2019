package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var (
	JSTLocation *time.Location
)

func init() {
	var err error
	JSTLocation, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Panic(err)
	}
}

type CedecTime struct {
	time.Time
}

const CedecTimeFormat = `"2006\/01\/02 15:04:05"`

func (t *CedecTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == `""` {
		*t = CedecTime{time.Unix(0, 0)}
		return nil
	}
	nt, err := time.ParseInLocation(CedecTimeFormat, string(data), JSTLocation)
	if err != nil {
		return err
	}
	*t = CedecTime{nt}
	return nil
}

func (t *CedecTime) MarshalJSON() ([]byte, error) {
	d := t.Format(CedecTimeFormat)
	return []byte(d), nil
}

type Session struct {
	Start         CedecTime `json:"start"`
	End           CedecTime `json:"end"`
	Format        string    `json:"format"`
	Title         string    `json:"title"`
	Speakers      []Speaker `json:"speakers"`
	Category      string    `json:"category"`
	SubCategory   []string  `json:"sub_category"`
	Platform      []string  `json:"platform"`
	KeywordTag    []string  `json:"keywordtag"`
	Difficulty    string    `json:"difficulty"`
	Duration      string    `json:"duration"`
	Description   string    `json:"description"`
	Takeaway      string    `json:"takeaway"`
	ExpectedSkill string    `json:"expected_skill"`
	PhotoOk       bool      `json:"photo_ok"`
	SnsOk         bool      `json:"sns_ok"`
	CEDiL         bool      `json:"CEDiL"`
	Note          []string  `json:"note"`
	URL           string    `json:"URL"`
}

type Speaker struct {
	Company string `json:"company"`
	Name    string `json:"name"`
	Profile string `json:"profile"`
	Message string `json:"message"`
}

func CalString(t *time.Time) string {
	return t.UTC().Format(`20060102T150405Z`)
}

func EscapeNL(s string) string {
	return strings.ReplaceAll(s, "\n", `\n`)
}

func SessionSpeaker(s Session) string {
	var ss []string
	for _, sp := range s.Speakers {
		ss = append(ss, fmt.Sprintf("%s(%s)", sp.Name, sp.Company))
	}
	return strings.Join(ss, ",")
}

func SessionDescription(s Session) string {
	return fmt.Sprint(
		s.URL, `\n`,
		// EscapeNL(s.Description), `\n`,
		"スピーカー:", SessionSpeaker(s), `\n`,
	)
}

func main() {
	data, err := ioutil.ReadFile("json_file.json")
	if err != nil {
		log.Panic(err)
	}

	var sessions []Session
	err = json.Unmarshal(data, &sessions)
	if err != nil {
		log.Panic(err)
	}

	const NL = "\r\n"

	fmt.Print("BEGIN:VCALENDAR", NL)
	fmt.Print("VERSION:2.0", NL)
	fmt.Print("PRODID:neguse/cedec_2019_calendar", NL)

	for _, session := range sessions {
		fmt.Print("BEGIN:VEVENT", NL)
		fmt.Printf("UID:%s%s", session.URL, NL) // セッションURLはイベントのIDとしてユニークに使える気がする
		fmt.Printf("URL:%s%s", session.URL, NL)
		fmt.Printf("SUMMARY:%s%s", EscapeNL(session.Title), NL)
		fmt.Printf("DESCRIPTION:%s%s", SessionDescription(session), NL)
		fmt.Printf("DTSTART:%s%s", CalString(&session.Start.Time), NL)
		fmt.Printf("DTEND:%s%s", CalString(&session.End.Time), NL)
		fmt.Print("END:VEVENT", NL)
	}
	fmt.Print("END:VCALENDAR", NL)
}
