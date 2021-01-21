# Generated by Django 3.1.3 on 2020-11-28 18:38

import django.core.validators
from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('user', '0003_auto_20201126_1638'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='profile',
            name='phone',
        ),
        migrations.AddField(
            model_name='profile',
            name='subscribtionId',
            field=models.IntegerField(default=0, validators=[django.core.validators.MinValueValidator(0)]),
        ),
        migrations.AddField(
            model_name='profile',
            name='telegram',
            field=models.CharField(default='0000000', max_length=300),
        ),
    ]