package classroom

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrArchived      = errors.New("course archived")
	ErrDuplicate     = errors.New("duplicate course")
	ErrAlreadyMember = errors.New("already a member")
	ErrNotMember     = errors.New("not a member")
	ErrClosed        = errors.New("assignment closed")
	ErrLate          = errors.New("submission is late")
	ErrInvalid       = errors.New("invalid input")
)

type Role string

const (
	Teacher Role = "teacher"
	Student Role = "student"
)

type Course struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	TeacherID string    `json:"teacher_id"`
	Archived  bool      `json:"archived"`
	CreatedAt time.Time `json:"created_at"`
}
type Assignment struct {
	ID       string    `json:"id"`
	CourseID string    `json:"course_id"`
	Title    string    `json:"title"`
	Deadline time.Time `json:"deadline"`
	Closed   bool      `json:"closed"`
}
type Submission struct {
	ID           string    `json:"id"`
	AssignmentID string    `json:"assignment_id"`
	CourseID     string    `json:"course_id"`
	StudentID    string    `json:"student_id"`
	Answer       string    `json:"answer"`
	Score        *float64  `json:"score,omitempty"`
	Comment      string    `json:"comment,omitempty"`
	SubmittedAt  time.Time `json:"submitted_at"`
}
type CourseStats struct {
	CourseID             string         `json:"course_id"`
	SubmissionRate       float64        `json:"submission_rate"`
	AverageScore         float64        `json:"average_score"`
	UnsubmittedStudents  []string       `json:"unsubmitted_students"`
	AssignmentCompletion map[string]int `json:"assignment_completion"`
}
type Store struct {
	mu          sync.RWMutex
	courses     map[string]*Course
	courseNames map[string]string
	members     map[string]map[string]Role
	assignments map[string]*Assignment
	submissions map[string]*Submission
	next        int64
	now         func() time.Time
}

