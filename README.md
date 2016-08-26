# Jafeed

Rss and atom feeds agregator using Rethinkdb for the storage, Django for the UI and a Go module for the aggregation
engine.

## Dependencies

- [Rethinkdb](http://rethinkdb.com)
- [Django](https://github.com/django/django) and [django-changefeed](https://github.com/synw/django-changefeed)

## Install

Clone the repo and `pip install rethinkdb`
and [install django-changefeed](http://django-changefeed.readthedocs.io/en/latest/src/install.html)

#### Database:

Create a `jafeed` database with a `feeds` table in Rethinkdb. Add a `Timestamp` secondary index to this table.

#### Django

Add `"changefeed",` and `"jafeed",` to INSTALLED_APPS

Add `url(r'^jafeed/', include('jafeed.urls')),` to urls.py

Run the migrations.

#### The Go module

Put the url of your domain in `jafeed/go/jafeed.config`. Default is set to `http://localhost:8000`. This is used by
the go worker to retrieve the feeds list from Django.

Set a cronjob for `jafeed/go/jafeed` that will retrieve the feeds.

## Usage

Create some feeds in the Django admin and run `jafeed/go/jafeed` to retrieve the data and store it in the database.

Go to `/jafeed/` and see your feeds.

## Realtime notifications

Displays a widget that notifies the number of new posts to the user.

Install [django-instant](https://github.com/synw/django-instant) and 
[Centrifugo](https://github.com/centrifugal/centrifugo)

Create a template ``instant/extra_handlers.js`` with this content:

  ```django
{% include "jafeed/js/handlers.js" %}
  ```

Add the counter widget in your templates where you want it:

  ```django
{% include "jafeed/js/handlers.js" %}
  ```

## Todo

- Categories
- More configuration info to use for the go worker