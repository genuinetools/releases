package main

const (
	releaseTmpl = `Below are easy install instructions by OS and Architecture. As always there are always the standard instructions in the [README.md](README.md) as well.

<< range $os, $v := . >>
## << $os  >>

<< range $arch, $r := $v >>
### << $arch >>

` + "```" + `console
# Export the sha256sum for verification.
$ export << $r.Repository.Name | ToUpper >>_SHA256="<< $r.BinarySHA256 >>"

# Download and check the sha256sum.
$ curl -fSL "<< $r.BinaryURL >>" -o "/usr/local/bin/<< $r.Repository.Name >>" \
	&& echo "` + "${" + `<< $r.Repository.Name | ToUpper >>_SHA256` + "}" + `  /usr/local/bin/<< $r.Repository.Name >>" | sha256sum -c - \
	&& chmod a+x "/usr/local/bin/<< $r.Repository.Name >>"

$ echo "<< $r.Repository.Name >> installed!"

# Run it!
$ << $r.Repository.Name >> -h
` + "```" + `
<<end>>
<<end>>
`
)
