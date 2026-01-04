# Hướng Dẫn Viết Unit Test - Mindset & Quy Trình

## 🧠 MINDSET - Tư Duy Viết Test

### 1. **Test là Documentation sống động**
- Test mô tả cách code hoạt động
- Test là ví dụ sử dụng code
- Test giúp người mới hiểu code nhanh hơn

### 2. **Test để tự tin refactor**
- Khi có test tốt, bạn có thể refactor mà không sợ break code
- Test giúp phát hiện bug sớm
- Test giúp thiết kế code tốt hơn (testable code = good code)

### 3. **Test cái gì?**
- ✅ **Test business logic** - Logic xử lý nghiệp vụ
- ✅ **Test edge cases** - Trường hợp biên, lỗi
- ✅ **Test happy path** - Luồng thành công
- ❌ **KHÔNG test** - Framework, library, code của người khác

### 4. **Nguyên tắc FIRST**
- **F**ast - Test chạy nhanh
- **I**ndependent - Test độc lập, không phụ thuộc nhau
- **R**epeatable - Chạy nhiều lần cho cùng kết quả
- **S**elf-validating - Tự động pass/fail
- **T**imely - Viết test cùng lúc với code

---

## 📋 QUY TRÌNH VIẾT TEST - AAA Pattern

### **Arrange → Act → Assert**

```
1. ARRANGE: Setup - Chuẩn bị dữ liệu, mocks
2. ACT: Execute - Chạy function cần test
3. ASSERT: Verify - Kiểm tra kết quả
```

---

## 🎯 VÍ DỤ CỤ THỂ - Phân Tích Từng Bước

### **Ví dụ 1: Test đơn giản - GetConversationByID**

```go
func TestConversationService_GetConversationByID(t *testing.T) {
    // ========== ARRANGE - Setup ==========
    // 1. Tạo mocks (giả lập dependencies)
    mockRepo := new(mocks.MockRepository)
    mockConversationRepo := new(mocks.MockConversationRepo)
    
    // 2. Tạo service với mock
    svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
    reqCtx := &model.RequestContext{}
    
    // 3. Setup expectations (mong đợi gì sẽ xảy ra)
    mockRepo.On("Conversation").Return(mockConversationRepo)
    mockConversationRepo.On("QueryOne", reqCtx, &model.Conversation{ID: 3}).
        Return(model.SuccessResponse(&model.Conversation{ID: 3}, "ok"))
    
    // ========== ACT - Execute ==========
    // Chạy function cần test
    resp := svc.GetConversationByID(reqCtx, 3)
    
    // ========== ASSERT - Verify ==========
    // Kiểm tra kết quả
    assert.True(t, resp.OK())                    // Response phải OK
    assert.Equal(t, uint(3), resp.Data.ID)      // ID phải đúng
    
    // Kiểm tra mocks đã được gọi đúng
    mockRepo.AssertExpectations(t)
    mockConversationRepo.AssertExpectations(t)
}
```

**Phân tích:**
- **Arrange**: Setup mocks và expectations
- **Act**: Gọi `GetConversationByID`
- **Assert**: Kiểm tra response và verify mocks

---

### **Ví dụ 2: Test validation - GetConversationByUserIDs**

```go
t.Run("requires at least two participants", func(t *testing.T) {
    // ========== ARRANGE ==========
    mockRepo := new(mocks.MockRepository)
    mockConversationRepo := new(mocks.MockConversationRepo)
    svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
    reqCtx := &model.RequestContext{}
    
    // ========== ACT ==========
    resp := svc.GetConversationByUserIDs(reqCtx, []uint{1}) // Chỉ 1 user
    
    // ========== ASSERT ==========
    assert.Equal(t, model.CodeBadRequest, resp.Code)  // Phải return BadRequest
    assert.False(t, resp.OK())                        // Không OK
    mockRepo.AssertNotCalled(t, "Conversation")       // KHÔNG được gọi repo
})
```

**Mindset:**
- Test validation logic
- Đảm bảo service không gọi repo khi input invalid
- Early return đúng cách

---

### **Ví dụ 3: Test error handling - CreateConversation fails**

