webhog
======

**webhog** is a package that stores and downloads a given URL (including js, css, and images) for offline use and uploads it to a given AWS-S3 account (more persistance options to come).

##Installation
`go get github.com/johnernaut/webhog`

##Usage
Make a `POST` request to `http://localhost:3000/scrape` with a header set to value `X-API-KEY: SCRAPEAPI`.  Pass in a JSON value of the URL you'd like to fetch: `{ "url": "http://facebook.com"}` (as an example).  You'll notice an `Ent dir: /blah/blah/blah` printed to the console - your assets are saved there.  To test, open the given `index.html` file.

##Configuration
Create a `webhog.yml` file in the running directory.  The following options are supported:
```
development:
  mongodb: mongodb://127.0.0.1:27017/webhog
  api_key: SCRAPEAPI
  aws_key: AWSKEY
  aws_secret: AWSSECRET
  bucket: mybucket
production:
staging:
```
The setting root-key is established via a `MARTINI_ENV` environment variable that you should set.