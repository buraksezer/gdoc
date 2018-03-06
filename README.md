# gdoc
[![Go Report Card](https://goreportcard.com/badge/github.com/buraksezer/gdoc)](https://goreportcard.com/report/github.com/buraksezer/gdoc)

gdoc is a small utility to search on GoDoc.org. It also has some handy functions to read package documents.

**Warning:** This tool is only tested on Linux and OSX. 

### Install
With a properly configured Go environment:

```
go get -u github.com/buraksezer/gdoc/cmd...
```

### Usage
#### Search something:

```
gdoc search <keyword>
```
gdoc lists the first 10 package as default. If you want to increase or decrease that number, use `-c/--count` parameter:

```
gdoc search -c 3 <keyword>
```
gdoc has an interactive mode. If you use the `-i/--interactive` parameter when you search something, gdoc lists packages with numbers and
await for your action to fetch its document from GoDoc.org.

Sample output:
```
gdoc search -i memberlist
==> (1) github.com/hashicorp/memberlist
==> imports: 192 stars: 804
memberlist is a library that manages cluster membership and member failure detection using a gossip based protocol.

==> (2) github.com/Nitro/memberlist
==> imports: 3 stars: 0
memberlist is a library that manages cluster membership and member failure detection using a gossip based protocol.

==> (3) github.com/journeymidnight/nentropy/memberlist
==> imports: 3 stars: 1
....
Give a number to read the document:
```
#### Read package documentation:

```
gdoc read <package path>
```

Sample usage:
```
gdoc read github.com/hashicorp/memberlist
```

Fetches the package document from GoDoc.org in text format and passes it to an available the pager. gdoc looks for `GDOC_PAGER` and `PAGER` environment 
variables to get pager command as respectively. If you want to disable paging, use `--disable-pager` parameter:
```
gdoc read --disable-pager <package path> 
```

gdoc supports aliases to access easily the frequently used documents. If you set an alias for a package previously, just use `-a/--alias` to get the 
document for that package.

```
gdoc read -a <alias>
```

#### Aliases:
In order to add an alias for a package:
```
gdoc alias set <short name> <package path>
```

Sample usage:
```
gdoc alias set memberlist github.com/hashicorp/memberlist
```
If you want to delete previously setted alias, just use `del` subcommand:
```
gdoc alias del <short name>
```
You can use `list` subcommand to list already setted aliases.

### Contribution
gdoc is a free software, you feel free to send PRs to improve gdoc.

### License
gdoc is licensed under the GNU General Public License v3.0 - see LICENSE for more details.

