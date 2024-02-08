package job

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestParseDate(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		str := "20230809 17:17:47.013"
		layout := "20060102 15:04:05.000" // Specify the layout matching the string format
		v, err := time.Parse(layout, str)
		assert.Equal(t, nil, err)
		fmt.Print(v)
	})

	t.Run("Success", func(t *testing.T) {
		str := "20230809 17:17:47.013"
		layout := "20060102 15:04:05.000" // Specify the layout matching the string format
		re := regexp.MustCompile(`^(.*?:.*?:.*?):(.*?)$`)

		// Replace the second colon with a period
		output := re.ReplaceAllString(str, "$1.$2")
		fmt.Println("output", output)
		v, err := time.Parse(layout, output)
		assert.Equal(t, nil, err)
		fmt.Print(v)
	})

}
