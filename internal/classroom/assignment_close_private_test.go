package classroom
import("testing";"time")
func TestClosingOneAssignmentDoesNotCloseAnother(t *testing.T){s:=NewStore();c,_:=s.CreateCourse("Math","teacher");a1,_:=s.CreateAssignment(c.ID,"one",time.Now().Add(time.Hour));a2,_:=s.CreateAssignment(c.ID,"two",time.Now().Add(time.Hour));if err:=s.CloseAssignment(c.ID,a1.ID);err!=nil{t.Fatal(err)};s.mu.RLock();closed2:=s.assignments[a2.ID].Closed;s.mu.RUnlock();if closed2{t.Fatal("unrelated assignment was closed")}}
