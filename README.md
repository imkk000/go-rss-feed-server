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

// implement using fmt.Println(v)
console.log(val)

// throw panic
exit("error message")
```

## Configuration

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
