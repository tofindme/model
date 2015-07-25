// GENERATED CODE - DO NOT EDIT
package Model

type employee_logins struct {
	id                   int       `db:"id"    json:"id"`
	employee_id          int       `db:"employee_id"    json:"employee_id"`
	IP                   varchar   `db:"IP"    json:"IP"`
	login_at             timestamp `db:"login_at"    json:"login_at"`
	login_status         tinyint   `db:"login_status"    json:"login_status"`
	employee_sessions_id int       `db:"employee_sessions_id"    json:"employee_sessions_id"`
}
