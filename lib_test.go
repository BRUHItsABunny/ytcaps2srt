package ytcaps2srt

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"io"
	"os"
	"testing"
)

func TestParseTimedText(t *testing.T) {
	f, _ := os.Open("timedtextauto.xml")
	sample, _ := io.ReadAll(f)
	result, err := ParseTimedText([]byte(sample))
	if err != nil {
		t.Error(err)
	}

	for _, text := range result.Content {
		fmt.Println(spew.Sdump(text))
	}
}

func TestTimedText_Beautify(t *testing.T) {
	f, _ := os.Open("timedtextauto.xml")
	sample, _ := io.ReadAll(f)
	result, err := ParseTimedText(sample)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(len(result.Content))
	// for _, text := range result.Content {
	// 	fmt.Println(spew.Sdump(text))
	// }

	rREsult, err2 := result.Beautify()
	if err2 != nil {
		t.Error(err)
	}

	for _, text := range rREsult {
		fmt.Println(spew.Sdump(text))
	}
	fmt.Println(len(rREsult))
}

func TestConvertToSRT(t *testing.T) {
	f, _ := os.Open("timedtextauto.xml")
	sample, _ := io.ReadAll(f)
	result, err := ParseTimedText(sample)
	if err != nil {
		t.Error(err)
	}

	rResult, err2 := result.Beautify()
	if err2 != nil {
		t.Error(err)
	}

	srtResult, err := ConvertToSRT(rResult)
	if err2 != nil {
		t.Error(err)
	}

	fmt.Println(srtResult)
}
