package chg

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangelogVersion(t *testing.T) {
	unreleased := &Version{Name: "Unreleased"}
	v123 := &Version{Name: "1.2.3"}

	c := NewChangelog()
	c.Versions = append(c.Versions, unreleased)
	c.Versions = append(c.Versions, v123)

	t.Run("version=unreleased", func(t *testing.T) {
		result := c.Version("unreleased")
		assert.Equal(t, unreleased, result)
	})

	t.Run("version=1.2.3", func(t *testing.T) {
		result := c.Version("1.2.3")
		assert.Equal(t, v123, result)
	})

	t.Run("version=unknown", func(t *testing.T) {
		result := c.Version("unknown")
		assert.Nil(t, result)
	})
}

func TestChangelogRelease(t *testing.T) {
	c := Changelog{
		Versions: []*Version{
			{
				Name: "Unreleased",
				Link: "http://example.com/1.0.0..HEAD",
				Changes: []*ChangeList{
					{
						Type: Added,
						Items: []*Item{
							{Description: "New feature"},
						},
					},
				},
			},
			{
				Name: "1.0.0",
				Link: "http://example.com/abcdef..1.0.0",
			},
			{
				Name: "0.2.0",
				Link: "http://example.com/abcdef..0.2.0",
			},
		},
	}

	t.Run("default", func(t *testing.T) {
		newVersion, err := c.Release(Version{Name: "2.0.0"})

		assert.Nil(t, err)
		assert.Equal(t, "2.0.0", newVersion.Name)
		// Make sure the changes were kept
		assert.Equal(t, 1, len(newVersion.Changes))
	})

	t.Run("explicit-compare-url", func(t *testing.T) {
		v := Version{Name: "2.0.0", Link: "https://localhost/<prev>..<next>"}
		newVersion, err := c.Release(v)

		assert.Equal(t, "2.0.0", newVersion.Name)

		unreleased := c.Version("Unreleased")
		assert.Equal(t, "https://localhost/2.0.0..HEAD", unreleased.Link)

		assert.Nil(t, err)
	})
}

func TestChangelogRenderLinks(t *testing.T) {
	unreleased := &Version{Name: "Unreleased", Link: "http://example.com/unreleased"}
	v123 := &Version{Name: "1.2.3", Link: "http://example.com/1.2.3"}
	v456 := &Version{Name: "4.5.6"}

	c := NewChangelog()
	c.Versions = append(c.Versions, unreleased)
	c.Versions = append(c.Versions, v123)
	c.Versions = append(c.Versions, v456)

	expected := "[Unreleased]: http://example.com/unreleased\n[1.2.3]: http://example.com/1.2.3\n"

	var buf bytes.Buffer
	c.RenderLinks(&buf)
	result := string(buf.Bytes())

	assert.Equal(t, expected, result)
}

func TestChangelogRender(t *testing.T) {
	c := Changelog{
		Preamble: `Any paragraph
to be inserted.
`,
	}

	t.Run("empty-versions", func(t *testing.T) {
		expected := `# Changelog

Any paragraph
to be inserted.
`
		var buf bytes.Buffer
		c.Render(&buf)
		result := buf.String()
		assert.Equal(t, expected, result)
	})

	t.Run("with-versions", func(t *testing.T) {
		c.Versions = []*Version{
			{Name: "1.0.0"},
			{Name: "2.0.0"},
		}

		expected := `# Changelog

Any paragraph
to be inserted.

## 1.0.0

## 2.0.0
`
		var buf bytes.Buffer
		c.Render(&buf)
		result := buf.String()
		assert.Equal(t, expected, result)
	})

	t.Run("sort-changes", func(t *testing.T) {
		c.Versions = []*Version{
			{
				Name: "1.0.0",
				Changes: []*ChangeList{
					{
						Type: Fixed,
					},
					{
						Type: Added,
					},
				},
			},
		}

		expected := `# Changelog

Any paragraph
to be inserted.

## 1.0.0
### Added

### Fixed
`

		var buf bytes.Buffer
		c.Render(&buf)
		result := buf.String()
		assert.Equal(t, expected, result)
	})
}
