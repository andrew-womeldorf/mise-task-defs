# Mise Task Definitions for Go

There's a little go code as well, for demonstration purposes, mostly.

All of the tasks should be runnable in this repository, even though the code
itself doesn't do much. This is primarily here for the purpose of copying task
definitions quickly.

Also, I included using zig for CGO compilation. I had issues if there were no
CGO files, so added some sqlite and mattn/sqlite3 as well so we could static
compile without errors related to there being no CGO code.
