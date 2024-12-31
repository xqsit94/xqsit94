package profile

import (
	"fmt"
	"time"
)

type Experience struct {
	Company   string
	Role      string
	StartDate time.Time
	EndDate   time.Time
	Location  string
	IsCurrent bool
}

type Education struct {
	Institution string
	Degree      string
	StartYear   int
	EndYear     int
}

func (e Experience) Duration() string {
	end := e.EndDate
	if e.IsCurrent {
		end = time.Now()
	}

	years := end.Year() - e.StartDate.Year()
	months := int(end.Month() - e.StartDate.Month())

	if months < 0 {
		years--
		months += 12
	}

	if years == 0 {
		return fmt.Sprintf("%d m", months)
	}
	if months == 0 {
		return fmt.Sprintf("%d yr", years)
	}
	return fmt.Sprintf("%d yr, %d m", years, months)
}

func (e Experience) DateRange() string {
	startDate := e.StartDate.Format("Jan 2006")
	if e.IsCurrent {
		return fmt.Sprintf("%s - Present", startDate)
	}
	return fmt.Sprintf("%s - %s", startDate, e.EndDate.Format("Jan 2006"))
}

func GetExperience() []Experience {
	return []Experience{
		{
			Company:   "Oreala B.V",
			Role:      "Full Stack Engineer",
			StartDate: parseDate("2022-04"),
			IsCurrent: true,
			Location:  "India",
		},
		{
			Company:   "Colan Infotech Private Limited",
			Role:      "Software Engineer",
			StartDate: parseDate("2018-07"),
			EndDate:   parseDate("2022-03"),
			Location:  "Chennai, Tamil Nadu, India",
		},
		{
			Company:   "Expose InfoTech India Pvt Ltd",
			Role:      "PHP Developer",
			StartDate: parseDate("2017-12"),
			EndDate:   parseDate("2018-06"),
			Location:  "Calicut Area, India",
		},
		{
			Company:   "Slogics Solutions",
			Role:      "Web Developer",
			StartDate: parseDate("2016-11"),
			EndDate:   parseDate("2017-12"),
			Location:  "Chennai Area, India",
		},
	}
}

func GetEducation() []Education {
	return []Education{
		{
			Institution: "Madha Engineering College",
			Degree:      "Bachelor of Engineering (B.E.), Computer Science",
			StartYear:   2012,
			EndYear:     2016,
		},
		{
			Institution: "Assisi Matriculation School - India",
			Degree:      "Primary and Secondary Examinations, General Studies",
			StartYear:   1997,
			EndYear:     2012,
		},
	}
}

func parseDate(date string) time.Time {
	t, _ := time.Parse("2006-01", date)
	return t
}

func CalculateTotalExperience() string {
	var totalMonths int
	now := time.Now()

	for _, exp := range GetExperience() {
		end := exp.EndDate
		if exp.IsCurrent {
			end = now
		}

		months := (end.Year()-exp.StartDate.Year())*12 + int(end.Month()-exp.StartDate.Month())
		totalMonths += months
	}

	years := totalMonths / 12
	return fmt.Sprintf("%d years", years)
}
