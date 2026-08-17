package classroom

import("context";"testing";"time")
func TestBasicCourseWorkflow(t *testing.T){s:=NewStore();c,err:=s.CreateCourse("Math","teacher-1");if err!=nil{t.Fatal(err)};if err=s.JoinCourse(c.ID,"student-1");err!=nil{t.Fatal(err)};a,err:=s.CreateAssignment(c.ID,"Quiz",time.Now().Add(time.Hour));if err!=nil{t.Fatal(err)};sub,err:=s.Submit(context.Background(),a.ID,"student-1","answer");if err!=nil{t.Fatal(err)};if err=s.Grade(sub.ID,90,"good");err!=nil{t.Fatal(err)}}
