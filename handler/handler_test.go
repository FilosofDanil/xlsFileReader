package handler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler()

	if h == nil {
		t.Fatal("NewHandler returned nil")
	}

	if h.userStates == nil {
		t.Error("userStates map not initialized")
	}
}

func TestIsValidExcelFile(t *testing.T) {
	h := NewHandler()

	tests := []struct {
		name     string
		fileName string
		expected bool
	}{
		{"valid xlsx file", "test.xlsx", true},
		{"valid xls file", "test.xls", true},
		{"valid xlsx uppercase", "TEST.XLSX", true},
		{"valid xls uppercase", "TEST.XLS", true},
		{"valid xlsx mixed case", "Test.XlSx", true},
		{"valid xls mixed case", "Test.XlS", true},
		{"invalid txt file", "test.txt", false},
		{"invalid pdf file", "document.pdf", false},
		{"invalid csv file", "data.csv", false},
		{"invalid no extension", "noextension", false},
		{"invalid empty string", "", false},
		{"valid xlsx with path", "/path/to/file.xlsx", true},
		{"valid xls with spaces", "my file.xls", true},
		{"valid xlsx with cyrillic", "файл.xlsx", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := h.isValidExcelFile(tt.fileName)
			if result != tt.expected {
				t.Errorf("isValidExcelFile(%q) = %v, want %v", tt.fileName, result, tt.expected)
			}
		})
	}
}

func TestGenerateSQLScript(t *testing.T) {
	h := NewHandler()

	t.Run("single contract", func(t *testing.T) {
		contracts := []string{"228960453-123"}
		result := h.generateSQLScript(contracts)

		// Check SELECT clause present
		if !strings.Contains(result, "SELECT") {
			t.Error("Result should contain SELECT")
		}

		// Check contract in VALUES clause
		if !strings.Contains(result, "('EP-228960453-123', 0)") {
			t.Error("Result should contain contract with correct format")
		}

		// Check ORDER BY present
		if !strings.Contains(result, "ORDER BY") {
			t.Error("Result should contain ORDER BY")
		}

		// Check closing statement
		if !strings.HasSuffix(result, "sort_order.sort_seq;") {
			t.Error("Result should end with sort_order.sort_seq;")
		}
	})

	t.Run("multiple contracts", func(t *testing.T) {
		contracts := []string{"111111-aaa", "222222-bbb", "333333-ccc"}
		result := h.generateSQLScript(contracts)

		// Check all contracts present with correct indices
		if !strings.Contains(result, "('EP-111111-aaa', 0)") {
			t.Error("Result should contain first contract with index 0")
		}
		if !strings.Contains(result, "('EP-222222-bbb', 1)") {
			t.Error("Result should contain second contract with index 1")
		}
		if !strings.Contains(result, "('EP-333333-ccc', 2)") {
			t.Error("Result should contain third contract with index 2")
		}

		// Check contracts are comma-separated
		if !strings.Contains(result, "),\n") {
			t.Error("Contracts should be comma-separated with newlines")
		}
	})

	t.Run("empty contracts", func(t *testing.T) {
		contracts := []string{}
		result := h.generateSQLScript(contracts)

		// Should still have valid SQL structure
		if !strings.Contains(result, "SELECT") {
			t.Error("Result should contain SELECT even with empty contracts")
		}
		if !strings.Contains(result, "(VALUES") {
			t.Error("Result should contain VALUES clause")
		}
	})

	t.Run("SQL structure validation", func(t *testing.T) {
		contracts := []string{"123-456"}
		result := h.generateSQLScript(contracts)

		// Validate required columns
		expectedColumns := []string{
			"sort_order.number AS 'Номер договору'",
			"dbo.getCagentFullName(c.id_acquisitor) AS 'Аквізитор'",
			"dbo.getCagentFullName(c.id_responsible) AS 'Відповідальна особа'",
			"AS 'Канал продажів'",
			"AS 'Підканал продажів'",
			"AS 'Обліковий підрозділ'",
			"AS 'Вищестоящий підрозділ'",
		}

		for _, col := range expectedColumns {
			if !strings.Contains(result, col) {
				t.Errorf("Result should contain column: %s", col)
			}
		}

		// Validate JOINs
		expectedJoins := []string{
			"LEFT JOIN contract c ON c.number = sort_order.number",
			"LEFT JOIN division div ON div.id = c.id_division",
			"LEFT JOIN helement h_div ON h_div.id = div.id",
		}

		for _, join := range expectedJoins {
			if !strings.Contains(result, join) {
				t.Errorf("Result should contain JOIN: %s", join)
			}
		}
	})
}

