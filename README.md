# RSS Feed Server

## Why?

I want to transform every content into atom feed without rebuilt the rss server.
I implemented own internal rss feed server before, and I have to rebuild the server when I add a new feed.
Actually, just for fun.

## Notes

- Custom feed using javascript
- Inject listening address on `0.0.0.0:80` on build docker
- For `development`, server run on `127.0.0.1:9000`
- New `sobek.vm` for every rss feed request from client
- Using `sobek` because of lightweight and it is supported by grafana k6

## API Specification

### GET /reload

- Reload configuration file
- Pre-compile program (js)
- Response no content 200

### GET /feed/:name

- Run js script
- Look up into program table with `:name` key

## JS Custom Function

```js
// http get method
// set user-agent into header automatically
{
    status: int,
    header: object,
    body: string,
} = get(url: string, { headers: object })

// fetch rss feed with gofeed
{
    Title: string
    Items: []*gofeed.Item
} = fetchFeed(url: string)
feeds.Items = feeds.Items.filter()

// convert gofeed to map[string]any (general model)
convertMapFeed(resp)

// parse html
// ignore error because want to run in chain
{
    length: int64
    each: callback func(fn)
    attr: func() string
    find: func(selector: string) object
} = parseHTML(html: string)

const html = parseHTML(content: string)
html.
    find(".container").
    each((i, elm) => console.log(elm.attr("src")))

// parse awesome markdown
{
    Topic: string
    Title: string
    Link: string
    Description: string
}: []Awesome = parseMarkdownAwesome(content: string)

// implement using fmt.Printf("%#v\n", v)
console.log(val: any)

// throw panic
exit("error message")
```

## Configuration

### From js file

- purpose: custom feed from difference source (json, html)
- config.yaml - Add feeds table

```yaml
feeds:
  myfeed:
    file: ./rules/myfeed.js
```

- rules/myfeed.js - JS file

```yaml
const url = "http://myfeed.com/atom"
const resp = get(url)
if (resp.status !== 200) exit("get bad response status");

const body = JSON.parse(resp.body);
const items = body.map((item) => ({
  title: item.name.trim(),
  link: item.link,
  author: item.author,
  created: new Date(),
  updated: new Date(),
}));

const result = {
  title: "My Feed",
  link: "http://myfeed.com/",
  items: items,
  author: "kk",
  created: new Date(),
  updated: new Date(),
};
result;
```

### From url + script

- purpose: filter out with keyword or regex
- require `url`, but `script` is optional
- config.yaml - Add feeds table

```yaml
feeds:
  myfeed:
    url: http://myfeed.com/atom
    script: >-
      feeds.Items =
      feeds.Items.filter(item => item.Title === "my feed 1");
```

- parse script into default template
- using go `text/template`

```yaml
const url = "{{.URL}}";
const url = "http://myfeed.com/atom";
let feeds = fetchFeed(url);

// inject here
{{.Script}};
feeds.Items = feeds.Items.filter(item => item.Title === "my feed 1");

convertMapFeed(feeds);
```
