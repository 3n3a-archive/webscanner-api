# Web Scanner API

This API will Scan a given hostname for server names

## Features

- [ ] Switch to labstack/echo/v5 instead of Gin
- [ ] Add scanning for other features than generator meta (Headers, Files, Sitemap Urla, Robots.Txt Url/)
- [x] Look at url being submitted and extract base url + schema (https:// + example.com)


## Deployment

The two GeoIP Databases need to be provided by the runtime or your local environment. I set it up as follows:

```
# Folder in this repo (only local)
geodb:
    - GeoLite2-ASN.mmdb
    - GeoLite2-City.mmdb
```

## Dev

### Coroutines

**ErrorGroup**

As described in [this lovely blog post](https://bostonc.dev/blog/go-errgroup)
.

There's also more docs on go's own page.


### Live-Reload

* [Air](https://github.com/cosmtrek/air)

Just use go install to get
