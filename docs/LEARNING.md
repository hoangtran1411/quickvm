# Hướng Dẫn Học & Trình Bày Project QuickVM (Learning & Presentation Guide)

Tài liệu này được thiết kế để giúp bạn:
1. **Truyền tải kiến thức** cho người mới học Go thông qua project này.
2. **Trình bày project** trong các buổi phỏng vấn (Interview) để gây ấn tượng mạnh.

---

## Phần 1: Dùng QuickVM để dạy lập trình Go (Teaching Modules)

Project này là một ví dụ tuyệt vời để dạy Go vì nó không quá phức tạp nhưng bao gồm đủ các khái niệm quan trọng. Bạn có thể chia bài giảng thành 5 module:

### Module 1: Cấu trúc dự án Go & Packages
*   **Mục tiêu**: Hiểu cách tổ chức code trong Go.
*   **Điểm nhấn**:
    *   `cmd/`: Chứa entry point của các lệnh (CLI logic).
    *   `hyperv/`: Chứa business logic (tách biệt logic nghiệp vụ khỏi giao diện).
    *   `ui/`: Chứa code giao diện TUI.
*   **Bài học**: Tại sao không để tất cả trong `main.go`? -> *Để dễ bảo trì và test từng phần riêng biệt.*

### Module 2: Xây dựng CLI với Cobra
*   **Mục tiêu**: Làm sao để tạo ra các lệnh như `quickvm start 1`.
*   **File tham khảo**: `cmd/start.go`, `cmd/root.go`.
*   **Bài học**:
    *   Cách khai báo `cobra.Command`.
    *   Cách đọc tham số (Arguments) và Cờ (Flags).
    *   Pattern "Command" để mở rộng tính năng dễ dàng.

### Module 3: Tương tác Hệ thống (System Calls)
*   **Mục tiêu**: Gọi lệnh PowerShell từ Go.
*   **File tham khảo**: `hyperv/powershell.go`.
*   **Bài học**:
    *   Sử dụng `os/exec` để gọi `powershell.exe`.
    *   Xử lý `CombinedOutput` để lấy kết quả hoặc lỗi.
    *   *Nâng cao*: Xử lý JSON trả về từ PowerShell để map vào struct Go.

### Module 4: Giao diện Terminal đẹp (TUI)
*   **Mục tiêu**: Làm app console trông "xịn xò" hơn.
*   **File tham khảo**: `ui/` directory.
*   **Bài học**:
    *   Mô hình ELM (Model - Update - View) trong thư viện `Bubble Tea`.
    *   Xử lý sự kiện bàn phím (Key bindings).

### Module 5: Testing & Interfaces
*   **Mục tiêu**: Viết code có thể test được (Testable Code).
*   **File tham khảo**: `hyperv/manager_test.go` (hoặc các file test tương tự).
*   **Bài học**:
    *   Tại sao dùng Interface (`ShellExecutor`) thay vì gọi trực tiếp? -> *Để Mock được hành vi khi chạy test mà không cần cài Hyper-V thật.*

---

## Phần 2: Mang Project đi Phỏng vấn (Interview Strategy)

Khi trả lời phỏng vấn, hãy dùng phương pháp **STAR** (Situation - Task - Action - Result).

### 1. Giới thiệu Project (Context)
> "Em/Mình đã xây dựng một công cụ CLI/TUI tên là QuickVM để quản lý máy ảo Hyper-V trên Windows."

### 2. Situation (Vấn đề gặp phải)
*   "Giao diện Hyper-V Manager mặc định của Windows khá chậm và nặng nề."
*   "Việc thao tác lặp đi lặp lại (bật/tắt VM để test) tốn nhiều thời gian click chuột."
*   "Cần một công cụ có thể script hóa (automation) được."

### 3. Task (Nhiệm vụ)
*   Xây dựng một công cụ bằng Go (Golang) vừa có CLI (để chạy lệnh nhanh/scripting) vừa có TUI (giao diện terminal) để dễ nhìn.
*   Yêu cầu: Nhanh, nhẹ, và dễ mở rộng.

### 4. Action (Giải pháp kỹ thuật - Phần quan trọng nhất)
Hãy nhấn mạnh vào các quyết định kỹ thuật (Engineering Decisions):

*   **Kiến trúc**: "Em chọn **Cobra** cho CLI vì nó là chuẩn công nghiệp (dùng bởi K8s, Docker). Em tách biệt logic `hyperv` ra khỏi `cmd` để code clean hơn."
*   **Tương tác hệ thống**: "Thay vì dùng CGO gọi Windows API (phức tạp, khó build), em chọn wrapper qua **PowerShell**. Để tối ưu, em parse output JSON từ PowerShell thẳng vào Go structs."
*   **Giao diện TUI**: "Sử dụng **Bubble Tea** để tạo trải nghiệm người dùng hiện đại, có tương tác bàn phím, update trạng thái realtime."
*   **Testability**: "Một thách thức là test code này trên môi trường không có Hyper-V (như CI/CD linux). Em đã giải quyết bằng cách dùng **Interface** để trừu tượng hóa lớp thực thi lệnh (ShellExecutor), cho phép mock kết quả trả về khi test."

### 5. Result (Kết quả)
*   "Quản lý VM nhanh hơn 50% so với dùng GUI."
*   "Dễ dàng tích hợp vào các script automation."
*   "Code coverage đạt mức tốt nhờ thiết kế testable."

---

## Các câu hỏi phỏng vấn thường gặp về Project này

Chuẩn bị sẵn câu trả lời cho các câu hỏi "xoáy" sau:

1.  **"Tại sao lại gọi PowerShell mà không dùng syscall của Windows?"**
    *   *Trả lời*: Dùng syscall/COM rất phức tạp và dễ lỗi bộ nhớ. PowerShell cung cấp sẵn các cmdlets ổn định của Microsoft. Trade-off là hiệu năng chậm hơn một chút nhưng an toàn và dễ bảo trì hơn nhiều.

2.  **"Làm sao để xử lý việc PowerShell chạy chậm?"**
    *   *Trả lời*: Em chỉ gọi PowerShell khi cần thiết. Với các tác vụ nặng, có thể dùng Goroutines để chạy song song (concurrency) hoặc cache lại danh sách VM nếu cần.

3.  **"Project này xử lý lỗi (Error Handling) như thế nào?"**
    *   *Trả lời*: Em wrap lỗi để ngữ cảnh rõ ràng (dùng `fmt.Errorf("start vm failed: %w", err)`). Điều này giúp debug dễ dàng hơn khi biết lỗi xuất phát từ lớp nào.

---

## Mẹo: In ấn tài liệu này (Convert to Word/Docx)

Nếu bạn muốn in tài liệu này hoặc gửi dưới dạng file Word, bạn có thể chuyển đổi trực tiếp bằng lệnh sau (yêu cầu máy có cài Node.js):

```powershell
# Chạy lệnh này tại thư mục gốc của dự án
npx -y markdown-docx --input docs/LEARNING.md --output docs/LEARNING.docx
```

Lệnh trên sẽ tự động tải công cụ về và tạo file `docs/LEARNING.docx` cho bạn mà không cần cài đặt thêm phần mềm gì.
