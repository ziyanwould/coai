package admin

import (
	"chat/channel"
	"chat/globals"
	"chat/utils"
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

// AuthLike is to solve the problem of import cycle.
type AuthLike struct {
	ID int64 `json:"id"`
}

func (a *AuthLike) GetID(_ *sql.DB) int64 {
	return a.ID
}

func (a *AuthLike) HitID() int64 {
	return a.ID
}

// Modify getUsersForm to accept isSubscribedFilter and isBannedFilter
func getUsersForm(db *sql.DB, page int64, search string, isSubscribedFilter *bool, isBannedFilter *bool, sortKey string) PaginationForm {
	var users []interface{}
	var total int64

	whereClause := "WHERE 1=1"
	searchQuery := "%" + search + "%"

	if search != "" {
		whereClause += " AND (auth.username LIKE ? OR auth.email LIKE ?)"
	}

	if isSubscribedFilter != nil {
		if *isSubscribedFilter {
			whereClause += " AND subscription.expired_at IS NOT NULL AND subscription.expired_at > NOW()"
		} else {
			whereClause += " AND (subscription.expired_at IS NULL OR subscription.expired_at <= NOW())"
		}
	}

	if isBannedFilter != nil {
		if *isBannedFilter {
			whereClause += " AND auth.is_banned = TRUE"
		} else {
			whereClause += " AND auth.is_banned = FALSE"
		}
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM auth LEFT JOIN subscription ON subscription.user_id = auth.id LEFT JOIN quota ON quota.user_id = auth.id %s", whereClause)
	queryArgsCount := []interface{}{}
	if search != "" {
		queryArgsCount = append(queryArgsCount, searchQuery, searchQuery)
	}
	if err := globals.QueryRowDb(db, countQuery, queryArgsCount...).Scan(&total); err != nil {
		return PaginationForm{
			Status:  false,
			Message: err.Error(),
		}
	}

	// 构建排序子句
	orderByClause := "ORDER BY auth.id" // 默认按 ID 升序
	switch sortKey {
	case "quota":
		orderByClause = "ORDER BY quota.quota DESC"
	case "used_quota":
		orderByClause = "ORDER BY quota.used DESC"
	case "total_month":
		orderByClause = "ORDER BY subscription.total_month DESC"
	case "expired_at":
		orderByClause = "ORDER BY subscription.expired_at DESC"
	}

	selectQuery := fmt.Sprintf(`
		SELECT
		    auth.id, auth.username, auth.email, auth.is_admin,
		    quota.quota, quota.used,
		    subscription.expired_at, subscription.total_month, subscription.enterprise, subscription.level,
		    auth.is_banned
		FROM auth
		LEFT JOIN quota ON quota.user_id = auth.id
		LEFT JOIN subscription ON subscription.user_id = auth.id
		%s
		%s
		LIMIT ? OFFSET ?
	`, whereClause, orderByClause) // 添加 orderByClause

	queryArgsSelect := []interface{}{}
	if search != "" {
		queryArgsSelect = append(queryArgsSelect, searchQuery, searchQuery)
	}
	queryArgsSelect = append(queryArgsSelect, pagination, page*pagination)

	rows, err := globals.QueryDb(db, selectQuery, queryArgsSelect...)
	if err != nil {
		return PaginationForm{
			Status:  false,
			Message: err.Error(),
		}
	}

	for rows.Next() {
		var user UserData
		var (
			email             sql.NullString
			expired           []uint8
			quota             sql.NullFloat64
			usedQuota         sql.NullFloat64
			totalMonth        sql.NullInt64
			isEnterprise      sql.NullBool
			subscriptionLevel sql.NullInt64
			isBanned          sql.NullBool
		)
		if err := rows.Scan(&user.Id, &user.Username, &email, &user.IsAdmin, &quota, &usedQuota, &expired, &totalMonth, &isEnterprise, &subscriptionLevel, &isBanned); err != nil {
			return PaginationForm{
				Status:  false,
				Message: err.Error(),
			}
		}
		if email.Valid {
			user.Email = email.String
		}
		if quota.Valid {
			user.Quota = float32(quota.Float64)
		}
		if usedQuota.Valid {
			user.UsedQuota = float32(usedQuota.Float64)
		}
		if totalMonth.Valid {
			user.TotalMonth = totalMonth.Int64
		}
		if subscriptionLevel.Valid {
			user.Level = int(subscriptionLevel.Int64)
		}
		stamp := utils.ConvertTime(expired)
		if stamp != nil {
			user.IsSubscribed = stamp.After(time.Now())
			user.ExpiredAt = stamp.Format("2006-01-02 15:04:05")
		}
		user.Enterprise = isEnterprise.Valid && isEnterprise.Bool
		user.IsBanned = isBanned.Valid && isBanned.Bool

		users = append(users, user)
	}

	return PaginationForm{
		Status: true,
		Total:  int(math.Ceil(float64(total) / float64(pagination))),
		Data:   users,
	}
}

// clearUserCache clears all cache keys starting with nio:user:
func clearUserCache(cache *redis.Client) error {
	ctx := context.Background()
	iter := cache.Scan(ctx, 0, "nio:user:*", 100).Iterator()
	for iter.Next(ctx) {
		if err := cache.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete cache key %s: %v", iter.Val(), err)
		}
	}
	return iter.Err()
}

func passwordMigration(db *sql.DB, cache *redis.Client, id int64, password string) error {
	password = strings.TrimSpace(password)
	if len(password) < 6 || len(password) > 36 {
		return fmt.Errorf("password length must be between 6 and 36")
	}
	hash_passwd := utils.Sha2Encrypt(password)

	// Update password in database
	_, err := globals.ExecDb(db, `
		UPDATE auth SET password = ? WHERE id = ?
	`, hash_passwd, id)

	if err != nil {
		return err
	}

	// Clear all user related cache
	if err := clearUserCache(cache); err != nil {
		return fmt.Errorf("failed to clear user cache: %v", err)
	}

	return nil
}

func emailMigration(db *sql.DB, id int64, email string) error {
	_, err := globals.ExecDb(db, `
		UPDATE auth SET email = ? WHERE id = ?
	`, email, id)

	return err
}

func setAdmin(db *sql.DB, id int64, isAdmin bool) error {
	_, err := globals.ExecDb(db, `
		UPDATE auth SET is_admin = ? WHERE id = ?
	`, isAdmin, id)

	return err
}

func banUser(db *sql.DB, id int64, isBanned bool) error {
	_, err := globals.ExecDb(db, `
		UPDATE auth SET is_banned = ? WHERE id = ?
	`, isBanned, id)

	return err
}

func quotaMigration(db *sql.DB, id int64, quota float32, override bool) error {
	// 获取操作前的配额和已用配额
	var quotaBefore, usedBefore float32
	err := globals.QueryRowDb(db, "SELECT quota, used FROM quota WHERE user_id = ?", id).Scan(&quotaBefore, &usedBefore)
	if err != nil && err != sql.ErrNoRows { // 检查是否是其他数据库错误
		return fmt.Errorf("failed to get quota before operation: %w", err)
	}
	// 如果是 sql.ErrNoRows，则 beforeQuota 和 usedBefore 保持为 0

	// 执行配额更新操作
	if override {
		_, err = globals.ExecDb(db, `
            INSERT INTO quota (user_id, quota, used) VALUES (?, ?, ?)
            ON DUPLICATE KEY UPDATE quota = ?
        `, id, quota, 0., quota)
	} else {
		_, err = globals.ExecDb(db, `
            INSERT INTO quota (user_id, quota, used) VALUES (?, ?, ?)
            ON DUPLICATE KEY UPDATE quota = quota + ?
        `, id, quota, 0., quota)
	}
	if err != nil {
		return fmt.Errorf("failed to update quota: %w", err)
	}

	// 获取操作后的配额和已用配额
	var quotaAfter, usedAfter float32
	err = globals.QueryRowDb(db, "SELECT quota, used FROM quota WHERE user_id = ?", id).Scan(&quotaAfter, &usedAfter)
	if err != nil {
		return fmt.Errorf("failed to get quota after operation: %w", err)
	}

	// 记录日志
	operation := "admin_increase"
	if override {
		operation = "admin_override"
	}
	quotaChange := quotaAfter - quotaBefore // 计算配额变化
	usedChange := usedAfter - usedBefore    //计算已用配额变化
	_, err = globals.ExecDb(db, `
        INSERT INTO quota_log (user_id, operation, quota_change, quota_before, quota_after, used_change, used_before, used_after)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, id, operation, quotaChange, quotaBefore, quotaAfter, usedChange, usedBefore, usedAfter)
	if err != nil {
		return fmt.Errorf("failed to record quota log: %w", err)
	}

	return nil
}

func subscriptionMigration(db *sql.DB, id int64, expired string) error {
	_, err := globals.ExecDb(db, `
		INSERT INTO subscription (user_id, expired_at) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE expired_at = ?
	`, id, expired, expired)
	return err
}

func subscriptionLevelMigration(db *sql.DB, id int64, level int64) error {
	if level < 0 || level > 3 {
		return fmt.Errorf("invalid subscription level")
	}

	_, err := globals.ExecDb(db, `
		INSERT INTO subscription (user_id, level) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE level = ?
	`, id, level, level)

	return err
}

func releaseUsage(db *sql.DB, cache *redis.Client, id int64) error {
	var level sql.NullInt64
	if err := globals.QueryRowDb(db, `
		SELECT level FROM subscription WHERE user_id = ?
	`, id).Scan(&level); err != nil {
		return err
	}

	if !level.Valid || level.Int64 == 0 {
		return fmt.Errorf("user is not subscribed")
	}

	u := &AuthLike{ID: id}

	plan := channel.PlanInstance.GetPlan(int(level.Int64))
	if !plan.ReleaseAll(u, cache) {
		return fmt.Errorf("cannot release usage")
	}

	return nil
}

func UpdateRootPassword(db *sql.DB, cache *redis.Client, password string) error {
	password = strings.TrimSpace(password)
	if len(password) < 6 || len(password) > 36 {
		return fmt.Errorf("password length must be between 6 and 36")
	}

	if _, err := globals.ExecDb(db, `
		UPDATE auth SET password = ? WHERE username = 'root'
	`, utils.Sha2Encrypt(password)); err != nil {
		return err
	}

	// Clear all user related cache
	if err := clearUserCache(cache); err != nil {
		return fmt.Errorf("failed to clear user cache: %v", err)
	}

	return nil
}

type QuotaLog struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"user_id"`
	Operation   string  `json:"operation"`
	QuotaChange float32 `json:"quota_change"`
	QuotaBefore float32 `json:"quota_before"`
	QuotaAfter  float32 `json:"quota_after"`
	UsedChange  float32 `json:"used_change"` //添加已用配额的改变
	UsedBefore  float32 `json:"used_before"`
	UsedAfter   float32 `json:"used_after"`
	CreatedAt   string  `json:"created_at"` // 使用字符串格式
	Username    string  `json:"username"`   // 新增 username 字段
}

// QuotaLogPaginationForm，用于返回分页数据
type QuotaLogPaginationForm struct {
	Status  bool        `json:"status"`
	Message string      `json:"message,omitempty"`
	Total   int         `json:"total"`
	Data    []*QuotaLog `json:"data"` // 使用 QuotaLog 指针的切片
}

// GetQuotaLogs 类似于 getUsersForm，用于获取 quota_log 数据
func getQuotaLogs(db *sql.DB, page int64, search string, userIDFilter *int64, sortKey string) QuotaLogPaginationForm {
	var logs []*QuotaLog // 使用 QuotaLog 指针的切片
	var total int64

	whereClause := "WHERE 1=1"
	searchQuery := "%" + search + "%"

	if search != "" {
		whereClause += " AND (operation LIKE ?)" // 假设只搜索 operation 字段
	}

	if userIDFilter != nil {
		whereClause += " AND quota_log.user_id = ?" // 修正: 使用 quota_log.user_id
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM quota_log %s", whereClause)
	queryArgsCount := []interface{}{}
	if search != "" {
		queryArgsCount = append(queryArgsCount, searchQuery)
	}
	if userIDFilter != nil {
		queryArgsCount = append(queryArgsCount, *userIDFilter)
	}

	if err := globals.QueryRowDb(db, countQuery, queryArgsCount...).Scan(&total); err != nil {
		return QuotaLogPaginationForm{
			Status:  false,
			Message: err.Error(),
		}
	}

	// 构建排序子句
	orderByClause := "ORDER BY id DESC" // 默认按 ID 倒序
	switch sortKey {
	case "user_id":
		orderByClause = "ORDER BY user_id"
	case "operation":
		orderByClause = "ORDER BY operation"
	case "quota_change":
		orderByClause = "ORDER BY quota_change" //可按变化量升序排列，如需降序排列，则加DESC
	case "created_at":
		orderByClause = "ORDER BY created_at" //按创建时间排序
		// 可以添加更多排序选项
	}

	selectQuery := fmt.Sprintf(`
        SELECT quota_log.id, quota_log.user_id, quota_log.operation, quota_log.quota_change, quota_log.quota_before, quota_log.quota_after, quota_log.used_change, quota_log.used_before, quota_log.used_after, quota_log.created_at, auth.username
        FROM quota_log
        LEFT JOIN auth ON quota_log.user_id = auth.id  
        %s
        %s
        LIMIT ? OFFSET ?
    `, whereClause, orderByClause)

	queryArgsSelect := []interface{}{}
	if search != "" {
		queryArgsSelect = append(queryArgsSelect, searchQuery)
	}
	if userIDFilter != nil {
		queryArgsSelect = append(queryArgsSelect, *userIDFilter)
	}
	queryArgsSelect = append(queryArgsSelect, pagination, page*pagination)

	rows, err := globals.QueryDb(db, selectQuery, queryArgsSelect...)
	if err != nil {
		return QuotaLogPaginationForm{
			Status:  false,
			Message: err.Error(),
		}
	}
	defer rows.Close() // 确保关闭 rows

	for rows.Next() {
		var log QuotaLog
		var createdAt []uint8       // 用于接收时间戳
		var username sql.NullString // 新增: 用于接收 username

		if err := rows.Scan(&log.ID, &log.UserID, &log.Operation, &log.QuotaChange,
			&log.QuotaBefore, &log.QuotaAfter, &log.UsedChange, &log.UsedBefore, &log.UsedAfter, &createdAt, &username); err != nil {
			return QuotaLogPaginationForm{
				Status:  false,
				Message: err.Error(),
			}
		}

		// 处理 username
		if username.Valid {
			log.Username = username.String
		}

		// 转换时间戳为字符串
		t := utils.ConvertTime(createdAt)
		if t != nil {
			log.CreatedAt = t.Format("2006-01-02 15:04:05")
		}

		logs = append(logs, &log) // 将 QuotaLog 指针添加到切片
	}

	return QuotaLogPaginationForm{
		Status: true,
		Total:  int(math.Ceil(float64(total) / float64(pagination))),
		Data:   logs, // 返回 QuotaLog 指针的切片
	}
}
