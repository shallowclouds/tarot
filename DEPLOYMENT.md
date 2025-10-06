# Tarot Card Reading - Deployment Guide

## Installation

The tarot card reading application has been deployed to your local environment.

### Binary Location
- Main binary: `~/.local/bin/tarot-divine`
- Wrapper script: `~/.local/bin/tarot`

### Usage

#### Using the full binary:
```bash
tarot-divine --thing "你的问题"
```

#### Using the convenient wrapper (shorter command):
```bash
tarot "你的问题"
```

If you don't provide a question, it will use the default: "我今天的运势如何？"

### DeepSeek API Configuration

如果想使用 DeepSeek API 生成真实的占卜解读（而不是预设回复），需要：

1. 获取 DeepSeek API 密钥：访问 https://platform.deepseek.com/
2. 设置环境变量：
```bash
export DEEPSEEK_API_KEY="your-api-key-here"
export DEEPSEEK_BASE_URL="https://api.deepseek.com"  # 可选，默认值
export DEEPSEEK_MODEL="deepseek-chat"  # 可选，默认值
```

3. 运行应用：
```bash
tarot-divine --reader deepseek --thing "你的问题"
```

或使用命令行参数：
```bash
tarot-divine --reader deepseek --api-key "your-key" --thing "你的问题"
```

**可用的 reader 类型：**
- `dumb`：使用预设回复（默认，无需 API key）
- `openai`：使用 OpenAI ChatGPT
- `deepseek`：使用 DeepSeek API

### Examples

```bash
# Ask about today's fortune (using default DumbGPTReader)
tarot "我今天的运势如何？"

# Ask about work (using DeepSeek)
tarot-divine --reader deepseek --thing "我的工作会顺利吗？"

# Ask about relationships (using DeepSeek with API key from environment)
export DEEPSEEK_API_KEY="sk-xxxxx"
tarot-divine --reader deepseek --thing "我的感情运势怎么样？"

# Use OpenAI instead
export OPENAI_API_KEY="sk-xxxxx"
tarot-divine --reader openai --thing "我这周的财运如何？"
```

### Output

The application will:
1. Display the selected cards in the terminal
2. Generate a tarot reading interpretation (Chinese text)
3. Save a beautiful image to `~/repos/tarot/assets/divine_results.jpg`

### Technical Details

- **Language**: Go 1.18+
- **Dependencies**: All embedded in binary
- **API Keys**: Not required (uses DumbGPTReader with preset responses)
- **Assets**: Embedded in binary via Go embed
- **Output Format**: JPG image with Chinese text

### Supported AI Models

The application now supports multiple AI readers:

1. **DumbGPTReader (Default)**: Uses preset responses, no API key needed
2. **DeepSeek**: Uses DeepSeek API for dynamic, contextual readings
3. **OpenAI ChatGPT**: Uses OpenAI's GPT models

To switch between readers, use the `--reader` flag or set the `DEEPSEEK_API_KEY` or `OPENAI_API_KEY` environment variables.

### Troubleshooting

**Command not found**:
If you get "command not found", make sure `~/.local/bin` is in your PATH:
```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

**Permission denied**:
Make sure the scripts are executable:
```bash
chmod +x ~/.local/bin/tarot-divine ~/.local/bin/tarot
```

## Repository Structure

- `cmd/divinetest/main.go` - Main application entry point
- `reader.go` - Core tarot reading logic
- `gpt_reader.go` - GPT integration (with dummy fallback)
- `cards.go` - Card definitions and asset management
- `assets/` - Embedded card images, fonts, and data
- `Makefile` - Build and lint commands

## Build from Source

```bash
cd ~/repos/tarot
go build -o tarot-divine cmd/divinetest/main.go
```
