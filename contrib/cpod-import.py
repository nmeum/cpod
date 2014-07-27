#!/usr/bin/env python2

import os
import sys
import xml.etree.ElementTree as ET

if len(sys.argv) <= 1:
    sys.exit("USAGE: %s [PATH]" % sys.argv[0])

tree = ET.parse(sys.argv[1])
root = tree.getroot()

urls = []
for outline in root.iter("outline"):
    urls.append(outline.get("xmlUrl"))

try:
    cpod_root = os.environ["XDG_CONFIG_HOME"]
except KeyError:
    cpod_root = os.path.join(os.environ["HOME"], ".config")

url_path = os.path.join(cpod_root, "cpod", "urls")
try:
    os.makedirs(os.path.dirname(url_path))
except FileExistsError:
    pass

url_list = open(url_path, "w")
for url in urls:
    url_list.write(url + "\n")
url_list.close()
