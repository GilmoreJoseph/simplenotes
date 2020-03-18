# simplenotes

simplenotes is a simple tool for taking and searching notes written in Go.

# Usage

go build

./simplenotes

This will start a webserver running on port 80. There it is a simple api to put and get notes.
to put a note simply send a PUT request with a Note field in the json. To search notes see below.

# Note format

Notes simply have the date in the format Jan 02, 2006, a newline, and tags (any word in square brackets []), and then the note.


# Note searching

To search notes use a GET request with query string parameters. The following are accepted.

start - Inclusive (Format Jan 02, 2006)

end - Inclusive (Formet Jan 02, 2006)

tag - Tag that note has, to search multiple add multiple &tag= to query

phrase - Exact phrase to be searched, can only include one.

# ios shortcut

https://www.icloud.com/shortcuts/5b401efd2a0e4afea092e792bd14f913

Replace the URL, the first action in the shorcut with the url of your server and make sure you include the note file you want at the end.

The shortcut supports writing and getting notes. Getting notes only works with tags, and dates. No phrases yet.

# future plans

-authorization

-automatically adding tags to top of file and a query to get a list of tags

-mutex file read/write
