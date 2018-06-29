package jwalk

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
		expect  interface{}
	}{
		{
			name: "simple object",
			data: []byte(`{"key":"value","key2":null,"key3":[1,2]}`),
			expect: interface{}(
				object{
					fields: []field{
						field{
							name:  "key",
							value: value{[]byte("\"value\"")},
						},
						field{
							name:  "key2",
							value: nil,
						},
						field{
							name:  "key3",
							value: value{[]byte("[1,2]")},
						},
					},
				},
			),
		},
		{
			name: "embed object",
			data: []byte(`{"key":{"key":"value"}}`),
			expect: interface{}(
				object{
					fields: []field{
						field{
							name: "key",
							value: object{
								fields: []field{
									field{
										name:  "key",
										value: value{[]byte("\"value\"")},
									},
								},
							},
						},
					},
				},
			),
		},
		{
			name: "objects key",
			data: []byte(`{"key":[{"key":"value"}]}`),
			expect: interface{}(
				object{
					fields: []field{
						field{
							name: "key",
							value: objects{
								object{
									fields: []field{
										field{
											name:  "key",
											value: value{[]byte("\"value\"")},
										},
									},
								},
							},
						},
					},
				},
			),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.data)

			if tt.wantErr && err == nil {
				assert.NotNil(t, err)
			}

			if !tt.wantErr && err != nil {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.expect, got)
		})
	}
}

func TestObjectIteratorMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		obj     ObjectIterator
		wantErr bool
		expect  []byte
	}{
		{
			name: "simple object",
			obj: object{
				fields: []field{
					field{
						name:  "key",
						value: value{[]byte("\"value\"")},
					},
				},
			},
			expect: []byte(`{"key":"value"}`),
		},
		{
			name: "embed object",
			obj: object{
				fields: []field{
					field{
						name: "key",
						value: object{
							fields: []field{
								field{
									name: "key",
									value: objects{
										object{
											fields: []field{
												field{
													name:  "key",
													value: value{[]byte("\"value\"")},
												},
											},
										},
										object{
											fields: []field{
												field{
													name:  "key",
													value: value{[]byte("\"value\"")},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expect: []byte(`{"key":{"key":[{"key":"value"},{"key":"value"}]}}`),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(tt.obj)

			if tt.wantErr && err == nil {
				assert.NotNil(t, err)
			}

			if !tt.wantErr && err != nil {
				assert.Nil(t, err)
			}

			assert.Equal(t, string(tt.expect), string(got))
		})
	}
}
