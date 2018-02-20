package model

// FetchCohortList returns an object of all cohorts in DB
func FetchCohortList() []Cohort {
	var cohorts []Cohort

	db.Find(&cohorts)

	return cohorts
}

// AddCohort adds a new cohort to the DB
func AddCohort(c Cohort) (Cohort, error) {
	db.Save(&c)

	return c, nil
}
