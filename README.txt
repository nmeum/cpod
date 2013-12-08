cpod
====

Yet another cron friendly podcatcher.

Introduction
------------

cpod is a small cron friendly podcatcher written in Go. It uses a tiny
json file to store your feeds but it doesn't track downloaded episodes.
When your are done with a podcast episode you can delete it and cpod
won't fetch it again.

Installation
------------

Install it using `go get`:

    $ go get github.com/nmeum/cpod

Usage
-----

If you don't pass any command line flags to cpod it will automatically
update all feeds and download all new episodes. The following command
line flags can be used to change this behaviour:

`-h`

    Display help and exit.

`-v`

    Display version and exit.

`-r <n>`

    Only download <n> most recent episodes.

`-c`

    Remove all episodes except the lastest ones.

`-u`

    Don't update feeds and don't download new episodes.

`-d`

    Don't download new episodes.

`-i <path>`

    Import feeds from opml file at <path>.

`-e <path>`

    Export all feeds as opml to <path>.

Examples
--------

Update all feeds and download new episodes:

    $ cpod

Update all feeds and download the latest episode:

    $ cpod -r 1

Update all feed without downloading episodes:

    $ cpod -d

Remove all episodes except the latest ones:

    $ cpod -u -c

Import a new opml file but don't download new episodes:

    $ cpod -d -i podcasts.opml

License
-------

This program is free software: you can redistribute it and/or modify it
under the terms of the GNU General Public License as published by the
Free Software Foundation, either version 3 of the License, or (at your
option) any later version.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
Public License for more details.

You should have received a copy of the GNU General Public License along
with this program. If not, see <http://www.gnu.org/licenses/>.
