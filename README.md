webhog
======

**webhog** is a package that stores and downloads a given URL (including js, css, and images) for offline use and uploads it to a given AWS-S3 account (more persistance options to come).

##Installation
`go get github.com/johnernaut/webhog`

You may also want to import the given SQL file into your database.

##Usage
Make a `POST` request to `http://localhost:3000/scrape` with a header set to value `X-API-KEY: SCRAPEAPI`.  Pass in a JSON value of the URL you'd like to fetch: `{ "url": "http://facebook.com"}` (as an example).  You'll notice an `Ent dir: /blah/blah/blah` printed to the console - your assets are saved there.  To test, open the given `index.html` file.

##Configuration
Create a `webhog.yml` file in the running directory.  The following options are supported:
```
development:
  mysql: 127.0.0.1:3306
  db_name: webhog_development
  api_key: SCRAPEAPI
  aws_key: AWSKEY
  aws_secret: AWSSECRET
production:
staging:
```

##TODO
* add in configuration options for S3 key, API key, etc...
* add configuration option to specify DB name / hosts
* finish ZIP / Upload to S3 functionality
