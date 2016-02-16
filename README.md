# `page2images_agent`

page2images.com is a web page screenshot service that converts a URL into a screenshot.

`page2images_agent` is a web server which hits page2images.com API according to the request URL and redirects when the image is ready.

## Usage

```
page2images_agent -api-key=xxxxxxxxxxxxxxxx
```

```
% page2images_agent -h
Usage of ./page2images_agent:
  -api-key string
    	API key for page2images.com
  -port string
    	http port to listen (default "8080")
  -url-prefix string
    	only URLs starting with this prefix are permitted (leave blank to permit any URLs)
```

## Parameters

All parameters supported by page2images.com. See http://www.page2images.com/code_wizard for more details.
The only exception is the `p2i_key`, which must be set as the command line argument.
