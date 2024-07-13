import os
import requests
from copy import deepcopy
import gzip
import json
import pprint
import time
from pypinyin import lazy_pinyin
from typing import List
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
from argparse import ArgumentParser

driver = webdriver.Chrome()
session = requests.Session()
parser = ArgumentParser(description="A simple scrapy to download the teacher profile and save as json.")
base_url = "https://faculty.sjtu.edu.cn/"
json_indent = 4
"""
type TeacherPO struct {
	gorm.Model
	Name       string `gorm:"index"`
	Code       string `gorm:"index;uniqueIndex"`
	Department string `gorm:"index"`
	Title      string
	Pinyin     string `gorm:"index"`
	PinyinAbbr string `gorm:"index"` // 拼音缩写
}
other attribute:
- profile_url
- head_image(url)
- mail
- profile_description
"""

def get_pinyin(chinese: str):
    return "".join(lazy_pinyin(chinese))
def get_pinyin_abbreviation(chinese: str):
    first_letters = [s[0] for s in lazy_pinyin(chinese)]
    return "".join(first_letters)
class Teacher:
    def __init__(self, name="", code="", department="", title="",
                 profile_url="", head_image="", mail="", profile_description="") -> None:
        self.name = name
        self.code = code
        self.department = department
        self.title = title
        self.pinyin = get_pinyin(name)
        self.pinyin_abbr = get_pinyin_abbreviation(name)
        self.profile_url = profile_url
        self.head_image = head_image
        self.mail = mail
        self.profile_description = profile_description
    def to_dict(self):
        return self.__dict__
    def to_json(self):
        return json.dumps(self)
def resp_to_teacher(resp: dict)->Teacher:
    teacher = Teacher()
    keys_map = dict()
    # ['teacherName', 'unit', 'url', 'prorank', 'picUrl', 'mail', 'profile']
    keys_map['teacherName'] = 'name'
    keys_map['unit']='department'
    keys_map['url']='profile_url'
    keys_map['picUrl']='head_image'
    keys_map['mail']='mail'
    keys_map['prorank']='title'
    keys_map['profile']='profile_description'
    keys_map['teacherId']='code' # NOTE
    for key in keys_map.keys():
        if key not in resp:
            print(f"Key '{key}({keys_map[key]})' is missing in the {resp}")
            teacher.__dict__[keys_map[key]] = ""
        if key == "teacherName":
            teacher.__dict__["pinyin"] = get_pinyin(resp[key])
            teacher.__dict__["pinyin_abbr"] = get_pinyin_abbreviation(resp[key])
            print(resp[key])
        if key == "picUrl":
            teacher.head_image = f"{base_url}{resp[key]}"
        else:
            teacher.__dict__[keys_map[key]] = resp[key]

    return teacher

def login(name:str, password:str):
    driver.get(url="https://i.sjtu.edu.cn/jaccountlogin")
    username_input = driver.find_element(By.ID, "input-login-user")
    password_input = driver.find_element(By.ID, "input-login-pass")
    username_input.send_keys("{}".format(name))
    password_input.send_keys("{}".format(password))
    time.sleep(10)
def set_cookie():
    cookies = driver.get_cookies()
    # 转化为 requests 库可以使用的格式
    for cookie in cookies:
        session.cookies.set(cookie['name'], cookie['value'])
def get_teacher_list(_print:bool=False, limit:int=10)->List[Teacher]:
    url = f"{base_url}/system/resource/tsites/advancesearch.jsp"
    query_param = {
        "collegeid": 0,
        "disciplineid": 0,
        "enrollid": 0,
        "pageindex": 1,
        "pagesize": limit,
        "rankid": 0,
        "degreeid": 0,
        "honorid": 0,
        "pinyin": "",
        "profilelen": 30,
        "teacherName": "",
        "searchDirection": "",
        "viewmode": 8,
        "viewid": 68237,
        "siteOwner": "1538087020",
        "viewUniqueId": "u11",
        "showlang": "zh_CN",
        "ispreview": False,
        "ellipsis": "...",
        "alignright": False,
        "productType": 0
    }
    resp = session.get(url=url, params=query_param)
    teacher_list = []
    if resp.status_code != 200:
        print(f"请求失败，状态码: {resp.status_code}")
        pprint.pprint(resp.text)
        return teacher_list
    print("请求成功")
    data:List[dict] = resp.json()['teacherData']
    if _print:
        pprint.pprint(data)
    for d in data:
        teacher_list.append(resp_to_teacher(d))
    return teacher_list
def save_json(data: List[Teacher], to="./data/teachers.json"):
    folder_path = os.path.dirname(to)
    if not os.path.exists(folder_path):
        os.makedirs(folder_path)
    with open(to, 'w') as f:
        json.dump([t.to_dict() for t in data], f, indent=json_indent)
def append_json(data: List[Teacher], to="./data/teachers.json"):
    with open(to, "r") as f:
        teachers = json.load(f)
    teachers.extend([t.to_dict() for t in data])
    with open(to, "w") as f:
        json.dump(teachers, f, indent=json_indent)
def automation():
    parser.add_argument("-n", "--name", type=str, action="store", help="jacount for login")
    parser.add_argument("-p", "--password", type=str, action="store", help="password for login")
    parser.add_argument("-d","--debug", action="store_true", help="In debug mode, the program will not save the result to the file.")
    parser.add_argument("-l", "--load_from", type=str, action="store", help="Load the result from the <load_from> json file.")
    parser.add_argument("-a", "--append_from", type=str, action="store", help="Append the result from the <append_from> json file.")
    parser.add_argument("-s", "--save_to", type=str, action="store", help="Save the result to the <save_to> json file.")
    parser.add_argument("--limit", type=int, help="The maximum number of request teachers")
    parser.add_argument("-v", "--verbose", action="store_true", help="Verbose mode")
def main():
    automation()
    args = parser.parse_args()
    init_teachers = []
    if not args.name or not args.password:
        print("jacount and password are required")
        exit(0)
    if args.load_from and os.path.exists(args.load_from):
        with open(args.load_from, "r") as f:
            init_teachers = json.load(f)
    if args.append_from and os.path.exists(args.append_from):
        with open(args.append_from, "r") as f:
            init_teachers.extend(json.load(f))

    login(name=args.name, password=args.password)
    set_cookie()
    new_teachers = get_teacher_list(args.verbose, args.limit)
    all_teachers = init_teachers + new_teachers
    print("当前教师数量：{}".format(len(all_teachers)))
    if not args.debug:
        if args.save_to and os.path.exists(args.save_to):
            save_json(all_teachers, to=args.save_to)
        else:
            save_json(all_teachers)

main()
driver.close()