package auth

import (
	"chat/channel"
	"chat/globals"
	"database/sql"
	"fmt"
)

func (u *User) CreateInitialQuota(db *sql.DB) bool {
	initialQuota := float32(channel.SystemInstance.GetInitialQuota())
	_, err := globals.ExecDb(db, `
		INSERT INTO quota (user_id, quota, used) VALUES (?, ?, ?)
	`, u.GetID(db), initialQuota, 0.)

	if err == nil {
		u.recordQuotaLog(db, "create", initialQuota, 0, initialQuota, 0, 0, 0, "")
	}
	return err == nil
}

func (u *User) GetQuota(db *sql.DB) float32 {
	var quota float32
	if err := globals.QueryRowDb(db, "SELECT quota FROM quota WHERE user_id = ?", u.GetID(db)).Scan(&quota); err != nil {
		return 0.
	}
	return quota
}

func (u *User) GetUsedQuota(db *sql.DB) float32 {
	var quota float32
	if err := globals.QueryRowDb(db, "SELECT used FROM quota WHERE user_id = ?", u.GetID(db)).Scan(&quota); err != nil {
		return 0.
	}
	return quota
}

func (u *User) SetQuota(db *sql.DB, quota float32) bool {
	used := u.GetUsedQuota(db)
	before := u.GetQuota(db)
	_, err := globals.ExecDb(db, `
		INSERT INTO quota (user_id, quota, used) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE quota = ?
	`, u.GetID(db), quota, 0., quota)

	if err == nil {
		u.recordQuotaLog(db, "set", quota-before, before, quota, -used, used, 0, "")
	}
	return err == nil
}

func (u *User) SetUsedQuota(db *sql.DB, used float32) bool {
	before := u.GetUsedQuota(db)
	quota := u.GetQuota(db)
	_, err := globals.ExecDb(db, `
		INSERT INTO quota (user_id, quota, used) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE used = ?
	`, u.GetID(db), 0., used, used)

	if err == nil {
		u.recordQuotaLog(db, "set_used", -quota, quota, 0, used-before, before, used, "")
	}
	return err == nil
}

func (u *User) IncreaseQuota(db *sql.DB, quota float32) bool {
	before := u.GetQuota(db)
	usedBefore := u.GetUsedQuota(db)
	_, err := globals.ExecDb(db, `
		INSERT INTO quota (user_id, quota, used) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE quota = quota + ?
	`, u.GetID(db), quota, 0., quota)
	after := u.GetQuota(db)
	if err == nil {
		u.recordQuotaLog(db, "increase", quota, before, after, 0, usedBefore, usedBefore, "")
	}

	return err == nil
}

func (u *User) IncreaseUsedQuota(db *sql.DB, used float32) bool {
	quotaBefore := u.GetQuota(db)
	usedBefore := u.GetUsedQuota(db)
	_, err := globals.ExecDb(db, `
		INSERT INTO quota (user_id, quota, used) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE used = used + ?
	`, u.GetID(db), 0., used, used)
	usedAfter := u.GetUsedQuota(db)
	if err == nil {
		u.recordQuotaLog(db, "increase_used", 0, quotaBefore, quotaBefore, used, usedBefore, usedAfter, "")
	}
	return err == nil
}

func (u *User) DecreaseQuota(db *sql.DB, quota float32) bool {
	before := u.GetQuota(db)
	usedBefore := u.GetUsedQuota(db)
	_, err := globals.ExecDb(db, `
		INSERT INTO quota (user_id, quota, used) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE quota = quota - ?
	`, u.GetID(db), quota, 0., quota)
	after := u.GetQuota(db)
	if err == nil {
		u.recordQuotaLog(db, "decrease", -quota, before, after, 0, usedBefore, usedBefore, "")
	}
	return err == nil
}

func (u *User) UseQuota(db *sql.DB, quota float32, modelName ...string) bool {
	if quota == 0 {
		return true
	}
	quotaBefore := u.GetQuota(db)
	usedBefore := u.GetUsedQuota(db)
	if !u.DecreaseQuota(db, quota) {
		return false
	}

	if !u.IncreaseUsedQuota(db, quota) {
		return false
	}

	quotaAfter := u.GetQuota(db)
	usedAfter := u.GetUsedQuota(db)

	usedModel := ""
	if len(modelName) > 0 {
		usedModel = modelName[0]
	}
	fmt.Println("usedModel", usedModel)

	u.recordQuotaLog(db, "use", quota, quotaBefore, quotaAfter, quota, usedBefore, usedAfter, usedModel) // 传递模型名称
	return true
}

func (u *User) PayedQuota(db *sql.DB, quota float32) bool {
	if quota == 0 {
		return true
	}
	quotaBefore := u.GetQuota(db)
	current := u.GetQuota(db)
	usedBefore := u.GetUsedQuota(db)
	if quota > current {
		return false
	}

	if !u.DecreaseQuota(db, quota) {
		return false
	}
	if !u.IncreaseUsedQuota(db, quota) {
		return false
	}
	quotaAfter := u.GetQuota(db)
	usedAfter := u.GetUsedQuota(db)
	u.recordQuotaLog(db, "payed", -quota, quotaBefore, quotaAfter, quota, usedBefore, usedAfter, "")
	return true
}

func (u *User) PayedQuotaAsAmount(db *sql.DB, amount float32) bool {
	return u.PayedQuota(db, amount*10)
}

func (u *User) recordQuotaLog(db *sql.DB, operation string, quotaChange, quotaBefore, quotaAfter, usedChange, usedBefore, usedAfter float32, modelName string) {
	_, err := globals.ExecDb(db, `
		INSERT INTO quota_log (user_id, operation, quota_change, quota_before, quota_after, used_change, used_before, used_after, model_name) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, u.GetID(db), operation, quotaChange, quotaBefore, quotaAfter, usedChange, usedBefore, usedAfter, modelName)
	if err != nil {
		globals.Warn(fmt.Sprintf("[quota] failed to record quota log: %s", err.Error()))
	}
}
