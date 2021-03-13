# FOM-Blackboard 
This Discord Bot can read from the Online-Campus Blackboard and write it's messages to an channel.

This is the basic application Design
```
+-------------------------------------------------+
| +-------------+ +-------------+ +-------------+ |
| |Auth|Context | |Auth Cookie  | |  MSG|Nodes  | |
| +-----^----+--+ +--^---+------+ +----+--------+ |
|       |    |       |   |             ^          |
|       |    +-------+   | +-----------+          |
|       |     |----------+ |                      |
+-------------v-----------------------------------+
        |POST Login        ^ GET with Session Cookie
+--------------+           |    +-----------------+
|              +-----------+    |                 |
|  Go Client   ^                |   Discord API   |
|              +---------------->                 |
+--------------+  Write Channel +-----------------+
```

We are using Webhooks to write our information into different channels. Currently we read out **Blackboard Messages** and **Zoom Notifications**. 

## Discord Server Config
- Create a channel #blackboard and create a webhook for `FOM_WEBHOOK` env-var
- Create a channel #zoom-notify and create a webhook for `FOM_WEBHOOK_COURSES` env-var

## Application Config 
- Use the env-Vars `FOM_USER` and `FOM_PWD` to set your login credentials. The programm needs a valid OC Login to authenticate against the Blackboard API. The Credentials must be encoded via base64 to stop Shoulder-Surfers from copying your valuable Online-Campus Credentials
- Use the env-var `FOM_WEBHOOK` to set the channe

Name | Use 
--- | --- 
FOM_WEBHOOK_COURSES | Course Notification Webhook
FOM_WEBHOOK | Blackboard Webhook
FOM_USER | Base64 User
FOM_WPD | Base64 Password

## Reverse Engineering Shizzle
In the /samples Folder some responses from the OC are saved. These can be used for testing and parsing

### Steps for Login:
- Get Login JSESSIONID
- Perform Login on Login.do with Username and Shit (you get a session cookie)
- Append session cookie and then perform get requests

## ToDo
- Parse non-Zoom Conference Invites from private messages
- Multi-Tenant Configuration for BWIT-Students
