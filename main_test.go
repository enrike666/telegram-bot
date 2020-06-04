package main

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseCheckInfo(t *testing.T) {
	testCases := []struct {
		name             string
		checkQueryString string
		wantCheckInfo    *CheckInfo
		wantErr          bool
	}{
		{
			name:             "positive",
			checkQueryString: "t=20200601T2122&s=213.07&fn=9282440300649733&i=12416&fp=2858316733&n=1",
			wantCheckInfo: &CheckInfo{
				DateTimeString: "20200601T2122",
				Sum:            "213.07",
				FN:             "9282440300649733",
				FD:             "12416",
				FP:             "2858316733",
				N:              "1",
			},
			wantErr: false,
		},
		{
			name:             "negative",
			checkQueryString: "t=20200601T2122&s=213.07&fn=9282440300649733&i=12416&fp=2858316733&$=$",
			wantErr:          true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actualCheckInfo, err := ParseCheckInfo(tt.checkQueryString)
			if tt.wantErr {
				t.Log(err)
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCheckInfo, actualCheckInfo)
			}
		})
	}
}
