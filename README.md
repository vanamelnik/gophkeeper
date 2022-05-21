# GophKeeper

GophKeeper is a client-server system that allows users securely store logins, passwords, binary data and other private information.

Interaction with the server occurs via the gRPC protocol, methods are described in ./proto/gophkeeper.proto

## User sessions

For signing up (and signing in) a pair of login + password is specified. The login and password hash is stored in the user account in the database.
Each user can log in from different devices, for each of which a session is created. A pair of tokens (access token and refresh token) is generated for each session. When the access token is expired it must be refreshed using the one-time refresh token.

## Data types

The service allows to store the following types of data:

- Text strings
- Passwords
- Binary data
- Credit cards
  Each data item can contain metadata - a set of key-value pairs with additional data, such as login, the name of the bank, etc. The metadata is transmitted and stored on the server as a JSON string. The types of such data are not strictly defined and must be handled on the client side.

## Data Synchronization Protocol

### Понятия

**Item** - Структура, в которой хранятся пользовательские данные. Данные могут быть следующих типов:

- Пароль
- Текст
- Бинарные данные
- Данные банковской карты

**Event**
Структура, представляющая событие, знаменующее обновление данных. Состоит из определения операции (создание, обновление, удаление) и элемента данных _Item_.

**Data Version**
Это поле хранится в пользовательской таблице в базе данных и представляет собой номер версии текущего слепка данных пользователя.

### Описание

Each local data update is wrapped in an event object. Unsent events are marked as pending to avoid confusion when data is synchronized. As events accumulate, they can be batched to the server for storage. Each time the user store is updated on the server, the "data version" field is incremented. The server refuses to accept new events from the client if the client version of the data is out of date. In this case, before sending updates to the server, you should get fresh data from the server.
When requesting data from the server, a snapshot of all user data is returned with the current version number. The client must independently handle disputed points when synchronizing data.
After the updates are successfully sent to the server, the data version on the client is updated and the pending flags are cleared.
