# GophKeeper

GophKeeper is a client-server system that allows users securely store logins, passwords, binary data and other private information.

Interaction with the server occurs via the gRPC protocol, methods are described in ./proto/gophkeeper.proto

### Definitions

**User** - service user. Can connect to the server from different clients.

**Client** - client program with local storage of user data. Data in local storage is periodically synchronized with user data stored on the server.

**Session** - session between client and server.
For signing up (and signing in) a pair of login + password is specified. The login and password hash is stored in the user account in the database.
Each user can log in from different devices, for each of which a session is created. A pair of tokens (access token and refresh token) is generated for each session. When the access token is expired it must be refreshed using the one-time refresh token.

**User data** - The service allows to store the following types of data:

- Text strings
- Passwords
- Binary data
- Credit cards

  Each data item can contain metadata - a set of key-value pairs with additional data, such as login, the name of the bank, etc. The metadata is transmitted and stored on the server as a JSON string. The types of such data are not strictly defined and must be handled on the client side.

**Item** - The structure that consist the user data. In addition to the payload, the structure contains the following fields:

- ID - resource identifier in UUID format
- Version - resource 
- CreatedAt - resource creation time
- DeletedAt - resource deletion time (if the resource is not deleted, then _nil_)

**Event** - a structure representing an event that marks the update of data. Consists of an operation definition (create, update, delete) and an _Item_ data element.

**Data Version**
This field is stored in the user table in the database and represents the version number of the user's current data snapshot. Updating data in storage on the server is possible only if the version of the data on the client before the update matches the value of the Data Version of this user. Each time the data in the storage on the server is successfully updated, the Data Version field is incremented.

## Data Synchronization Protocol

### Storing and updating data on the client

#### Storing

Each _Item_ stored on the client is wrapped in an **Entry** structure that has a **Pending** field indicating that the item has been locally created or modified and awaits confirmation from the server that the changes have been saved. There is also an **OldVersion** field, which is a pointer to the version of this _Item_ from the data snapshot with the Data Version relevant for the local storage

#### Update

Each local data update (creation or update item) is wrapped in an _Event_ structure. An updated (or created) _item_ is marked as _pending_. If this is the first data change before an update confirmation is received from the server, then a copy of this _item_ is created corresponding to the last confirmed version of the data for possible conflict resolution. A pointer to the copy is placed in the _OldVersion_ field. If such a copy has already been created (the _OldVersion_ field is not equal to _nil_) and the confirmation has not yet been received from the server, the change is simply wrapped in a new _Event_.

#### Delete
Deleted items mark with timestamp DeletedAt. 

### Sending updates to the server

All _events_ are sent to the server as a batch with a certain frequency. Together with the event package, the latest up-to-date _Data Version_ is sent. If it matches the given user's _Data Version_ on the server, the changes are accepted. Updated data comes from the server as a response to the next update request.

### Synchronizing data with the server

A request to update the data is sent with a certain frequency. The request specifies the current Data Version of the client. If it matches the Data Version on the server, the error _"data version is up to date"_ is returned. If the Data Version of the client does not match the Data Version of the server, then a complete snapshot of the user data of the latest version is downloaded from the server.

#### Parsing received data

Data parsing blocks receiving and sending updates.
When a new snapshot of data is received from the server, the received _items_ are merged with those stored locally. Each _item_ received from the server is analyzed:

- if there is no _item_ with such ID in the local storage, then this item is considered created in a session on another client and is added to the existing storage _entries_.
- if there is an _item_ with such ID in local storage, but it is not marked as _pending_, this _item_ is considered updated in a session on another client and is replaced by the item received from the server.
- if there is an _item_ with such ID marked as _pending_ in local storage, then the content of the _items_ is compared:
- - if the contents match, the update is considered confirmed, the _pending_ flag is unset, the _OldVersion_ pointer is set to _nil_.
- - if the contents do not match, then the update is compared with a copy of the _OldVersion_ item. If the contents of the item received from the server match the _OldVersion_ copy, the local update is considered to be pending acknowledgment, but not yet processed by the server. Nothing is done, the new version of the data for this _item_ is considered committed.
- in all other cases, a **conflict resolution** procedure is performed.

#### Conflict resolution

In disputable cases, the user is shown the _item_ data that came with the latest update from the server, the _item_ data stored locally, and (if available) a copy of this _item_ data of the latest confirmed version of the data. The user is given the choice of which version of the data to accept as valid.

- if the user selects _item_ data from the server, the _item_ stored locally is replaced.
- if the user selects the local _item_ data, then the _item_ that came from the server is considered accepted, but is placed in the _OldVersion_ field, and the local _item_ is wrapped in a new _event_, which is waiting to be sent to the server. Thus, the data from the latest Data Version from the server does not conflict with the data stored locally.
- if the user selects data from the OldVersion copy, they are wrapped in a new _event_, and the _item_ received from the server is considered accepted and placed in the OldVersion.
