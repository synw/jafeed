# -*- coding: utf-8 -*-

from django.contrib import admin
from jafeed.models import Feed


@admin.register(Feed)
class Feedadmin(admin.ModelAdmin):
    list_display = ['title', 'url', 'is_active']
    list_display_links = ['title', 'url']
