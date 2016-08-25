# Jafeed

Rss and atom feeds agregator using Rethinkdb for the storage, Django for the UI and a Go module for the aggregation.

## Dependencies

- [Rethinkdb](http://rethinkdb.com)
- [Django](https://github.com/django/django) and [django-changefeed](https://github.com/synw/django-changefeed)

## Install

1. Database:

Create a `jafeed` database with a `feeds` table in Rethinkdb. Add a `Timestamp` secondary index to this table.

2. Django

Add `"changefeed",` and `"jafeed",` to INSTALLED_APPS

Add `url(r'^jafeed/', include('jafeed.urls')),` to urls.py

3. The Go module

Put the url of your domain in `jafeed/go/jafeed.config`. Default is set to `http://localhost:8000`.

Set a cronjob for `jafeed/go/jafeed` that will retrieve the feeds.

## Usage

Create some feeds in the Django admin and run `jafeed/go/jafeed` to retrieve the data and store it in the database.

Go to `/jafeed/` and see your feeds.

## Todo

- Changefeeds realtime notifications
- Categories
- More configuration info to use for the go worker