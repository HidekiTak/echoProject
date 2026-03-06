package asyncjobs

import (
	"context"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	//    config = asyncjobs.db.mysql.config(config, host = "localhost:13306", schema = "unittest_unittest", user = "root", password = "root")
	ctx := context.Background()
	controller, err := NewJobController[string](ctx, "unit_test", "mysql", "root:root@tcp(localhost:13306)/unittest_unittest")
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, controller)
	_, err = controller.Kick(ctx, "tenant", "job_type", "job_id", "meta", "example@example.com", 12)
	if err != nil {
		t.Error(err)
	}
}
