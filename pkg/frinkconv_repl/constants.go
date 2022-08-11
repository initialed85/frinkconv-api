package frinkconv_repl

const (
	delimiter = "\n\n\n\n"         // spam newlines after each command so we've got something consistent to read until
	pattern   = `(\d+\.\d+)|(\d+)` // always grab the last thing that looks like a number from the result
)
