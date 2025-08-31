package course

type FilterAggregateItem struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

type CourseFilter struct {
	Departments []FilterAggregateItem `json:"departments"`
	Credits     []FilterAggregateItem `json:"credits"`
	Semesters   []FilterAggregateItem `json:"semesters"`
	Categories  []FilterAggregateItem `json:"categories"`
}

type TeacherFilter struct {
	Departments []FilterAggregateItem `json:"departments"`
	Titles      []FilterAggregateItem `json:"titles"`
}

type TrainingPlanFilter struct {
	Departments []FilterAggregateItem `json:"departments"`
	Degrees     []FilterAggregateItem `json:"degrees"`
	EntryYears  []FilterAggregateItem `json:"entry_years"`
}

type Pinyin struct {
	Full string `json:"full"`
	Abbr string `json:"abbr"`
}
