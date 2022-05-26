# GophKeeper

GophKeeper is a client-server system that allows users securely store logins, passwords, binary data and other private information.

Interaction with the server occurs via the gRPC protocol, methods are described in ./proto/gophkeeper.proto

### Definitions

**User** - service user. Can connect to the server from different clients.

**Client** - client program with local repository of user data. Data in local repository is periodically synchronized with user data stored on the server.

**Server** is a server application that handles client requests for synchronization and updating of user data, as well as user registration and authentication. The server communicates with the Postgresql database, which stores user information and user data.

**User data** - The service allows to store the following types of data:

- Text strings
- Passwords
- Binary data
- Credit cards

Each data item can contain metadata - a set of key-value pairs with additional data, such as login, the name of the bank, etc. The metadata is transmitted and stored on the server as a JSON string. The types of such data are not strictly defined and must be handled on the client side.

**Item** - the structure that consist the user data. In addition to the payload, the structure contains the following fields:

- ID - resource identifier in UUID format
- Version - resource current version. When the new item is created, _Version_ must be set to zero. This field should be changed ONLY on the server!
- CreatedAt - resource creation time
- DeletedAt - resource deletion time (if the resource is not deleted, then _nil_)

**Entry** - the object that consist an _item_ at local repository. Has a _pending_ flag.

**Event** - a structure representing an event that marks the update of user data. Consists of an operation definition (create, update, delete) and the payload (_item_)

**Data Version**
This field is stored in the user table in the database and represents the version number of the user's current data snapshot. Updating data in storage on the server is possible only if the version of the data on the client before the update matches the value of the Data Version of this user. Each time any item in the storage on the server is successfully updated, the Data Version field is incremented.

## User registration and authentication

...

## Data Synchronization Protocol

### Storing and updating data on the client

#### Storing

Each _Item_ stored on the client is wrapped in an **Entry** structure that has a **Pending** field indicating that the item has been locally created or modified and awaits confirmation from the server that the changes have been saved.

#### Update

Each local data update (creation or update item) is sent to the server wrapped in an _Event_ structure. An updated (or created) _item_ is marked as _pending_. For locally created items _Version_ field must be set to 0 (_Version_ field will be incremented on the server).

#### Delete

Deleted items are marked with timestamp DeletedAt. Payload should be erased.

### Sending updates to the server

All _events_ are sent to the server as a batch with a certain frequency. Together with the event package, the latest up-to-date _Data Version_ is sent. If it matches the given user's _Data Version_ on the server, the changes are accepted. Updated data comes from the server as a response to the next update request.

### Synchronizing data with the server

A _WhatsNew_ request is sent with a certain frequency. The request specifies the current Data Version of the client. If it matches the Data Version on the server, the OK status is returned. Otherwise, "_download the updates_" error is returned. In that case client invokes _DownloadUpdates_ method with JSON objectwhich contains a table <item ID>: <item version> for all local items. The server analyses the table and sends all new or modified items to the client in the response.

#### Parsing received data

Data parsing blocks receiving and sending updates.
When a new set of data is received from the server, the received _items_ are merged with those stored locally. Each _item_ received from the server is analyzed:

- if there is no _item_ with such ID in the local repository, then this item is considered created in a session on another client and is added to the existing storage _entries_.
- if there is the _item_ with such ID in local repository and the _Version_ of the received _item_ is newer:
- - if local _item_ has _Pending_ flag set and the payload is equal to the received item's payload, the update is considered approved by the server, _Pending_ flag is unset and the local _Version_ field is updates.
- - if local _item_ has _Pending_ flag unset, the item considered changed on another client, payload is replaced and _Version_ field is renewed.
- in all other cases, a **conflict resolution** procedure is performed.

TODO: It would be good to implement a mechanism that tracks long-term pending items and communicates them to the user. User can command "send again".

#### Conflict resolution

In disputable cases, the user is shown the _item_ payload that came with the latest update from the server, and the _item_ stored locally. The user is given the choice of which version of the data to accept as valid.

- if the user selects _item_ data from the server, the _item_ stored locally is replaced.
- if the user prefers the local _item_, then the _Version_ of the _item_ that came from the server is accepted, but the local _item_ is sent to the server within a new _event_. Thus, the data from the latest Data Version from the server does not conflict with the data stored locally.
