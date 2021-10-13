# Counter

This is a simple web counter written in Go. It uses a SQLite database for its data; the path
is counter.sqlite.

It's accessible on configured port if not cofigured it will default to 9776.

The counter responds to three possible paths:

`/`: The root, will respond with a 200 and a small JSON message.

`/count`: The text counter, with a required `path` query parameter for the path you're interested in. Does not increment on requests. For example: `http://localhost:9776/count?path=example.com/about`

`/count/counter.jpg`: The counter image, with a required `path` query parameter for the path you're interested in. DOES increment on requests. For example: `http://localhost:9776/count/counter.jpg?path=example.com/about`

## Configuration

Use environment variables to configure this.

* `COUNTER_FONT_DIR`: The directory path where the font file lives. Defaults to the same directory as the binary.
* `COUNTER_FONT_FILE`: The filename for the font. Defaults to `Berylium.ttf`.
* `COUNTER_IMAGE_WIDTH`: The width of the counter image. Defaults to `200`.
* `COUNTER_IMAGE_HEIGHT`: The height of the counter image. Defaults to `50`.
* `COUNTER_FONT_SIZE`: The font size for the counter image. Defaults to `32`.
* `COUNTER_BG_COLOR`: The hex color for the background. Defaults to `"#000000"`.
* `COUNTER_FONT_COLOR`: The hex color for the font. Defaults to `"#FFFFFF"`.
* `COUNTER_LOG_LEVEL`: The log level to use. You can set it to `debug`, otherwise it uses `info`.
* `COUNTER_PORT`: The port for the web counter. Defaults to `9776`.

## Running in Docker

If you run this via Docker, then remember to bind a volume for `counter.sqlite`!
