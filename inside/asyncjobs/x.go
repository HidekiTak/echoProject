package asyncjobs

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type JobParameter[T any] struct {
}

type JobInfo struct {
	kickAt  int64  `db:"kick_at"`
	subAt   int32  `db:"sub_at"`
	tenant  string `db:"tenant"`
	jobType string `db:"job_type"`
	jobId   string `db:"job_id"`
	jobHash string `db:"job_hash"`
}

func (ji JobInfo) KickAt() int64 {
	return ji.kickAt
}

func (ji JobInfo) SubAt() int32 {
	return ji.subAt
}

func (ji *JobInfo) Tenant() string {
	return ji.tenant
}

func (ji *JobInfo) JobType() string {
	return ji.jobType
}

func (ji *JobInfo) JobId() string {
	return ji.jobId
}

func (ji *JobInfo) JobHash() string {
	return ji.jobHash
}

type JobController[T any] interface {
	Kick(
		ctx context.Context,
		tenant string,
		jobType string,
		jobId string,
		meta T,
		kickedBy string,
		kickedAt int64,
	) (*JobInfo, error)
}

func NewJobController[T any](ctx context.Context, jobGroup string, driver string, datasourceName string) (JobController[T], error) {
	db, _ := sqlx.Open(driver, datasourceName)
	defer db.Close()
	exists, err := checkTable(ctx, db, jobGroup)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = createTable(ctx, db, jobGroup)
		if err != nil {
			return nil, err
		}
	}
	return &jobController[T]{driver: driver, datasourceName: datasourceName, tableBase: jobGroup}, nil
}

const sql_mysql string = `CREATE TABLE ${JOB_GROUP}_async_jobs(
kicked_at bigint                                                NOT NULL,
sub_at    int                                                   NOT NULL,
tenant    varchar(255) CHARACTER SET ascii COLLATE ascii_bin    NOT NULL,
job_type  varchar(255) CHARACTER SET ascii COLLATE ascii_bin    NOT NULL,
job_id    varchar(255) CHARACTER SET ascii COLLATE ascii_bin    NOT NULL,
job_hash  varchar(255) CHARACTER SET ascii COLLATE ascii_bin    NOT NULL,
retry     smallint                                              NOT NULL DEFAULT '0',
meta      text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
status    varchar(255) CHARACTER SET ascii COLLATE ascii_bin    NOT NULL,
result    text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
kicked_by varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  DEFAULT NULL,
exec_at   bigint                                                         DEFAULT NULL,
PRIMARY KEY (kicked_at, sub_at, retry),
KEY ix_${JOB_GROUP}_async_jobs_tenant (tenant ASC, job_type ASC, job_id ASC)
)`

func createTable(ctx context.Context, db *sqlx.DB, jobGroup string) error {
	createSQL := strings.Replace(sql_mysql, "${JOB_GROUP}", jobGroup, -1)
	_, err := db.ExecContext(ctx, createSQL)
	return err
}

func checkTable(ctx context.Context, db *sqlx.DB, jobGroup string) (bool, error) {
	records, err := db.QueryContext(ctx, "SHOW TABLES LIKE '"+jobGroup+"_async_jobs'")
	defer records.Close()
	if err != nil {
		return false, err
	}
	err = records.Err()
	if err != nil {
		return false, err
	}
	if !records.Next() {
		return false, nil
	}
	var tableName string
	err = records.Scan(&tableName)
	if err != nil {
		return false, err
	}
	return true, nil
}

type jobController[T any] struct {
	driver         string
	datasourceName string
	tableBase      string
}

const sqlSelect string = `
SELECT kick_at,sub_at,tenant,job_type,job_id,job_hash
FROM %s_async_jobs
WHERE tenant = ? AND job_type = ? AND job_id = ? AND job_hash=?
ORDER BY exec_at DESC limit 1`

func (jc *jobController[T]) Kick(
	ctx context.Context,
	tenant string,
	jobType string,
	jobId string,
	meta T,
	kickedBy string,
	kickedAt int64,
) (*JobInfo, error) {
	metaJson, _ := json.Marshal(meta)
	jobHash := string(md5.New().Sum(metaJson))

	con, err := sqlx.Open(jc.driver, jc.datasourceName)
	if err != nil {
		return nil, err
	}
	defer con.Close()
	tx, err := con.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			println("_ = tx.Commit()")
			_ = tx.Commit()
		} else {
			println("_ = tx.Rollback()")
			_ = tx.Rollback()
		}
	}()
	jobInfo := JobInfo{}
	err = con.GetContext(
		ctx, &jobInfo,
		fmt.Sprintf(sqlSelect, jc.tableBase),
		tenant, jobType, jobId, jobHash)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, err
		}
	} else {
		return &jobInfo, nil
	}

	var subAt int32
	subAt = 0
	sql_insert := fmt.Sprintf(
		"INSERT IGNORE INTO `%s_async_jobs`(`kicked_at`,`sub_at`,`tenant`,`job_type`,`job_id`,`job_hash`,`meta`,`status`,`kicked_by`)VALUES(?,?,?,?,?,?,?,?,?)",
		jc.tableBase,
	)
X:
	result, err := con.ExecContext(ctx,
		sql_insert,
		kickedAt, subAt,
		tenant, jobType, jobId, jobHash,
		metaJson, "prepare", kickedBy)
	if err != nil {
		return nil, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if 0 < count {
		return &JobInfo{
			kickAt:  kickedAt,
			subAt:   subAt,
			tenant:  tenant,
			jobType: jobType,
			jobId:   jobId,
			jobHash: jobHash,
		}, nil
	}
	subAt = subAt + 1
	if 100 < subAt {
		return nil, errors.New("internal server error")
	}
	goto X
}
