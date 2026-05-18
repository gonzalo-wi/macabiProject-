package expensesstorage

import "testing"

func TestResolveStorageURL(t *testing.T) {
	s := &SupabaseSigner{BaseURL: "https://example.supabase.co"}

	tests := []struct {
		in   string
		want string
	}{
		{
			in:   "/object/sign/expense-receipts/projects/p1/expenses/e1/receipt.png?token=abc",
			want: "https://example.supabase.co/storage/v1/object/sign/expense-receipts/projects/p1/expenses/e1/receipt.png?token=abc",
		},
		{
			in:   "https://cdn.example.com/file.png",
			want: "https://cdn.example.com/file.png",
		},
		{
			in:   "/storage/v1/object/public/bucket/x.png",
			want: "https://example.supabase.co/storage/v1/object/public/bucket/x.png",
		},
	}

	for _, tc := range tests {
		got := s.resolveStorageURL(tc.in)
		if got != tc.want {
			t.Fatalf("resolveStorageURL(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
