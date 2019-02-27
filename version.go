//go:generate sh -c "sed -i \"s/const Version = \\\".*\\\"/const Version = \\\"`git describe --always --abbrev=12`\\\"/\" version.go"

package ledgerserver

// Version git version number
const Version = "1f52fd9b4446"
