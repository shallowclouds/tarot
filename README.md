# 🔮 Tarot Card Reading

A tarot card reading application powered by AI, supporting both command-line interface and web interface.

## Features

- 🎴 **Random Card Selection**: Draws three tarot cards for divination
- 🤖 **AI-Powered Readings**: Generate interpretations using GPT models (OpenAI, DeepSeek, or fallback)
- 🖼️ **Beautiful Result Images**: Auto-generates visual representations of readings
- 💻 **CLI Tool**: Quick readings from the command line
- 🌐 **Web Interface**: User-friendly browser-based interface with responsive design
- 🔌 **Flexible AI Providers**: Support for OpenAI, DeepSeek, or dummy mode

## Quick Start

### Web Interface

1. Set your API key (optional, defaults to dummy mode):
```bash
export DEEPSEEK_API_KEY="your-api-key-here"
export READER_TYPE="deepseek"  # or "openai" or "dumb"
```

2. Run the web server:
```bash
cd ~/repos/tarot
go run cmd/tarotweb/main.go
```

3. Open http://localhost:8080 in your browser

### CLI Tool

```bash
# Using DeepSeek (recommended)
tarot-divine --reader deepseek --api-key "your-key" --thing "我今天的运势如何？"

# Using environment variables
export DEEPSEEK_API_KEY="your-api-key"
tarot-divine --reader deepseek --thing "我的工作会顺利吗？"

# Using dummy mode (no API key needed)
tarot-divine --thing "今天能写完方案吗？"
```

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/shallowclouds/tarot.git
cd tarot

# Build CLI tool
go build -o tarot-divine cmd/divinetest/main.go

# Build web server
go build -o tarotweb cmd/tarotweb/main.go

# Optional: Install to PATH
cp tarot-divine ~/.local/bin/
cp tarotweb ~/.local/bin/
```

## Configuration

### Environment Variables

**For DeepSeek API:**
- `DEEPSEEK_API_KEY`: Your DeepSeek API key (required for DeepSeek reader)
- `DEEPSEEK_BASE_URL`: API endpoint (default: `https://api.deepseek.com`)
- `DEEPSEEK_MODEL`: Model name (default: `deepseek-chat`)

**For OpenAI API:**
- `OPENAI_API_KEY`: Your OpenAI API key (required for OpenAI reader)

**For Web Server:**
- `READER_TYPE`: AI reader type - `deepseek`, `openai`, or `dumb` (default: `dumb`)
- `PORT`: Web server port (default: `8080`)

### Reader Types

- **`deepseek`**: Uses DeepSeek AI for dynamic, contextual readings
- **`openai`**: Uses OpenAI ChatGPT for readings
- **`dumb`**: Uses preset responses (no API key needed)

## Usage Examples

### Web Interface

The web interface provides a beautiful, responsive UI for tarot readings:

1. Enter your question in the text area
2. Click "开始占卜" (Start Reading)
3. View your drawn cards, AI interpretation, and generated image
4. Click "再次占卜" (Read Again) to start over

Features:
- Purple gradient design with smooth animations
- Loading states during API calls
- Error handling and user feedback
- Mobile-responsive layout
- Base64-encoded images displayed inline

### CLI Tool Options

```bash
# All available options
tarot-divine --help

# Specify reader type
tarot-divine --reader deepseek --thing "question"

# Provide API key inline (not recommended for security)
tarot-divine --reader deepseek --api-key "sk-xxx" --thing "question"

# Custom output file
tarot-divine --thing "question"  # Saves to assets/divine_results.jpg
```

## Example Output

![Tarot Card Divine Results](assets/divine_results.jpg)

Each reading includes:
- Three randomly selected tarot cards with position (Upright/Reversed)
- AI-generated interpretation in Chinese
- Beautiful visual representation with card images

## API Integration

### DeepSeek API (Recommended)

Get your API key from [DeepSeek Platform](https://platform.deepseek.com/)

```bash
export DEEPSEEK_API_KEY="your-api-key-here"
```

### OpenAI API

Get your API key from [OpenAI Platform](https://platform.openai.com/)

```bash
export OPENAI_API_KEY="sk-..."
```

## Development

### Project Structure

```
tarot/
├── cmd/
│   ├── divinetest/     # CLI tool entry point
│   └── tarotweb/       # Web server entry point
├── static/             # Web interface files
│   ├── index.html
│   ├── style.css
│   └── script.js
├── assets/             # Embedded card images and fonts
├── cards.go            # Card definitions and management
├── reader.go           # Core tarot reading logic
├── gpt_reader.go       # AI reader implementations
├── DEPLOYMENT.md       # Detailed deployment guide
└── README.md           # This file
```

### Building

```bash
# Format code
make fmt

# Build CLI
go build -o tarot-divine cmd/divinetest/main.go

# Build web server
go build -o tarotweb cmd/tarotweb/main.go
```

### Running Tests

```bash
go test ./...
```

## Security Notes

- Never commit API keys to the repository
- Use environment variables for sensitive configuration
- The web server has no built-in authentication - deploy behind a reverse proxy if needed
- Consider rate limiting for production deployments

## License

See LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
