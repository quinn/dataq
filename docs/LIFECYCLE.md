## Data Lifecyle

### ExtractRequest
An extract request contains a `Kind` that has meaning to plugin. Fields:
* Hash: the content address of the object responsible for creating the Extract. Should reference either a `Permanode` or content linked to an `ExtractResponse`.
* Kind: the operation to be performed that will produce a piece of data. 

*example*
One `Kind: initial` request is made to the Gmail Plug-in. 
A `Kind: get_message` request is made, to get and store the api response for a given message ID.
A `Kind: get_page` with a page token. 

### ExtractResponse
Fields:
- Hash: the content address of the data produced by the extract. 
- Kind: the kind of content being stored
- RequestHash: the address of the request. 
- Transforms: a list of transform requests to be created. 
	- Kind: the name of the transform. So far, this has been 1:1 with the kind of data. 
	- metadata: Any additional information necessary for the transform. So far, the content itself is all that is necessary. This will always be empty for all current use-cases. 

*example*
From gmail, a `Kind: page` response will be created. linked to api page data. It will also contain a single transform of type `page`. 
A `Kind: message` will also have a single transform of type `message`. 
### TransformRequest:
Any action or behavior performed on a piece of data is considered a transform. Each transform step should reference a hash for the piece of data in the CAS.
Fields:
* Hash: the address of the data to transform
* Kind: the kind of transform to be applied to the data. So far, this has been 1:1 with the kind of the data. 

*example*
a `Kind: page` with the address to a page api response from gmail.

### TransformResponse
The response from a transform. Enumerates next steps for the host. 
* Hash: value can be copied from RequestHash.
* Kind: N/A. Can be copied from Request. 
* RequestHash: address of request
* Extracts: A list of extracts.
	* Kind: the operation to perform
	* metadata: any data necessary to perform the operation
* Permanodes: a list of nodes to be managed by this plugin. These have semantic meaning to DataQ, in terms of the higher level domain of the application (it provides tooling for interacting with emails as a concept separate from a particular plugin)
	* Kind: the type of object within the DataQ schema
	* DataSources:
		* Plugin: the unique identifier for the plugin
		* Key: A value that can be used to uniquely identify the plugin's internal representation of the permanode
	* Payload: The data, structured to match the expected Kind. 

*example*
transform will scan the page for two things: the next page to extract, as well as a list of messages to extract. the extracts will be structured like this:
* `get_page`: includes page cursor
* `get_message`: includes message ID
In the case of transforming a `Kind: message` for gmail, This will produce a permanode of `Kind: email` with the fields necessary to populate an email object. 

### Permanode
Fields:
- nonce: random data. 
- Kind: the type of object within the DataQ schema.

### PermanodeVersion
Fields:
- PermanodeHash: the hash of the permanode
- Timestamp: used for versioning the data of the permanode
- Payload: the data of the permanode
- Deleted: bool

*example*
the data of an email, formatted into the DataQ schema. 
### DataSource
Fields:
* PermanodeHash: The hash of the permanode
* TransformResponseHash: The hash of the TransformResponse that created this DS. 
* Plugin: the plugin identifier
* Key: A value that can be used to uniquely identify the plugin's internal representation of the permanode

*example*
A `Kind: email` permanode is created. A DS is created That links the permanode 
