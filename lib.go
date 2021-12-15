package ytcaps2srt

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type TimedText struct {
	Content []*Text `xml:"text"`
}

type Text struct {
	Duration int    `xml:"d,attr"`
	Time     int    `xml:"t,attr"`
	Value    string `xml:",chardata"`

	// Old caps have this set to 1, new caps don't have the attribute
	Append int `xml:"append,attr"`
}

func (t *TimedText) Beautify() ([]*Text, error) {
	var (
		err             error
		result          = make([]*Text, 0)
		tmpText         *Text
		tmpEnd, tmpCEnd int
	)

	for i, text := range t.Content {
		if tmpText == nil {
			tmpText = &Text{
				Duration: text.Duration,
				Time:     text.Time,
				Value:    strings.TrimSpace(text.Value),
			}
			tmpEnd = tmpText.Time + tmpText.Duration
		} else {
			tmpCEnd = text.Time

			if tmpCEnd < tmpEnd {
				tmpText.Value += " " + strings.TrimSpace(text.Value)
			} else {
				tmpText.Value = strings.ReplaceAll(tmpText.Value, "  ", " ")
				result = append(result, tmpText)
				tmpText = &Text{
					Duration: text.Duration,
					Time:     text.Time,
					Value:    strings.TrimSpace(text.Value),
				}
				tmpEnd = tmpText.Time + tmpText.Duration
			}
		}

		if i == len(t.Content)-1 {
			tmpText.Value = strings.ReplaceAll(tmpText.Value, "  ", " ")
			result = append(result, tmpText)
		}
	}

	return result, err
}

func ParseTimedText(content []byte) (*TimedText, error) {
	var (
		err    error
		result = new(TimedText)
	)

	err = xml.Unmarshal(content, result)

	return result, err
}

func ConvertToSRT(in []*Text) (string, error) {
	var (
		err    error
		result = new(strings.Builder)
	)

	timeFMT := "15:04:05,000"

	for i, text := range in {
		// N
		_, err = result.WriteString(strconv.Itoa(i+1) + "\n")
		// 00:00:49,514 --> 00:00:50,427
		from := time.Time{}.Add(time.Millisecond * time.Duration(text.Time)).Format(timeFMT)
		until := time.Time{}.Add(time.Millisecond * time.Duration(text.Time+text.Duration)).Format(timeFMT)
		_, err = result.WriteString(fmt.Sprintf("%s --> %s\n", from, until))
		// DATA
		// Per line 47 chars, see why: https://stackoverflow.com/a/69772685
		for ii, letter := range text.Value {
			_, err = result.WriteRune(letter)
			if ii+1%47 == 0 {
				_, err = result.WriteString("\n\n")
			}
		}
		if len(text.Value)%47 != 0 {
			_, err = result.WriteString("\n\n")
		}
	}

	return result.String(), err
}