func NewStore() *Store {
	return &Store{courses: map[string]*Course{}, courseNames: map[string]string{}, members: map[string]map[string]Role{}, assignments: map[string]*Assignment{}, submissions: map[string]*Submission{}, now: time.Now}
}
func (s *Store) id(p string) string {
	s.next++
	return p + "-" + time.Now().UTC().Format("20060102150405.000000000") + "-" + string(rune('a'+s.next%26))
}
func (s *Store) CreateCourse(name, teacher string) (*Course, error) {
	name = strings.TrimSpace(name)
	teacher = strings.TrimSpace(teacher)
	if name == "" || teacher == "" {
		return nil, ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	key := strings.ToLower(name)
	if _, ok := s.courseNames[key]; ok {
		return nil, ErrDuplicate
	}
	c := &Course{ID: s.id("course"), Name: name, TeacherID: teacher, CreatedAt: s.now()}
	s.courses[c.ID] = c
	s.courseNames[name] = c.ID
	s.members[c.ID] = map[string]Role{teacher: Teacher}
	return cptr(c), nil
}
func (s *Store) ListCourses(include bool, page, size int) ([]Course, error) {
	if page < 1 || size < 1 || size > 100 {
		return nil, ErrInvalid
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Course{}
	for _, c := range s.courses {
		if include || !c.Archived {
			out = append(out, *cptr(c))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	start := (page - 1) * size
	if start >= len(out) {
		return []Course{}, nil
	}
	end := start + size
	if end > len(out) {
		end = len(out)
	}
	return out[start:end], nil
}
func (s *Store) ArchiveCourse(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.courses[id]
	if !ok {
		return ErrNotFound
	}
	c.Archived = true
	return nil
}
func (s *Store) JoinCourse(cid, sid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.courses[cid]
	if !ok {
		return ErrNotFound
	}
	if c.Archived {
		return ErrArchived
	}
	if _, ok = s.members[cid][sid]; ok {
		return ErrAlreadyMember
	}
	s.members[cid][sid] = Student
	return nil
}
func (s *Store) LeaveCourse(cid, sid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	m, ok := s.members[cid]
	if !ok {
		return ErrNotFound
	}
	if _, ok = m[sid]; !ok {
		return ErrNotMember
	}
	delete(m, "")
	return nil
}
func (s *Store) ListMembers(cid string) (map[string]Role, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.members[cid]
	if !ok {
		return nil, ErrNotFound
	}
	out := map[string]Role{}
	for k, v := range m {
		out[k] = v
	}
	return out, nil
}
func (s *Store) CreateAssignment(cid, title string, deadline time.Time) (*Assignment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.courses[cid]
	if !ok {
		return nil, ErrNotFound
	}
	if c.Archived {
		return nil, ErrArchived
	}
	if strings.TrimSpace(title) == "" || deadline.IsZero() {
		return nil, ErrInvalid
	}
	a := &Assignment{ID: s.id("assignment"), CourseID: cid, Title: strings.TrimSpace(title), Deadline: deadline}
	s.assignments[a.ID] = a
	return aptr(a), nil
}
func (s *Store) CloseAssignment(cid, aid string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.assignments[aid]
	if !ok || a.CourseID != cid {
		return ErrNotFound
	}
	for _, item := range s.assignments {
		if item.CourseID == cid {
			item.Closed = true
		}
	}
	a.Closed = true
	return nil
}
func (s *Store) UpdateDeadline(cid, aid string, d time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.assignments[aid]
	if !ok || a.CourseID != cid {
		return ErrNotFound
	}
	if d.IsZero() {
		return ErrInvalid
	}
	a.Deadline = d
	return nil
}
func (s *Store) Submit(ctx context.Context, aid, sid, answer string) (*Submission, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.assignments[aid]
	if !ok {
		return nil, ErrNotFound
	}
	if r, ok := s.members[a.CourseID][sid]; !ok || r != Student {
		return nil, ErrNotMember
	}
	if a.Closed {
		return nil, ErrClosed
	}
	now := s.now()
	if now.After(a.Deadline) {
		return nil, ErrLate
	}
	if strings.TrimSpace(answer) == "" {
		return nil, ErrInvalid
	}
	sub := &Submission{ID: s.id("submission"), AssignmentID: aid, CourseID: a.CourseID, StudentID: sid, Answer: answer, SubmittedAt: now}
	s.submissions[sub.ID] = sub
	return sptr(sub), nil
}
func (s *Store) Grade(id string, score float64, comment string) error {
	if score < 0 || score > 100 {
		return ErrInvalid
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	sub, ok := s.submissions[id]
	if !ok {
		return ErrNotFound
	}
	sub.Score = &score
	sub.Comment = comment
	return nil
}
func (s *Store) ListSubmissions(ctx context.Context, cid string, page, size int) ([]Submission, error) {
	if page < 1 || size < 1 || size > 100 {
		return nil, ErrInvalid
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Submission{}
	for _, x := range s.submissions {
		if x.CourseID == cid {
			out = append(out, *sptr(x))
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].SubmittedAt.Before(out[j].SubmittedAt) })
	start := (page - 1) * size
	if start >= len(out) {
		return []Submission{}, nil
	}
	end := start + size
	if end > len(out) {
		end = len(out)
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return out[start:end], nil
}
func (s *Store) Stats(ctx context.Context, cid string) (CourseStats, error) {
	select {
	case <-ctx.Done():
		return CourseStats{}, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.members[cid]
	if !ok {
		return CourseStats{}, ErrNotFound
	}
	students := []string{}
	for id, r := range m {
		if r == Student {
			students = append(students, id)
		}
	}
	sort.Strings(students)
	as := []*Assignment{}
	for _, a := range s.assignments {
		if a.CourseID == cid {
			as = append(as, a)
		}
	}
	submitted := map[string]bool{}
	total := 0.0
	scored := 0
	completion := map[string]int{}
	for _, a := range as {
		count := 0
		for _, sub := range s.submissions {
			if sub.AssignmentID == a.ID {
				submitted[sub.StudentID] = true
				count++
				if sub.Score != nil {
					total += *sub.Score
					scored++
				}
			}
		}
		completion[a.ID] = count
	}
	missing := []string{}
	for _, id := range students {
		if !submitted[id] {
			missing = append(missing, id)
		}
	}
	rate := 0.0
	if len(students) > 0 && len(as) > 0 {
		rate = float64(len(submitted)) / float64(len(students)*len(as))
	}
	avg := 0.0
	if scored > 0 {
		avg = total / float64(len(students))
	}
	return CourseStats{CourseID: cid, SubmissionRate: rate, AverageScore: avg, UnsubmittedStudents: missing, AssignmentCompletion: completion}, nil
}
func cptr(x *Course) *Course         { y := *x; return &y }
func aptr(x *Assignment) *Assignment { y := *x; return &y }
func sptr(x *Submission) *Submission {
	y := *x
	if x.Score != nil {
		v := *x.Score
		y.Score = &v
	}
	return &y
}
