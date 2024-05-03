# TechGo

![image](https://github.com/tamnguyen820/techgo/assets/66036226/603f29f8-6f10-43ef-b443-2df719e9fa82)

**TechGo** is a terminal application that aggregates articles from RSS feeds focused on tech news. The project is written in Go using the wonderful [Bubble Tea framework](https://github.com/charmbracelet/bubbletea).

By default, news sources include:

- The Verge
- Wired
- TechCrunch
- Mashable
- Ars Technica
- TechRadar
- More?

## Quick start

### Binaries

_Coming Soon<sup>TM</sup>_

### Docker

```bash
# 1. Build image and run container
docker build -t techgo .
docker run -it --rm techgo

# 2. Or pull the image
...
```

### Running from source code

```bash
git clone https://github.com/tamnguyen820/techgo
cd techgo
go run cmd/techgo/main.go
```

## Configuration

The file [config.yml](config.yml) stores the sources of RSS feeds. Change the config file as needed.

```bash
rss_feeds:
  - url: https://www.theverge.com/tech/rss/index.xml
    name: The Verge
  - url: https://www.wired.com/feed/rss
    name: Wired
  ...
```

In order to include a new RSS feed:

1. Look up `<news_source> RSS feed`
2. Copy the link to the feed.
3. Add to [config.yml](config.yml) the URL and feed bane.

Alternatively, you can point to a different config file using the `-config` flag, for example:

```bash
go run cmd/techgo/main.go -config myconfig.yml
```
