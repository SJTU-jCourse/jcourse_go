import requests
from typing import Callable
from selenium.webdriver.common.by import By
from argparse import ArgumentParser
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
from selenium import webdriver
from selenium.webdriver.chrome.service import Service
class Automator:
    timeout = 60
    def __init__(self, description:str="") -> None:
        # automate downloading chrome driver if
        # 1.you haven't installed chrome
        # 2.imcompatible with your browser
        service = Service()
        options = webdriver.ChromeOptions()
        self.driver = webdriver.Chrome(service=service, options=options)
        self.parser = ArgumentParser(description=description)
        self.parser.add_argument("-n", "--name", type=str, action="store", help="jacount for login")
        self.parser.add_argument("-p", "--password", type=str, action="store", help="password for login")
        self.parser.add_argument("-d","--debug", action="store_true", help="In debug mode, the program will not save the result.")
        self.parser.add_argument("-v", "--verbose", action="store_true", help="Verbose mode")
        self.has_login = False
    def login(self, name:str, password:str):
        login_url = "https://i.sjtu.edu.cn/jaccountlogin"
        self.driver.implicitly_wait(time_to_wait=self.timeout)
        self.driver.get(url=login_url)
        username_input = self.driver.find_element(By.ID, "input-login-user")
        password_input = self.driver.find_element(By.ID, "input-login-pass")
        username_input.send_keys("{}".format(name))
        password_input.send_keys("{}".format(password))
        # driver -> bool
        cond:Callable[[object], bool] = EC.url_contains("https://i.sjtu.edu.cn/xtgl")
        WebDriverWait(self.driver, self.timeout).until(cond)
        session = requests.Session()
        # 获取登录后的 Cookie
        cookies = self.driver.get_cookies()

        # 将 Cookie 转换为 requests 库可以使用的格式
        for cookie in cookies:
            session.cookies.set(cookie['name'], cookie['value'])
        self.session = session
        self.has_login = True
    def gen_driver(self):
        return self.driver

    def end(self):
        self.driver.close()

    def get_session(self):
        if not self.login:
            raise Exception("Please login first")
        return self.session
