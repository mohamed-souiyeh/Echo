package db

import	sqlc "echo/db/sqlc_generated"


var hashingCost int = 13

type User = sqlc.User