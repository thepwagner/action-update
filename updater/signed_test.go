package updater_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thepwagner/action-update/updater"
)

var (
	testKey = []byte{1, 2, 3, 4}
)

func TestNewSignedUpdateGroup(t *testing.T) {
	var (
		update1 = updater.Update{Path: "github.com/foo/bar", Previous: "v1.0.0", Next: "v1.1.0"}
		update2 = updater.Update{Path: "github.com/foo/baz", Previous: "v1.0.0", Next: "v2.0.0"}
	)

	// XXX: hardcoded signatures are intentionally fragile - try not to break users
	cases := []struct {
		signature string
		group     string
		updates   []updater.Update
	}{
		{
			signature: "93aZ4eM9U9p3iPKQXcE8YMwfeASEN67G+Bu9Tu0ynC5807cJ7Lldr8BTdRsXEAsqQhqpVFlC1cX590Yca+RplQ==",
			updates:   []updater.Update{update1},
		},
		{
			signature: "pzpXqoxExE6Ot4ZLRWwvtxU1roCuA0njRPY3fiw0QSzDoESeHXLjVfxAsSqm1oZDc0vOSy7jBfYxj+aWOxYMbA==",
			updates:   []updater.Update{update1, update2},
		},
		{
			signature: "pzpXqoxExE6Ot4ZLRWwvtxU1roCuA0njRPY3fiw0QSzDoESeHXLjVfxAsSqm1oZDc0vOSy7jBfYxj+aWOxYMbA==",
			updates:   []updater.Update{update2, update1},
		},
		{
			signature: "9GyHBP/SKi3jDaGiY9Z6z9F5Z2S9mhASrSuf1I4sy0pTWNiHnUTc+ogeNBNvOQZTrYDiQ78hSy7BWbr2ze55nQ==",
			group:     "test",
			updates:   []updater.Update{update1},
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc.updates), func(t *testing.T) {
			ug := updater.NewUpdateGroup(tc.group, tc.updates...)
			signed, err := updater.NewSignedUpdateGroup(testKey, ug)
			require.NoError(t, err)

			buf, err := json.Marshal(&signed)
			require.NoError(t, err)
			t.Log(string(buf))

			assert.Equal(t, tc.group, signed.Updates.Name)
			assert.Equal(t, tc.updates, signed.Updates.Updates)
			assert.Equal(t, tc.signature, base64.StdEncoding.EncodeToString(signed.Signature))

			verified, err := updater.VerifySignedUpdateGroup(testKey, signed)
			require.NoError(t, err)
			assert.Equal(t, tc.group, verified.Name)
			assert.Equal(t, tc.updates, verified.Updates)
		})
	}
}

func TestVerifySignedUpdateGroup_Invalid(t *testing.T) {
	_, err := updater.VerifySignedUpdateGroup([]byte{}, updater.SignedUpdateGroup{})
	assert.EqualError(t, err, "invalid signature")
}
