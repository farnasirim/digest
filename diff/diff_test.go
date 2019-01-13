package diff

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/farnasirim/digest"
)

func TestDiffInterface(t *testing.T) {
	var service digest.DiffService = NewPlainTextDiff()
	assert.Equal(t, service.DiffDirsHtml(".", "."), "")
}

func TestDiffDirSimpe(t *testing.T) {
	diffService := NewPlainTextDiff()

	diffService.DiffDirsHtml(
		"/home/colonelmo/.digest/data/2019-01-13T02:11:27-05:00",
		"/home/colonelmo/.digest/data/2019-01-13T02:14:02-05:00")
}

func TestDiffFileSimple(t *testing.T) {
	oldContent := `hi
hello
wow

to be deleted

some more lines
`
	newContent := `hi
he
wow
insert


some more liness


a totally new section
1 2 3
going on and on
`
	expectedDiff := `<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size=4px color="darkred">
 hello
 <br>
</font>
</div>
<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size="4px" color="darkgreen">
 he
 <br>
</font>
</div>
<br>
<br>

<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size="4px" color="darkgreen">
 insert
 <br> 
</font>
</div>
<br><br>

<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size=4px color="darkred">
 to be deleted
 <br>
</font>
</div>
<br><br>

<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size=4px color="darkred">
 some more lines
 <br>
</font>
</div>
<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size="4px" color="darkgreen">
 some more liness
 <br> 
 <br> 
 <br> a totally new section
 <br> 1 2 3
 <br> going on and on
 <br>
</font>
</div>
<br>
<br>
`

	dirName, err := ioutil.TempDir("", "some_prefix")
	assert.Nil(t, err)
	defer os.RemoveAll(dirName)

	oldFileAddr := path.Join(dirName, "a")
	f, err := os.OpenFile(oldFileAddr, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	assert.Nil(t, err)
	f.WriteString(oldContent)
	defer f.Close()

	newFileAddr := path.Join(dirName, "new-a")
	f, err = os.OpenFile(newFileAddr, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	assert.Nil(t, err)
	f.WriteString(newContent)
	defer f.Close()

	diffService := NewPlainTextDiff()
	assert.Equal(t, strings.TrimSpace(expectedDiff),
		strings.TrimSpace(diffService.DiffFilesHtml(oldFileAddr, newFileAddr)),
	)
}
