package d3

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// Skip live tests unless explicitly enabled with RUN_LIVE_TESTS=1 and a real API key
func shouldRunLiveTests() bool {
	runLive := os.Getenv("RUN_LIVE_TESTS")
	apiKey := os.Getenv("D3_API_KEY")
	return runLive == "1" && apiKey != ""
}

func TestClient_LiveAPI_UploadConvertPollDownload(t *testing.T) {
	if !shouldRunLiveTests() {
		t.Skip("Skipping live API tests. Set RUN_LIVE_TESTS=1 and D3_API_KEY to run.")
	}

	apiKey := os.Getenv("D3_API_KEY")
	apiBase := os.Getenv("D3_BASE_URL")
	if apiBase == "" {
		apiBase = "https://api-dev.dragdropdo.com"
	}

	if apiKey == "" {
		t.Fatal("D3_API_KEY is required for live tests")
	}

	// Create client
	client, err := NewDragdropdo(Config{
		APIKey:  apiKey,
		BaseURL: apiBase,
		Timeout: 120 * time.Second,
	})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Create temporary test file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "d3-live-test.txt")
	defer os.Remove(tmpFile)

	content := []byte("hello world")
	if err := os.WriteFile(tmpFile, content, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Log("[live-test] Uploading file...")
	upload, err := client.UploadFile(UploadFileOptions{
		File:     tmpFile,
		FileName: "hello.txt",
		MimeType: "text/plain",
		Parts:    1,
	})
	if err != nil {
		t.Fatalf("Upload failed: %v", err)
	}
	t.Logf("[live-test] Upload result: file_key=%s, upload_id=%s", upload.FileKey, upload.UploadID)

	t.Log("[live-test] Starting convert...")
	operation, err := client.Convert([]string{upload.FileKey}, "png", nil)
	if err != nil {
		t.Fatalf("Convert failed: %v", err)
	}
	t.Logf("[live-test] Operation: main_task_id=%s", operation.MainTaskID)

	t.Log("[live-test] Polling status...")
	status, err := client.PollStatus(PollStatusOptions{
		StatusOptions: StatusOptions{
			MainTaskID: operation.MainTaskID,
		},
		Interval: 3 * time.Second,
		Timeout:  60 * time.Second,
	})
	if err != nil {
		t.Fatalf("Poll status failed: %v", err)
	}
	t.Logf("[live-test] Final status: operation_status=%s", status.OperationStatus)

	if status.OperationStatus != "completed" {
		t.Errorf("Expected operation_status 'completed', got '%s'", status.OperationStatus)
	}

	if len(status.FilesData) == 0 {
		t.Fatal("Expected at least one file in files_data")
	}

	link := status.FilesData[0].DownloadLink
	if link == "" {
		t.Fatal("Expected download_link in files_data[0]")
	}

	t.Log("[live-test] Downloading output...")
	resp, err := http.Get(link)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	t.Logf("[live-test] Downloaded bytes: %d", len(body))
	if len(body) == 0 {
		t.Error("Expected non-empty download")
	}
}

