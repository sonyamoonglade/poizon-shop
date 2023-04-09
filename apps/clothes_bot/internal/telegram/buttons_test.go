package telegram

import (
	"math"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInjectAndParseCallback(t *testing.T) {
	t.Run("inject message ids and parse", func(t *testing.T) {
		t.Parallel()
		type testcase struct {
			msgIDs   []int
			callback int
		}
		var testcases []testcase
		for i := int(-math.MaxInt16); i < math.MaxInt16; i++ {
			testcases = append(testcases, testcase{msgIDs: []int{i, i + 1, i + 2}, callback: int(math.MaxInt16 - i)})
		}

		for _, tc := range testcases {
			data := injectMessageIDs(tc.callback, tc.msgIDs...)
			require.NotZero(t, data)
			gotMsgIDs, callback, err := parseCallbackData(data)
			intMsgIDs, ok := gotMsgIDs.([]int)
			require.True(t, ok)
			require.NoError(t, err)
			require.EqualValues(t, tc.msgIDs, intMsgIDs)
			require.Equal(t, tc.callback, callback)
		}
	})

	t.Run("inject string data and parse", func(t *testing.T) {
		t.Parallel()
		type testcase struct {
			payload  string
			callback int
		}
		var testcases []testcase
		for i := int(math.MaxInt16); i < math.MaxInt16; i++ {
			testcases = append(testcases, testcase{payload: strconv.Itoa(i), callback: int(math.MaxInt16 - i)})
		}

		for _, tc := range testcases {
			data := injectStringData(tc.callback, tc.payload)
			require.NotZero(t, data)
			out, callback, err := parseCallbackData(data)
			require.NoError(t, err)
			require.Equal(t, tc.callback, callback)
			strOut, ok := out.(string)
			require.True(t, ok)
			require.Equal(t, tc.payload, strOut)
		}
	})
}
