package logging

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// return how many seconds time zone difference there is between the current location and
// Pacific Time Zone
func secondsFromPST(t *testing.T) int64 {
	pst, err := time.LoadLocation("America/Los_Angeles")
	assert.NoError(t, err)
	timeHere, err := time.ParseInLocation(time.ANSIC, "Mon Jan 2 15:04:05 2006", time.Now().Location())
	assert.NoError(t, err)
	timePST, err := time.ParseInLocation(time.ANSIC, "Mon Jan  2 15:04:05 2006", pst)
	assert.NoError(t, err)
	return int64(timeHere.Sub(timePST).Seconds())
}

func TestFormatFromString(t *testing.T) {
	assert.Equal(t, FormatFromString("FuLl"), FULL, "formats are case insensitive")
	assert.Equal(t, FormatFromString("SimplE"), SIMPLE, "formats are case insensitive")
	assert.Equal(t, FormatFromString("MinimalTagged"), MINIMALTAGGED, "formats are case insensitive")
	assert.Equal(t, FormatFromString("Minimal"), MINIMAL, "formats are case insensitive")
	assert.Equal(t, FormatFromString("foo"), SIMPLE, "default is simple")
}

func TestFormatGetFormatter(t *testing.T) {
	// Note: we can't do function equality, so we need to hack around it by temporarily
	// swapping out the functions
	fullFormatOriginal := fullFormat
	simpleFormatOriginal := simpleFormat
	minimalWithTagsFormatOriginal := minimalWithTagsFormat
	minimalFormatOriginal := minimalFormat
	defer func() {
		// put everything back
		fullFormat = fullFormatOriginal
		simpleFormat = simpleFormatOriginal
		minimalWithTagsFormat = minimalWithTagsFormatOriginal
		minimalFormat = minimalFormatOriginal
	}()

	fullFormat = func(level LogLevel, tags []string, message string, t time.Time, original time.Time) string {
		return "FULL FORMAT"
	}
	simpleFormat = func(level LogLevel, tags []string, message string, t time.Time, original time.Time) string {
		return "SIMPLE FORMAT"
	}
	minimalWithTagsFormat = func(level LogLevel, tags []string, message string, t time.Time, original time.Time) string {
		return "MINIMAL WITH TAGS FORMAT"
	}
	minimalFormat = func(level LogLevel, tags []string, message string, t time.Time, original time.Time) string {
		return "MINIMAL FORMAT"
	}

	assert.Equal(t, "FULL FORMAT", GetFormatter(FULL)(WARN, nil, "", time.Now(), time.Now()), "should be full")
	assert.Equal(t, "SIMPLE FORMAT", GetFormatter(SIMPLE)(WARN, nil, "", time.Now(), time.Now()), "should be simple")
	assert.Equal(t, "MINIMAL WITH TAGS FORMAT", GetFormatter(MINIMALTAGGED)(WARN, nil, "", time.Now(), time.Now()), "should be minimal tagged")
	assert.Equal(t, "MINIMAL FORMAT", GetFormatter(MINIMAL)(WARN, nil, "", time.Now(), time.Now()), "should be minimal")
	assert.Equal(t, "SIMPLE FORMAT", GetFormatter(LogFormat("foo"))(WARN, nil, "", time.Now(), time.Now()), "should be simple")
}

func TestFormatFull(t *testing.T) {
	// the test was written for PST, so we need to apply the offset
	timeZoneOffset := secondsFromPST(t)

	at := time.Unix(1000+timeZoneOffset, 0)
	original := at.AddDate(0, 0, 1)

	expected := "[Dec 31 16:16:40.000] [INFO] [one two] [replayed from Jan  1 16:16:40.000] hello"
	assert.Equal(t, fullFormat(INFO, []string{"one", "two"}, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))

	expected = "[Dec 31 16:16:40.000] [INFO] [replayed from Jan  1 16:16:40.000] hello"
	assert.Equal(t, fullFormat(INFO, nil, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))

	expected = "[Dec 31 16:16:40.000] [INFO] [one two] hello"
	assert.Equal(t, fullFormat(INFO, []string{"one", "two"}, "hello", at, at), expected, fmt.Sprintf("should equal %s", expected))

	expected = "[Dec 31 16:16:40.000] [INFO] [one two] [replayed from Jan  1 16:16:40.000] hello"
	assert.Equal(t, fullFormat(INFO, []string{"one", "two"}, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))
}

func TestFormatSimple(t *testing.T) {
	// the test was written for PST, so we need to apply the offset
	timeZoneOffset := secondsFromPST(t)

	at := time.Unix(1000+timeZoneOffset, 0)
	original := at.AddDate(0, 0, 1)

	expected := "[Dec 31 16:16:40] [INFO] hello"
	assert.Equal(t, simpleFormat(INFO, []string{"one", "two"}, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))
	assert.Equal(t, simpleFormat(INFO, nil, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))
	assert.Equal(t, simpleFormat(INFO, []string{"one", "two"}, "hello", at, at), expected, fmt.Sprintf("should equal %s", expected))
	assert.Equal(t, simpleFormat(INFO, []string{"one", "two"}, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))
}

func TestFormatMinimal(t *testing.T) {

	at := time.Unix(1000, 0)
	original := at.AddDate(0, 0, 1)

	expected := "hello"
	assert.Equal(t, minimalFormat(INFO, []string{"one", "two"}, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))
	assert.Equal(t, minimalFormat(INFO, nil, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))
	assert.Equal(t, minimalFormat(INFO, []string{"one", "two"}, "hello", at, at), expected, fmt.Sprintf("should equal %s", expected))
	assert.Equal(t, minimalFormat(INFO, []string{"one", "two"}, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))
}

func TestFormatMinimalWithTags(t *testing.T) {

	at := time.Unix(1000, 0)
	original := at.AddDate(0, 0, 1)

	expected := "[INFO] [one two] hello"
	assert.Equal(t, minimalWithTagsFormat(INFO, []string{"one", "two"}, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))

	expected = "[INFO] hello"
	assert.Equal(t, minimalWithTagsFormat(INFO, nil, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))

	expected = "[INFO] [one two] hello"
	assert.Equal(t, minimalWithTagsFormat(INFO, []string{"one", "two"}, "hello", at, at), expected, fmt.Sprintf("should equal %s", expected))

	expected = "[INFO] [one two] hello"
	assert.Equal(t, minimalWithTagsFormat(INFO, []string{"one", "two"}, "hello", at, original), expected, fmt.Sprintf("should equal %s", expected))
}
