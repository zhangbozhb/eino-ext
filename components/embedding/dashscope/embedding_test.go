/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dashscope

import (
	"context"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/meguminnnnnnnnn/go-openai"
)

func TestEmbedding(t *testing.T) {
	expectedRequest := openai.EmbeddingRequest{
		Input:          []string{"input"},
		Model:          "mock_model",
		EncodingFormat: openai.EmbeddingEncodingFormatFloat,
		Dimensions:     1024,
	}
	mockResponse := openai.EmbeddingResponse{
		Object: "object",
		Data: []openai.Embedding{
			{
				Embedding: []float32{0.1, 0.2},
			},
			{
				Embedding: []float32{0.3, 0.4},
			},
		},
		Model: "embedding",
		Usage: openai.Usage{
			PromptTokens:     1,
			CompletionTokens: 2,
			TotalTokens:      3,
		},
	}

	t.Run("full param", func(t *testing.T) {
		ctx := context.Background()
		emb, err := NewEmbedder(ctx, &EmbeddingConfig{
			APIKey:  "mock_key",
			Timeout: time.Second * 5,
			Model:   "mock_model",
		})
		if err != nil {
			t.Fatal(err)
		}

		defer mockey.Mock((*openai.Client).CreateEmbeddings).To(func(ctx context.Context, conv openai.EmbeddingRequestConverter) (res openai.EmbeddingResponse, err error) {
			if !reflect.DeepEqual(conv.Convert(), expectedRequest) {
				t.Fatal("openai embedding request is unexpected")
				return
			}
			return mockResponse, nil
		}).Build().UnPatch()

		handler := callbacks.NewHandlerBuilder().
			OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
				nOutput := embedding.ConvCallbackOutput(output)
				if nOutput.TokenUsage.PromptTokens != 1 {
					t.Fatal("PromptTokens is unexpected")
				}
				if nOutput.TokenUsage.CompletionTokens != 2 {
					t.Fatal("CompletionTokens is unexpected")
				}
				if nOutput.TokenUsage.TotalTokens != 3 {
					t.Fatal("TotalTokens is unexpected")
				}
				return ctx
			})

		ctx = callbacks.InitCallbacks(ctx, &callbacks.RunInfo{}, handler.Build())
		result, err := emb.EmbedStrings(ctx, []string{"input"})
		if err != nil {
			t.Fatal(err)
		}
		expectedResult := [][]float64{{0.1, 0.2}, {0.3, 0.4}}
		for i := range result {
			for j := range result[i] {
				if math.Abs(result[i][j]-expectedResult[i][j]) > 1e-7 {
					t.Fatal("result is unexpected")
				}
			}
		}
	})
}
