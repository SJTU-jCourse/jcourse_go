from copy import deepcopy
import os
from typing import List
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common.action_chains import ActionChains
import requests
import time
import argparse

class major:
    main_class_table:dict = {
        "01": "哲学",
        "02": "经济学",
        "03": "法学",
        "04": "教育学",
        "05": "文学",
        "06": "历史学",
        "07": "理学",
        "08": "工学",
        "09": "农学",
        "10": "医学",
        "11": "军事学",
        "12": "管理学",
        "13": "艺术学",
    }
    # ref: https://yzb.sjtu.edu.cn/2022.pdf
    department_table:dict = {
        "010": "船舶海洋与建筑工程学院",
        "030": "电子信息与电气工程学院",
        "020": "机械与动力工程学院",
        "050": "材料科学与工程学院（含塑性研究院）",
        "071": "数学科学学院",
        "072": "物理与天文学院（含李政道研究所）",
        "080": "生命科学技术学院（含系统生物医学研究院）",
        "082": "生物医学工程学院（含Med-X研究院）",
        "090": "人文学院",
        "110": "化学化工学院",
        "120": "安泰经济与管理学院",
        "130": "国际与公共事务学院",
        "140": "外国语学院",
        "150": "农业与生物学院",
        "160": "环境科学与工程学院",
        "170": "药学院",
        "190": "凯原法学院",
        "200": "媒体与传播学院",
        "230": "马克思主义学院",
        "251": "体育系",
        "260": "上海交大-巴黎高科卓越工程师学院",
        "270": "上海交大-南加州大学文化创意产业学院",
        "280": "中英国际低碳学院",
        "350": "教育学院",
        "351": "中美物流研究院",
        "370": "密西根学院",
        "380": "上海高级金融学院",
        "413": "航空航天学院",
        "430": "设计学院",
        "440": "海洋学院",
        "450": "智慧能源创新学院",
        "700": "医学院",
        "710": "基础医学院",
        "711": "公共卫生学院",
        "712": "护理学院",
        "720": "瑞金医院",
        "721": "仁济医院",
        "722": "新华医院",
        "723": "第九人民医院",
        "724": "第一人民医院",
        "725": "第六人民医院",
        "727": "儿童医院",
        "728": "胸科医院",
        "729": "精神卫生中心",
        "730": "国际和平妇幼保健院",
        "731": "儿童医学中心",
        "732": "同仁医院",
    }
    def __init__(self, name:str, code:str, zyh_id:str, jxzxjhxx_id:str, xz:str, min_points:str, jg_id:str):
        self.name = name
        self.code = code
        self.zyh_id = zyh_id # 专业号
        self.jxzxjhxx_id = jxzxjhxx_id
        self.total_year = xz # 学制
        self.min_points = min_points # 最低学分
        self.major_class = self.get_main_class() # 专业类
        self.set_department(jg_id)
    def __csv__(self):
        return f"{self.zyh_id},{self.jxzxjhxx_id},{self.name},{self.code},{self.total_year},{self.min_points},{self.major_class},{self.department_name}"
    def __str__(self):
        return self.__csv__()
    def get_main_class(self)->str:
        if len(self.code) < 2 or self.code[:2] not in major.main_class_table:
            return "未知类别"
        return major.main_class_table[self.code[:2]]
    def set_department(self, id:str)->None:
        if len(id) < 3 or id[:3] not in major.department_table:
            self.department_name = "未知院系"
            return
        self.department_name = major.department_table[id[:3]]
class course:
    def __init__(self, code:str, name:str, credit:str, suggest_year:str, suggest_semester:str, department:str) -> None:
        self.code = code
        self.name = name
        self.credit = credit
        self.suggest_year = suggest_year
        self.suggest_semester = suggest_semester
        self.department = department
    def __str__(self):
        return f"{self.code},{self.name},{self.credit},{self.suggest_year},{self.suggest_semester},{self.department}"
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
session = requests.Session()
def login(jaccount_name:str = "", password: str = ""):
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
    for cookie in cookies:
        session.cookies.set(cookie['name'], cookie['value'])


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
        "queryModel.showCount": "30",
        "queryModel.currentPage": "1",
        "queryModel.sortName": "",
        "queryModel.sortOrder": "des",
        "time": "1"
    }
    def generate_training_plan(tp:dict):
        required_keys = ['zyh', 'zymc', 'zyh_id', 'jxzxjhxx_id', 'jg_id']
        min_points = "0"
        if 'zdxf' in tp:
            min_points = tp['zdxf']
        xz = '0'
        if 'xz' in tp:
            xz = tp['xz']
        for key in required_keys:
            if key not in tp:
                print(f"Key '{key}' is missing in the {tp}")
                return
        _major = major(code=tp['zyh'], name=tp['zymc'], zyh_id=tp['zyh_id'],
                       jxzxjhxx_id=tp['jxzxjhxx_id'], min_points=min_points, jg_id=tp['jg_id'], xz=xz)
        tps.append(trainingPlan(deepcopy(_major), tp['njmc']))
    callback = generate_training_plan
    abstract_postlist(url, data, callback)
"""
响应格式
trainingplan:
"zyh": "03010011", // 专业代码
"zyh_id": "DE7CE9928A87B81CE055F8163EE1DCCC", // 这个不知道用来干嘛的
"zymc": "工科平台(信息类)",
"jxzxjhxx_id":"FAD4594F4EF48A96E055F8163EE1DCCC"  // 这个是用于查询的
"jg_id" //学院代码
"xz" //学制 例如4（年制）
"zdxf" // 最低学分(未必有)

course:
"kch" //课程号
"kcmc" //课程名称
"xf" //学分(不一定有)
"kkbmmc" //开课学院名称(不一定有)
"jyxdxnm" //建议修读学年
"jyxdxqm" //建议修读学期
""
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
            required_keys = ['kch', 'kcmc', 'jyxdxnm', 'jyxdxqm', 'kkbmmc', 'xf']
            if 'xf' not in d:
                d['xf'] = '0'
            for key in required_keys:
                if key not in d:
                    print(f"Key '{key}' is missing in the {d}")
                    return
            tp.add_course(course(code=d['kch'], name=d['kcmc'], credit=d['xf']
            , suggest_year=d['jyxdxnm'], suggest_semester=d['jyxdxqm'], department=d['kkbmmc']))
        abstract_postlist(url, data, callback)

def save_csv(tps:List[trainingPlan], to="./data/trainingPlan.txt", append=False):
    dir = os.path.dirname(to)
    if not os.path.exists(dir):
        os.makedirs(dir)
    mode = 'w'
    if append:
        mode = 'a'
    with open(to, mode=mode) as f:
        for tp in tps:
            f.write(tp.__csv__() + "\n")

parser = argparse.ArgumentParser()
parser.add_argument('-n',"--name", help="jaccount name", default="")
parser.add_argument('-p',"--password", help="password", default="")
args = parser.parse_args()
TrainingPlans:List[trainingPlan] = []
login(args.name, args.password)
get_training_plans(TrainingPlans)
get_training_plan_courses(TrainingPlans)
save_csv(TrainingPlans)
# 关闭 Selenium WebDriver
driver.quit()