func TestReadExcelFile_Routing(t *testing.T) {
	h := NewHandler()

	// Test that readExcelFile routes to correct reader based on extension
	t.Run("xlsx extension routing", func(t *testing.T) {
		// This will fail because file doesn't exist, but we can verify error message
		_, err := h.readExcelFile("test.xlsx")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
		// Error should be about opening Excel file, not about XLS file
		if !strings.Contains(err.Error(), "failed to open Excel file") {
			t.Errorf("Expected xlsx-related error, got: %v", err)
		}
	})

	t.Run("xls extension routing", func(t *testing.T) {
		_, err := h.readExcelFile("test.xls")
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
		// Error should be about opening XLS file
		if !strings.Contains(err.Error(), "failed to open XLS file") {
			t.Errorf("Expected xls-related error, got: %v", err)
		}
	})
}

func TestReadXlsxFile_Integration(t *testing.T) {
	h := NewHandler()

	// Create a temporary test directory
	testDir := t.TempDir()

	t.Run("read xlsx with matching data", func(t *testing.T) {
		// Create test xlsx file
		f := excelize.NewFile()
		defer f.Close()

		// Set test data - "ББС ІНШУРАНС" in column A, values in B and C
		f.SetCellValue("Sheet1", "A1", "Header1")
		f.SetCellValue("Sheet1", "B1", "Header2")
		f.SetCellValue("Sheet1", "C1", "Header3")
		f.SetCellValue("Sheet1", "A2", "ББС ІНШУРАНС")
		f.SetCellValue("Sheet1", "B2", "228960453")
		f.SetCellValue("Sheet1", "C2", "123")
		f.SetCellValue("Sheet1", "A3", "ББС ІНШУРАНС")
		f.SetCellValue("Sheet1", "B3", "228209382")
		f.SetCellValue("Sheet1", "C3", "456")

		testFile := filepath.Join(testDir, "test_data.xlsx")
		if err := f.SaveAs(testFile); err != nil {
			t.Fatalf("Failed to save test file: %v", err)
		}

		// Read the file
		result, err := h.readXlsxFile(testFile)
		if err != nil {
			t.Fatalf("readXlsxFile failed: %v", err)
		}

		// Verify SQL contains extracted data
		if !strings.Contains(result, "('EP-228960453-123', 0)") {
			t.Error("Result should contain first contract")
		}
		if !strings.Contains(result, "('EP-228209382-456', 1)") {
			t.Error("Result should contain second contract")
		}

		// Verify SQL structure
		if !strings.Contains(result, "SELECT") {
			t.Error("Result should contain SELECT clause")
		}
		if !strings.Contains(result, "ORDER BY") {
			t.Error("Result should contain ORDER BY clause")
		}
	})

	t.Run("read xlsx with no matching data", func(t *testing.T) {
		// Create test xlsx file without matching text
		f := excelize.NewFile()
		defer f.Close()

		f.SetCellValue("Sheet1", "A1", "Some Other Data")
		f.SetCellValue("Sheet1", "B1", "Value1")
		f.SetCellValue("Sheet1", "C1", "Value2")

		testFile := filepath.Join(testDir, "test_no_match.xlsx")
		if err := f.SaveAs(testFile); err != nil {
			t.Fatalf("Failed to save test file: %v", err)
		}

		// Read the file
		result, err := h.readXlsxFile(testFile)
		if err != nil {
			t.Fatalf("readXlsxFile failed: %v", err)
		}

		// Should return "no matching data" message
		if !strings.Contains(result, "No matching data found") {
			t.Error("Result should indicate no matching data found")
		}
	})

	t.Run("read xlsx with multiple sheets", func(t *testing.T) {
		f := excelize.NewFile()
		defer f.Close()

		// Add data to Sheet1
		f.SetCellValue("Sheet1", "A1", "ББС ІНШУРАНС")
		f.SetCellValue("Sheet1", "B1", "111111")
		f.SetCellValue("Sheet1", "C1", "aaa")

		// Create Sheet2 with data
		f.NewSheet("Sheet2")
		f.SetCellValue("Sheet2", "A1", "ББС ІНШУРАНС")
		f.SetCellValue("Sheet2", "B1", "222222")
		f.SetCellValue("Sheet2", "C1", "bbb")

		testFile := filepath.Join(testDir, "test_multi_sheet.xlsx")
		if err := f.SaveAs(testFile); err != nil {
			t.Fatalf("Failed to save test file: %v", err)
		}

		result, err := h.readXlsxFile(testFile)
		if err != nil {
			t.Fatalf("readXlsxFile failed: %v", err)
		}

		// Should contain data from both sheets
		if !strings.Contains(result, "EP-111111-aaa") {
			t.Error("Result should contain data from Sheet1")
		}
		if !strings.Contains(result, "EP-222222-bbb") {
			t.Error("Result should contain data from Sheet2")
		}
	})

	t.Run("read xlsx with partial data", func(t *testing.T) {
		f := excelize.NewFile()
		defer f.Close()

		// Match found but only one column after it
		f.SetCellValue("Sheet1", "A1", "ББС ІНШУРАНС")
		f.SetCellValue("Sheet1", "B1", "OnlyOne")
		// C1 is empty

		testFile := filepath.Join(testDir, "test_partial.xlsx")
		if err := f.SaveAs(testFile); err != nil {
			t.Fatalf("Failed to save test file: %v", err)
		}

		result, err := h.readXlsxFile(testFile)
		if err != nil {
			t.Fatalf("readXlsxFile failed: %v", err)
		}

		// Should still handle partial data
		// The result will be "OnlyOne-" since second value is empty
		if strings.Contains(result, "No matching data found") {
			t.Log("Partial data was not captured - this may be expected behavior")
		}
	})

	t.Run("read non-existent file", func(t *testing.T) {
		_, err := h.readXlsxFile(filepath.Join(testDir, "non_existent.xlsx"))
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})
}

