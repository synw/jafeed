# -*- coding: utf-8 -*-

from django.db import models
from django.utils.translation import ugettext_lazy as _


class Feed(models.Model):
    title = models.CharField(_(u"Title"), max_length=255, blank=True)
    url = models.URLField(_(u"Url"), unique=True)
    is_active = models.BooleanField(_(u'Is active'), default=True)
    
    class Meta:
        verbose_name = _(u"Feed")
        verbose_name_plural = _(u"Feeds")
        ordering = ('title', 'url',)
    
    def __unicode__(self):
        return self.url