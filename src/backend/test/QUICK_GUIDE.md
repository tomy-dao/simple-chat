# 🚀 Hướng Dẫn Viết Unit Test - Tóm Tắt Nhanh

## 🎯 MINDSET (3 điều cần nhớ)

1. **Test = Documentation** - Test mô tả cách code hoạt động
2. **Test để tự tin refactor** - Có test thì sửa code không sợ
3. **Test business logic** - Không test framework/library

---

## 📝 QUY TRÌNH (3 bước - AAA)

```
ARRANGE → ACT → ASSERT
```

### **ARRANGE** - Chuẩn bị
- Tạo mocks (giả lập dependencies)
- Setup expectations (mong đợi gì sẽ xảy ra)

### **ACT** - Chạy
- Gọi function cần test

### **ASSERT** - Kiểm tra
- Kiểm tra kết quả (result)
- Kiểm tra behavior (mocks đã được gọi đúng chưa)

---

## 💡 VÍ DỤ ĐƠN GIẢN

```go
func TestGetConversationByID(t *testing.T) {
    // ========== ARRANGE ==========
    mockRepo := new(mocks.MockRepository)
    mockConversationRepo := new(mocks.MockConversationRepo)
    svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
    
    // Setup: Khi gọi Conversation() → return mockConversationRepo
    mockRepo.On("Conversation").Return(mockConversationRepo)
    
    // Setup: Khi gọi QueryOne() → return success
    mockConversationRepo.On("QueryOne", reqCtx, &model.Conversation{ID: 3}).
        Return(model.SuccessResponse(&model.Conversation{ID: 3}, "ok"))
    
    // ========== ACT ==========
    resp := svc.GetConversationByID(reqCtx, 3)
    
    // ========== ASSERT ==========
    assert.True(t, resp.OK())              // Kết quả đúng?
    assert.Equal(t, uint(3), resp.Data.ID) // Data đúng?
    mockRepo.AssertExpectations(t)         // Đã gọi đúng chưa?
}
```

---

## 🎓 5 ĐIỀU CẦN NHỚ

### 1. **Mỗi test case test 1 thing**
```go
t.Run("validation fails", ...)  // ✅
t.Run("create fails", ...)       // ✅
t.Run("success", ...)            // ✅
```

### 2. **Test name mô tả rõ**
```go
t.Run("returns error when create fails", ...)  // ✅
t.Run("test1", ...)                             // ❌
```

### 3. **Setup expectations trước khi chạy**
```go
// ✅ GOOD
mockRepo.On("Conversation").Return(mockConversationRepo)
mockConversationRepo.On("Create", ...).Return(...)

// ❌ BAD - Thiếu setup
mockRepo.On("Conversation").Return(mockConversationRepo)
// Thiếu setup Create → test sẽ fail
```

### 4. **Verify cả result VÀ behavior**
```go
assert.True(t, resp.OK())        // ✅ Verify result
mockRepo.AssertExpectations(t)   // ✅ Verify behavior
```

### 5. **Test cả error cases**
```go
t.Run("validation fails", ...)  // ✅
t.Run("create fails", ...)       // ✅
t.Run("success", ...)            // ✅
```

---

## 🔄 QUY TRÌNH 4 BƯỚC

```
1. ĐỌC CODE → Hiểu logic
2. XÁC ĐỊNH CASES → Happy path + Error cases
3. VIẾT TEST (AAA) → Arrange → Act → Assert
4. CHẠY TEST → Pass? Done ✅
```

---

## ❓ FAQ NGẮN

**Q: Mock là gì?**
A: Giả lập dependencies (như repo, database) để test độc lập.

**Q: AssertExpectations để làm gì?**
A: Kiểm tra code đã gọi đúng dependencies chưa (không chỉ kiểm tra kết quả).

**Q: Test bao nhiêu là đủ?**
A: Test tất cả branches (if/else, error paths). Aim >80% coverage.

**Q: Khi nào cần mock?**
A: Khi test service → mock repo. Khi test repo → dùng test DB.

---

## 🎯 TEMPLATE NHANH

```go
func TestFunctionName(t *testing.T) {
    // ARRANGE
    mockRepo := new(mocks.MockRepository)
    mockXRepo := new(mocks.MockXRepo)
    svc := service.NewService(&common.Params{Repo: mockRepo})
    
    mockRepo.On("X").Return(mockXRepo)
    mockXRepo.On("Method", ...).Return(...)
    
    // ACT
    resp := svc.FunctionName(...)
    
    // ASSERT
    assert.True(t, resp.OK())
    assert.Equal(t, expected, resp.Data)
    mockRepo.AssertExpectations(t)
    mockXRepo.AssertExpectations(t)
}
```

---

**Tóm lại: AAA (Arrange → Act → Assert) + Verify cả result và behavior!**



