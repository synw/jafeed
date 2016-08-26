# -*- coding: utf-8 -*-
# Generated by Django 1.9.8 on 2016-08-26 08:45
from __future__ import unicode_literals

from django.db import migrations, models


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Feed',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('title', models.CharField(blank=True, max_length=255, verbose_name='Title')),
                ('url', models.URLField(unique=True, verbose_name='Url')),
                ('is_active', models.BooleanField(default=True, verbose_name='Is active')),
            ],
            options={
                'ordering': ('title', 'url'),
                'verbose_name': 'Feed',
                'verbose_name_plural': 'Feeds',
            },
        ),
    ]
