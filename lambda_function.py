from wsgiref import headers
import requests
import os
from bs4 import BeautifulSoup
from datetime import datetime
from bs4.element import Comment
import json


def lambda_handler(event, context):
    today = datetime.today().date()

    print(f"Creating OC HTTP-Session and parsing Blackboard HTML")
    session = create_oc_session()
    posts = get_blackboard_posts(session)
    print("Found posts:",posts)
    
    # Only select posts which were created today
    new_posts = [
        get_single_post(session=session, post=p) for p in posts if p["date"] == today
    ]

    print(f"Sending {len(new_posts)} new posts to discord channel")
    resp_codes = [send_discord_message(post=post,session=session) for post in new_posts]
    return resp_codes

def send_discord_message(post,session):
    webhook = os.environ["DISCORD_CHANNEL"]

    message = {
        "username" : "FOM-Blackboard",
        "embeds" : [
            {
                "title" : post['title'],
                "url" : "https://campus.bildungscentrum.de" + post['link'],
                "description" : post['text'],
                "color" : "3066993"
            }
        ]
    }

    resp = session.post(webhook,data=json.dumps(message),headers={'content-type' : 'application/json'})
    if (resp.status_code) == 204:
        print("Webhook executed successfully")
    else: 
        print("Something went wrong with the webhook",resp.status_code)
        print(resp.text)
    return resp.status_code

def get_single_post(session, post):
    single_post = session.get(
        "https://campus.bildungscentrum.de" + post["link"]
    ).content

    soup = BeautifulSoup(single_post, "html.parser")

    content = soup.find("div", {"id": "content"})

    post['title'] = content.find("b").text

    post['text'] = " ".join(
        text_from_html(content)
        .strip("| Meine Hochschule")
        .strip("zurück zur Übersicht")
        .split()
    )
    return post


def text_from_html(soup):
    texts = soup.findAll(text=True)
    visible_texts = filter(tag_visible, texts)
    return " ".join(t.strip() for t in visible_texts)


def tag_visible(element):
    # also filters h1 and b tags
    if element.parent.name in [
        "style",
        "script",
        "head",
        "title",
        "meta",
        "[document]",
        "h1",
        "b",
    ]:
        return False
    if isinstance(element, Comment):
        return False
    return True


def get_blackboard_posts(session):
    msg_size = 10

    page = session.get(
        f"https://campus.bildungscentrum.de/nfcampus/startapi/blackboardsite?page=0&size={msg_size}&tab=Blackboard&tab_=&n=5003&semester=&newstyp="
    ).json()

    return [parse_blackboard_html_msg(msg) for msg in page["data"]]


def parse_blackboard_html_msg(msg):
    soup = BeautifulSoup(msg[1], "html.parser")

    msg_date = soup.find("span", {"class": "date"}).text
    return {
        "date": datetime.strptime(msg_date.split(" ")[1], "%Y-%m-%d").date(),
        "link": soup.find("a").attrs["href"],
    }


def create_oc_session():
    OC_USER = os.environ["OC_USER"]
    OC_PWD = os.environ["OC_PWD"]

    with requests.session() as c:
        payload = {"quelle": "LoginForm-BCW", "name": OC_USER, "password": OC_PWD}
        c.post("https://campus.bildungscentrum.de/nfcampus/Login.do", params=payload)

        return c
