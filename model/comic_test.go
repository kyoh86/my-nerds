package model_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kyoh86/my-nerds/model"
)

func TestParseComicFuzzy(t *testing.T) {
	for _, testcase := range []struct {
		title string
		input string
		want  model.Comic
	}{
		{
			title: "append",
			input: "(append1) [作者] タイトル 第3巻 (append2)",
			want: model.Comic{
				Author: "作者",
				Title:  "タイトル",
				Number: 3,
				Append: "append1 append2",
			},
		},
		{
			title: "old order",
			input: "[作者] タイトル 第3巻",
			want: model.Comic{
				Author: "作者",
				Title:  "タイトル",
				Number: 3,
			},
		},
		{
			title: "subtitle",
			input: "タイトル 第3巻 サブ編 [作者]",
			want: model.Comic{
				Author: "作者",
				Title:  "タイトル",
				Volume: "サブ編",
				Number: 3,
			},
		},
		{
			title: "standard",
			input: "タイトル 第3巻 [作者]",
			want: model.Comic{
				Author: "作者",
				Title:  "タイトル",
				Number: 3,
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			got, err := model.ParseComicFuzzy(testcase.input)
			if err != nil {
				t.Fatal(err)
			}
			if got == nil {
				t.Fatal("no response")
			}
			if diff := cmp.Diff(testcase.want, *got); diff != "" {
				t.Errorf("mismatched. -want, +got:\n%s", diff)
			}
		})
	}
}

func TestParseComicName(t *testing.T) {
	for _, testcase := range []struct {
		title string
		input string
		want  model.Comic
	}{
		{
			title: "standard",
			input: "タイトル [作者]",
			want: model.Comic{
				Author: "作者",
				Title:  "タイトル",
			},
		},
		{
			title: "with number",
			input: "タイトル 3巻 [作者]",
			want: model.Comic{
				Author: "作者",
				Title:  "タイトル",
				Number: 3,
			},
		},
		{
			title: "with volume",
			input: "タイトル 3巻 サブ編 [作者]",
			want: model.Comic{
				Author: "作者",
				Title:  "タイトル",
				Number: 3,
				Volume: "サブ編",
			},
		},
	} {
		t.Run(testcase.title, func(t *testing.T) {
			got, err := model.ParseComicName(testcase.input)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(testcase.want, *got); diff != "" {
				t.Errorf("mismatched. -want, +got:\n%s", diff)
			}
		})
	}
}
