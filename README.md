# GrumpDumpInfo Aggregator
## About
This project collects either some or all of the metadata for videos
uploaded by a given YouTube channel and stores that metadata in mongoDB

## Modes
**NOTE**: This section describes incomplete functionality

A given mode can be selected using `aggregator <mode name>` 
the modes are as follows:
* `build`  - downloads all video metadata for all listed channels and exits
* `update` - downloads the latest video metadata for all listed channels and exits
* `server` - hosts a server that will download the latest video metadata when it receives a request of any kind.

`server` mode will perform at most n updates per minute (where n is the number of channels in `config.yaml`)

## Configuration
Below is an example configuration
```
APIKey: "abcdefghijklmnopqrstuvwxyz-1234567890"
ServerPort: ":8080"
TargetPlaylists:
  - "UU9CuvdOVfMPvKCiwdGKL3cQ"
  - "UUAQ0o3l-H3y_n56C3yJ9EHA"
  - "UUXq2nALoSbxLMehAvYTxt_A"
```
APIKey is the Google developer console key used to access the YouTube API

ServerPort is the desired port to host on when running in `server` mode

TargetPlaylists is the list of playlists/channels to download metadata from