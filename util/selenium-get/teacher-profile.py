import os
import requests
import json
import pprint
from pypinyin import lazy_pinyin
from typing import List, Self

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
    def __eq__(self, other:Self):
        return self.to_dict() == other.to_dict()
    def __hash__(self):
        t = tuple(self.to_dict().values())
        return hash(t)
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

def get_teacher_list(session:requests.Session, _print:bool=False, limit:int=1000)->List[Teacher]:
    url = f"{base_url}/system/resource/tsites/advancesearch.jsp"
    page_size = 20
    teacher_list = []
    # ATTENTION: Recantly(2024.8.14), a request with a large limit (like 1000) will sometimes be rejected
    # So have to throttle the requests by paging now (though quite slower)
    end = limit // page_size + 1
    page_index = 1
    while True:
        print(f"第 {page_index} 页")
        query_param = {
            "collegeid": 0,
            "disciplineid": 0,
            "enrollid": 0,
            "pageindex": page_index,
            "pagesize": page_size,
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
        if resp.status_code != 200:
            print(f"请求失败，状态码: {resp.status_code}")
            pprint.pprint(resp.text)
            break
        print("请求成功")
        if page_index == 1:
            end = min(end, resp.json()['totalpage'] + 1)
        data = resp.json()['teacherData']
        if not data or len(data) == 0:
            break
        if _print:
            pprint.pprint(data)
        for d in data:
            teacher_list.append(resp_to_teacher(d))
        page_index += 1
        if page_index >= end:
            break
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

if __name__=="__main__":
    description = "A simple scrapy to download the teacher profile and save as json."
    if __package__ is None:
        import sys
        from os import path
        sys.path.append(path.dirname(path.dirname( path.abspath(__file__))))
        from common import Automator
    else:
        from .common import Automator
    automator = Automator(
        description=description
    )
    base_url = "https://faculty.sjtu.edu.cn/"
    parser = automator.parser
    parser.add_argument("-l", "--load_from", type=str, action="store", help="Load the result from the <load_from> json file.")
    parser.add_argument("-a", "--append_from", type=str, action="store", help="Append the result from the <append_from> json file.")
    parser.add_argument("-s", "--save_to", type=str, action="store", help="Save the result to the <save_to> json file.")
    parser.add_argument("--limit", type=int, help="The maximum number of request teachers")
    args = parser.parse_args()
    init_teachers = []
    if not args.name or not args.password:
        print("jacount and password are required")
        exit(0)
    automator.login(args.name, args.password)
    session = automator.get_session()
    if args.load_from and os.path.exists(args.load_from):
        with open(args.load_from, "r") as f:
            init_teachers = json.load(f)
    if args.append_from and os.path.exists(args.append_from):
        with open(args.append_from, "r") as f:
            init_teachers.extend(json.load(f))
    print("开始查询，请稍等.......")
    if args.limit:
        new_teachers = get_teacher_list(session=session,_print=args.verbose, limit=args.limit)
    else:
        new_teachers = get_teacher_list(session=session,_print=args.verbose)
    all_teachers = set(init_teachers).union(set(new_teachers))
    # data from jwc MAY be duplicated (e.g. 唐晟媚) for unknown reason, need to remove them
    all_teachers = list(all_teachers)
    print("当前教师数量：{}".format(len(all_teachers)))

    if not args.debug:
        if args.save_to and os.path.exists(args.save_to):
            save_json(all_teachers, to=args.save_to)
        else:
            save_json(all_teachers)
    automator.end()

