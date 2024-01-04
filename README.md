Basic automation tool for expiring some feed items from miniflux

Example config file

```yaml
miniflux:
  url: https://miniflux.url
  token: xxxxxx
  feeds_expire:
    62: 24h
    63: 24h
    9: 48h
    14: 48h
    11: 48h
    12: 48h
    40: 120h
```
