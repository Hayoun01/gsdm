<div align="center">
	<h1>GSDM (GO Simple Download Manager)</h1>
	<h4 align="center">
		GSDM is a CLI tool for downloading files from URLs. It provides a simple and efficient way to download files.
	</h4>
</div>

<p align="center">
	<a href="#installation">Installation</a> ❘
	<a href="#cli-usage">CLI Usage</a> ❘
    <a href="#how-it-works">How it works?</a> ❘
	<a href="#license">License</a> ❘
	<a href="#contribute">Contribute</a>
</p>


## Installation

To install GSDM, you can use the `go get` command:
```bash
go get github.com/Hayoun01/gsdm
```

## CLI Usage
GSDM provides the following command-line options:

* **-o**: Specify the output filename for the downloaded file.
* **-w**: Number of goroutines for concurrent downloading (Default: 4).
* **-v**: Enable verbose mode to display detailed information during the download process.

Here's an example of how to use GSDM:
```bash
gsdm -o vid.mp4 -w 10 -v https://example.com/video.mp4
```
> **Note:** that the args goes before the link since The [flag](https://pkg.go.dev/flag#:~:text=Flag%20parsing%20stops%20just%20before%20the%20first%20non%2Dflag%20argument%20(%22%2D%22%20is%20a%20non%2Dflag%20argument)%20or%20after%20the%20terminator%20%22%2D%2D%22.) package doesn't adhere to GNU parsing rules.

## How it works?
GSDM follows [RFC 7233](https://datatracker.ietf.org/doc/html/rfc7233), which defines the standard for HTTP range requests. This allows for efficient handling of partial content downloads, enabling features like resuming interrupted downloads and downloading files in chunks.

## License
This project is licensed under the [MIT License](https://github.com/Hayoun01/gsdm/blob/master/LICENSE) ©️ Mohammed Hayyoun.

## Contribute
Contributions are welcome! Feel free to open an issue or submit a pull request to contribute to this project.