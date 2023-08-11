# TFD (Telegram File Downloader) 
============================================ 
## Overview 
TFD is a tool for downloading files from Telegram. 
It will allow you to download content such as documents, videos, audios, and links from the Telegram.

## Usage 
To use the TFD tool, you will need to build it:
```bash
go build -o tfd
```

then run it:
```bash
./tfd -d
```

This utility will provide you with the necessary commands and configuration settings to download files and content from Telegram.

## Configuration 
By default, the main configuration file is stored in the user's home directory as `.tfd.yaml`. 
The file is used to store essential access and usage information for the TFD tool. It can also be configured by passing the `--config` flag to the `tfd` utility.

## License
TFD is open source and released under the MIT License. 
For more information, please read the `LICENSE` file located in the root of the project.

