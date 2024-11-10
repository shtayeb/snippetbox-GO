package ui

import "embed"

//go:embed "html" "static"
var Files embed.FS

// The important line here is //go:embed "html" "static".
// This looks like a comment, but it is actually a special comment
// directive. When our application is compiled, this comment directive
// instructs Go to store the files from our ui/html and ui/static
// folders in an embed.FS embedded filesystem referenced by the global variable Files

// go:embed "static" "html" embeds the directories ui/static
// and ui/html from our project.
