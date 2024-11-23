I am working on a project called dataq, a general purpose tool for archiving and organizing data from various sources.

I would like to start with a protobuf based system for data extraction. Basically, I would like to be able to register plugins in a config, and then the plugins return the data they are responsible for. For example, a plugin could be:

1. something that scans folders for files
2. integrates with Gmail or Outlook for collecting emails
3. Connecting to icloud for collecting personal data
4. a web scraper for a product or service which does not provide an API

So, it should be fairly general purpose. I think, to start, the plugins should just return raw API data they are responsible for collecting.

Lets get started with something simple here, and the rest can be fleshed out later.
