from copy import deepcopy
import gzip
import json
import pprint
from typing import List
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
import requests
import time
class major:
    def __init__(self, name, code, zyh_id, jxzxjhxx_id):
        self.name = name
        self.code = code
        self.zyh_id = zyh_id
        self.jxzxjhxx_id = jxzxjhxx_id
    def __csv__(self):
        return f"{self.zyh_id},{self.jxzxjhxx_id},{self.name},{self.code}"
    def __str__(self):
        return self.__csv__()
class course:
    def __init__(self, code, name, credit) -> None:
        self.code = code
        self.name = name
        self.credit = credit
    def __str__(self):
        return f"{self.code},{self.name},{self.credit}"
    def __csv__(self):
        return self.__str__()
class trainingPlan:
    def __init__(self, major: major, year):
        self.major = major
        self.year = year
        self.courses = []
    def __csv__(self):
        courses_csv = "\n".join([str(c) for c in self.courses])
        return f"{self.major.__csv__()},{self.year}\n{courses_csv}\n"
    def add_course(self, course):
        self.courses.append(deepcopy(course))
    def __str__(self) -> str:
        return self.__csv__()
# 初始化 Selenium WebDriver
driver = webdriver.Chrome()

jaccount_name = "" # HINT: your jaccount
password = "" # HINT: your password
driver.get(url="https://i.sjtu.edu.cn/jaccountlogin")
username_input = driver.find_element(By.ID, "input-login-user")
password_input = driver.find_element(By.ID, "input-login-pass")
username_input.send_keys("{}".format(jaccount_name))
password_input.send_keys("{}".format(password))
time.sleep(10)
# 手输一下验证码

# 获取登录后的 Cookie
cookies = driver.get_cookies()

# 将 Cookie 转换为 requests 库可以使用的格式
session = requests.Session()
for cookie in cookies:
    session.cookies.set(cookie['name'], cookie['value'])

# 关闭 Selenium WebDriver
driver.quit()
def abstract_postlist(url:str, post_data:dict, callback):
    response = session.post(url, data=post_data)
    if response.status_code == 200:
        print("请求成功")
        data = response.json()['items']
        for d in data:
            callback(deepcopy(d))
    else:
        print(f"请求失败，状态码: {response.status_code}")
        print(response.text)
def get_training_plans(tps:List[trainingPlan]):
    url = "https://i.sjtu.edu.cn/jxzxjhgl/jxzxjhck_cxJxzxjhckIndex.html?doType=query&gnmkdm=N153540"
    data = {
        "jg_id": "",
        "njdm_id": "",
        "dlbs": "",
        "zyh_id": "",
        "currentPage_cx": "",
        "_search": "false",
        "nd": "",
        "queryModel.showCount": "5000",
        "queryModel.currentPage": "1",
        "queryModel.sortName": "",
        "queryModel.sortOrder": "asc",
        "time": "1"
    }
    def generate_training_plan(tp:dict):
        required_keys = ['zyh', 'zymc', 'zyh_id', 'jxzxjhxx_id']
        for key in required_keys:
            if key not in tp:
                print(f"Key '{key}' is missing in the {tp}")
                return
        _major = major(code=tp['zyh'], name=tp['zymc'], zyh_id=tp['zyh_id'], jxzxjhxx_id=tp['jxzxjhxx_id'])
        tps.append(trainingPlan(deepcopy(_major), tp['njmc']))
    callback = generate_training_plan
    abstract_postlist(url, data, callback)
"""
响应格式
items
"zyh": "03010011",
"zyh_id": "DE7CE9928A87B81CE055F8163EE1DCCC", // 这个不知道用来干嘛的
"zymc": "工科平台(信息类)",
"jxzxjhxx_id":"FAD4594F4EF48A96E055F8163EE1DCCC"  // 这个是用于查询的
"""
def get_training_plan_courses(tps: List[trainingPlan]):
    # 构建请求的 URL 和数据
    url = "https://i.sjtu.edu.cn/jxzxjhgl/jxzxjhkcxx_cxJxzxjhkcxxIndex.html?doType=query&gnmkdm=N153540"
    for tp in tps:
        data = {
            "jyxdxnm": "",
            "jyxdxqm": "",
            "yxxdxnm": "",
            "yxxdxqm": "",
            "shzt": "",
            "kch": "",
            "jxzxjhxx_id": str(tp.major.jxzxjhxx_id),
            "xdlx": "",
            "_search": "false",
            "nd": "",
            "queryModel.showCount": "100",
            "queryModel.currentPage": "1",
            "queryModel.sortName": "jyxdxnm,jyxdxqm,kch",
            "queryModel.sortOrder": "asc",
            "time": "1"
        }
        # err handle
        def callback(d:dict):
            required_keys = ['kch', 'kcmc', 'xf']
            for key in required_keys:
                if key not in d:
                    print(f"Key '{key}' is missing in the {d}")
                    return
            tp.add_course(course(code=d['kch'], name=d['kcmc'], credit=d['xf']))
        abstract_postlist(url, data, callback)

def save_csv(tps:List[trainingPlan], to="trainingPlan.csv", append=False):
    mode = 'w'
    if append:
        mode = 'a'
    with open(to, mode=mode) as f:
        for tp in tps:
            f.write(tp.__csv__() + "\n")


TrainingPlans:List[trainingPlan] = []
get_training_plans(TrainingPlans)
get_training_plan_courses(TrainingPlans)
save_csv(TrainingPlans)




