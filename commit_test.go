package gitobj

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCommitReturnsCorrectObjectType(t *testing.T) {
	assert.Equal(t, CommitObjectType, new(Commit).Type())
}

func TestCommitEncoding(t *testing.T) {
	author := &Signature{Name: "John Doe", Email: "john@example.com", When: time.Now()}
	committer := &Signature{Name: "Jane Doe", Email: "jane@example.com", When: time.Now()}

	c := &Commit{
		Author:    author,
		Committer: committer,
		ParentIds: [][]byte{
			[]byte("aaaaaaaaaaaaaaaaaaaa"), []byte("bbbbbbbbbbbbbbbbbbbb"),
		},
		TreeId:  []byte("cccccccccccccccccccc"),
		Message: "initial commit",
	}

	buf := new(bytes.Buffer)

	_, err := c.Encode(buf)
	assert.Nil(t, err)

	assertLine(t, buf, "tree 6363636363636363636363636363636363636363")
	assertLine(t, buf, "parent 6161616161616161616161616161616161616161")
	assertLine(t, buf, "parent 6262626262626262626262626262626262626262")
	assertLine(t, buf, "author %s", author)
	assertLine(t, buf, "committer %s", committer)
	assertLine(t, buf, "")
	assertLine(t, buf, "initial commit")

	assert.Equal(t, 0, buf.Len())
}

func TestCommitDecoding(t *testing.T) {
	when := time.Unix(1494258422, 0)
	author := &Signature{Name: "John Doe", Email: "john@example.com", When: when}
	committer := &Signature{Name: "Jane Doe", Email: "jane@example.com", When: when}

	p1 := []byte("aaaaaaaaaaaaaaaaaaaa")
	p2 := []byte("bbbbbbbbbbbbbbbbbbbb")
	treeId := []byte("cccccccccccccccccccc")

	from := new(bytes.Buffer)
	fmt.Fprintf(from, "author %s\n", author)
	fmt.Fprintf(from, "committer %s\n", committer)
	fmt.Fprintf(from, "parent %s\n", hex.EncodeToString(p1))
	fmt.Fprintf(from, "parent %s\n", hex.EncodeToString(p2))
	fmt.Fprintf(from, "tree %s\n", hex.EncodeToString(treeId))
	fmt.Fprintf(from, "initial commit\n")

	flen := from.Len()

	commit := new(Commit)
	n, err := commit.Decode(from, int64(flen))

	assert.Nil(t, err)
	assert.Equal(t, flen, n)

	assert.Equal(t, author.String(), commit.Author.String())
	assert.Equal(t, committer.String(), commit.Committer.String())
	assert.Equal(t, [][]byte{p1, p2}, commit.ParentIds)
	assert.Equal(t, "initial commit", commit.Message)
}

func assertLine(t *testing.T, buf *bytes.Buffer, wanted string, args ...interface{}) {
	got, err := buf.ReadString('\n')

	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf(wanted, args...), strings.TrimSuffix(got, "\n"))
}
