import requests
import os
from bs4 import BeautifulSoup
from datetime import datetime
from bs4.element import Comment


def handler(event, context):
    today = datetime.today().date()

    print(f"Creating OC HTTP-Session and parsing Blackboard HTML")
    session = create_oc_session()
    posts = get_blackboard_posts(session)
    posts[9]["date"] = datetime.today().date()

    # Only select posts which were created today
    new_posts = [
        get_single_post(session=session, post=p) for p in posts if p["date"] == today
    ]

    print(f"Sending {len(new_posts)} new posts to discord channel")
    [send_discord_message(text) for text in new_posts]


def send_discord_message(text):
    print(text)


def get_single_post(session, post):
    single_post = session.get(
        "https://campus.bildungscentrum.de" + post["link"]
    ).content

    soup = BeautifulSoup(single_post, "html.parser")

    content = soup.find("div", {"id": "content"})
    return " ".join(
        text_from_html(content)
        .strip("| Meine Hochschule")
        .strip("zurück zur Übersicht")
        .split()
    )


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
