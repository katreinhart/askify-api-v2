package model

// FetchCohortList returns an object of all cohorts in DB
func FetchCohortList() []Cohort {
	var cohorts []Cohort

	db.Find(&cohorts)

	return cohorts
}

// AddCohort adds a new cohort to the DB
func AddCohort(c Cohort, uid string) (Cohort, error) {
	var u UserModel
	db.First(&u, "id = ?", uid)

	if u.ID == 0 || u.Admin == false {
		return Cohort{}, ErrorForbidden
	}

	db.Save(&c)

	return c, nil
}