```go
t.Run("returns error when create fails", func(t *testing.T) {
    // ========== ARRANGE ==========
    mockRepo.ExpectedCalls = nil  // Clear previous expectations
    mockConversationRepo.ExpectedCalls = nil
    
    // Setup: Mock Create sẽ return error
    mockRepo.On("Conversation").Return(mockConversationRepo)
    mockConversationRepo.On("Create", reqCtx, mock.AnythingOfType("*model.Conversation")).
        Return(model.BadRequest[*model.Conversation]("Failed to create conversation"))
    
    // ========== ACT ==========
    resp := svc.CreateConversation(reqCtx, []uint{1, 2})
    
    // ========== ASSERT ==========
    assert.Equal(t, model.CodeBadRequest, resp.Code)  // Phải return error
    mockRepo.AssertExpectations(t)                    // Verify đã gọi Conversation()
    mockConversationRepo.AssertExpectations(t)        // Verify đã gọi Create()
    // LƯU Ý: Không expect Participant() vì Create fail → early return
})
```

**Mindset:**
- Test error path
- Đảm bảo service xử lý lỗi đúng
- Đảm bảo không gọi thêm methods không cần thiết khi lỗi

---

### **Ví dụ 4: Test success path - CreateConversation thành công**

```go
t.Run("creates conversation and returns full conversation", func(t *testing.T) {
    // ========== ARRANGE ==========
    mockRepo.ExpectedCalls = nil
    mockConversationRepo.ExpectedCalls = nil
    mockParticipantRepo.ExpectedCalls = nil
    
    // Setup tất cả steps sẽ thành công
    mockRepo.On("Conversation").Return(mockConversationRepo)
    mockRepo.On("Participant").Return(mockParticipantRepo)
    
    created := &model.Conversation{ID: 9}
    // Step 1: Create conversation
    mockConversationRepo.On("Create", reqCtx, mock.AnythingOfType("*model.Conversation")).
        Return(model.SuccessResponse(created, "created"))
    
    // Step 2: Add participants (2 users)
    mockParticipantRepo.On("AddParticipantToConversation", reqCtx, uint(9), uint(1)).
        Return(model.SuccessResponse(&model.ConversationParticipant{}, "added"))
    mockParticipantRepo.On("AddParticipantToConversation", reqCtx, uint(9), uint(2)).
        Return(model.SuccessResponse(&model.ConversationParticipant{}, "added"))
    
    // Step 3: Query conversation để lấy full data
    mockConversationRepo.On("QueryOne", reqCtx, &model.Conversation{ID: 9}).
        Return(model.SuccessResponse(&model.Conversation{ID: 9}, "ok"))
    
    // ========== ACT ==========
    resp := svc.CreateConversation(reqCtx, []uint{1, 2})
    
    // ========== ASSERT ==========
    assert.True(t, resp.OK())              // Phải thành công
    assert.Equal(t, uint(9), resp.Data.ID) // ID đúng
    
    // Verify tất cả steps đã được gọi
    mockConversationRepo.AssertExpectations(t)
    mockParticipantRepo.AssertExpectations(t)
})
```

**Mindset:**
- Test happy path (luồng thành công)
- Setup expectations cho tất cả steps
- Verify từng step đã được thực thi đúng

---

## 🔍 QUY TRÌNH CHI TIẾT - Từng Bước

### **Bước 1: Đọc và hiểu code cần test**

```go
// Code cần test
func (svc *conversationService) CreateConversation(reqCtx *model.RequestContext, userIds []uint) model.Response[*model.Conversation] {
    // Validation
    if len(userIds) < 2 {
        return model.BadRequest[*model.Conversation]("At least 2 participants are required")
    }
    
    // Create conversation
    conversation := &model.Conversation{...}
    createResponse := svc.repo.Conversation().Create(reqCtx, conversation)
    if !createResponse.OK() {
        return createResponse  // Early return nếu fail
    }
    
    // Add participants
    for _, userID := range userIds {
        participantResponse := svc.repo.Participant().AddParticipantToConversation(...)
        if !participantResponse.OK() {
            return model.BadRequest[...]("Failed to add participant")
        }
    }
    
    // Return full conversation
    queryResponse := svc.repo.Conversation().QueryOne(...)
    return queryResponse
}
```

