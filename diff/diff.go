package diff

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
)

type PlainTextDiff struct {
	unchangedGroupFormat string
	newLineFormat        string
	oldLineFormat        string
	newGroupFromat       string
	changedGroupFormat   string
	oldGroupFormat       string
}

func (d *PlainTextDiff) DiffFilesHtml(older, newer string) string {
	cmd := exec.Command("diff",
		"--unchanged-group-format", d.unchangedGroupFormat,
		"--new-line-format", d.newLineFormat,
		"--old-line-format", d.oldLineFormat,
		"--new-group-format", d.newGroupFromat,
		"--changed-group-format", d.changedGroupFormat,
		"--old-group-format", d.oldGroupFormat,
		older, newer)

	buffer := make([]byte, 0)
	buf := bytes.NewBuffer(buffer)
	cmd.Stdout = buf
	cmd.Stderr = buf

	_ = cmd.Run()

	diff := buf.String()
	return diff
}

func filesMap(dir string) map[string]string {
	mp := make(map[string]string)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, file := range files {
		mp[file.Name()] = path.Join(dir, file.Name())
	}
	return mp
}

func (d *PlainTextDiff) DiffDirsHtml(older, newer string) string {
	ret := ""

	olderFiles := filesMap(older)
	newerFiles := filesMap(newer)

	for name, newFileName := range newerFiles {
		oldFile := "/dev/null"
		if oldFileName, ok := olderFiles[name]; ok {
			oldFile = oldFileName
			delete(olderFiles, name)
		}
		diff := d.DiffFilesHtml(oldFile, newFileName)
		if diff != "" {
			ret += fmt.Sprintf(`<h1> %s </h1>`, name)
			ret += diff
		}
	}

	for name, oldFileName := range olderFiles {
		diff := d.DiffFilesHtml(oldFileName, "/dev/null")
		if diff != "" {
			ret += fmt.Sprintf(`<h1> %s </h1>`, name)
			ret += diff
		}
	}

	return ret
}

func NewPlainTextDiff() *PlainTextDiff {
	return &PlainTextDiff{
		unchangedGroupFormat: "",
		newLineFormat:        " %L <br>",
		oldLineFormat:        " %L <br>",
		newGroupFromat: `
<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size="4px" color="darkgreen">
%> 
</font>
</div>
<br><br>
`,
		changedGroupFormat: `
<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size=4px color="darkred">
%<
</font>
</div>
<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size="4px" color="darkgreen">
%>
</font>
</div>
<br>
<br>
`,
		oldGroupFormat: `
<div style="word-wrap: break-word; width:700px; font-family: monospace;">
<font size=4px color="darkred">
%<
</font>
</div>
<br><br>
`,
	}
}
