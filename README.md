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
            migrate
                -path database/migrations
                -database "$DATABASE_URL"
                up

            migrate
                -path database/migrations
                -database "$DATABASE_URL"
                down 1

    5. Structure
        Init App -> init App with real-struct
        Every module-struct must be Init in each provider inside module

    6. Install Redis via Docker
        Os linux + open Docker Deskhop dowload img Redis: docker pull redis:7.4-alpine
        In DockerDeskhop -> build img to container + setup params
                    name: myRedis + port: 6379
        Install go-redis/v9 : go get github.com/redis/go-redis/v9
        Contruct module Redis and trigger in provider.go -> register for serviec in each modules need:
                ex: service := &Service{
                        Repo: repo,
                        Redis: redisClient,
                    }
    7. Create UUID-Value and store in Redis -> with TTL 48h.
            Set cookies for Login-request for first-time
            Next time , the request has cookies inside to valid whom request
            Check request with cookies inside to confirm who is requested from browser
            Key was created via 'github.com/google/uuid' : go get github.com/google/uuid


## CORS cho UI gọi API bằng cookie

Khai báo trong `.env` (local), hoặc biến môi trường khi `ENV=prod`:

```dotenv
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,https://app.example.com
```

Thay `https://app.example.com` bằng origin UI thực tế. Mỗi origin gồm scheme,
host và port nếu có; không kèm path hoặc dấu `/` cuối. Không dùng `*`.
Danh sách là các địa chỉ UI gọi API, không phải URL API đích.
Local mặc định cho phép hai origin localhost trên khi chưa cấu hình;
production không cho phép cross-origin nếu danh sách trống.

Frontend cần `credentials: "include"` cho cả login và các request sau:

```js
fetch("https://api.example.com/login", {
  method: "POST",
  credentials: "include",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({ name, password }),
});
```

Middleware trả origin khớp request, cho phép credentials và xử lý preflight
OPTIONS. CORS không thay thế chính sách cookie: UI/API khác site cần
`SameSite=None; Secure` và vẫn chịu chính sách cookie bên thứ ba của browser.
