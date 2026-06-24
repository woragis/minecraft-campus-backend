package models

const (
	AffiliationStudent = "student"
	AffiliationStaff   = "staff"
	AffiliationGuest   = "guest"
	AffiliationAlumni  = "alumni"
)

func ValidAffiliationType(t string) bool {
	switch t {
	case AffiliationStudent, AffiliationStaff, AffiliationGuest, AffiliationAlumni:
		return true
	default:
		return false
	}
}

func (p *Player) IsGuest() bool {
	return p != nil && p.AffiliationType == AffiliationGuest
}

type University struct {
	Slug     string `gorm:"primaryKey" json:"slug"`
	Name     string `gorm:"not null" json:"name"`
	ColorHex string `gorm:"not null;column:color_hex" json:"colorHex"`
}

func (University) TableName() string { return "universities" }

type Faculty struct {
	Slug           string `gorm:"primaryKey" json:"slug"`
	UniversitySlug string `gorm:"not null;column:university_slug" json:"universitySlug"`
	Name           string `gorm:"not null" json:"name"`
	ShortAbbr      string `gorm:"not null;column:short_abbr" json:"shortAbbr"`
	ColorHex       string `gorm:"not null;column:color_hex" json:"colorHex"`
}

func (Faculty) TableName() string { return "faculties" }

type Course struct {
	Slug      string `gorm:"primaryKey" json:"slug"`
	FacultySlug string `gorm:"not null;column:faculty_slug" json:"facultySlug"`
	Name      string `gorm:"not null" json:"name"`
	ShortAbbr string `gorm:"not null;column:short_abbr" json:"shortAbbr"`
	ColorHex  string `gorm:"not null;column:color_hex" json:"colorHex"`
}

func (Course) TableName() string { return "courses" }
