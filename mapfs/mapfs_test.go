package mapfs_test

import (
	"io/fs"
	"testing"

	"oss.terrastruct.com/util-go/assert"
	"oss.terrastruct.com/util-go/mapfs"
)

func TestMapFS(t *testing.T) {
	t.Parallel()

	m := map[string]string{
		"index":                       "<Espy_on_crack> I installed 'Linux 6.1', doesn't that make me a unix",
		"d2/imports":                  "Do your part to help preserve life on Earth -- by trying to preserve your own.",
		"d2/globs":                    "I'm going to raise an issue and stick it in your ear.",
		"nested/nested/nested/nested": "Yuppie Wannabes",
	}

	mapfs, err := mapfs.New(m)
	assert.Success(t, err)
	t.Cleanup(func() {
		err := mapfs.Close()
		assert.Success(t, err)
	})

	for p, s := range m {
		b, err := fs.ReadFile(mapfs, p)
		assert.Success(t, err)
		assert.Equal(t, s, string(b))
	}

	_, err = fs.ReadFile(mapfs, "../escape")
	assert.ErrorString(t, err, "stat ../escape: invalid argument")
	_, err = fs.ReadFile(mapfs, "/root")
	assert.ErrorString(t, err, "stat /root: invalid argument")
}
