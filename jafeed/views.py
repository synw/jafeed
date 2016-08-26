# -*- coding: utf-8 -*-

from django.views.generic import TemplateView
from rethinkdb import r
from changefeed.orm import R
from jafeed.models import Feed


class ConfigView(TemplateView):
    template_name = 'jafeed/index.cfg'
    
    def get_context_data(self, **kwargs):
        context = super(ConfigView, self).get_context_data(**kwargs)
        context['feeds'] = Feed.objects.filter(is_active=True)
        return context
    
    
class IndexView(TemplateView):
    template_name = 'jafeed/index.html'
    
    def get_context_data(self, **kwargs):
        context = super(IndexView, self).get_context_data(**kwargs)
        r_query = r.db("jafeed").table("feeds").order_by(index=r.desc("Timestamp")).slice(0,30)
        cursor = R.run_query(r_query)
        feeds = []
        for document in cursor:
            feeds.append("<h3>"+document['Title'].encode('utf-8')+"</h3>"+document['Description'].encode('utf-8'))
        context['feeds'] = feeds
        return context