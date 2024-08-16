from copy import deepcopy
import os
from typing import List, Dict, Callable

import requests


class Major:
    main_class_table: Dict[str, str] = {
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
    """
    date: 2024.7
    ref:
    - https://yzb.sjtu.edu.cn/2022.pdf
    - https://i.sjtu.edu.cn/  上课课程查询 表单 kkbm_id 项
    - SJTU规定学院代码为前三位，但有些单位在“开课学院”会填写更细粒度的5位代码，
    即下面的special_department, 需要以 https://i.sjtu.edu.cn/ 为准，新学院需要手动更新
    原则上，为XXX00的代码可以放入department_table，其他不匹配的代码需要放入special_department_table
    """
    department_table: Dict[str, str] = {
        "010": "船舶海洋与建筑工程学院",
        "030": "电子信息与电气工程学院",
        "020": "机械与动力工程学院",
        "050": "材料科学与工程学院",
        "071": "数学科学学院",
        "072": "物理与天文学院",
        "080": "生命科学技术学院",
        "082": "生物医学工程学院",
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
        "220": "继续教育学院",
        "230": "马克思主义学院",
        "240": "致远学院",
        "251": "体育系",
        "260": "上海交大-巴黎高科卓越工程师学院",
        "270": "上海交大-南加州大学文化创意产业学院",
        "280": "中英国际低碳学院",
        "330": "图书馆",
        "350": "教育学院",
        "351": "中美物流研究院",
        "370": "密西根学院",
        "380": "上海高级金融学院",
        "390": "创业学院",
        "400": "上海中医药大学",
        "401": "网络信息中心",
        "403": "档案文博管理中心",
        "404": "分析测试中心",
        "413": "航空航天学院",
        "410": "学生创新中心",
        "430": "设计学院",
        "440": "海洋学院",
        "450": "智慧能源创新学院",
        "490": "溥渊未来技术学院",
        "531": "人工智能学院",
        "602": "教务处",
        "603": "研究生院",
        "605": "国际交流与合作处",
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
    specail_department_table: Dict[str, str] = {
        "50500": "学指委（学生处、团委、人武部）合署",
        "50501": "军事教研室",
        "50511": "学生就业服务与职业发展中心",
        "80510": "校医院",
        "50520": "共青团上海交通大学委员会",
        "50514": "心理健康教育与咨询中心",
        "CF242766FA996146E055F8163EE1DCCC": "心理与行为科学研究院",
        "0E7A4795A03286A1E065F8163EE1DCCC": "张江高等研究院",

    }

    def __init__(self, name: str, code: str, zyh_id: str, jxzxjhxx_id: str,
                 xz: str, min_points: str, jg_id: str):
        self.name = name
        self.code = code
        self.degree = "本科"  # 学位
        self.zyh_id = zyh_id  # 专业号
        self.jxzxjhxx_id = jxzxjhxx_id
        self.total_year = xz  # 学制
        self.min_points = min_points  # 最低学分
        self.major_class = self.get_main_class()  # 专业类
        self.set_department(jg_id)

    def __csv__(self):
        return f"{self.zyh_id},{self.jxzxjhxx_id},\
{self.name},{self.code},{self.total_year},\
{self.min_points},{self.major_class},\
{self.department_name},{self.degree}"

    def __str__(self):
        return self.__csv__()

    def get_main_class(self) -> str:
        if len(self.code) < 2 or self.code[:2] not in Major.main_class_table:
            return "未知类别"
        return Major.main_class_table[self.code[:2]]

    def set_department(self, id: str) -> None:
        if len(id) < 3:
            self.department_name = "未知院系"
        elif id[:3] in Major.department_table:
            self.department_name = Major.department_table[id[:3]]
            return
        elif id in Major.specail_department_table:
            self.department_name = Major.specail_department_table[id]
            return
        self.department_name = "未知院系"


class Course:
    def __init__(self, code: str, name: str, credit: str, suggest_year: str,
                 suggest_semester: str, department: str) -> None:
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


class TrainingPlan:
    def __init__(self, major: Major, year):
        self.major = major
        self.year = year
        self.courses:List[Course] = []

    def __csv__(self):
        courses_csv = "\n".join([str(c) for c in self.courses])
        return f"{self.major.__csv__()},{self.year}\n{courses_csv}\n"

    def add_course(self, course: Course):
        self.courses.append(deepcopy(course))

    def __str__(self) -> str:
        return self.__csv__()

def abstract_postlist(session:requests.Session, url: str, post_data: Dict, callback: Callable[[Dict], None]):
    response = session.post(url, data=post_data)
    if response.status_code == 200:
        print("请求成功")
        try:
            data = response.json()['items']
            for d in data:
                callback(deepcopy(d))
        except Exception as e:
            print(response.text)
            print(e)

    else:
        print(f"请求失败，状态码: {response.status_code}")
        print(response.text)


def get_training_plans(session: requests.Session, tps: List[TrainingPlan]):
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
        "queryModel.sortOrder": "des",
        "time": "1"
    }

    def generate_training_plan(tp: Dict)->None:
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
        _major = Major(code=tp['zyh'], name=tp['zymc'], zyh_id=tp['zyh_id'],
                       jxzxjhxx_id=tp['jxzxjhxx_id'], min_points=min_points,
                       jg_id=tp['jg_id'], xz=xz)
        tps.append(TrainingPlan(deepcopy(_major), tp['njmc']))

    callback = generate_training_plan
    abstract_postlist(session=session, url=url, post_data=data, callback=callback)


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
"jyxdxnm" //建议修读学年 (2018-2019格式)
"jyxdxqm" //建议修读学期 (1,2,3)
""
"""


def get_training_plan_courses(session: requests.Session, tps: List[TrainingPlan]):
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
        def callback(d: Dict):
            required_keys = ['kch', 'kcmc', 'jyxdxnm', 'jyxdxqm', 'kkbmmc',
                             'xf']
            if 'xf' not in d:
                d['xf'] = '0'
            for key in required_keys:
                if key not in d:
                    print(f"Key '{key}' is missing in the {d}")
                    return
            tp.add_course(Course(code=d['kch'], name=d['kcmc'], credit=d['xf']
                                 , suggest_year=d['jyxdxnm'],
                                 suggest_semester=d['jyxdxqm'],
                                 department=d['kkbmmc']))

        abstract_postlist(session=session, url=url, post_data=data, callback=callback)


def save_csv(tps: List[TrainingPlan], to:str="./data/trainingPlan.txt",
             append=False):
    dir = os.path.dirname(to)
    if not os.path.exists(dir):
        os.makedirs(dir)
    mode = 'w'
    if append:
        mode = 'a'
    with open(to, mode=mode) as f:
        for tp in tps:
            f.write(tp.__csv__() + "\n")

if __name__ == "__main__":
    if __package__ is None:
        import sys
        from os import path
        sys.path.append(path.dirname( path.dirname(path.abspath(__file__))))
        from common import Automator
    else:
        from .common import Automator
    automator = Automator(description="Get training-plans")
    parser = automator.parser
    args = parser.parse_args()
    trainingPlans: List[TrainingPlan] = []
    automator.login(args.name, args.password)
    session = automator.get_session()
    get_training_plans(session=session, tps=trainingPlans)
    get_training_plan_courses(session=session,tps=trainingPlans)
    if not args.debug:
        save_csv(trainingPlans)
    else:
        print(trainingPlans[:10])
    automator.end()
