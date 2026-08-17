package main

import (
	"encoding/json"
	"errors"
	"github.com/zhangkui/go-classroom/internal/classroom"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	s := classroom.NewStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var in struct {
				Name      string `json:"name"`
				TeacherID string `json:"teacher_id"`
			}
			if json.NewDecoder(r.Body).Decode(&in) != nil {
				http.Error(w, "invalid json", 400)
				return
			}
			c, e := s.CreateCourse(in.Name, in.TeacherID)
			if e != nil {
				http.Error(w, e.Error(), status(e))
				return
			}
			write(w, 201, c)
			return
		}
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		size, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		if size == 0 {
			size = 20
		}
		list, e := s.ListCourses(false, page, size)
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		write(w, 200, list)
	})
	mux.HandleFunc("/courses/", func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 2 {
			http.NotFound(w, r)
			return
		}
		cid := parts[1]
		if len(parts) == 2 && r.Method == http.MethodDelete {
			if e := s.ArchiveCourse(cid); e != nil {
				http.Error(w, e.Error(), status(e))
				return
			}
			w.WriteHeader(204)
			return
		}
		if len(parts) == 3 && parts[2] == "members" {
			if r.Method == http.MethodPost {
				var in struct {
					StudentID string `json:"student_id"`
				}
				json.NewDecoder(r.Body).Decode(&in)
				if e := s.JoinCourse(cid, in.StudentID); e != nil {
					http.Error(w, e.Error(), status(e))
					return
				}
				w.WriteHeader(204)
				return
			}
			members, e := s.ListMembers(cid)
			if e != nil {
				http.Error(w, e.Error(), status(e))
				return
			}
			write(w, 200, members)
			return
		}
		http.NotFound(w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
func status(e error) int {
	if errors.Is(e, classroom.ErrNotFound) {
		return 404
	}
	if errors.Is(e, classroom.ErrDuplicate) {
		return 409
	}
	if errors.Is(e, classroom.ErrArchived) {
		return 410
	}
	return 400
}

var _ = time.Now
