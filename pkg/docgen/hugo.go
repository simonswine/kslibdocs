package docgen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gobuffalo/plush"
	"github.com/ksonnet/kslibdocs/pkg/docgen/template"

	"github.com/google/go-jsonnet/ast"
	"github.com/ksonnet/ksonnet-lib/ksonnet-gen/printer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	reFirstSentence = regexp.MustCompile(`^.*?[\.](?:\s|$)`)
)

type hugo struct {
	root string
}

func newHugo(root string) (*hugo, error) {

	h := &hugo{
		root: root,
	}

	if err := h.cleanContents(); err != nil {
		return nil, err
	}

	return h, nil
}

func (h *hugo) makePath(path ...string) string {
	return filepath.Join(append([]string{h.root}, path...)...)
}

func (h *hugo) mkdir(path ...string) error {
	dirName := h.makePath(path...)
	if err := os.MkdirAll(dirName, 0750); err != nil {
		return errors.Wrapf(err, "create directory %s", dirName)
	}

	return nil
}

func (h *hugo) writeProperty(prepend, group, version, kind string, property []string, fm *hugoProperty) error {
	category := []string{group, version, kind}
	for i := range property {
		if i == len(property)-1 {
			break
		}

		category = append(category, property[i])
	}

	id := property[len(property)-1]
	return h.writeDoc(prepend, category, id, fm)
}

func (h *hugo) writeGroup(prepend, group string, fm *hugoGroup) error {
	return h.writeDoc(prepend, []string{"groups"}, group, fm)
}

func (h *hugo) writeKind(prepend, group, kind string, fm *hugoKind) error {
	return h.writeDoc(prepend, []string{group}, kind, fm)
}

func (h *hugo) writeVersionedKind(prepend, group, version, kind, comment string) error {
	category := []string{group, version}
	fm := newHugoVersionedKind(group, version, kind, comment)
	return h.writeDoc(prepend, category, kind, fm)
}

func (h *hugo) writeDoc(prepend string, category []string, name string, fm frontMatterer) error {
	logrus.WithFields(logrus.Fields{
		"category": strings.Join(category, "/"),
		"name":     name,
	}).Debug("writing doc")

	if prepend != "" {
		category = append([]string{prepend}, category...)
	}

	parentPath := append([]string{"content"}, category...)
	if err := h.mkdir(parentPath...); err != nil {
		return errors.Wrapf(err, "create %s dir", category)
	}

	var buf docbuf

	b, err := json.MarshalIndent(fm.FrontMatter(), "", "  ")
	if err != nil {
		return err
	}

	if _, err = buf.Write(b); err != nil {
		return err
	}

	buf.WriteString("\n")

	content, err := fm.Content()
	if err != nil {
		return err
	}

	buf.WriteString(content)

	path := h.makePath(append(parentPath, fm.Filename())...)

	if err = buf.Err(); err != nil {
		return err
	}

	return ioutil.WriteFile(path, buf.Bytes(), 0644)
}

func (h *hugo) cleanContents() error {
	path := h.makePath("content")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return errors.Wrap(err, "check content path")
	}

	return os.RemoveAll(path)
}

type propertyFrontMatter struct {
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Draft    bool      `json:"draft"`
	Weight   int       `json:"weight"`
	K8SKind  string    `json:"k8s_kind"`
	Current  string    `json:"current"`
	DataType string    `json:"data_type"`
}

type hugoProperty struct {
	group        string
	version      string
	kind         string
	comment      string
	property     []string
	weight       int
	function     *ast.Function
	propertyType propertyType
	fieldType    fieldType
	typeDef      string
}

var _ frontMatterer = (*hugoProperty)(nil)

func newHugoProperty(group, version, kind, comment string, property []string, pt propertyType) *hugoProperty {
	return &hugoProperty{
		group:        group,
		version:      version,
		kind:         kind,
		comment:      comment,
		property:     property,
		weight:       100,
		propertyType: pt,
	}
}

func (hp *hugoProperty) FrontMatter() interface{} {
	kind := []string{hp.group, hp.version, hp.kind}
	cur := append(kind, hp.property...)
	for i := range hp.property {
		if i == len(hp.property)-1 {
			break
		}

		kind = append(kind, hp.property[i])
	}
	return &propertyFrontMatter{
		Title:    hp.name(),
		Date:     time.Now().UTC(),
		Draft:    false,
		Weight:   hp.weight,
		K8SKind:  strings.Join(kind, "."),
		Current:  strings.Join(cur, "."),
		DataType: hp.propertyType.String(),
	}
}

func (hp *hugoProperty) name() string {
	return hp.property[len(hp.property)-1]
}

func (hp *hugoProperty) Filename() string {
	return hp.name() + ".md"
}

