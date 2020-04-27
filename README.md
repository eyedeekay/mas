# mas

Make a Site/Makefile Abuse Stopper: Just generate a site, as simply as possible.

This is the world's simplest static site generator. Maybe. Maybe that's an
overstatement. What it is for sure is super damned simple. You give it a
directory, which is full of files. If any of those directories contain ```.md```
files, they are processed in alphabetical order and turned into a single
index.html page. The general idea is "One directory, one page." HTML is allowed
in the markdown files as well.

Anything in the ./css ./js or ./images folder is copied directly to the root of
the site, and all Javascript and CSS is shared by all generated pages.

The output of these processes, page generation and static file preparation, is
eventually placed in the ./site directory, which you can copy to any web server.