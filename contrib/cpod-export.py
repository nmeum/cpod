#!/usr/bin/env python2

import os
import sys
import datetime
import feedparser
import xml.etree.ElementTree as ET

if len(sys.argv) <= 1:
    sys.exit("USAGE: %s [PATH]" % sys.argv[0])

root = ET.Element("opml")
root.set("version", "2.0")

head = ET.SubElement(root, "head")
body = ET.SubElement(root, "body")

title = ET.SubElement(head, "title")
title.text = "Podcast subscriptions"

created = ET.SubElement(head, "dateCreated")
created.text = str(datetime.datetime.now().time()) # FIXME

try:
    cpod_root = os.environ["XDG_CONFIG_HOME"]
except KeyError:
    cpod_root = os.path.join(os.environ["HOME"], ".config")

url_list = open(os.path.join(cpod_root, "cpod", "urls"), "r")
urls = url_list.readlines()

for url in urls:
    url = url.replace("\n", "").encode("utf-8")
    feed = feedparser.parse(url)

    element = ET.SubElement(body, "outline")
    element.set("xmlUrl", url)
    element.set("text", feed['feed']['title'])

tree = ET.ElementTree(root)
tree.write(sys.argv[1])