func (hp *hugoProperty) Content() (string, error) {
	var buf docbuf

	content := fmt.Sprintf("%s/%s/%s - %s",
		hp.group, hp.version, hp.kind, strings.Join(hp.property, "."))

	if hp.function != nil {
		buf.WriteString("<div>")
		buf.WriteString(`{{< highlight js >}}`)
		buf.WriteString(hp.name())
		if err := printer.Fprint(&buf, hp.function); err != nil {
			logrus.WithError(err).WithField("id", content).Fatal("print property")
		}
		buf.WriteString(`{{< /highlight >}}`)
		buf.WriteString("</div>")
	}

	sections := strings.Split(hp.comment, "\n")
	for _, section := range sections {
		buf.WriteString("<p>")
		buf.WriteString(section)
		buf.WriteString("</p>")
	}

	if hp.fieldType == ftArray {
		at := argumentType{
			propertyName: hp.name(),
			typeName:     typeName(hp.name()),
			typeDef:      hp.typeDef,
		}

		out, err := at.ToDoc()
		if err != nil {
			return "", errors.Wrap(err, "build array template")
		}

		buf.WriteString(out)
	}

	if err := buf.Err(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type versionedKindFrontMatter struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Draft   bool      `json:"draft"`
	Layout  string    `json:"layout"`
	Type    string    `json:"type"`
	Matcher string    `json:"matcher"`
}

type hugoVersionedKind struct {
	group   string
	version string
	kind    string
	comment string
}

var _ frontMatterer = (*hugoVersionedKind)(nil)

func newHugoVersionedKind(group, version, kind, comment string) *hugoVersionedKind {
	return &hugoVersionedKind{
		group:   group,
		version: version,
		kind:    kind,
		comment: comment,
	}
}

func (hvk *hugoVersionedKind) FrontMatter() interface{} {
	return &versionedKindFrontMatter{
		Title:   hvk.kind,
		Date:    time.Now().UTC(),
		Draft:   false,
		Layout:  "kind",
		Type:    "kind",
		Matcher: fmt.Sprintf("%s.%s.%s", hvk.group, hvk.version, hvk.kind),
	}
}

func (hvk *hugoVersionedKind) Filename() string {
	return hvk.kind + ".md"
}

func (hvk *hugoVersionedKind) Content() (string, error) {
	return hvk.comment, nil
}

func newGroupFrontMatter(name string) *hugoGroup {
	return &hugoGroup{
		Title: name,
		Name:  name,
	}
}

type groupFrontMatter struct {
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	Draft     bool      `json:"draft"`
	GroupName string    `json:"group_name"`
}

type hugoGroup struct {
	Title string
	Name  string
}

var _ frontMatterer = (*hugoGroup)(nil)

func (hg *hugoGroup) FrontMatter() interface{} {
	return &groupFrontMatter{
		Title:     hg.Title,
		Date:      time.Now().UTC(),
		Draft:     false,
		GroupName: hg.Name,
	}
}

func (hg *hugoGroup) Filename() string {
	return hg.Name + ".md"
}

func (hg *hugoGroup) Content() (string, error) {
	return fmt.Sprintf("Group %s", hg.Name), nil
}

type kindFrontMatter struct {
	Title       string    `json:"title"`
	Date        time.Time `json:"date"`
	Draft       bool      `json:"draft"`
	KindName    string    `json:"kind_name"`
	Versions    []string  `json:"versions"`
	ParentGroup string    `json:"parent_group"`
}

type hugoKind struct {
	Title        string
	group        string
	Name         string
	Versions     []string
	descriptions map[string]string
}

var _ frontMatterer = (*hugoKind)(nil)

func newKindFrontMatter(group, name string, versions []string, descriptions map[string]string) *hugoKind {
	return &hugoKind{
		Title:        name,
		group:        group,
		Name:         name,
		Versions:     versions,
		descriptions: descriptions,
	}
}

func (hk *hugoKind) FrontMatter() interface{} {
	return &kindFrontMatter{
		Title:       hk.Title,
		Date:        time.Now().UTC(),
		Draft:       false,
		KindName:    hk.Name,
		Versions:    hk.Versions,
		ParentGroup: hk.group,
	}
}

func (hk *hugoKind) Content() (string, error) {
	return kindSummary(hk.group, hk.Name, hk.Versions, hk.descriptions)
}

func kindSummary(group, kind string, versions []string, descriptions map[string]string) (string, error) {
	var out []map[string]interface{}

	for i, v := range versions {
		classes := []string{"kind-summary"}
		if i == 0 {
			classes = append(classes, "default")
		}
		m := map[string]interface{}{
			"classes": classes,
			"version": fmt.Sprintf("%s-%s-%s", group, v, kind),
			"summary": summarize(descriptions[v]),
		}
		out = append(out, m)
	}

	ctx := plush.NewContext()
	ctx.Set("classFor", template.ClassFor)
	ctx.Set("versions", out)

	return plush.Render(tmplKindSummary, ctx)

}

func (hk *hugoKind) Filename() string {
	return hk.Name + ".md"
}

type frontMatterer interface {
	FrontMatter() interface{}
	Filename() string
	Content() (string, error)
}

func summarize(s string) string {
	out := reFirstSentence.FindAllString(strings.TrimSpace(s), 1)

	if len(out) == 0 {
		return ""
	}

	return out[0]
}

var tmplKindSummary = `<%= for (version) in versions { %>
<div id="<%= version["version"] %>" class="<%= classFor(version["classes"]) %>">
	<%= version["summary"] %>
</div>
<% } %>
`