**Phân tích:**
- Có validation: `len(userIds) < 2`
- Có error handling: `if !createResponse.OK()`
- Có loop: `for _, userID := range userIds`
- Có multiple steps: Create → Add participants → Query

### **Bước 2: Xác định test cases**

```
Test cases cần cover:
1. ✅ Validation fail (userIds < 2)
2. ✅ Create conversation fail
3. ✅ Add participant fail
4. ✅ Success - tất cả steps thành công
```

### **Bước 3: Viết test theo AAA**

**Test case 1: Validation**

```go
t.Run("validates minimum participants", func(t *testing.T) {
    // ARRANGE
    mockRepo := new(mocks.MockRepository)
    svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
    reqCtx := &model.RequestContext{}
    
    // ACT
    resp := svc.CreateConversation(reqCtx, []uint{1})
    
    // ASSERT
    assert.Equal(t, model.CodeBadRequest, resp.Code)
    mockRepo.AssertNotCalled(t, "Conversation") // Không gọi repo
})
```

**Test case 2: Create fail**

```go
t.Run("returns error when create fails", func(t *testing.T) {
    // ARRANGE
    mockRepo := new(mocks.MockRepository)
    mockConversationRepo := new(mocks.MockConversationRepo)
    svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
    reqCtx := &model.RequestContext{}
    
    // Setup: Create sẽ fail
    mockRepo.On("Conversation").Return(mockConversationRepo)
    mockConversationRepo.On("Create", reqCtx, mock.AnythingOfType("*model.Conversation")).
        Return(model.BadRequest[*model.Conversation]("Failed"))
    
    // ACT
    resp := svc.CreateConversation(reqCtx, []uint{1, 2})
    
    // ASSERT
    assert.Equal(t, model.CodeBadRequest, resp.Code)
    mockRepo.AssertExpectations(t)
    mockConversationRepo.AssertExpectations(t)
    // KHÔNG expect Participant vì Create fail → early return
})
```

**Test case 3: Success**

```go
t.Run("creates conversation and returns full conversation", func(t *testing.T) {
    // ARRANGE
    mockRepo := new(mocks.MockRepository)
    mockConversationRepo := new(mocks.MockConversationRepo)
    mockParticipantRepo := new(mocks.MockParticipantRepo)
    svc := conversation.NewConversationService(&common.Params{Repo: mockRepo})
    reqCtx := &model.RequestContext{}
    
    // Setup tất cả steps
    mockRepo.On("Conversation").Return(mockConversationRepo)
    mockRepo.On("Participant").Return(mockParticipantRepo)
    
    created := &model.Conversation{ID: 9}
    mockConversationRepo.On("Create", reqCtx, mock.AnythingOfType("*model.Conversation")).
        Return(model.SuccessResponse(created, "created"))
    
    // Setup cho loop (2 participants)
    mockParticipantRepo.On("AddParticipantToConversation", reqCtx, uint(9), uint(1)).
        Return(model.SuccessResponse(&model.ConversationParticipant{}, "added"))
    mockParticipantRepo.On("AddParticipantToConversation", reqCtx, uint(9), uint(2)).
        Return(model.SuccessResponse(&model.ConversationParticipant{}, "added"))
    
    mockConversationRepo.On("QueryOne", reqCtx, &model.Conversation{ID: 9}).
        Return(model.SuccessResponse(&model.Conversation{ID: 9}, "ok"))
    
    // ACT
    resp := svc.CreateConversation(reqCtx, []uint{1, 2})
    
    // ASSERT
    assert.True(t, resp.OK())
    assert.Equal(t, uint(9), resp.Data.ID)
    mockConversationRepo.AssertExpectations(t)
    mockParticipantRepo.AssertExpectations(t)
})
```

---

## 🎓 BEST PRACTICES

### 1. **Mỗi test case test 1 thing**
```go
// ✅ GOOD
t.Run("validates minimum participants", ...)
t.Run("returns error when create fails", ...)
t.Run("creates conversation successfully", ...)

// ❌ BAD
t.Run("test everything", ...) // Quá nhiều assertions
```

### 2. **Test name mô tả rõ ràng**
```go
// ✅ GOOD
t.Run("returns error when create fails", ...)
t.Run("requires at least two participants", ...)

// ❌ BAD
t.Run("test1", ...)
t.Run("test create", ...)
```

