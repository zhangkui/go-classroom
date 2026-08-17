package classroom
import("testing";"errors")
func TestCourseNamesAreUniqueAfterNormalization(t *testing.T){s:=NewStore();if _,err:=s.CreateCourse(" Algebra ","t1");err!=nil{t.Fatal(err)};if _,err:=s.CreateCourse("algebra","t2");!errors.Is(err,ErrDuplicate){t.Fatalf("expected duplicate, got %v",err)}}
