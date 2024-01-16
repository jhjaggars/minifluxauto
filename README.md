Basic automation tool for expiring some feed items from miniflux

You can expire items from feeds and categories by id. You can get the id from the url of the feed or category.

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
  categories_expire:
    7: 24h
    16: 48h
```
