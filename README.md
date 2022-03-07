# FOM-Blackboard 

FOM-Blackboard is an AWS Lambda Function to parse the posts from FOM University of Applied Sciences for Economics and Management OnlineCampus. It can provide your local discord channel with the texts from the OnlineCampus, so users can manage their study efforts in one single place.

The Lambda Function runs once a day, then sends all the messages from today into the configured channel.

## Configuration
Set these 3 enviroment variables to access the OnlineCampus and the Discord Channel

```bash
export OC_USER = ""
export OC_PWD = ""
export DISCORD_CHANNEL = ""
```


## Deployment 
Run commands and upload the .zip folder to AWS

```bash
cd .env/lib/python3.8/site-packages/
zip -r ../../../../libs-deployment-package.zip .
cd ../../../../
 zip -g libs-deployment-package.zip lambda_function.py 
```