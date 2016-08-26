# -*- coding: utf-8 -*-

import rethinkdb as r
from instant import broadcast


def feed_handlers(database, table, change):
    _type = change["type"]
    if _type == "add":
        print "Jafeed handler ->"
        print ".. broadcasting message"
        broadcast(message='Update', event_class="__jafeed__")
    return

def r_query():
    return r.db("jafeed").table("feeds").changes(include_types=True)
