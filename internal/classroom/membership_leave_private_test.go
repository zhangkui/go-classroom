package classroom
import "testing"
func TestLeavingCourseRemovesOnlyRequestedStudent(t *testing.T){s:=NewStore();c,_:=s.CreateCourse("Math","teacher");_ = s.JoinCourse(c.ID,"alice");_ = s.JoinCourse(c.ID,"bob");if err:=s.LeaveCourse(c.ID,"alice");err!=nil{t.Fatal(err)};m,_:=s.ListMembers(c.ID);if _,ok:=m["alice"];ok{t.Fatal("alice is still a member")};if _,ok:=m["bob"];!ok{t.Fatal("bob was removed")}}
