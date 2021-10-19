migrateup:
	migrate -path migrations -database $(DATABASE_URL) -verbose up

migratedown:
	migrate -path migrations -database $(DATABASE_URL) -verbose down

.PHONY: migrateup migratedown