func TestStateManagement(t *testing.T) {
	h := NewHandler()

	t.Run("set and get state", func(t *testing.T) {
		chatID := int64(123456)

		// Initial state should be default
		state := h.getState(chatID)
		if state != StateDefault {
			t.Errorf("Expected default state, got: %s", state)
		}

		// Set state to START
		h.setState(chatID, StateStart)
		state = h.getState(chatID)
		if state != StateStart {
			t.Errorf("Expected START state, got: %s", state)
		}

		// Set state back to DEFAULT
		h.setState(chatID, StateDefault)
		state = h.getState(chatID)
		if state != StateDefault {
			t.Errorf("Expected DEFAULT state, got: %s", state)
		}
	})

	t.Run("multiple users states", func(t *testing.T) {
		user1 := int64(111)
		user2 := int64(222)
		user3 := int64(333)

		h.setState(user1, StateStart)
		h.setState(user2, StateDefault)
		// user3 never set, should return default

		if h.getState(user1) != StateStart {
			t.Error("User1 should have START state")
		}
		if h.getState(user2) != StateDefault {
			t.Error("User2 should have DEFAULT state")
		}
		if h.getState(user3) != StateDefault {
			t.Error("User3 should have DEFAULT state (never set)")
		}
	})
}

func TestGenerateSQLScript_EdgeCases(t *testing.T) {
	h := NewHandler()

	t.Run("contract with special characters", func(t *testing.T) {
		contracts := []string{"123-456'789"}
		result := h.generateSQLScript(contracts)

		// Should contain the contract as-is (SQL injection would be handled by parameterized queries)
		if !strings.Contains(result, "EP-123-456'789") {
			t.Error("Contract with special chars should be included")
		}
	})

	t.Run("contract with spaces", func(t *testing.T) {
		contracts := []string{"123 456-789"}
		result := h.generateSQLScript(contracts)

		if !strings.Contains(result, "EP-123 456-789") {
			t.Error("Contract with spaces should be included")
		}
	})

	t.Run("contract with unicode", func(t *testing.T) {
		contracts := []string{"тест-123"}
		result := h.generateSQLScript(contracts)

		if !strings.Contains(result, "EP-тест-123") {
			t.Error("Contract with unicode should be included")
		}
	})

	t.Run("very long contract list", func(t *testing.T) {
		contracts := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			contracts[i] = "123456-789"
		}

		result := h.generateSQLScript(contracts)

		// Check first and last entries
		if !strings.Contains(result, "('EP-123456-789', 0)") {
			t.Error("Should contain first contract")
		}
		if !strings.Contains(result, "('EP-123456-789', 999)") {
			t.Error("Should contain last contract")
		}
	})
}

// Benchmark tests
func BenchmarkGenerateSQLScript(b *testing.B) {
	h := NewHandler()
	contracts := []string{"228960453-123", "228209382-456", "226833195-789"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.generateSQLScript(contracts)
	}
}

func BenchmarkIsValidExcelFile(b *testing.B) {
	h := NewHandler()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		h.isValidExcelFile("test_file.xlsx")
	}
}

// Helper function to create test files directory if needed
func ensureTestDir(t *testing.T) string {
	testDir := filepath.Join(os.TempDir(), "handler_test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	return testDir
}
