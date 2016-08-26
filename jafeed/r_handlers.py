# -*- coding: utf-8 -*-

import rethinkdb as r


def feed_handlers(database, table, change):
    print "Jafeed handler ->"
    print 'DB: '+str(database)
    print 'Table: '+str(table)
    print 'Old values: '+str(change['old_val'])
    print 'New values: '+str(change['new_val'])
    return

def r_query():
    return r.db("jafeed").table("feeds").changes()
