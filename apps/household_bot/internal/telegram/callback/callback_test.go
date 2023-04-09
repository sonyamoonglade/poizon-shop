package callback

import (
	"strconv"
	"testing"

	f "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func TestInjectAndParseData(t *testing.T) {

	t.Run("raw callback", func(t *testing.T) {
		for cb := 0; cb < 10000; cb++ {
			data := Inject(Callback(cb))
			require.Equal(t, rawCallbackPrefix+strconv.Itoa(cb), data)
		}
	})

	t.Run("with data < 64 bytes", func(t *testing.T) {
		for i := 1; i < 3; i++ {
			values := make([]string, i)
			f.Slice(&values)
			data := Inject(Callback(i), values...)
			parsedCb, parsedValues, err := ParseButtonData(data)
			require.NoError(t, err)
			require.Equal(t, Callback(i), parsedCb)
			require.EqualValues(t, values, parsedValues)
		}
	})
}
