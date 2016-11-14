import urllib
from datetime import datetime
from zipfile import ZipFile

import os

start = datetime.now()
with open('input/img_urls.txt', 'r') as input, ZipFile('py_out/py_out.zip', 'w') as zipper:

    for idx, line in enumerate(input.readlines()):
        line = line.replace('\n', '')
        filename = 'py_out/image_{}.jpg'.format(idx)
        if line:
            urllib.urlretrieve(line, filename)
            zipper.write(filename)
            os.remove(filename)

print "total time: {}".format((datetime.now() - start))
