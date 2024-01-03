#!/bin/bash

execute() {
	# echo "Execute sql-migrate $1 to $2"
	cd /app/databases/$2 && sql-migrate $1 $3
}

status() {
	# echo "Display sql-migrate status from $1"
	cd /app/databases/$1 && sql-migrate status
}

main() {
	case "$1" in
	up | down | redo | status)
		# generate
		;;
	*)
		echo "Usage: entrypoint.sh COMMAND <database>"
		echo
		echo "Commands:"
		echo "  up <database>        Migrates the database to the most recent version available"
		echo "  down <database>      Undo a database migration"
		echo "  redo <database>      Reapply the last migration"
		echo "  status <database>    Show migration status"
		exit 1
		;;
	esac

	if [[ $1 == "status" ]]; then
		status $2
	else
		execute $1 $2 $3
	fi

}

main "$@"
