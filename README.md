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

## Invite Link for the Bot
https://discord.com/oauth2/authorize?client_id=780869817921962027&scope=bot
- Bot is currently set to public. Needs to be changed after it's deployed on the Discord Server

## Discord Server Config
- Create a channel #blackboard and paste the channel id into the `FOM_CHANNEL` env-var
- Create a permission role for the bot. The bot needs to have write/read permissions for the channel. Permission Integer is 190464

## Application Config 
- Use the env-Vars `FOM_USER` and `FOM_PWD` to set your login credentials. The programm needs a valid OC Login to authenticate against the Blackboard API. 
- Use the env-var `FOM_DTOKEN` to set the authentication token for the discord bot.
- User the env-Var `FOM_CHANNEL` to set binding of channel. `export FOM_CHANNEL=780873287126220850`

## Reverse Engineering Shizzle
In the /samples Folder some responses from the OC are saved. These can be used for testing and parsing

### Steps for Login:
- Get Login JSESSIONID
- Perform Login on Login.do with Username and Shit (you get a session cookie)
- Append session cookie and then perform get requests