### 3. **Setup mocks rõ ràng**
```go
// ✅ GOOD - Clear expectations
mockRepo.On("Conversation").Return(mockConversationRepo)
mockConversationRepo.On("Create", reqCtx, mock.AnythingOfType("*model.Conversation")).
    Return(model.SuccessResponse(created, "created"))

// ❌ BAD - Không rõ ràng
mockRepo.On("Conversation").Return(mockConversationRepo)
// Thiếu setup Create expectation
```

### 4. **Verify cả behavior và result**
```go
// ✅ GOOD
assert.True(t, resp.OK())                    // Verify result
mockRepo.AssertExpectations(t)               // Verify behavior

// ❌ BAD - Chỉ verify result
assert.True(t, resp.OK())
// Không verify mocks → không biết code có gọi đúng dependencies không
```

### 5. **Test edge cases**
```go
// ✅ GOOD
t.Run("empty userIDs", ...)
t.Run("single userID", ...)
t.Run("duplicate userIDs", ...)

// ❌ BAD - Chỉ test happy path
t.Run("success", ...)
```

### 6. **Isolate tests**
```go
// ✅ GOOD - Clear expectations mỗi test
t.Run("test1", func(t *testing.T) {
    mockRepo.ExpectedCalls = nil  // Clear
    // Setup fresh
})

// ❌ BAD - Dùng chung expectations
mockRepo.On("Conversation").Return(...) // Setup ở ngoài
t.Run("test1", ...) // Có thể bị ảnh hưởng bởi test khác
```

---

## 🚀 QUY TRÌNH TỔNG QUÁT

```
1. ĐỌC CODE
   ↓
2. PHÂN TÍCH
   - Input/Output là gì?
   - Có validation không?
   - Có error handling không?
   - Có dependencies gì?
   ↓
3. XÁC ĐỊNH TEST CASES
   - Happy path
   - Error cases
   - Edge cases
   ↓
4. VIẾT TEST (AAA)
   - Arrange: Setup mocks, data
   - Act: Execute function
   - Assert: Verify result + behavior
   ↓
5. CHẠY TEST
   - Pass? → Done ✅
   - Fail? → Fix code hoặc fix test
   ↓
6. REFACTOR (nếu cần)
   - Test vẫn pass sau refactor? → Good ✅
```

---

## 💡 TIPS & TRICKS

### 1. **Dùng table-driven tests cho nhiều cases tương tự**
```go
tests := []struct {
    name     string
    input    []uint
    expected int
}{
    {"empty", []uint{}, model.CodeBadRequest},
    {"single", []uint{1}, model.CodeBadRequest},
    {"two", []uint{1, 2}, model.CodeSuccess},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // Test logic
    })
}
```

### 2. **Dùng helper functions**
```go
func setupMocks() (*mocks.MockRepository, *mocks.MockConversationRepo) {
    mockRepo := new(mocks.MockRepository)
    mockConversationRepo := new(mocks.MockConversationRepo)
    mockRepo.On("Conversation").Return(mockConversationRepo)
    return mockRepo, mockConversationRepo
}
```

### 3. **Clear expectations giữa các tests**
```go
t.Run("test1", func(t *testing.T) {
    mockRepo.ExpectedCalls = nil
    // Fresh setup
})
```

---

## ❓ FAQ

**Q: Khi nào dùng mock?**
A: Khi test service layer, mock repository. Khi test repository layer, dùng test DB.

**Q: Test bao nhiêu là đủ?**
A: Cover tất cả branches (if/else, loops, error paths). Aim for >80% coverage.

**Q: Test có cần test private functions không?**
A: Không cần. Test qua public interface. Nếu private function phức tạp, có thể tách thành function riêng để test.

**Q: AssertExpectations có cần thiết không?**
A: Có! Đảm bảo code gọi đúng dependencies. Không chỉ test result, mà còn test behavior.

---

## 📚 TÀI LIỆU THAM KHẢO

- [Testify Mock Documentation](https://github.com/stretchr/testify#mock-package)
- [Go Testing Best Practices](https://golang.org/doc/effective_go#testing)
- [Unit Testing Principles](https://martinfowler.com/bliki/UnitTest.html)



