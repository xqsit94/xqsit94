package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Experience struct {
	Company   string    `json:"company"`
	Role      string    `json:"role"`
	StartDate time.Time `json:"-"`
	EndDate   time.Time `json:"-"`
	Location  string    `json:"location"`
	IsCurrent bool      `json:"isCurrent"`
	StartDateStr string `json:"startDate"`
	EndDateStr   string `json:"endDate,omitempty"`
}

type Education struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	StartYear   int    `json:"startYear"`
	EndYear     int    `json:"endYear"`
}

type experienceData struct {
	Experience []Experience `json:"experience"`
}

type educationData struct {
	Education []Education `json:"education"`
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

func loadExperience() ([]Experience, error) {
	data, err := os.ReadFile("data/experience.json")
	if err != nil {
		return nil, fmt.Errorf("error reading experience data: %w", err)
	}

	var expData experienceData
	if err := json.Unmarshal(data, &expData); err != nil {
		return nil, fmt.Errorf("error unmarshaling experience data: %w", err)
	}

	// Convert date strings to time.Time
	for i := range expData.Experience {
		expData.Experience[i].StartDate = parseDate(expData.Experience[i].StartDateStr)
		if expData.Experience[i].EndDateStr != "" {
			expData.Experience[i].EndDate = parseDate(expData.Experience[i].EndDateStr)
		}
	}

	return expData.Experience, nil
}

func loadEducation() ([]Education, error) {
	data, err := os.ReadFile("data/education.json")
	if err != nil {
		return nil, fmt.Errorf("error reading education data: %w", err)
	}

	var eduData educationData
	if err := json.Unmarshal(data, &eduData); err != nil {
		return nil, fmt.Errorf("error unmarshaling education data: %w", err)
	}

	return eduData.Education, nil
}

func GetExperience() []Experience {
	experience, err := loadExperience()
	if err != nil {
		return []Experience{}
	}
	return experience
}

func GetEducation() []Education {
	education, err := loadEducation()
	if err != nil {
		return []Education{}
	}
	return education
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
