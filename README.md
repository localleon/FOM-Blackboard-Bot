# FOM-Blackboard 
This Discord Bot can read from the Online-Campus Blackboard and write it's messages to an channel.


## Config 
Use the env-Vars `FOM_USER` and `FOM_PWD` to set your login credentials. The programm needs a valid OC Login to authenticate against the Blackboard API

## Reverse Engineering Shizzle
In the /samples Folder some responses from the OC are saved. These can be used for testing and parsing

### Steps for Login:
- Get Login JSESSIONID
- Perform Login on Login.do with Username and Shit (you get a session cookie)
- Append session cookie and then perform get requests

