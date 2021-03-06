.. Autogenerated by Gollum RST generator (docs/generator/*.go)

ElasticSearch
=============

The ElasticSearch producer sends messages to elastic search using the bulk
http API. The producer expects a json payload.




Parameters
----------

**Enable** (default: true)

  Switches this plugin on or off.
  

**Retry/Count**

  Set the amount of retries before a Elasticsearch request
  fail finally.
  By default this parameter is set to "3".
  
  

**Retry/TimeToWaitSec**

  This value denotes the time in seconds after which a
  failed dataset will be  transmitted again.
  By default this parameter is set to "3".
  
  

**SetGzip**

  This value enables or disables gzip compression for Elasticsearch
  requests (disabled by default). This option is used one to one for the library
  package. See http://godoc.org/gopkg.in/olivere/elastic.v5#SetGzip
  By default this parameter is set to "false".
  
  

**Servers**

  This value defines a list of servers to connect to.
  
  

**User**

  This value used as the username for the elasticsearch server.
  By default this parameter is set to "".
  
  

**Password**

  This value used as the password for the elasticsearch server.
  By default this parameter is set to "".
  
  

**StreamProperties**

  This value defines the mapping and settings for each stream.
  As index use the stream name here.
  
  

**StreamProperties/<streamName>/Index**

  The value defines the Elasticsearch
  index used for the stream.
  
  

**StreamProperties/<streamName>/Type**

  This value defines the document type
  used for the stream.
  
  

**StreamProperties/<streamName>/TimeBasedIndex**

  This value can be set to "true"
  to append the date of the message to the index as in "<index>_<TimeBasedFormat>".
  NOTE: This setting incurs a performance penalty because it is necessary to
  check if an index exists for each message!
  By default this parameter is set to "false".
  
  

**StreamProperties/<streamName>/TimeBasedFormat**

  This value can be set to a valid
  go time format string to be used with DayBasedIndex.
  By default this parameter is set to "2006-01-02".
  
  

**StreamProperties/<streamName>/Mapping**

  This value is a map which is used
  for the document field mapping. As document type, the already defined type is
  reused for the field mapping. See
  https://www.elastic.co/guide/en/elasticsearch/reference/5.4/indices-create-index.html#mappings
  
  

**StreamProperties/<streamName>/Settings**

  This value is a map which is used
  for the index settings. See
  https://www.elastic.co/guide/en/elasticsearch/reference/5.4/indices-create-index.html#mappings
  
  

Parameters (from core.BatchedProducer)
--------------------------------------

**Batch/MaxCount** (default: 8192)

  Defines the maximum number of messages per batch. If this
  limit is reached a flush is always triggered.
  By default this parameter is set to 8192.
  
  

**Batch/FlushCount** (default: 4096)

  Defines the minimum number of messages required to flush
  a batch. If this limit is reached a flush might be triggered.
  By default this parameter is set to 4096.
  
  

**Batch/TimeoutSec** (default: 5, unit: sec)

  Defines the maximum time in seconds messages can stay in
  the internal buffer before being flushed.
  By default this parameter is set to 5.
  
  

Parameters (from core.SimpleProducer)
-------------------------------------

**Streams**

  Defines a list of streams the producer will receive from. This
  parameter is mandatory. Specifying "*" causes the producer to receive messages
  from all streams except internal internal ones (e.g. _GOLLUM_).
  By default this parameter is set to an empty list.
  
  

**FallbackStream**

  Defines a stream to route messages to if delivery fails.
  The message is reset to its original state before being routed, i.e. all
  modifications done to the message after leaving the consumer are removed.
  Setting this paramater to "" will cause messages to be discared when delivery
  fails.
  
  

**ShutdownTimeoutMs** (default: 1000, unit: ms)

  Defines the maximum time in milliseconds a producer is
  allowed to take to shut down. After this timeout the producer is always
  considered to have shut down.  Decreasing this value may lead to lost
  messages during shutdown. Raising it may increase shutdown time.
  
  

**Modulators**

  Defines a list of modulators to be applied to a message when
  it arrives at this producer. If a modulator changes the stream of a message
  the message is NOT routed to this stream anymore.
  By default this parameter is set to an empty list.
  
  

Examples
--------

This example starts a simple twitter example producer for local running ElasticSearch:

.. code-block:: yaml

	 producerElasticSearch:
	   Type: producer.ElasticSearch
	   Streams: tweets_stream
	   SetGzip: true
	   Servers:
	     - http://127.0.0.1:9200
	   StreamProperties:
	     tweets_stream:
	       Index: twitter
	       TimeBasedIndex: true
	       Type: tweet
	       Mapping:
	         # index mapping for payload
	         user: keyword
	         message: text
	       Settings:
	         number_of_shards: 1
	         number_of_replicas: 1





