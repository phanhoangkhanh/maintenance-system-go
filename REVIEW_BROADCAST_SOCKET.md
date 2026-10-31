# Ghi chú tiếp tục review broadcastSocket/service.go

Ngày lưu: 2026-10-30 (Asia/Bangkok).

Đây là bản tóm tắt cuộc trao đổi và trạng thái công việc, không phải bản chép nguyên văn hội thoại.

## Yêu cầu của người dùng

- Review `broadcastSocket/service.go`, giải thích thiếu sót và hướng sửa.
- Chỉ hướng dẫn, không tự chỉnh code. Người dùng tự sửa rồi yêu cầu review lại.
- Tiếp tục giữ giới hạn này trừ khi người dùng cho phép chỉnh code.
- Các nhận xét dưới đây dựa trên đọc mã nguồn; chưa chạy kiểm thử.

## Bối cảnh

- `tradingSocket` gửi `MarketPriceInfo` vào `BridgeChannel`.
- Broadcaster nhân bản message sang `msgSend` của từng client, buffer 1000.
- Mục tiêu: queue của một client đầy thì đóng client đó, không làm ngừng các client còn lại.
- Ban đầu có data race trên slice clients, gửi channel có thể block toàn bộ broadcaster, và writer có thể chờ mãi sau disconnect.

## Những thay đổi người dùng đã làm

- Tạo `ClientList` có mutex bảo vệ danh sách.
- Dùng `select/default` để thử gửi vào `msgSend` mà không chờ queue trống.
- Thêm `signalDone` và callback `callCloseNow` vào `Client`.
- Callback dùng `sync.Once` riêng cho mỗi connection để đóng socket, đóng `signalDone`, xóa client và ghi log.
- Không giữ `client.lock` khi đóng socket.
- Reader và writer đều defer callback cleanup.
- Writer dùng vòng `for` bao quanh `select`, nhận message hoặc tín hiệu dừng.
- Đã bỏ các lệnh gửi giá trị vào `signalDone`; chỉ đóng channel trong cleanup.
- Không đóng `msgSend`, tránh gửi vào channel đã đóng khi broadcaster vẫn giữ tham chiếu client.

## Lỗi quan trọng nhất còn lại ở lần review cuối

Broadcaster vẫn giữ `clientList.lock` khi gọi `client.callCloseNow(client)` ở nhánh queue đầy.

Luồng gây deadlock trong phiên bản đã đọc:

```text
multipleMsgToClients
  → Lock(clientList.lock)          dòng 107
  → queue đầy
  → client.callCloseNow(client)    dòng 116
  → removeClientFromList(client)
  → Lock(clientList.lock)          dòng 180
  → chờ mãi
```

Cleanup chờ khóa do broadcaster đang giữ; broadcaster chờ cleanup hoàn tất. Việc broadcast và thêm/xóa client có thể bị kẹt.

Hướng sửa đã hướng dẫn, chưa áp dụng vào source:

```go
clientList.lock.Lock()
snapshot := append([]*Client(nil), clientList.clients...)
clientList.lock.Unlock()

for _, client := range snapshot {
    select {
    case <-client.signalDone:
        continue
    case client.msgSend <- msg:
    default:
        client.callCloseNow(client)
    }
}
```

Phải sao chép phần tử sang mảng slice mới; chỉ gán `snapshot := clientList.clients` vẫn dùng chung mảng. Gửi message và cleanup phải diễn ra sau khi nhả khóa danh sách.

## Những điểm cần kiểm tra tiếp

1. `select` không ưu tiên nhánh `signalDone`. Khi signal đã đóng và queue còn chỗ, vẫn có thể enqueue thêm message cho client đã dừng. Hiện không gây panic vì `msgSend` không bị đóng; không nên giải quyết bằng cách đóng `msgSend` tùy tiện.
2. Ghi phản hồi login (dòng 143 của bản đã đọc) đang bỏ qua lỗi rồi đặt `isReady = true`. Cần kiểm tra lỗi, nhả khóa và thoát reader khi ghi thất bại.
3. Chưa có write deadline cho ghi giá và phản hồi login. Nếu ghi bị kẹt nhưng nguồn không gửi thêm đủ dữ liệu để đầy queue, cơ chế ngắt client chậm chưa kích hoạt.
4. `ReadJSON` đang chờ không tự bị ngắt khi `signalDone` đóng. Cleanup chủ động đóng socket là phần giải phóng thao tác đọc/ghi; giữ cơ chế đó.
5. Các cải thiện chưa xử lý từ review đầu: lỗi `ListenAndServe`, read limit, ping/pong và read deadline, thời hạn đăng nhập, xác thực cấu hình/token thay secret hardcode, quản lý shutdown server và WebSocket.
6. Comment đầu file còn nói gửi alarm vào channel; cập nhật cho khớp cơ chế cleanup mới khi người dùng sửa tiếp.

## Các tình huống kiểm thử đề xuất

- Queue của một client đầy: client đó đóng, client khác vẫn nhận dữ liệu.
- Client disconnect khi queue rỗng: cả reader và writer đều thoát.
- Reader/writer cùng lỗi hoặc broadcaster cùng yêu cầu dừng: không panic, cleanup chỉ chạy một lần.
- Nhiều client connect/disconnect trong lúc broadcast: dùng race detector với kiểm thử thực sự thực thi các đường chạy đó.
- Client không đọc hoặc mất mạng: deadline/heartbeat kết thúc kết nối trong thời gian giới hạn.

## Cách tiếp tục

Đọc lại file hiện tại trước khi review vì người dùng có thể đã sửa sau khi lưu ghi chú. Tập trung kiểm tra deadlock ở nhánh queue đầy trước, rồi lỗi login và deadline. Không tự sửa source nếu chưa được yêu cầu.
