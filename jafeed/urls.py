# -*- coding: utf-8 -*-

from django.conf.urls import url
from jafeed.views import ConfigView, IndexView

urlpatterns = [
    url(r'^config/$', ConfigView.as_view(), name="jafeed-config"),
    url(r'^$', IndexView.as_view(), name="jafeed-index"),
    ]