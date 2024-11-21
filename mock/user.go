package mock

import (
	"jcourse_go/model/po"
	"jcourse_go/util"
	"strconv"
	"time"

	"syreclabs.com/go/faker"
)

// MockMajors fake的专业列表
var MockMajors = []string{
	"计算机科学与技术",
	"软件工程",
	"信息安全",
	"网络工程",
	"物联网工程",
	"数据科学与大数据技术",
	"人工智能",
	"数字媒体技术",
	"数字艺术设计",
	"机械工程",
	"电子信息工程",
	"通信工程",
	"自动化",
	"电气工程及其自动化",
	"电子信息科学与技术",
	"电子科学与技术",
	"微电子科学与工程",
	"法学",
	"社会学",
	"心理学",
	"交通运输",
	"工业工程",
	"工商管理",
	"市场营销",
	"会计学",
	"财务管理",
	"金融学",
	"国际经济与贸易",
	"经济学",
	"船舶与海洋工程",
	"航海技术",
	"土木工程",
	"建筑学",
	"城乡规划",
	"环境工程",
	"环境科学",
	"风景园林",
	"食品科学与工程",
	"马克思主义理论",
	"中国语言文学",
	"新闻传播学",
	"历史学",
	"考古学",
	"暂无",
}

var timeSeedStr = "2024-11-22 11:22:33"
var TimeSeed, _ = util.ParseTime(timeSeedStr)

func MockUsers(gen MockDBGenerator, n int) ([]po.UserPO, error) {
	users := make([]po.UserPO, n)
	emails, err := GenerateUniqueSet(n, func() string {
		return faker.Internet().Email()
	})
	if err != nil {
		return nil, err
	}
	userNames, err := GenerateUniqueSet(n, func() string {
		return faker.Internet().UserName()
	})
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		users[i] = po.UserPO{
			Username:   userNames[i],
			Password:   faker.Internet().Password(8, 14),
			Email:      emails[i],
			Avatar:     faker.Avatar().String(),
			Department: GenDepartment(gen),
			Type:       "学生",
			Major:      MockMajors[gen.Rand.Intn(len(MockMajors))],
			Degree:     "本科",
			Grade:      strconv.Itoa(gen.Rand.Intn(5) + 2020),
			Bio:        faker.Lorem().Sentence(10),
			Points:     int64(gen.Rand.Intn(1000)),
			LastSeenAt: TimeSeed.Add(-time.Duration(gen.Rand.Intn(1000)) * time.Hour),
		}
	}
	err = gen.db.CreateInBatches(users, gen.batchSize).Error
	return users, err
}
