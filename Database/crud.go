package Database

func PushUser(user, group string) error {
	_, err := db.Exec(`insert into "user" values($1, $2) on conflict do update set "group"=$2`, user, group)
	return err
}

func ReadUser(user string) (string, error) {
	var group string
	err := db.QueryRow(`select "group" from "user" where "id"=$1`, user).Scan(&group)
	return group, err
}

func DeleteUser(user string) error {
	_, err := db.Exec(`delete from "user" where "id"=$1`, user)
	return err
}
