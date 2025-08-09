# Bubly CLI

A CLI tool to download video, audio, and transcriptions from YouTube.

<img width="800" src="./preview.gif" />

## Features

- Download YouTube videos
- Download audio only from YouTube videos
- Download video subtitles
- Format selection for audio and video downloads
- Language selection for subtitles
- Pagination for long lists
- Detailed logging to output.log for debugging
- Automatic installation of yt-dlp and ffmpeg
- Clean terminal interface with auto-clear on startup

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/AbdelilahOu/Bubly-cli-app.git
   ```
2. Navigate to the project directory:
   ```bash
   cd Bubly-cli-app
   ```
3. Run the application:
   ```bash
   go run main.go
   ```
   
   Or use the Makefile:
   ```bash
   make run
   ```

## Troubleshooting

If you encounter any issues, check the `output.log` file for detailed error information from yt-dlp.

The application will automatically prompt to install yt-dlp and ffmpeg if they are not found. You can also manually install them:

- yt-dlp: https://github.com/yt-dlp/yt-dlp
- ffmpeg: https://ffmpeg.org/download.html
