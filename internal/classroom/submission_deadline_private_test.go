package classroom
import("context";"testing";"time")
func TestSubmissionAtDeadlineIsRejected(t *testing.T){s:=NewStore();now:=time.Now().UTC();s.now=func()time.Time{return now};c,_:=s.CreateCourse("Math","teacher");_ = s.JoinCourse(c.ID,"student");a,_:=s.CreateAssignment(c.ID,"one",now);if _,err:=s.Submit(context.Background(),a.ID,"student","answer");err!=ErrLate{t.Fatalf("expected late error, got %v",err)}}
