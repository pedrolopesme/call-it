# Call-It Usage Examples

## 🎨 Default TUI Mode (Interactive)
```bash
# Start the beautiful TUI interface (default)
call-it

# This will open an interactive interface where you can:
# - Enter URL, attempts, and concurrent calls
# - See real-time progress with spinners
# - View beautiful results with colored status codes
```

## 🖥️ Classic CLI Mode
```bash
# Use the original CLI interface
call-it --cli https://httpbin.org/status/200 10 5

# Or with arguments (automatically uses CLI mode)
call-it https://example.com 20 10
```

## 📋 Config File Mode
```bash
# Use a configuration file
call-it -c

# This reads from config.json in the current directory
# See examples/config.json for format
```

## 🎯 Quick Examples

### TUI Mode (Default)
```bash
# Simply run call-it and enjoy the interactive experience
call-it
```

### CLI Mode for Scripts
```bash
# Great for automation and scripts
call-it --cli https://api.example.com/health 5 2
call-it --cli https://httpbin.org/delay/1 3 1
```

### Batch Testing with Config
```bash
# Multiple endpoints in one go
call-it -c
```

## 🚀 Pro Tips

- **TUI Mode**: Perfect for interactive testing and exploration
- **CLI Mode**: Ideal for automation, scripts, and CI/CD pipelines  
- **Config Mode**: Best for complex test scenarios with multiple endpoints
- Use `--help` to see all available options
- Press `Ctrl+C` in TUI mode to exit safely