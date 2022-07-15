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
							start: 7,
							end:   14,
						},
						field{
							name:  "key2",
							value: value{[]byte("null")},
							start: 22,
							end:   25,
						},
						field{
							name:  "key3",
							value: value{[]byte("[1,2]")},
							start: 34,
							end:   39,
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
										start: 14,
										end:   21,
									},
								},
							},
							start: 7,
							end:   22,
						},
					},
				},
			),
		},
		{
			name: "objects key",
			data: []byte(`{"key":[{"key2":"value","key3":[{"key4":"value2"}]}]}`),
			expect: interface{}(
				object{
					fields: []field{
						field{
							name: "key",
							value: objects{
								object{
									fields: []field{
										field{
											name:  "key2",
											value: value{[]byte("\"value\"")},
											start: 16,
											end:   23,
										},
										field{
											name: "key3",
											value: objects{
												object{
													[]field{
														field{
															name:  "key4",
															value: value{[]byte("\"value2\"")},
															start: 40,
															end:   48,
														},
													},
												},
											},
											start: 31,
											end:   50,
										},
									},
								},
							},
							start: 7,
							end:   52,
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

func TestObjectWalkerMarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		obj     ObjectWalker
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
