# maintenance-system-go

## Connect database postgress :

    1. Migrate:
        github.com/golang-migrate/migrate/v4/cmd/migrate@v4.18.3

    2. Driver:
        go get gorm.io/gorm@v1.31.1
        go get gorm.io/driver/postgres@v1.6.0

    3. Database:
        sudo -i -u postgres -> psql ->
        CREATE DATABASE maintenance_system ->
        GRANT CONNECT ON DATABASE maintenance_system TO admin_ets;
        ALTER DATABASE maintenance_system OWNER TO admin_ets;
    4. Set Export params:
        export DATABASE_URL='postgres://admin_ets:123456@localhost:5432/maintenance_system?sslmode=disable'

        4.1. Run migrate with export params:
            migrate \
                -path database/migrations \
                -database "$DATABASE_URL" \
                up

            migrate \
                -path database/migrations \
                -database "$DATABASE_URL" \
                down 1
