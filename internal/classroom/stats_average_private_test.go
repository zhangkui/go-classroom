package classroom
import("context";"testing";"time")
func TestAverageScoreUsesGradedSubmissions(t *testing.T){s:=NewStore();c,_:=s.CreateCourse("Math","teacher");_ = s.JoinCourse(c.ID,"one");_ = s.JoinCourse(c.ID,"two");a,_:=s.CreateAssignment(c.ID,"one",time.Now().Add(time.Hour));sub,_:=s.Submit(context.Background(),a.ID,"one","answer");_ = s.Grade(sub.ID,80,"ok");stats,err:=s.Stats(context.Background(),c.ID);if err!=nil{t.Fatal(err)};if stats.AverageScore!=80{t.Fatalf("expected average 80, got %v",stats.AverageScore)}}
